package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
	"tinygrpc"
)

func startServer(addr chan string) {
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	tinygrpc.Accept(l)
}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)
	client, _ := tinygrpc.Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)

	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ { // 实现五个 rpc 调用，参数和返回值都是 string
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("tinygrpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
}
