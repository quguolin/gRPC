package main

import (
	"context"
	"fmt"
	"time"

	pb "server/api/v1"

	"google.golang.org/grpc"
)

func main()  {
	conn, err := grpc.Dial("127.0.0.1:9001",grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	for{
		msg,err := client.SayHello(context.Background(),&pb.HelloRequest{
			Name:"jack",
		})
		if err != nil{
			panic(err)
		}
		fmt.Println(msg.Message)
		time.Sleep(time.Second)
	}
}