package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func LaunchHttpRequest(method, url string, params map[string]string, body interface{}, headerContentType, contentType string) (map[string]interface{}, error) {
	//logrus.Info("LaunchHttpRequest url: ", url)
	request, err := newRequest(url, method, params, body, headerContentType, contentType)
	if err != nil {
		return
	}

	return parseHttpResponseDefault(http.DefaultClient.Do(request))
}

func parseHttpResponseDefault(resp *http.Response) (map[string]interface{}, error) {
	var response map[string]interface{}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = json.Unmarshal(b, &response)
		return response, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(b))
	}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return response, fmt.Errorf("error: %v\n response in bytes: %v\n", err, b)
	}
	return response, err
}

func newRequest(url, method string, params map[string]string, body interface{}, headerContentType, contentType string) (req *http.Request, err error) {
	var tempBody []byte

	if tempBody, err = json.Marshal(body); err != nil {
		return nil, err
	}

	req, err = http.NewRequest(method, url, bytes.NewBuffer(tempBody))
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set(headerContentType, contentType)

	return req, nil
}
