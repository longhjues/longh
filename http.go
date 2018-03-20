package longh

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// HTTPDo 请求URL并且解析JSON格式的返回数据
func HTTPDo(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// HTTPGet 请求URL
func HTTPGet(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// HTTPGetValue 请求URL 附带参数
func HTTPGetValue(URL string, params url.Values) ([]byte, error) {
	resp, err := http.Get(fmt.Sprint(URL, "?", params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// HTTPPost 请求URL
func HTTPPost(URL string, params url.Values) ([]byte, error) {
	resp, err := http.PostForm(URL, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// HTTPPostToJSON POST请求 BODY为json格式
func HTTPPostToJSON(URL string, v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// HTTPCheckNotNil 检查HTTP参数是否为空
func HTTPCheckNotNil(r *http.Request, args ...string) error {
	if args == nil || r == nil {
		return nil
	}

	switch r.Method {
	case "GET":
		query := r.URL.Query()
		for _, v := range args {
			if strings.TrimSpace(query.Get(v)) == "" {
				return fmt.Errorf("param(%s) is invalid", v)
			}
		}
	case "POST":
		for _, v := range args {
			if strings.TrimSpace(r.PostFormValue(v)) == "" {
				return fmt.Errorf("param(%s) is invalid", v)
			}
		}
	default:
		return errors.New("r.Method is not GET or POST")
	}
	return nil
}

// IsStringEmpty 判断是否有值为空或null或(null)
func IsStringEmpty(s ...string) bool {
	var str string
	for _, v := range s {
		str = strings.TrimSpace(v)
		if v == "" || strings.EqualFold(str, "(null)") || strings.EqualFold(str, "null") {
			return true
		}
	}
	return false
}

// HTTPWriteToJSON 写入json字符串
func HTTPWriteToJSON(w io.Writer, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
