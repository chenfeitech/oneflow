package jsonrpc_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
)

var ErrNullResult = errors.New("result is null")

func Request(url string, method string, args interface{}, reply interface{}) error {
	body, err := EncodeClientRequest(method, args)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	return DecodeClientResponse(resp.Body, reply)
}

// EncodeClientRequest encodes parameters for a JSON-RPC client request.
func EncodeClientRequest(method string, args interface{}) ([]byte, error) {
	c := &ClientRequest{
		Version: "2.0",
		Method:  method,
		Params:  args,
		Id:      uint64(rand.Int31()),
	}
	return json.Marshal(c)
}

// DecodeClientResponse decodes the response body of a client request into
// the interface reply.
func DecodeClientResponse(r io.Reader, reply interface{}) error {
	var c ClientResponse
	if err := json.NewDecoder(r).Decode(&c); err != nil {
		return err
	}
	if c.Error != nil {
		jsonErr := &Error{}
		if err := json.Unmarshal(*c.Error, jsonErr); err != nil {
			return &Error{
				Code:    -32700,
				Message: string(*c.Error),
			}
		}
		return jsonErr
	}

	if c.Result == nil {
		return ErrNullResult
	}

	return json.Unmarshal(*c.Result, reply)
}
