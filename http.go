package library

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	WARN = iota
	ERROR
	DEBUG
	INFO
)

type RespParams struct {
	Severity    int
	URL         string
	Section     string
	ErrorCode   string
	Description string
	Error       error
	Input       interface{}
}

type NetAdaptor struct {
	Client HttpClient
}

type HTTPResponse struct {
	Status      bool   `json:"status"`
	ErrorCode   string `json:"error_code"`
	Description string `json:"description"`
	Token       string `json:"token"`
	Data        string `json:"data"`
}

type HttpClient struct {
	Client *http.Client
}

func GoodResponse(c *gin.Context, data interface{}) {
	returnData, _ := json.Marshal(data)
	response := HTTPResponse{
		Token:  c.GetString("token"),
		Status: true,
		Data:   string(returnData),
	}
	c.JSON(http.StatusOK, response)
}

func BadResponse(LOG *zap.Logger, c *gin.Context, rp RespParams) {
	switch rp.Severity {
	case DEBUG:
		LOG.Debug(rp.Section,
			zap.String("connection", rp.URL),
			zap.Any("parameters", rp.Input),
			zap.String("description", rp.Description),
			zap.Error(rp.Error))
	case WARN:
		LOG.Warn(rp.Section,
			zap.String("connection", rp.URL),
			zap.Any("parameters", rp.Input),
			zap.String("description", rp.Description),
			zap.Error(rp.Error))
	case ERROR:
		LOG.Error(rp.Section,
			zap.String("connection", rp.URL),
			zap.Any("parameters", rp.Input),
			zap.String("description", rp.Description),
			zap.Error(rp.Error))
	}
	response := HTTPResponse{
		Status:      false,
		Description: rp.Description,
		ErrorCode:   rp.ErrorCode,
	}

	c.JSON(http.StatusBadRequest, response)
}

func ExtractBody(c *gin.Context) ([]byte, error) {
	var bodyBytes []byte
	var err error
	if c.Request.Body != nil {
		bodyBytes, err = ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	// Restore the io.ReadCloser to its original state
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes, nil
}

func (adaptor NetAdaptor) GET(log *zap.Logger, token, uri string, data interface{}) (HTTPResponse, error) {
	handleErr := func(err error) (HTTPResponse, error) {
		return HTTPResponse{}, fmt.Errorf("get %s  : %w", uri, err)
	}

	baseUrl, err := url.Parse(uri)
	if err != nil {
		return handleErr(err)
	}

	// Add a Path Segment (Path segment is automatically escaped)
	params, err := StructToUrlValue(data)
	if err != nil {
		return handleErr(err)
	}

	// Add Query Parameters to the URL
	baseUrl.RawQuery = params.Encode() // Escape Query Parameters

	log.Debug("http request",
		zap.String("method", "GET"),
		zap.String("url", baseUrl.String()))

	result, err := adaptor.Client.GET(makeHeaders(token), baseUrl.String())
	if err != nil {
		return handleErr(fmt.Errorf("http process (%w)", err))
	}

	rr := HTTPResponse{}
	if err := json.Unmarshal(result, &rr); err != nil {
		return handleErr(fmt.Errorf("unmarshal response , %s (%w)", string(result), err))
	}

	return rr, nil
}

func (adaptor NetAdaptor) POST(log *zap.Logger, token, url string, data interface{}) (HTTPResponse, error) {
	handleErr := func(err error) (HTTPResponse, error) {
		return HTTPResponse{}, fmt.Errorf("post %s  : %w", url, err)
	}

	message, err := json.Marshal(data)
	if err != nil {
		return handleErr(err)
	}

	log.Debug("http request",
		zap.String("method", "POST"),
		zap.String("url", url),
		zap.Any("data", data))

	result, err := adaptor.Client.POST(makeHeaders(token), url, message)
	if err != nil {
		return handleErr(fmt.Errorf("http process (%w)", err))
	}

	rr := HTTPResponse{}
	if err := json.Unmarshal(result, &rr); err != nil {
		return handleErr(fmt.Errorf("unmarshal response , %s (%w)", string(result), err))
	}

	return rr, nil
}

func (adaptor NetAdaptor) EXTPOST(log *zap.Logger, token, url string, data interface{}) ([]byte, error) {
	handleErr := func(err error) ([]byte, error) {
		return nil, fmt.Errorf("post %s  : %w", url, err)
	}

	message, err := json.Marshal(data)
	if err != nil {
		return handleErr(err)
	}

	log.Debug("http request",
		zap.String("method", "POST"),
		zap.String("url", url),
		zap.Any("data", data))

	return adaptor.Client.POST(makeHeaders(token), url, message)
}
func (adaptor NetAdaptor) PUT(log *zap.Logger, token, url string, data interface{}) (HTTPResponse, error) {
	handleErr := func(err error) (HTTPResponse, error) {
		return HTTPResponse{}, fmt.Errorf("put %s  : %w", url, err)
	}

	message, err := json.Marshal(data)
	if err != nil {
		return handleErr(err)
	}

	log.Debug("http request",
		zap.String("method", "PUT"),
		zap.String("url", url),
		zap.Any("data", data))

	result, err := adaptor.Client.PUT(makeHeaders(token), url, message)
	if err != nil {
		return handleErr(fmt.Errorf("http process (%w)", err))
	}

	rr := HTTPResponse{}
	if err := json.Unmarshal(result, &rr); err != nil {
		return handleErr(fmt.Errorf("unmarshal response , %s (%w)", string(result), err))
	}

	return rr, nil
}

func (adaptor NetAdaptor) DELETE(log *zap.Logger, token, url string, data interface{}) (HTTPResponse, error) {
	handleErr := func(err error) (HTTPResponse, error) {
		return HTTPResponse{}, fmt.Errorf("delete %s  : %w", url, err)
	}

	message, err := json.Marshal(data)
	if err != nil {
		return handleErr(err)
	}

	log.Debug("http request",
		zap.String("method", "DELETE"),
		zap.String("url", url),
		zap.Any("data", data))

	result, err := adaptor.Client.DELETE(makeHeaders(token), url, message)
	if err != nil {
		return handleErr(fmt.Errorf("http process (%w)", err))
	}

	rr := HTTPResponse{}
	if err := json.Unmarshal(result, &rr); err != nil {
		return handleErr(fmt.Errorf("unmarshal response , %s (%w)", string(result), err))
	}

	return rr, nil
}

func (hc HttpClient) POST(header http.Header, url string, load []byte) ([]byte, error) {
	handleErr := func(err error) ([]byte, error) {
		return nil, err
	}

	load, err := request(hc.Client, url, "POST", header, load)
	if err != nil {
		return handleErr(err)
	}

	return load, nil
}

func (hc HttpClient) PUT(header http.Header, url string, load []byte) ([]byte, error) {
	handleErr := func(err error) ([]byte, error) {
		return nil, err
	}

	load, err := request(hc.Client, url, "PUT", header, load)
	if err != nil {
		return handleErr(err)
	}

	return load, nil
}

func (hc HttpClient) DELETE(header http.Header, url string, load []byte) ([]byte, error) {
	handleErr := func(err error) ([]byte, error) {
		return nil, err
	}

	load, err := request(hc.Client, url, "DELETE", header, load)
	if err != nil {
		return handleErr(err)
	}

	return load, nil
}

func request(client *http.Client, url, rtype string, headers http.Header, load []byte) ([]byte, error) {
	handleErr := func(err error) ([]byte, error) {
		return nil, fmt.Errorf("http call : %w", err)
	}

	reqBody := bytes.NewBuffer(load)

	req, err := http.NewRequest(rtype, url, reqBody)
	if err != nil {
		return handleErr(fmt.Errorf("prepare request (%w)", err))
	}

	req.Close = true
	req.Header = headers

	resp, err := client.Do(req)
	if err != nil {
		return handleErr(fmt.Errorf("do request (%w)", err))
	}

	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return handleErr(fmt.Errorf("read response (%w)", err))
	}

	return response, nil
}

