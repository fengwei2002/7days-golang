package singleflight

import "sync"

/*

几个缓存相关的概念：

缓存雪崩：缓存在同一个时刻全部失效，造成 瞬间 的 DB 请求量增加，压力剧增，引起雪崩
		缓存血本通常因为服务器宕机，缓存的 key 设置了相同的过期时间等引起
缓存击穿：一个存在的 key 在缓存过期的那一刻， 同时存在大量的请求，造成 DB 压力变大
缓存穿透：查询一个不存在的数据，因为不存在则不会写到缓存中，所以每次都会请求 DB 如果瞬间流量过大，穿透到 DB 导致宕机
*/

// call 表示正在进行中，或者已经结束的请求，使用 sync waitGroup 避免重入
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group 是防止缓存击穿的主要数据结构，管理不同的 key 的请求 call
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized
}

// Do 接收两个参数
// 第一个参数是 key 第二个参数是 fn
// do 的作用是 针对相同的 key 无论 Do 调用多少次，fn 函数都只会被调用一次
// 等待 fn 的调用结束了，返回返回值或者是错误
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()         // 如果请求正在进行中，则等待
		return c.val, c.err // 请求结束，返回结果
	}
	c := new(call)
	c.wg.Add(1)  // 发起请求之前加锁
	g.m[key] = c // 添加 g.m 表明 key 已经有对应的 请求在处理
	g.mu.Unlock()

	c.val, c.err = fn() // 调用 fn 发起请求
	c.wg.Done()         // 请求结束

	g.mu.Lock()
	delete(g.m, key) // 更新 g.m
	g.mu.Unlock()

	return c.val, c.err // 返回结果
}
