
grpc 代码生成
```shell
protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative  .\user_growth.proto
```

下载 gin
```shell
go get -u github.com/gin-gonic/gin
```

官网<https://github.com/grpc-ecosystem/grpc-gateway>

引入 grpc-getway
```shell
go get github.com/grpc-ecosystem/grpc-gateway/v2
```


生成 gw 代码
```shell
protoc -I . --grpc-gateway_out ./  --grpc-gateway_opt logtostderr=true  --grpc-gateway_opt paths=source_relative     --grpc-gateway_opt generate_unbound_methods=true  user_growth.proto
```


