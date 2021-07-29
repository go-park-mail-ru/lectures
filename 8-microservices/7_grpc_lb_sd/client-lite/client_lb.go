package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	registry "github.com/liyue201/grpc-lb/registry/consul"

	"github.com/go-park-mail-ru/lectures/8-microservices/4_grpc/session"

	consulapi "github.com/hashicorp/consul/api"
)

var (
	consulAddr = flag.String("addr", "127.0.0.1:8500", "consul addr (8500 in original consul)")
)

func main() {
	flag.Parse()

	var err error

	fmt.Println("resolver", *consulAddr)

	registry.RegisterResolver("consul", &consulapi.Config{Address: *consulAddr}, "session-api")

	grcpConn, err := grpc.Dial(
		"consul:///",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	sessManager := session.NewAuthCheckerClient(grcpConn)

	ctx := context.Background()
	step := 1
	for {
		// проверяем несуществуюущую сессию
		// потому что сейчас между сервисами нет общения
		// получаем загшулку
		sess, err := sessManager.Check(ctx,
			&session.SessionID{
				ID: "not_exist_" + strconv.Itoa(step),
			})
		fmt.Println("get sess", step, sess, err)

		time.Sleep(1500 * time.Millisecond)
		step++
	}
}

