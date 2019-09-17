package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	pb "../counter"
	"google.golang.org/grpc"
)

// The client is intended to be used as a command line tool.
// Make sure there is a corresponding server instance running.
// To increment a key:
// counter_client inc key val (where val is an int32)
// To get a key's value:
// counter_client get key

func main() {
	const use = "use: counter_client [inc $key $val] | [get $key]"

	if len(os.Args) < 3 {
		log.Fatal(use)
	}

	// Connect to the server.
	conn, err := grpc.Dial("localhost:55555", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewCounterClient(conn)

	if os.Args[1] == "inc" {

		if len(os.Args) < 4 {
			log.Fatal(use)
		}
		key := os.Args[2]
		valInt, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		val := int32(valInt)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		valueResponse, err := client.IncrementCounter(ctx, &pb.Increment{
			Countername: key,
			Increment:   val,
		})
		if err != nil {
			log.Fatal(err)
		}

		// Use generated getters to extract values.
		countervalue := valueResponse.GetCountervalue()
		countername := countervalue.GetCountername()
		newValue := countervalue.GetValue()
		log.Printf("%v %v\n", countername, newValue)

	} else if os.Args[1] == "get" {

		key := os.Args[2]

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err = client.ReadCounter(ctx, &pb.Read{
			Countername: key,
		})
		if err != nil {
			log.Fatal(err)
		}

		valueResponse, err := client.ReadCounter(ctx, &pb.Read{
			Countername: key,
		})
		if err != nil {
			log.Fatal(err)
		}

		// Use generated getters to extract values.
		countervalue := valueResponse.GetCountervalue()
		countername := countervalue.GetCountername()
		newValue := countervalue.GetValue()
		log.Printf("%v %v\n", countername, newValue)

	} else {
		log.Fatal(use)
	}
}
