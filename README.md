# gRPC


## 生成proto文件
 ```protoc --go_out=plugins=grpc:. *.proto```
 
 ## consul
 ```bash
    brew install consul
    consul agent -dev
    http://127.0.0.1:8500/ui/dc1/services
```