func makeHeaders(token string) http.Header {
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	h.Set("Authorization", token)
	return h
}

func HTTPRequest(reqType string, headers http.Header, url string, load []byte) ([]byte, error) {
	errHandle := func(err error) ([]byte, error) {
		return nil, err
	}

	var request *http.Request
	var err error

	if load == nil {
		request, err = http.NewRequest(reqType, url, nil)
		if err != nil {
			return errHandle(err)
		}
	} else {
		request, err = http.NewRequest(reqType, url, bytes.NewBuffer(load))
		if err != nil {
			return errHandle(err)
		}
	}

	request.Header = headers
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return errHandle(err)
	}
	defer resp.Body.Close()

	//unmarshal load
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errHandle(err)
	}

	return body, nil
}

func HttpGet(net NetAdaptor, log *zap.Logger, url, token string, input interface{}) (HTTPResponse, error) {
	errHandle := func(err error) (HTTPResponse, error) {
		return HTTPResponse{}, err
	}

	result, err := net.GET(log, token, url, input)
	if err != nil {
		return errHandle(fmt.Errorf("request to service: %w", err))
	}
	if !result.Status {
		return errHandle(fmt.Errorf("fail request: %s", result.Description))
	}

	return result, nil
}

func HttpPost(net NetAdaptor, log *zap.Logger, url, token string, input interface{}) (HTTPResponse, error) {
	errHandle := func(err error) (HTTPResponse, error) {
		return HTTPResponse{}, err
	}

	result, err := net.POST(log, token, url, input)
	if err != nil {
		return errHandle(fmt.Errorf("request to service: %w", err))
	}
	if !result.Status {
		return errHandle(fmt.Errorf("fail request: %s", result.Description))
	}

	return result, nil
}

func HttpPut(net NetAdaptor, log *zap.Logger, url, token string, input interface{}) (HTTPResponse, error) {
	errHandle := func(err error) (HTTPResponse, error) {
		return HTTPResponse{}, err
	}

	result, err := net.PUT(log, token, url, input)
	if err != nil {
		return errHandle(fmt.Errorf("request to service: %w", err))
	}
	if !result.Status {
		return errHandle(fmt.Errorf("fail request: %s", result.Description))
	}

	return result, nil
}

func HttpDelete(net NetAdaptor, log *zap.Logger, url, token string, input interface{}) (HTTPResponse, error) {
	errHandle := func(err error) (HTTPResponse, error) {
		return HTTPResponse{}, err
	}

	result, err := net.DELETE(log, token, url, input)
	if err != nil {
		return errHandle(fmt.Errorf("request to service: %w", err))
	}
	if !result.Status {
		return errHandle(fmt.Errorf("fail request: %s", result.Description))
	}

	return result, nil
}
