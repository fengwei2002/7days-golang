## HTTP web framework written in Go (koo) 

> file history: 

- [01: 基础封装 http](https://github.com/fengwei2002/7days-golang/tree/23dcbfa6779a8973ecc39dd86a0f0dfaa52ce3ff) 使用 engine 截获路由进行处理
- [02: 创建上下文 context](https://github.com/fengwei2002/my-gin/tree/93e0fd4f8909bdaf7f75e02887ad08aa7b3854c9) 封装 context
![konng0120-README-2022-05-17-15-53-08](https://raw.githubusercontent.com/psychonaut1f/a/main/img/konng0120-README-2022-05-17-15-53-08.png)
- [03: 使用 trie 树管理路由](https://github.com/fengwei2002/my-gin/commit/3c3791d02e0552da3518a3c287875902e9890932) 将简单的 map 映射改为 trie 存储，支持两种模式 `/:name` 和 `/*filepath`
![konng0120-README-2022-05-17-21-41-46](https://raw.githubusercontent.com/psychonaut1f/a/main/img/konng0120-README-2022-05-17-21-41-46.png)
- [04: 实现路由分组控制(Route Group Control)]()