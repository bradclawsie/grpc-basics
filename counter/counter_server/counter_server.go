package main

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"

	pb "../counter"
	"google.golang.org/grpc"
)

// server will implement CounterServer interface as described
// in ../counter/counter.pb.go. We'll also use it to encapsulate
// the internal state - a map and a mutex.
type server struct {
	m   map[string]int32
	mut *sync.RWMutex
}

// increment a key countername with value val.
// increment is not directly tied to gRPC so it can be tested
// independently.
func (s *server) increment(countername string, val int32) int32 {
	s.mut.Lock()
	defer s.mut.Unlock()
	if _, ok := s.m[countername]; !ok {
		s.m[countername] = val
	} else {
		s.m[countername] += val
	}
	return s.m[countername]
}

// IncrementCounter increments a key countername with value val.
// IncrementCounter implements one function in the CounterServer interface
// generated in ../counter/counter.pb.go.
func (s *server) IncrementCounter(ctx context.Context, in *pb.Increment) (*pb.ValueResponse, error) {
	countername := in.GetCountername()
	val := in.GetIncrement()
	newVal := s.increment(countername, val)
	resp := &pb.ValueResponse{
		Countervalue: &pb.Value{
			Countername: countername, Value: newVal,
		},
	}
	return resp, nil
}

// read the val corresponding to key countername or return an error.
// read is not directly tied to gRPC so it can be tested
// independently.
func (s *server) read(countername string) (int32, error) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	v, ok := s.m[countername]
	if !ok {
		return 0, errors.New("key not set")
	}
	return v, nil
}

// ReadCounter reads the val corresponding to key countername or return an error.
// ReadCounter implements one function in the CounterServer interface
// generated in ../counter/counter.pb.go.
func (s *server) ReadCounter(ctx context.Context, in *pb.Read) (*pb.ValueResponse, error) {
	countername := in.GetCountername()
	val, err := s.read(countername)
	if err != nil {
		return nil, err
	}
	resp := &pb.ValueResponse{
		Countervalue: &pb.Value{
			Countername: countername, Value: val,
		},
	}
	return resp, nil
}

func main() {

	// Instantiate our server instance.
	s := &server{m: make(map[string]int32), mut: &sync.RWMutex{}}

	// You'll need to remember to have clients bind to port 55555.
	listener, err := net.Listen("tcp", ":55555")
	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpc.NewServer()

	// s must satisfy the CounterServer interface defined in
	// ../counter/counter.pb.go.
	pb.RegisterCounterServer(grpcSrv, s)

	if err := grpcSrv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
