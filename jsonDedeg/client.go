package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

type Reply struct {
	res []byte
}

func init() {
	encoding.RegisterCodec(JSON{
		Marshaler: jsonpb.Marshaler{
			EmitDefaults: true,
			OrigName:     true,
		},
	})
}


func main()  {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(JSON{}.Name())),
	}
	conn, err := grpc.Dial("127.0.0.1:9003", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	var reply Reply
	str := `{"name":"guolin"}`
	if err = conn.Invoke(context.Background(),"/helloworld.v1.Greeter/SayHello",[]byte(str), &reply);err != nil{
		panic(err)
	}
	fmt.Println(string(reply.res))
}

// JSON is impl of encoding.Codec
type JSON struct {
	jsonpb.Marshaler
	jsonpb.Unmarshaler
}

// Name is name of JSON
func (j JSON) Name() string {
	return "json"
}

// Marshal is json marshal
func (j JSON) Marshal(v interface{}) (out []byte, err error) {
	return v.([]byte), nil
}

// Unmarshal is json unmarshal
func (j JSON) Unmarshal(data []byte, v interface{}) (err error) {
	v.(*Reply).res = data
	return nil
}

