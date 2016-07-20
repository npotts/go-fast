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
)

var matchjs = regexp.MustCompile("/app-[0-9a-f]{6}.js")
var matchToken = regexp.MustCompile(`token:"[a-zA-Z]{32}"`)

//script returns the js script that contains the token
func (gofast) script() (jscript string, err error) {
	var resp *http.Response
	resp, err = http.Get("https://www.fast.com")
	defer resp.Body.Close()
	if err != nil {
		return
	}
	err = fmt.Errorf("Could not find script")
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

func (gf *gofast) getToken() (token string, err error) {
	err = fmt.Errorf("Could not find token")
	if token, err = gf.script(); err != nil {
		return
	}
	if resp, e := http.Get("https://fast.com" + token); e == nil {
		defer resp.Body.Close()
		if body, e := ioutil.ReadAll(resp.Body); e == nil {
			if matchToken.Match(body) {
				return string(matchToken.Find(body)[7:39]), nil
			}
		}
	}
	return
}

func (gf *gofast) getURLs(count int) (urls []string, err error) {
	err = fmt.Errorf("Unable to get URLs")
	if token, err := gf.getToken(); err == nil {
		url := fmt.Sprintf("http://api.fast.com/netflix/speedtest?https=true&token=%s&urlCount=%d", token, count)

		if resp, err := http.Get(url); err == nil {
			defer resp.Body.Close()
			if body, e := ioutil.ReadAll(resp.Body); e == nil {
				var v []map[string]string
				if err = json.Unmarshal(body, &v); err == nil {
					for _, m := range v {
						urls = append(urls, m["url"])
					}
					return urls, nil
				}
			}
		}
	}
	return
}
