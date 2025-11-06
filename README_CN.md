<a name="readme-top"></a>

[[英文介绍]](https://github.com/go-plugify/example)  |  [[Join Discord]](https://discord.gg/B3FwBSQq)

# go-plugify 例子

这是一个简单的展示`go-plugify`怎么使用的例子。包括服务端跟客户端两部分。

## 依赖（如果使用go原生plugin）

- 需要 `linux` 或 `mac` 系统
- 需要 golang 1.23.10 版本及以上
- 需要 CGO 支持

## 运行

<img alt="example" src="https://github.com/go-plugify/example/blob/main/example.gif?raw=true" width="651">

### 服务端

进入 `server` 目录，直接运行 `go run main.go`，看到启动成功打印出路由信息即可。
服务端程序挂载了几个web路由，并且对 `go-plugify` 进行了初始化。挂载了两个组件：`ginengine` 与 `bookService`。提供给插件调用。

```go
func setupRouter() *gin.Engine {
	r := gin.Default()

	ginRouter := ginadapter.NewHttpRouter(r)

	bookService := service.NewBookService()

	// 初始化插件管理器，并挂载组件
	plugManager := goplugify.InitPluginManagers("default",
		goplugify.ComponentWithName("ginengine", ginRouter),
		goplugify.ComponentWithName("bookService", bookService),
	)

	registerCoreRoutes(r, bookService)

	// 挂载路由，这里可以根据您的需要增加认证或其他权限拦截中间件
	goplugify.InitHTTPServer(plugManager).RegisterRoutes(ginRouter, "/api/v1")

	return r
}
```

### 客户端

进入 `client` 目录，中可以看到有两种支持的方式：

- yaegi 解释器（推荐）
- go 原生 plugin

#### yaegi

进入 `yaegi` 目录，运行：`make init`，即可。
代码逻辑在 `main.go` 中。脚本均需要实现下面三个函数或变量：

- `Run(input map[string]any) (any, error)`
- `Methods() map[string]func(any) any{}`
- `Destroy(map[string]any) error`

```go
func Run(input map[string]any) (any, error) {
	plugify.Logger.Info("Example plugin is running")
	plugify.BookService.AddBook(plugify.ServiceBook{ID: 1, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"})
	plugify.BookService.AddBook(plugify.ServiceBook{ID: 2, Title: "Pride and Prejudice", Author: "Jane Austen"})
	plugify.BookService.DeleteBook(1)
	plugify.Logger.Info("Books in the service: %+v", plugify.BookService.ListBooks())
	plugify.Logger.Info("Example plugin finished execution")
	return map[string]any{
		"message": "Plugin executed successfully",
		"books":   plugify.BookService.ListBooks(),
	}, nil
}

func Methods() map[string]func(any) any{} {
	return nil
}

func Destroy(map[string]any) error {
	return nil
}
```

执行：`make help` 可以查看更多的命令。

#### 原生 plugin

进入 `go_plugin` 目录，运行：`make run`，即可。
客户端程序的代码在 `src` 目录下。他的逻辑是修改了服务端中 `http://localhost:8080/` 路径对应的处理方法，以及调用服务端的程序进行处理计算逻辑返回。
如 `client/src/plugin.go` 的代码：

```go
// Run is called when the plugin is executed.
func (p Plugin) Run(args any) (any, error) {
	p.Logger().Info("Plugin %s is running", p.Name)
	// 动态修改路由方法逻辑
	p.GetGinEngine().ReplaceHandler("GET", "/", func(ctx context.Context) {
		ctx.(HttpContext).JSON(200, map[string]string{"message": "Hello from example plugin !!!"})
	})
	// 获取宿主组件，并调用其方法
	bookService := NewBookService(p.Component("bookService"), p)
	bookService.AddBook(Book{ID: 1, Title: "1984", Author: "George Orwell"})
	bookService.AddBook(Book{ID: 2, Title: "To Kill a Mockingbird", Author: "Harper Lee"})
	bookService.DeleteBook(1)
	return map[string]any{
		"message":  "Plugin executed successfully",
		"load pkg": pkg.SayHello(),
		"books":    bookService.ListBooks(),
	}, nil
}
```

您可以修改体验下这里的使用方式。

执行：`make help` 可以查看更多的命令。

### 部署运行

客户端通过HTTP请求服务端接口进行部署运行，在 `Makefile` 中可以看到具体的 `curl` 命令。

```shell
curl --insecure -s -X POST \
	-F "file=@build/patch.so" \
	-F "meta={\"id\": \"example\", \"name\": \"Example Plugin\", \"description\": \"An example plugin\", \"author\": \"Your Name\", \"version\": \"v0.0.1\", \"loader\": \"native_plugin_http\"}"  \
	"http://localhost:8080/api/v1/plugin/init"
```

以表单形式传输，其中 `file` 字段为具体逻辑代码，在原生golang模式中是 `.so` 文件，在 yaegi 模式中则是对应的脚本。
`meta` 字段为插件的元信息，包括字段：

| 字段 | 类型 | 解释 | 例子 |
| --- | --- | --- | --- |
| id | string | 唯一的插件ID | example |
| name | string | 插件名字 | Example Plugin |
| description | string |  插件描述 | An example plugin |
| author | string |  插件作者 | Jack |
| version | string |  插件版本，格式为：vx.x.x | v0.0.1 |
| loader | enum |  插件加载器，目前支持：native_plugin_http/yaegi_http | native_plugin_http |
| components | JSON | 依赖的类型等 | [{"pkg_path": "example.com/server/service","name": "Book"},{"pkg_path": "example.com/server/service","name": "FictionBook"}] |

最后，接口地址是服务端中注册的地址，域名+监听端口+前缀+`/plugin/init`。
