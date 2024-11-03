package base

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

type APIType int

const (
	Public APIType = iota
	Private
)

type params map[string]interface{}

type Request struct {
	Method      string
	Endpoint    string
	Params      url.Values
	ParamString string
	Header      http.Header
	Body        io.Reader
	ApiType     APIType
	Url         string
}

func NewRequest(method string, endpoint string, isPublic bool) *Request {
	apiType := Public
	if !isPublic {
		apiType = Private
	}
	return &Request{
		Method:   method,
		Endpoint: endpoint,
		ApiType:  apiType,
	}
}

// AddParam  add param with key/value to params string
func (r *Request) AddParam(key string, value interface{}) *Request {
	if r.Params == nil {
		r.Params = url.Values{}
	}
	r.Params.Add(key, fmt.Sprintf("%v", value))
	return r
}

// SetParam  set param with key/value to params string
func (r *Request) SetParam(key string, value interface{}) *Request {
	if r.Params == nil {
		r.Params = url.Values{}
	}

	if reflect.TypeOf(value).Kind() == reflect.Slice {
		v, err := json.Marshal(value)
		if err == nil {
			value = string(v)
		}
	}

	r.Params.Set(key, fmt.Sprintf("%v", value))
	return r
}

// SetParams  set params with key/values to params string
func (r *Request) SetParams(m params) *Request {
	for k, v := range m {
		r.SetParam(k, v)
	}
	return r
}

func (r *Request) EncodeParams() {
	if r.Params == nil {
		r.Params = url.Values{}
	}
	r.ParamString = r.Params.Encode()
}
