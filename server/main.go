package main

import (
	"net/http"
	"strconv"
	"sync"

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

	bookService := NewBookService()

	plugManager := goplugify.Init("default",
		goplugify.ComponentWithName("ginengine", ginRouter),
		goplugify.ComponentWithName("bookService", bookService),
	)

	registerCoreRoutes(r, bookService)

	plugManager.RegisterRoutes(ginRouter, "/api/v1")

	return r
}

func registerCoreRoutes(r *gin.Engine, svc *BookService) {
	r.GET("/", func(c *gin.Context) {
		c.String(200, "📚 Welcome to Go-Plugify Book Manager!")
	})

	r.GET("/api/v1/books", func(c *gin.Context) {
		c.JSON(http.StatusOK, svc.ListBooks())
	})

	r.POST("/api/v1/books", func(c *gin.Context) {
		var book Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		svc.AddBook(book)
		c.JSON(http.StatusOK, gin.H{"msg": "book added"})
	})

	r.PUT("/api/v1/books/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var book Book
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

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type BookService struct {
	mu    sync.Mutex
	books map[int]Book
	next  int
}

func NewBookService() *BookService {
	return &BookService{
		books: map[int]Book{},
		next:  1,
	}
}

func (s *BookService) ListBooks() []Book {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]Book, 0, len(s.books))
	for _, b := range s.books {
		result = append(result, b)
	}
	return result
}

func (s *BookService) AddBook(b Book) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b.ID = s.next
	s.next++
	s.books[b.ID] = b
}

func (s *BookService) UpdateBook(id int, b Book) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.books[id]; !ok {
		return false
	}
	b.ID = id
	s.books[id] = b
	return true
}

func (s *BookService) DeleteBook(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.books[id]; !ok {
		return false
	}
	delete(s.books, id)
	return true
}
