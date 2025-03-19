package main

import (
	"context"
	"fmt"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/resource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func main() {
	// 验证grpc服务
	conn, err := grpc.NewClient("127.0.0.1:18010", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := resource.NewRpcClient(conn)
	set, err := client.Search(context.Background(), resource.NewSearchRequest())
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Println(set)
}
