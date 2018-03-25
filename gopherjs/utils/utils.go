//build +js

package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"

	fmt "github.com/cathalgarvey/fmtless"
	"github.com/cathalgarvey/fmtless/encoding/json"
	"github.com/cathalgarvey/fmtless/net/url"
	"github.com/gopherjs/gopherjs/js"
	"github.com/shutej/protobuf/jsonpb"
	"github.com/shutej/protobuf/proto"
	"honnef.co/go/js/xhr"

	"github.com/shutej/twirp/gopherjs/header"
	"github.com/shutej/twirp/gopherjs/twirp"
)

// XHRClient is the interface used by generated clients to send HTTP requests.
//
// TODO(shutej): XHRClient implementations should not follow redirects.
// Package xhr does, however, which needs to be fixed.
type XHRClient interface {
	Do(req *XHRRequest) (*XHRResponse, error)
}

type xhrClient struct{}

func NewXHRClient() XHRClient {
	return &xhrClient{}
}

type XHRRequest struct {
	Method string
	URL    *url.URL
	Header header.Header
	Body   io.ReadCloser
}

type XHRResponse struct {
	StatusCode int
	StatusText string
	Body       io.ReadCloser
}

func NewXHRRequest(method, url_ string, body io.Reader) (req *XHRRequest, err error) {
	var u *url.URL
	if u, err = url.Parse(url_); err != nil {
		return
	}
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	req = &XHRRequest{
		Method: method,
		URL:    u,
		Header: nil,
		Body:   rc,
	}
	return
}

func (_ xhrClient) Do(request *XHRRequest) (response *XHRResponse, err error) {
	xhrRequest := xhr.NewRequest(request.Method, request.URL.String())
	xhrRequest.ResponseType = xhr.ArrayBuffer
	var requestData []byte
	if requestData, err = ioutil.ReadAll(request.Body); err != nil {
		return
	}
	if err = xhrRequest.Send(requestData); err != nil {
		return
	}

	responseData := bytes.NewBuffer(js.Global.Get("Uint8Array").New(xhrRequest.Response).Interface().([]byte))
	response = &XHRResponse{
		StatusCode: xhrRequest.Status,
		StatusText: fmt.Sprintf("%d %s", xhrRequest.Status, xhrRequest.StatusText),
		Body:       ioutil.NopCloser(responseData),
	}
	return
}

// URLBase helps ensure that addr specifies a scheme. If it is unparsable
// as a URL, it returns addr unchanged.
func URLBase(addr string) string {
	// If the addr specifies a scheme, use it. If not, default to
	// http. If url.Parse fails on it, return it unchanged.
	url, err := url.Parse(addr)
	if err != nil {
		return addr
	}
	if url.Scheme == "" {
		url.Scheme = "http"
	}
	return url.String()
}

// newRequest makes an XHRRequest from a client, adding common headers.
func newRequest(url string, reqBody io.Reader, contentType string) (*XHRRequest, error) {
	req, err := NewXHRRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", contentType)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Twirp-Version", "v5.3.0")
	return req, nil
}

// JSON serialization for errors
type twerrJSON struct {
	Code string            `json:"code"`
	Msg  string            `json:"msg"`
	Meta map[string]string `json:"meta,omitempty"`
}

// marshalErrorToJSON returns JSON from a twirp.Error, that can be used as HTTP error response body.
// If serialization fails, it will use a descriptive Internal error instead.
func marshalErrorToJSON(twerr twirp.Error) []byte {
	// make sure that msg is not too large
	msg := twerr.Msg()
	if len(msg) > 1e6 {
		msg = msg[:1e6]
	}

	tj := twerrJSON{
		Code: string(twerr.Code()),
		Msg:  msg,
		Meta: twerr.MetaMap(),
	}

	buf, err := json.Marshal(&tj)
	if err != nil {
		buf = []byte("{\"type\": \"" + twirp.Internal + "\", \"msg\": \"There was an error but it could not be serialized into JSON\"}") // fallback
	}

	return buf
}

// errorFromResponse builds a twirp.Error from a non-200 HTTP response.
// If the response has a valid serialized Twirp error, then it's returned.
// If not, the response status code is used to generate a similar twirp
// error. See twirpErrorFromIntermediary for more info on intermediary errors.
func errorFromResponse(resp *XHRResponse) twirp.Error {
	// TODO(shutej): This used to deal with redirects, but currently ignores them.
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return clientError("failed to read server error response body", err)
	}
	var tj twerrJSON
	if err := json.Unmarshal(respBodyBytes, &tj); err != nil {
		// Invalid JSON response; it must be an error from an intermediary.
		msg := fmt.Sprintf("Error from intermediary with HTTP status code %d %q", resp.StatusCode, resp.StatusText)
		return twirpErrorFromIntermediary(resp.StatusCode, msg, string(respBodyBytes))
	}

	errorCode := twirp.ErrorCode(tj.Code)
	if !twirp.IsValidErrorCode(errorCode) {
		msg := "invalid type returned from server error response: " + tj.Code
		return twirp.InternalError(msg)
	}

	twerr := twirp.NewError(errorCode, tj.Msg)
	for k, v := range tj.Meta {
		twerr = twerr.WithMeta(k, v)
	}
	return twerr
}

