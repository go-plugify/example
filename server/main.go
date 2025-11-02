package main

import (
	"net/http"
	"strconv"

	"example.com/server/service"

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

	ginRouter := ginadapter.NewHttpRouter(r)

	bookService := service.NewBookService()

	plugManager := goplugify.InitPluginManagers("default",
		goplugify.ComponentWithName("ginengine", ginRouter),
		goplugify.ComponentWithName("bookService", bookService),
		goplugify.ComponentWithName("allKindBook", new(service.AllKindBook)),
	)

	registerCoreRoutes(r, plugManager, bookService)

	goplugify.InitHTTPServer(plugManager).RegisterRoutes(ginRouter, "/api/v1")

	return r
}

func registerCoreRoutes(r *gin.Engine, plugManager goplugify.PluginManagers, svc *service.BookService) {
	r.GET("/", func(c *gin.Context) {
		c.String(200, "📚 Welcome to Go-Plugify Book Manager!")
	})

	r.GET("/hello", func(c *gin.Context) {
		examplePlug, err := plugManager["default"].GetPlugin("example")
		if err != nil {
			c.String(500, "Error retrieving plugin: %v", err)
			return
		}
		helloFn, ok := examplePlug.Method("hello")
		if !ok {
			c.String(500, "Plugin method 'hello' not found")
			return
		}
		resp := helloFn(nil)
		c.JSON(200, gin.H{"message": resp})
	})

	r.GET("/api/v1/books", func(c *gin.Context) {
		c.JSON(http.StatusOK, svc.ListBooks())
	})

	r.POST("/api/v1/books", func(c *gin.Context) {
		var book service.Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		svc.AddBook(book)
		c.JSON(http.StatusOK, gin.H{"msg": "book added"})
	})

	r.PUT("/api/v1/books/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var book service.Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if svc.UpdateBook(id, book) {
			c.JSON(http.StatusOK, gin.H{"msg": "book updated"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"msg": "book not found"})
		}
	})

	r.DELETE("/api/v1/books/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if svc.DeleteBook(id) {
			c.JSON(http.StatusOK, gin.H{"msg": "book deleted"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"msg": "book not found"})
		}
	})
}
