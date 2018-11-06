package utils

import (
	"io"
	"encoding/json"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"time"
	"errors"
)

var (
	httpClient = &http.Client {
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     true,
		},
	}
	errNotFound = errors.New("resource not found")
)

func HttpGet(url string, out interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("NewRequest: %s", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %s", req.URL.String(), err)
	}
	defer resp.Body.Close()

	var cnt bytes.Buffer
	_, err = io.Copy(&cnt, resp.Body)
	if err != nil {
		return fmt.Errorf("Copy: %s", err)
	}

	if resp.StatusCode == 404 {
		return errNotFound
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("%s: status code=%d, body=%s", req.URL.String(), resp.StatusCode, cnt.String())
	}

	fmt.Println(string(cnt.Bytes()))

	if err := json.Unmarshal(cnt.Bytes(), out); err != nil {
		return fmt.Errorf("Unmarshal: %s", err)
	}

	return nil
}