// twirpErrorFromIntermediary maps HTTP errors from non-twirp sources to twirp errors.
// The mapping is similar to gRPC: https://github.com/grpc/grpc/blob/master/doc/http-grpc-status-mapping.md.
// Returned twirp Errors have some additional metadata for inspection.
func twirpErrorFromIntermediary(status int, msg string, bodyOrLocation string) twirp.Error {
	var code twirp.ErrorCode
	if isHTTPRedirect(status) { // 3xx
		code = twirp.Internal
	} else {
		switch status {
		case 400: // Bad Request
			code = twirp.Internal
		case 401: // Unauthorized
			code = twirp.Unauthenticated
		case 403: // Forbidden
			code = twirp.PermissionDenied
		case 404: // Not Found
			code = twirp.BadRoute
		case 429, 502, 503, 504: // Too Many Requests, Bad Gateway, Service Unavailable, Gateway Timeout
			code = twirp.Unavailable
		default: // All other codes
			code = twirp.Unknown
		}
	}

	twerr := twirp.NewError(code, msg)
	twerr = twerr.WithMeta("http_error_from_intermediary", "true") // to easily know if this error was from intermediary
	twerr = twerr.WithMeta("status_code", strconv.Itoa(status))
	if isHTTPRedirect(status) {
		twerr = twerr.WithMeta("location", bodyOrLocation)
	} else {
		twerr = twerr.WithMeta("body", bodyOrLocation)
	}
	return twerr
}

func isHTTPRedirect(status int) bool {
	return status >= 300 && status <= 399
}

// wrappedError implements the github.com/pkg/errors.Causer interface, allowing errors to be
// examined for their root cause.
type wrappedError struct {
	msg   string
	cause error
}

func wrapErr(err error, msg string) error { return &wrappedError{msg: msg, cause: err} }
func (e *wrappedError) Cause() error      { return e.cause }
func (e *wrappedError) Error() string     { return e.msg + ": " + e.cause.Error() }

// clientError adds consistency to errors generated in the client
func clientError(desc string, err error) twirp.Error {
	return twirp.InternalErrorWith(wrapErr(err, desc))
}

// badRouteError is used when the twirp server cannot route a request
func badRouteError(msg string, method, url string) twirp.Error {
	err := twirp.NewError(twirp.BadRoute, msg)
	err = err.WithMeta("twirp_invalid_route", method+" "+url)
	return err
}

// DoProtobufRequest is common code to make a request to the remote twirp service.
func DoProtobufRequest(client XHRClient, url string, in, out proto.Message) (err error) {
	reqBodyBytes, err := proto.Marshal(in)
	if err != nil {
		return clientError("failed to marshal proto request", err)
	}
	reqBody := bytes.NewBuffer(reqBodyBytes)

	req, err := newRequest(url, reqBody, "application/protobuf")
	if err != nil {
		return clientError("could not build request", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return clientError("failed to do request", err)
	}

	defer func() {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = clientError("failed to close response body", cerr)
		}
	}()

	if resp.StatusCode != 200 {
		return errorFromResponse(resp)
	}

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return clientError("failed to read response body", err)
	}

	if err = proto.Unmarshal(respBodyBytes, out); err != nil {
		return clientError("failed to unmarshal proto response", err)
	}
	return nil
}

// DoJSONRequest is common code to make a request to the remote twirp service.
func DoJSONRequest(client XHRClient, url string, in, out proto.Message) (err error) {
	reqBody := bytes.NewBuffer(nil)
	marshaler := &jsonpb.Marshaler{OrigName: true}
	if err = marshaler.Marshal(reqBody, in); err != nil {
		return clientError("failed to marshal json request", err)
	}

	req, err := newRequest(url, reqBody, "application/json")
	if err != nil {
		return clientError("could not build request", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return clientError("failed to do request", err)
	}

	defer func() {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = clientError("failed to close response body", cerr)
		}
	}()

	if resp.StatusCode != 200 {
		return errorFromResponse(resp)
	}

	unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err = unmarshaler.Unmarshal(resp.Body, out); err != nil {
		return clientError("failed to unmarshal json response", err)
	}

	return nil
}
