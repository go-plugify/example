# go-plugify 例子

这是一个简单的展示`go-plugify`怎么使用的例子。包括服务端跟客户端两部分。

## 依赖

- 需要 `linux` 或 `mac` 系统
- 需要 golang 1.23.10 版本及以上

## 运行

<img alt="example" src="https://github.com/go-plugify/example/blob/main/example.gif?raw=true" width="651">

### 服务端

进入 `server` 目录，直接运行 `go run main.go`，看到启动成功打印出路由信息即可。
服务端程序挂载了几个web路由，并且对 `go-plugify` 进行了初始化。

### 客户端

进入 `client` 目录，运行：`make run`，即可运行。
客户端程序的代码在 `src` 目录下。他的逻辑是修改了服务端中 `http://localhost:8080/` 路径对应的处理方法，以及调用服务端的程序进行处理计算逻辑返回。
如 `client/src/plugin.go` 的代码：

```go
// Run is called when the plugin is executed.
func (p Plugin) Run(args any) {
	ctx := args.(HttpContext)
	p.Logger().Info("Plugin %s is running, ctx %+v", p.Name, ctx)
	p.Component("ginengine").(HttpRouter).ReplaceHandler("GET", "/", func(ctx context.Context) {
		ctx.(HttpContext).JSON(200, map[string]string{"message": "Hello from plugin !!!"})
	})
	cal := p.Component("calculator").(Calculator)
	ctx.JSON(200, map[string]any{
		"message":      "Plugin executed successfully",
		"load pkg":     pkg.SayHello(),
		"1 + 5 * 5 = ": cal.Add(1, cal.Mul(5, 5)),
	})
}
```

您可以修改体验下这里的使用方式。
