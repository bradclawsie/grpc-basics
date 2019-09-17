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

const (
	port = ":55555"
)

// server will implement CounterServer interface.
type server struct {
	m   map[string]int32
	mut *sync.RWMutex
}

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

func (s *server) read(countername string) (int32, error) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	v, ok := s.m[countername]
	if !ok {
		return 0, errors.New("key not set")
	}
	return v, nil
}

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
	s := &server{m: make(map[string]int32), mut: &sync.RWMutex{}}
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	grpcSrv := grpc.NewServer()
	pb.RegisterCounterServer(grpcSrv, s)
	if err := grpcSrv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
