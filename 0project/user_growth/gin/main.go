package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net/http"
	"time"
	"user_growth/conf"
	"user_growth/dbhelper"
	"user_growth/pb"
	"user_growth/ugserver"
)

func main() {
	grpcGateway()
}

var AllowOrigin = map[string]bool{
	"http://a.site.com": true,
	"http://b.site.com": true,
	"http://web.com":    true,
}

func mainGin() {
	// 连接到grpc服务的客户端
	conn, err := grpc.Dial("localhost:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	clientCoin := pb.NewUserCoinClient(conn)
	clientGrade := pb.NewUserGradeClient(conn)

	router := gin.New()
	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})
	// 用户积分服务的方法
	v1Group := router.Group("/v1", func(ctx *gin.Context) {
		// prometheus 指标
		//MetricAdd()

		// 支持跨域
		origin := ctx.GetHeader("Origin")
		if AllowOrigin[origin] {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTION")
			ctx.Header("Access-Control-Allow-Headers", "*")
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}
		ctx.Next()
	})
	gUserCoin := v1Group.Group("/UserGrowth.UserCoin")
	gUserCoin.GET("/ListTasks", func(ctx *gin.Context) {
		out, err := clientCoin.ListTasks(ctx, &pb.ListTasksRequest{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    2,
				"message": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, out)
		}
	})
	gUserCoin.POST("/UserCoinChange", func(ctx *gin.Context) {
		body := &pb.UserCoinChangeRequest{}
		err := ctx.BindJSON(body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    2,
				"message": err.Error(),
			})
		} else if out, err := clientCoin.UserCoinChange(ctx, body); err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    2,
				"message": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, out)
		}
		ctx.JSON(http.StatusOK, nil)
	})
	// 用户等级服务的方法
	gUserGrade := v1Group.Group("/UserGrowth.UserGrade")
	gUserGrade.GET("/ListGrades", func(ctx *gin.Context) {
		out, err := clientGrade.ListGrades(ctx, &pb.ListGradesRequest{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    2,
				"message": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, out)
		}
	})

	// prometheus client Create non-global registry.
	//MetricInit(router)

	// 为http/2配置参数
	h2Handler := h2c.NewHandler(router, &http2.Server{})
	// 配置http服务
	server := &http.Server{
		Addr:    ":8080",
		Handler: h2Handler,
	}
	// 启动http服务
	server.ListenAndServe()
}
func initDb() {
	// default UTC time location
	time.Local = time.UTC
	// Load global config
	conf.LoadConfigs()
	// Initialize db
	dbhelper.InitDb()
}
func grpcGateway() {
	// 初始化数据库实例
	initDb()

	s := grpc.NewServer()
	//注册服务
	pb.RegisterUserCoinServer(s, &ugserver.UgCoinServer{})   //用户积分
	pb.RegisterUserGradeServer(s, &ugserver.UgGradeServer{}) //用户等级

	//grpc gateway
	mux := runtime.NewServeMux()
	ctx := context.Background()
	err := pb.RegisterUserCoinHandlerServer(ctx, mux, &ugserver.UgCoinServer{})
	if err != nil {
		log.Printf("failed to register UgCoinServer: %v", err)
	}

	err = pb.RegisterUserGradeHandlerServer(ctx, mux, &ugserver.UgGradeServer{})
	if err != nil {
		log.Printf("failed to register UgGradeServer: %v", err)
	}

	reflection.Register(s)

	httpMux := http.NewServeMux()
	httpMux.Handle("/v1/UserGrowth", mux)

	server := &http.Server{
		Addr: ":8081",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("http.HandlerFunc url=%s", r.URL)
			mux.ServeHTTP(w, r)
		}),
	}

	//启动服务
	log.Printf("服务启动中 %v\n", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to listen and serve %v", err)
	}
}
