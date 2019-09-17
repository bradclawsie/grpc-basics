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

func main() {
	const use = "use: client [inc $key $val] | [get $key]"
	if len(os.Args) < 3 {
		log.Fatal(use)
	}

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

		countervalue := valueResponse.GetCountervalue()
		countername := countervalue.GetCountername()
		newValue := countervalue.GetValue()
		log.Printf("%v %v\n", countername, newValue)

	} else {
		log.Fatal(use)
	}
}
