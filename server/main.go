package main

import (
	"github.com/gin-gonic/gin"
	goplugify "github.com/go-plugify/go-plugify"
	ginadapter "github.com/go-plugify/webadapters/gin"
)

func main() {
	router := setupRouter()
	router.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	ginrouters := ginadapter.NewHttpRouter(r)
	plugManager := goplugify.Init("default",
		goplugify.ComponentWithName("ginengine", ginrouters),
		goplugify.ComponentWithName("calculator", &Caclulator{}),
	)
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})
	r.GET("/add", func(c *gin.Context) {
		// call plugin method "add"
		examplePlug, exist := plugManager["default"].GetPlugins().Get("example")
		if exist {
			c.JSON(200, gin.H{"result": examplePlug.CallMethod("add", []int{5, 10})})
			return
		}
		c.JSON(200, gin.H{"result": 1 + 1})
	})
	plugManager.RegisterRoutes(ginrouters, "/api/v1")
	return r
}

// example custom component

type Caclulator struct {
}

func (c *Caclulator) Add(a, b int) int {
	return a + b
}

func (c *Caclulator) Sub(a, b int) int {
	return a - b
}

func (c *Caclulator) Mul(a, b int) int {
	return a * b
}

func (c *Caclulator) Div(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}
