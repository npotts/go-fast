/*
Copyright (c) 2016 Nicholas Potts

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the "Software"), to deal in the Software without restriction, including without
limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT
LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package gofast

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var matchjs = regexp.MustCompile("/app-[0-9a-f]{6}.js")
var matchToken = regexp.MustCompile(`token:"[a-zA-Z]{32}"`)

func (gf gofast) timeoutGet(url string) (resp *http.Response, err error) {
	client := http.Client{Timeout: gf.cfg.Network}
	ch := make(chan bool)
	go func() {
		resp, err = client.Get(url)
		ch <- true
	}()
	select {
	case <-ch: //read
		return
	case <-time.After(3 * time.Second):
		return nil, fmt.Errorf("Unable to get data within timeout")
	}
	return
}

//script returns the js script that contains the token
func (gf gofast) script() (jscript string, err error) {
	var resp *http.Response
	if resp, err = gf.timeoutGet("https://www.fast.com"); err != nil {
		err = errors.Wrap(err, "Unable to retrieve HTML to extract script")
		return
	}
	err = fmt.Errorf("Could not find script")
	defer resp.Body.Close()
	izer := html.NewTokenizer(resp.Body)
	for {
		if izer.Next() == html.StartTagToken {
			if tok := izer.Token(); tok.Data == "script" {
				for _, src := range tok.Attr {
					if src.Key == "src" && matchjs.MatchString(src.Val) {
						return src.Val, nil
					}
				}
			}
		}
	}
	return
}

func (gf gofast) getToken() (token string, err error) {
	var resp *http.Response
	body := []byte{}
	err = fmt.Errorf("Could not find token")
	if token, err = gf.script(); err != nil {
		err = errors.Wrap(err, "Unable to get a token")
		return
	}
	if resp, err = gf.timeoutGet("https://fast.com" + token); err != nil {
		err = errors.Wrap(err, "Unable to retrieve from fast.com with a token")
		return
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		err = errors.Wrap(err, "Unable to Readall")
		return
	}
	if matchToken.Match(body) {
		return string(matchToken.Find(body)[7:39]), nil
	}
	return
}

func (gf *gofast) getURLs(count int) (urls []string, err error) {
	token := ""
	err = fmt.Errorf("Unable to get URLs")
	if token, err = gf.getToken(); err != nil {
		err = errors.Wrap(err, "Unable to retrieve URLs without a token")
		return
	}
	murls := map[string]string{}
	url := fmt.Sprintf("http://api.fast.com/netflix/speedtest?https=true&token=%s&urlCount=%d", token, count)
	enought := false
	for {
		var resp *http.Response
		body := []byte{}
		if resp, err = gf.timeoutGet(url); err != nil {
			err = errors.Wrap(err, "Cannot get JSON body")
			return
		}
		defer resp.Body.Close()
		if body, err = ioutil.ReadAll(resp.Body); err != nil {
			err = errors.Wrap(err, "ReadAll on body failed")
			return
		}
		var v []map[string]string
		if err = json.Unmarshal(body, &v); err != nil {
			err = errors.Wrapf(err, "Unable to unmarshal %q", body)
			return
		}
		for _, m := range v {
			enought = (len(murls) >= count)
			if !enought {
				murls[m["url"]] = "" //Make sure URLs are unique
			}
		}
		if enought {
			break
		}
	}
	for url := range murls {
		urls = append(urls, url)
	}
	return urls, nil
}
