# Example of go-plugify

This is a simple example demonstrating how to use go-plugify, including both the server and client parts.

## Dependencies

- Requires `Linux` or `macOS`
- Requires `Golang 1.23.10` or higher

## Run

### Server

Enter the `server` directory and run `go run main.go`.
Once you see the startup logs showing route information, it means the server has started successfully.
The server program mounts several web routes and initializes `go-plugify`.

### Client

Enter the `client` directory and run `make run`.
The client program’s code is located in the `src` directory.
Its logic modifies the handler function corresponding to `http://localhost:8080/` on the server and calls the server to perform processing and return results.
For example, in `client/src/plugin.go`:


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

You can modify the code to explore and experience how it works.