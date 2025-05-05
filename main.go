package main

import (
	"log"
	"net"
	"net/http"

	"nova-panel/internal/grpcserver"
	"nova-panel/internal/webserver"

	"google.golang.org/grpc"
	"nova-panel/pb"
)

func main() {
	// 启动 gRPC 服务
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("gRPC 监听失败: %v", err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterVpsServer(grpcServer, &grpcserver.Server{})
		log.Println("gRPC 启动成功: :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC 服务异常: %v", err)
		}
	}()

	// 启动 Gin Web 服务
	r := webserver.InitRouter()
	log.Println("Web 接口启动: :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Web 接口异常: %v", err)
	}
}
