package main

/*
$ curl http://localhost:8080/_cache/scores/Tom
630
$ curl http://localhost:8080/_cache/scores/kkk
kkk not exist
*/

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"tiny-cache/tinyCache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *tinyCache.Group {
	return tinyCache.NewGroup("scores", 2<<10, tinyCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, gee *tinyCache.Group) {
	peers := tinyCache.NewHTTPPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("tinyCache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *tinyCache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("frontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "cache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], addrs, gee)
}

// 							是
// 接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
// 				|  否                         是
// 				|-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
// 							| 否
// 							|-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶

// 流程 2 细化
// 使用一致性哈希选择节点        是                                    是
//      |-----> 是否是远程节点 -----> HTTP 客户端访问远程节点 --> 成功？-----> 服务端返回返回值
//             |  否                                    ↓  否
//             |----------------------------> 回退到本地节点处理。
