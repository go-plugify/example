# go-plugify Example

This is a simple example demonstrating how to use `go-plugify`, including both the server and client parts.

## Dependencies (for native Go plugin mode)

- Requires `Linux` or `macOS`
- Requires Go version 1.23.10 or above
- Requires CGO support

## Run

<img alt="example" src="https://github.com/go-plugify/example/blob/main/example.gif?raw=true" width="651">

### Server

Go to the `server` directory and run `go run main.go`. Once you see the startup message and route information printed, it means the server has started successfully.  
The server mounts several web routes and initializes `go-plugify` with two components: `ginengine` and `bookService`, which are exposed for plugins to call.

```go
func setupRouter() *gin.Engine {
	r := gin.Default()

	ginRouter := ginadapter.NewHttpRouter(r)
	bookService := service.NewBookService()

	// Initialize plugin manager and mount components
	plugManager := goplugify.InitPluginManagers("default",
		goplugify.ComponentWithName("ginengine", ginRouter),
		goplugify.ComponentWithName("bookService", bookService),
	)

	registerCoreRoutes(r, bookService)

	// Register plugin routes — you can add authentication or middleware here as needed
	goplugify.InitHTTPServer(plugManager).RegisterRoutes(ginRouter, "/api/v1")

	return r
}
```

### Client

Go to the `client` directory. Two modes are supported:

- **Yaegi interpreter** (recommended)
- **Native Go plugin**

#### Yaegi

Enter the `yaegi` directory and run `make run`.  
The logic is in `main.go`. Each script must implement the following three functions or variables:

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

Run `make help` to see more commands.

#### Native plugin

Enter the `go_plugin` directory and run `make run`.  
The client-side logic is in the `src` directory. It modifies the handler for `http://localhost:8080/` and interacts with the server components for computation and response.  
For example, in `client/src/plugin.go`:

```go
// Run is called when the plugin is executed.
func (p Plugin) Run(args any) (any, error) {
	p.Logger().Info("Plugin %s is running", p.Name)
	// Dynamically replace HTTP route handler
	p.GetGinEngine().ReplaceHandler("GET", "/", func(ctx context.Context) {
		ctx.(HttpContext).JSON(200, map[string]string{"message": "Hello from example plugin !!!"})
	})
	// Access host component and call its methods
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

You can freely modify and experiment with these usage patterns.

Run `make help` to see more commands.

### Deployment & Execution

The client deploys and runs plugins via HTTP requests to the server.  
You can find the corresponding `curl` command in the `Makefile`:

```shell
curl --insecure -s -X POST \
	-F "file=@build/patch.so" \
	-F "meta={\"id\": \"example\", \"name\": \"Example Plugin\", \"description\": \"An example plugin\", \"author\": \"Your Name\", \"version\": \"v0.0.1\", \"loader\": \"native_plugin_http\"}"  \
	"http://localhost:8080/api/v1/plugin/init"
```

The request uses multipart form data:
- `file` contains the plugin logic (a `.so` file for native Go mode or a script for Yaegi mode).
- `meta` contains the plugin metadata with the following fields:

| Field | Type | Description | Example |
| --- | --- | --- | --- |
| id | string | Unique plugin ID | example |
| name | string | Plugin name | Example Plugin |
| description | string | Plugin description | An example plugin |
| author | string | Plugin author | Jack |
| version | string | Plugin version, format: vx.x.x | v0.0.1 |
| loader | enum | Plugin loader type (supported: native_plugin_http / yaegi_http) | native_plugin_http |

Finally, the API endpoint corresponds to the registered route on the server:  
`<domain>:<port>/<prefix>/plugin/init`.