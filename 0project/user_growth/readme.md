
grpc 代码生成
```shell
protoc 
  --go_out=. --go_opt=paths=source_relative 
  --go-grpc_out=. --go-grpc_opt=paths=source_relative  .\user_growth.proto
```

下载 gin
```shell
go get -u github.com/gin-gonic/gin
```




