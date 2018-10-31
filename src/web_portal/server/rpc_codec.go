package server

import (
	"bufio"
	"bytes"
	log "github.com/cihub/seelog"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"
)

var (
	request_seq int64 = 0
)

// NewCodec returns a new JSON Codec.
func NewCodec() *Codec {
	return &Codec{json.NewCodec()}
}

// NewCodec returns a new JSON Codec.
func NewCodecRequest(s int64, proxy rpc.CodecRequest) *CodecRequest {
	return &CodecRequest{seq: s, proxy_codec_request: proxy}
}

// Codec creates a CodecRequest to process each request.
type Codec struct {
	proxy_codec rpc.Codec
}

type CodecRequest struct {
	seq                 int64
	proxy_codec_request rpc.CodecRequest
	response_writer     http.ResponseWriter
}

// Method returns the RPC method for the current request.
//
// The method uses a dotted notation as in "Service.Method".
func (c *CodecRequest) Method() (string, error) {
	return c.proxy_codec_request.Method()
}

// ReadRequest fills the request object for the RPC method.
func (c *CodecRequest) ReadRequest(args interface{}) error {
	return c.proxy_codec_request.ReadRequest(args)
}

// WriteResponse encodes the response and writes it to the ResponseWriter.
//
// The err parameter is the error resulted from calling the RPC method,
// or nil if there was no error.
func (c *CodecRequest) WriteResponse(w http.ResponseWriter, reply interface{}, methodErr error) error {
	c.response_writer = w
	return c.proxy_codec_request.WriteResponse(c, reply, methodErr)
}

func (c *CodecRequest) Header() http.Header {
	return c.response_writer.Header()
}

func (c *CodecRequest) Write(content []byte) (int, error) {
	log.Debug("Response[", c.seq, "]:", (string)(content))
	return c.response_writer.Write(content)
}

func (c *CodecRequest) WriteHeader(h int) {
	c.response_writer.WriteHeader(h)
}

// NewRequest returns a CodecRequest.
func (c *Codec) NewRequest(r *http.Request) rpc.CodecRequest {
	seq := atomic.AddInt64(&request_seq, 1)
	buffer := bytes.Buffer{}
	buffer_writer := bufio.NewWriter(&buffer)
	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, buffer_writer))
	ret := c.proxy_codec.NewRequest(r)
	buffer_writer.Flush()
	log.Debug("Accept request[", seq, "]:", buffer.String())
	return NewCodecRequest(seq, ret)
}
