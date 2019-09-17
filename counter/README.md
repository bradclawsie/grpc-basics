To generated gRPC bindings:

`make protoc`

To build server:

`cd counter_server; go build`

To build client:

`cd counter_client; go build`

Once server is running, try client:

`counter_client inc key 2`

`counter_client get key`

`counter_client get not-there`
