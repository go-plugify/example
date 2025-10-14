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
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})
	ginrouters := ginadapter.NewHttpRouter(r)
	goplugify.Init("default", goplugify.ComponentWithName("ginengine", ginrouters)).RegisterRoutes(ginrouters, "/api/v1")
	return r
}