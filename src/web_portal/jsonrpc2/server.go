package jsonrpc2

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// ----------------------------------------------------------------------------
// Codec
// ----------------------------------------------------------------------------

// Codec creates a CodecRequest to process each request.
type Codec interface {
	NewRequest(*http.Request) CodecRequest
	NewRequestWithBytes(r []byte) CodecRequest
}

// CodecRequest decodes a request and encodes a response using a specific
// serialization scheme.
type CodecRequest interface {
	// Reads the request and returns the RPC method name.
	Method() (string, error)
	// Reads the request filling the RPC method args.
	ReadRequest(interface{}) error
	// Writes the response using the RPC method reply.
	WriteResponse(io.Writer, interface{})
	// Writes an error produced by the server.
	WriteError(w io.Writer, status int, err error)
}

// ----------------------------------------------------------------------------
// Server
// ----------------------------------------------------------------------------

// NewServer returns a new RPC server.
func NewServer() *Server {
	return &Server{
		codecs:   make(map[string]Codec),
		services: new(serviceMap),
	}
}

// Server serves registered RPC services using registered codecs.
type Server struct {
	codecs   map[string]Codec
	services *serviceMap
}

// RegisterCodec adds a new codec to the server.
//
// Codecs are defined to process a given serialization scheme, e.g., JSON or
// XML. A codec is chosen based on the "Content-Type" header from the request,
// excluding the charset definition.
func (s *Server) RegisterCodec(codec Codec, contentType string) {
	s.codecs[strings.ToLower(contentType)] = codec
}

// RegisterService adds a new service to the server.
//
// The name parameter is optional: if empty it will be inferred from
// the receiver type name.
//
// Methods from the receiver will be extracted if these rules are satisfied:
//
//    - The receiver is exported (begins with an upper case letter) or local
//      (defined in the package registering the service).
//    - The method name is exported.
//    - The method has three arguments: *http.Request, *args, *reply.
//    - All three arguments are pointers.
//    - The second and third arguments are exported or local.
//    - The method has return type error.
//
// All other methods are ignored.
func (s *Server) RegisterService(receiver interface{}, name string) error {
	return s.services.register(receiver, name)
}

// HasMethod returns true if the given method is registered.
//
// The method uses a dotted notation as in "Service.Method".
func (s *Server) HasMethod(method string) bool {
	if _, _, err := s.services.get(method); err == nil {
		return true
	}
	return false
}

// ServeHTTP
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteError(w, 405, "rpc: POST method required, received "+r.Method)
		return
	}
	contentType := r.Header.Get("Content-Type")
	idx := strings.Index(contentType, ";")
	if idx != -1 {
		contentType = contentType[:idx]
	}
	codec := s.codecs[strings.ToLower(contentType)]
	if codec == nil {
		WriteError(w, 415, "rpc: unrecognized Content-Type: "+contentType)
		return
	}
	// Create a new codec request.
	codecReq := codec.NewRequest(r)
	// Get service method to be called.
	method, errMethod := codecReq.Method()
	if errMethod != nil {
		codecReq.WriteError(w, 400, errMethod)
		return
	}
	serviceSpec, methodSpec, errGet := s.services.get(method)
	if errGet != nil {
		codecReq.WriteError(w, 400, errGet)
		return
	}
	// Decode the args.
	args := reflect.New(methodSpec.argsType)
	if errRead := codecReq.ReadRequest(args.Interface()); errRead != nil {
		codecReq.WriteError(w, 400, errRead)
		return
	}
	// Call the service method.
	reply := reflect.New(methodSpec.replyType)
	errValue := methodSpec.method.Func.Call([]reflect.Value{
		serviceSpec.rcvr,
		reflect.ValueOf(r),
		args,
		reply,
	})
	// Cast the result to error if needed.
	var errResult error
	errInter := errValue[0].Interface()
	if errInter != nil {
		errResult = errInter.(error)
	}
	// Prevents Internet Explorer from MIME-sniffing a response away
	// from the declared content-type
	w.Header().Set("x-content-type-options", "nosniff")
	// Encode the response.
	if errResult == nil {
		codecReq.WriteResponse(w, reply.Interface())
	} else {
		codecReq.WriteError(w, 400, errResult)
	}
}

// ServeBytes
func (s *Server) ServeBytes(request *http.Request, body []byte) []byte {
	w := &bytes.Buffer{}
	codec := s.codecs["application/json"]
	if codec == nil {
		WriteError(w, 415, "rpc: unrecognized Content-Type: application/json")
		return w.Bytes()
	}
	// Create a new codec request.
	codecReq := codec.NewRequestWithBytes(body)
	// Get service method to be called.
	method, errMethod := codecReq.Method()
	if errMethod != nil {
		codecReq.WriteError(w, 400, errMethod)
		return w.Bytes()
	}
	serviceSpec, methodSpec, errGet := s.services.get(method)
	if errGet != nil {
		codecReq.WriteError(w, 400, errGet)
		return w.Bytes()
	}
	// Decode the args.
	args := reflect.New(methodSpec.argsType)
	if errRead := codecReq.ReadRequest(args.Interface()); errRead != nil {
		codecReq.WriteError(w, 400, errRead)
		return w.Bytes()
	}
	// Call the service method.
	reply := reflect.New(methodSpec.replyType)
	errValue := methodSpec.method.Func.Call([]reflect.Value{
		serviceSpec.rcvr,
		reflect.ValueOf(request),
		args,
		reply,
	})
	// Cast the result to error if needed.
	var errResult error
	errInter := errValue[0].Interface()
	if errInter != nil {
		errResult = errInter.(error)
	}

	// Encode the response.
	if errResult == nil {
		codecReq.WriteResponse(w, reply.Interface())
	} else {
		codecReq.WriteError(w, 400, errResult)
	}
	return w.Bytes()
}

func WriteError(w io.Writer, status int, msg string) {
	if hw, ok := w.(http.ResponseWriter); ok {
		hw.WriteHeader(status)
		hw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
	fmt.Fprint(w, msg)
}
