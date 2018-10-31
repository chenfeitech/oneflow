package jsonrpc_client

import (
	"encoding/json"
	"fmt"
)

// ClientRequest represents a JSON-RPC request sent by a client.
type ClientRequest struct {
	// JSON-RPC protocol.
	Version string `json:"jsonrpc"`

	// A String containing the name of the method to be invoked.
	Method string `json:"method"`

	// Object to pass as request parameter to the method.
	Params interface{} `json:"params"`

	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	Id uint64 `json:"id"`
}

// ClientResponse represents a JSON-RPC response returned to a client.
type ClientResponse struct {
	Version string           `json:"jsonrpc"`
	Result  *json.RawMessage `json:"result"`
	Error   *json.RawMessage `json:"error"`
	Id      uint64           `json:"id"`
}

type Error struct {
	// A Number that indicates the error type that occurred.
	Code int `json:"code"` /* required */

	// A String providing a short description of the error.
	// The message SHOULD be limited to a concise single sentence.
	Message string `json:"message"` /* required */

	// A Primitive or Structured value that contains additional information about the error.
	Data interface{} `json:"data"` /* optional */
}

func (j *Error) Error() string {
	return fmt.Sprint("(", j.Code, ")", j.Message)
}
