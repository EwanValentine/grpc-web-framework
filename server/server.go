package server

import (
	"context"
	"io"
	"net/http"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func New() *Server {
	r := http.NewServeMux()

	// TODO: we probably need a mutex here to manage concurrent access to the endpoints map
	return &Server{
		router:    r,
		endpoints: make(map[string]map[string]Endpoint),
		mu:        &sync.RWMutex{},
	}
}

// Server adds additional functionality to a `http.ServeMux` to turn incoming http requests
// into gRPC requests for matching registered endpoints
type Server struct {
	router    *http.ServeMux
	endpoints map[string]map[string]Endpoint
	mu        *sync.RWMutex
}

// RegisterEndpoint registers a handler to a method and path
func (s *Server) RegisterEndpoint(method, path string, handler Endpoint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.endpoints[path] = make(map[string]Endpoint)
	s.endpoints[path][method] = handler
}

// InvokeEndpoint invokes an endpoint by path and method
func (s *Server) InvokeEndpoint(method, path string, body []byte) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.endpoints[path][method](body)
}

// ServeHTTP introspects the incoming request, and looks for a registered endpoint
// for this path, and http method
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if endpoint, ok := s.endpoints[r.URL.Path][r.Method]; ok {

		// TODO: map headers to grpc metadata ctx?
		// TODO: implement some kind of middleware architecture
		// TODO: handle path params, this will likely require some fuzzy checking for the path=>endpoint mapping

		var b []byte
		if r.Method == http.MethodPost {
			b, _ = io.ReadAll(r.Body)
		} else {
			// TODO: map get request params to a map[string]string and pass it to the grpc request
			b = []byte(`{}`)
		}

		body, err := endpoint(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
		return
	}
	s.router.ServeHTTP(w, r)
}

// ClientMethod represents a client method on a gRPC service
type ClientMethod[RQ, RS proto.Message] func(context.Context, RQ, ...grpc.CallOption) (RS, error)

// Endpoint represents what happens when a route is called
// note: endpoint is no longer the best name for this...
type Endpoint func([]byte) ([]byte, error)

// Marshaller is a mapper callback which is used to map incoming request bytes into
// a gRPC proto.Message type, to be sent to the server
type Marshaller[RQ proto.Message] func([]byte) (RQ, error)

// NewHandler creates a handler, which encapsulates the client method, the marshaller used to
// create the gRPC request type, and returns an endpoint ready to be called via http
func NewHandler[RQ, RS proto.Message](clientMethod ClientMethod[RQ, RS], marshaller Marshaller[RQ]) Endpoint {
	return func(incoming []byte) ([]byte, error) {
		grpcRequest, err := marshaller(incoming)
		if err != nil {
			return nil, err
		}

		response, err := clientMethod(context.Background(), grpcRequest)

		b, _ := protojson.Marshal(response)

		return b, err
	}
}
