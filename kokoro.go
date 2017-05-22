package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type kokoro struct {
	mutex *sync.Mutex

	w io.Writer

	BaseUrl string

	Logger *log.Logger
}

type requestMessage struct {
	JSONRPC string         `json:"jsonrpc"`
	Id      interface{}    `json:"id"`
	Method  string         `json:"method"`
	Params  *requestParams `json:"params"`
}

type requestParams struct {
	Type    string            `json:"type"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Data    interface{}       `json:"data"`
	Timeout time.Duration     `json:"timeout"`
}

type response struct {
	JSONRPC string      `json:"jsonrpc"`
	Id      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type errorObject struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (self *kokoro) Start(c context.Context, r io.Reader) error {
	errCh := make(chan error, 1)
	go func(r *bufio.Reader) {
		// read line as jsonlines
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				errCh <- err
				break
			}
			msg := new(requestMessage)
			if err := json.Unmarshal(line, msg); err != nil {
				self.Logger.Errorf("Can't parse request message: %s", err)
				continue
			}

			// verify jsonrpc request
			if msg.JSONRPC != "2.0" {
				self.Logger.Warnf("Spec error, the \"jsonrpc\" key wants \"2.0\", but \"%v\"", msg.JSONRPC)
				continue
			}
			if msg.Id == "" {
				self.Logger.Warn("Spec error, ths \"id\" key must be specified")
				continue
			}

			var (
				successData interface{}
				errorData   *errorObject
			)
			switch msg.Method {
			case "http":
				if successData, errorData, err = self.TranslateToREST(msg); err != nil {
					self.Logger.Errorf("Error: %s", err)
				}
			case "websocket":
				self.Logger.Warn("Sorry, not implemented yet")
			}
			reply := &response{
				JSONRPC: msg.JSONRPC,
				Id:      msg.Id,
			}
			if errorData != nil {
				reply.Error = errorData
			} else {
				reply.Result = successData
			}
			if err = self.WriteJSON(reply); err != nil {
				self.Logger.Errorf("Can't write output: %s", err)
			}
		}
	}(bufio.NewReader(r))
	select {
	case <-c.Done():
		return c.Err()
	case err := <-errCh:
		return err
	}
}

func (self *kokoro) WriteJSON(v interface{}) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	return json.NewEncoder(self.w).Encode(v)
}

func (self *kokoro) TranslateToREST(msg *requestMessage) (interface{}, *errorObject, error) {
	client := new(http.Client)

	requestUrl, err := url.ParseRequestURI(msg.Params.Url)
	if err != nil {
		return nil, nil, fmt.Errorf("Can't parse request url: %s", err)
	}
	var requestBody io.Reader
	if strings.ToUpper(msg.Params.Type) == "POST" {
		data, err := json.Marshal(msg.Params.Data)
		if err != nil {
			return nil, nil, fmt.Errorf("Can't create request payload: %s", err)
		}
		requestBody = bytes.NewReader(data)
	} else {
		if queryParams, ok := msg.Params.Data.(map[string]interface{}); ok {
			for k, v := range queryParams {
				requestUrl.Query().Set(k, fmt.Sprintf("%v", v))
			}
		}
	}

	req, err := http.NewRequest(msg.Params.Type, requestUrl.String(), requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("Can't create REST API Request: %s", err)
	}
	if msg.Params.Headers != nil {
		for k, v := range msg.Params.Headers {
			req.Header.Set(k, v)
		}
	}

	client.Timeout = msg.Params.Timeout
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("Error: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("Can't read response body: %s", err)
	}

	var (
		successData map[string]interface{} = nil
		errorData   *errorObject           = nil
	)
	if resp.StatusCode < 300 {
		successData = map[string]interface{}{
			"status": resp.StatusCode,
			"body":   string(body),
		}
	} else {
		errorData = &errorObject{
			Code:    resp.StatusCode,
			Message: resp.Status,
			Data:    string(body),
		}
	}
	return successData, errorData, nil
}
