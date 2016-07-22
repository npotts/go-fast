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
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var matchjs = regexp.MustCompile("/app-[0-9a-f]{6}.js")
var matchToken = regexp.MustCompile(`token:"[a-zA-Z]{32}"`)

func (gofast) getWithTimeout(url string, timeout time.Duration) (*http.Response, error) {
	client := http.Client{Timeout: timeout}
	return client.Get(url)
}

//script returns the js script that contains the token
func (gf gofast) script() (jscript string, err error) {
	var resp *http.Response
	if resp, err = gf.getWithTimeout("https://www.fast.com", gf.cfg.Network); err != nil {
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
		return
	}
	if resp, err = http.Get("https://fast.com" + token); err != nil {
		return
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
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
		return
	}
	url := fmt.Sprintf("http://api.fast.com/netflix/speedtest?https=true&token=%s&urlCount=%d", token, count)
	for {
		var resp *http.Response
		body := []byte{}
		if resp, err = http.Get(url); err != nil {
			return
		}
		defer resp.Body.Close()
		if body, err = ioutil.ReadAll(resp.Body); err != nil {
			return
		}
		var v []map[string]string
		if err = json.Unmarshal(body, &v); err != nil {
			return
		}
		for _, m := range v {
			if len(urls) >= count {
				return urls, nil
			}
			urls = append(urls, m["url"])
		}
	}
}
