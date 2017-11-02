package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type kokoro struct {
	mutex *sync.Mutex

	w io.Writer

	Insecure bool

	Host string

	wsConn *websocket.Conn

	AccessToken string

	Logger *log.Logger
}

type requestMessage struct {
	JSONRPC string      `json:"jsonrpc"`
	Id      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type httpRequestParams struct {
	Type    string            `json:"type"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Data    interface{}       `json:"data"`
	Timeout time.Duration     `json:"timeout"`
}

type websocketRequestParams interface{}

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
	go func() {
		if err := self.startWS(c); err != nil {
			self.Logger.Errorf("WebSocket start failed: %s", err)
		}
	}()
	go func(r *bufio.Reader) {
		// read line as jsonlines
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				if err == io.EOF {
					errCh <- nil
				} else {
					errCh <- err
				}
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

			rawParams, err := json.Marshal(msg.Params)
			if err != nil {
				self.Logger.Errorf("Can't marshal json: %s", err)
				continue
			}

			var (
				successData interface{}
				errorData   *errorObject
			)
			switch msg.Method {
			case "http":
				params, err := self.UnmarshalRequestParamsForREST(rawParams)
				if err != nil {
					self.Logger.Errorf("Can't unmarshal params: %s", err)
					continue
				}
				msg.Params = params
				if successData, errorData, err = self.TranslateToREST(msg); err != nil {
					self.Logger.Errorf("REST API Error: %s", err)
				}
			case "websocket":
				params, err := self.UnmarshalRequestParamsForWebSocket(rawParams)
				if err != nil {
					self.Logger.Errorf("Can't unmarshal params: %s", err)
					continue
				}
				msg.Params = params
				if successData, errorData, err = self.TranslateToWebSocket(msg); err != nil {
					self.Logger.Errorf("WebSocket Error: %s", err)
				}
			}
			// no reply when result and error are empty
			if successData != nil || errorData != nil {
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
		}
	}(bufio.NewReader(r))
	select {
	case <-c.Done():
		return c.Err()
	case err := <-errCh:
		return err
	}
}

func (self *kokoro) startWS(c context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		wsEndpoint := new(url.URL)
		if self.Insecure {
			wsEndpoint.Scheme = "ws"
		} else {
			wsEndpoint.Scheme = "wss"
		}
		wsEndpoint.Host = self.Host
		wsEndpoint.Path = "/cable"

		self.Logger.Infof("Connecting to %s", wsEndpoint.String())
		wsHeader := make(http.Header)
		wsHeader.Set("X-Access-Token", self.AccessToken)
		if conn, wsResp, err := websocket.DefaultDialer.Dial(wsEndpoint.String(), wsHeader); err != nil {
			if wsResp != nil {
				defer wsResp.Body.Close()
				self.Logger.Errorf("WS Response: %v - %v", wsResp.StatusCode, wsResp.Status)
				for k, v := range wsResp.Header {
					self.Logger.Errorf("WS Response Header: %v: %v", k, v)
				}
				body, _ := ioutil.ReadAll(wsResp.Body)
				self.Logger.Errorf("WS Response: %v", string(body))
			}
			errCh <- fmt.Errorf("Can't upgrade websocket connection: %s", err)
			return
		} else {
			defer conn.Close()
			self.wsConn = conn
		}

		for {
			mt, r, err := self.wsConn.NextReader()
			if err != nil {
				errCh <- err
				break
			}
			if mt != websocket.TextMessage {
				self.Logger.Infof("Non text message frame, ignored")
				continue
			}
			b, err := ioutil.ReadAll(r)
			if err != nil {
				self.Logger.Errorf("Can't read websocket: %s", err)
				continue
			}
			self.Logger.Debugf("websocket = %v", string(b))
			data := make(map[string]interface{})
			if err = json.Unmarshal(b, &data); err != nil {
				self.Logger.Errorf("Can't unmarshal websocket message: %s", err)
				continue
			}
			switch s := data["type"].(type) {
			case string:
				if s == "ping" {
					self.Logger.Debugf("Ignore ping message: %v", data)
					continue
				}
			}
			notif := &response{
				JSONRPC: "2.0",
				Id:      "__websocket__",
				Result:  data,
			}
			if err = self.WriteJSON(notif); err != nil {
				self.Logger.Errorf("Can't write websocket output: %s", err)
			}
		}
	}()
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

func (self *kokoro) UnmarshalRequestParamsForREST(b []byte) (interface{}, error) {
	params := new(httpRequestParams)
	return params, json.Unmarshal(b, params)
}

func (self *kokoro) TranslateToREST(msg *requestMessage) (interface{}, *errorObject, error) {
	client := new(http.Client)
	params, ok := msg.Params.(*httpRequestParams)
	if !ok {
		self.Logger.Panicf("Type mismatch: wants %v, but %v", reflect.TypeOf(new(httpRequestParams)), reflect.TypeOf(msg.Params))
	}

	requestMethod := strings.ToUpper(params.Type)
	if requestMethod == "" {
		requestMethod = http.MethodGet
	}

	requestUrl, err := url.ParseRequestURI(params.Url)
	if err != nil {
		return nil, nil, fmt.Errorf("Can't parse request url: %s", err)
	}
	if requestUrl.Scheme == "" {
		if self.Insecure {
			requestUrl.Scheme = "http"
		} else {
			requestUrl.Scheme = "https"
		}
	}
	if requestUrl.Host == "" {
		requestUrl.Host = self.Host
	}
	var requestBody io.Reader
	if requestMethod == "POST" {
		data, err := json.Marshal(params.Data)
		if err != nil {
			return nil, nil, fmt.Errorf("Can't create request payload: %s", err)
		}
		requestBody = bytes.NewReader(data)
	} else {
		if queryParams, ok := params.Data.(map[string]interface{}); ok {
			for k, v := range queryParams {
				requestUrl.Query().Set(k, fmt.Sprintf("%v", v))
			}
		}
	}

	req, err := http.NewRequest(requestMethod, requestUrl.String(), requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("Can't create REST API Request: %s", err)
	}
	if params.Headers != nil {
		for k, v := range params.Headers {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("X-Access-Token", self.AccessToken)

	self.Logger.Debugf("Send http request with timeout %v: %v", params.Timeout, req)
	client.Timeout = params.Timeout
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

func (self *kokoro) UnmarshalRequestParamsForWebSocket(b []byte) (interface{}, error) {
	params := new(websocketRequestParams)
	return params, json.Unmarshal(b, params)
}

func (self *kokoro) TranslateToWebSocket(msg *requestMessage) (interface{}, *errorObject, error) {
	if self.wsConn == nil {
		self.Logger.Panicf("WebSocket connection is nil")
	}
	params, ok := msg.Params.(websocketRequestParams)
	if !ok {
		self.Logger.Panicf("Type mismatch: wants %v, but %v", reflect.TypeOf(new(websocketRequestParams)), reflect.TypeOf(msg.Params))
	}

	w, err := self.wsConn.NextWriter(websocket.TextMessage)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create ws writer: %s", err)
	}
	defer w.Close()

	self.Logger.Debugf("Send data via websocket: %v", params)
	if err = json.NewEncoder(w).Encode(params); err != nil {
		return nil, nil, fmt.Errorf("Can't marshal json: %s", err)
	}
	return nil, nil, nil
}
