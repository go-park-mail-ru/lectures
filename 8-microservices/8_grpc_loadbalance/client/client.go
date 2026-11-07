package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"

	"github.com/go-park-mail-ru/lectures/8-microservices/5_grpc/session"

	consulapi "github.com/hashicorp/consul/api"
)

var (
	consulAddr = flag.String("addr", "127.0.0.1:8500", "consul addr (8500 in original consul)")

	consul *consulapi.Client
)

func main() {
	flag.Parse()

	var err error
	config := consulapi.DefaultConfig()
	config.Address = *consulAddr
	consul, err = consulapi.NewClient(config)

	services, _, err := consul.Health().Service("session-api", "", false, nil)
	if err != nil {
		log.Fatalf("cant get alive services")
	}

	servers := make([]resolver.Address, 0, len(services))
	for _, item := range services {
		addr := fmt.Sprintf("%s:%d", item.Service.Address, item.Service.Port)
		servers = append(servers, resolver.Address{Addr: addr})
	}

	nameResolver := manual.NewBuilderWithScheme("myservice")
	nameResolver.InitialState(resolver.State{Addresses: servers})

	// grpclog.SetLogger(log.New(os.Stdout, "", log.LstdFlags)) // add logging
	grpcClient, err := grpc.NewClient(
		nameResolver.Scheme()+":///",
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithResolvers(nameResolver),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc: %v", err)
	}
	defer grpcClient.Close()

	sessManager := session.NewAuthCheckerClient(grpcClient)

	// тут мы будем периодически опрашивать консул на предмет изменений
	go runOnlineServiceDiscovery(nameResolver)

	ctx := context.Background()
	step := 1
	for {
		// проверяем несуществующую сессию
		// потому что сейчас между сервисами нет общения
		// получаем заглушку
		sess, err := sessManager.Check(ctx,
			&session.SessionID{
				ID: "not_exists_" + strconv.Itoa(step),
			})
		fmt.Println("get sess", step, sess, err)

		time.Sleep(500 * time.Millisecond)
		step++
	}
}

func runOnlineServiceDiscovery(nameResolver *manual.Resolver) {
	ticker := time.Tick(5 * time.Second)
	for range ticker {
		services, _, err := consul.Health().Service("session-api", "", false, nil)
		if err != nil {
			log.Fatalf("cant get alive services")
		}

		servers := make([]resolver.Address, 0, len(services))
		for _, item := range services {
			addr := item.Service.Address +
				":" + strconv.Itoa(item.Service.Port)
			servers = append(servers, resolver.Address{Addr: addr})
		}

		nameResolver.UpdateState(resolver.State{Addresses: servers})
	}
}
