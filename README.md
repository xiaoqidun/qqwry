# QQWry [![Go Reference](https://pkg.go.dev/badge/github.com/xiaoqidun/qqwry.svg)](https://pkg.go.dev/github.com/xiaoqidun/qqwry)

Golang QQWry，高性能纯真IP查询库。

# 使用须知

1. 仅支持ipv4查询。

# 使用说明

```go
package main

import (
	"fmt"
	"github.com/xiaoqidun/qqwry"
)

func main() {
	// 从文件加载IP数据库
	if err := qqwry.LoadFile("qqwry.dat"); err != nil {
		panic(err)
	}
	// 从内存或缓存查询IP
	location, err := qqwry.QueryIP("119.29.29.29")
	if err != nil {
		fmt.Printf("错误：%v\n", err)
		return
	}
	fmt.Printf("国家：%s，省份：%s，城市：%s，区县：%s，运营商：%s\n",
		location.Country,
		location.Province,
		location.City,
		location.District,
		location.ISP,
	)
}
```

# IP数据库

- DAT格式：[https://aite.xyz/share-file/qqwry/qqwry.dat](https://aite.xyz/share-file/qqwry/qqwry.dat)
- IPDB格式：[https://aite.xyz/share-file/qqwry/qqwry.ipdb](https://aite.xyz/share-file/qqwry/qqwry.ipdb)

# 编译说明

1. 下载IP数据库并放置于assets目录中。
2. client和server需要go1.16的内嵌资源特性。
3. 作为库使用，请直接引包，并不需要go1.16+才能编译。

# 数据更新

- 由于qqwry.dat缺乏更新，官方czdb格式又难以获得和分发，建议使用ipdb格式。
- 这里的ipdb格式指metowolf提供的官方czdb格式转换而来的ipdb格式（纯真格式原版）。

# 服务接口

1. 自行根据需要调整server下源码。
2. 可以通过-listen参数指定http服务地址。
3. json api：curl http://127.0.0.1/ip/119.29.29.29

# 特别感谢

- 感谢[纯真IP库](https://www.cz88.net/)一直坚持为大家提供免费IP数据库。
- 感谢[yinheli](https://github.com/yinheli)的[qqwry](https://github.com/yinheli/qqwry)项目，为我提供纯真ip库解析算法参考。
- 感谢[metowolf](https://github.com/metowolf)的[qqwry.ipdb](https://github.com/metowolf/qqwry.ipdb)项目，提供纯真czdb转ipdb数据库。

# 授权说明

使用本类库你唯一需要做的就是把LICENSE文件往你用到的项目中拷贝一份。