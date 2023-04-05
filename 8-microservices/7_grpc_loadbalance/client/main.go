package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"

	"github.com/go-park-mail-ru/lectures/8-microservices/4_grpc/session"

	consulapi "github.com/hashicorp/consul/api"
)

var (
	consulAddr = flag.String("addr", "127.0.0.1:8500", "consul addr (8500 in original consul)")
)

var (
	consul *consulapi.Client
)

func main() {
	flag.Parse()

	var err error
	config := consulapi.DefaultConfig()
	config.Address = *consulAddr
	consul, err = consulapi.NewClient(config)

	health, _, err := consul.Health().Service("session-api", "", false, nil)
	if err != nil {
		log.Fatalf("cant get alive services")
	}

	servers := make([]resolver.Address, 0, len(health))
	for _, item := range health {
		addr := item.Service.Address +
			":" + strconv.Itoa(item.Service.Port)
		servers = append(servers, resolver.Address{Addr: addr})
	}

	if len(servers) == 0 {
		panic("no alive session-api servers")
	}

	nameResolver := manual.NewBuilderWithScheme("myservice")
	nameResolver.InitialState(resolver.State{Addresses: servers})

	// grpclog.SetLogger(log.New(os.Stdout, "", log.LstdFlags)) // add logging
	grpcConn, err := grpc.Dial(
		nameResolver.Scheme()+":///",
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithResolvers(nameResolver),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc: %v", err)
	}
	defer grpcConn.Close()

	sessManager := session.NewAuthCheckerClient(grpcConn)

	// тут мы будем периодически опрашивать консул на предмет изменений
	go runOnlineServiceDiscovery(nameResolver)

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

		time.Sleep(500 * time.Millisecond)
		step++
	}
}

func runOnlineServiceDiscovery(nameResolver *manual.Resolver) {
	ticker := time.Tick(5 * time.Second)
	for _ = range ticker {
		health, _, err := consul.Health().Service("session-api", "", false, nil)
		if err != nil {
			log.Fatalf("cant get alive services")
		}

		servers := make([]resolver.Address, 0, len(health))
		for _, item := range health {
			addr := item.Service.Address +
				":" + strconv.Itoa(item.Service.Port)
			servers = append(servers, resolver.Address{Addr: addr})
		}
		nameResolver.CC.NewAddress(servers)
	}
}
