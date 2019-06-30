package main

import (
	"context"
	"fmt"
	"time"

	pb "server/api/v1"
	"server/consule"

	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

const (
	target = "consul://127.0.0.1:8500/hello"
)

func main()  {
	resolver.Register(consule.New())
	ctx, _ := context.WithTimeout(context.Background(),10*time.Second)
	conn, err := grpc.DialContext(ctx,target, grpc.WithBlock(), grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
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
			fmt.Println("service is error",err.Error())
			time.Sleep(time.Second)
			continue
		}
		fmt.Println(msg.Message)
		time.Sleep(time.Second)
	}
}