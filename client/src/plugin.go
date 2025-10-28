package main

import (
	"context"

	"plugin20251014202027/pkg"
)

func init() {
	ExportPlugin.Name = "example"
	ExportPlugin.Description = "An example plugin"
	ExportPlugin.Version = "v0.0.1"
}

// Run is called when the plugin is executed.
func (p Plugin) Run(args any) {
	ctx := args.(HttpContext)
	p.Logger().Info("Plugin %s is running, ctx %+v", p.Name, ctx)
	p.GetGinEngine().ReplaceHandler("GET", "/", func(ctx context.Context) {
		ctx.(HttpContext).JSON(200, map[string]string{"message": "Hello from example plugin !!!"})
	})
	bookService := NewBookService(p.Component("bookService"), p)
	bookService.AddBook(Book{ID: 1, Title: "1984", Author: "George Orwell"})
	bookService.AddBook(Book{ID: 2, Title: "To Kill a Mockingbird", Author: "Harper Lee"})
	bookService.DeleteBook(1)
	ctx.JSON(200, map[string]any{
		"message":  "Plugin executed successfully",
		"load pkg": pkg.SayHello(),
		"books":    bookService.ListBooks(),
	})
}

// Methods returns a map of method names to functions that can be called on the plugin.
func (p Plugin) Methods() map[string]func(any) any {
	return map[string]func(any) any{
		"add": func(args any) any {
			nums := args.([]int)
			if len(nums) != 2 {
				return 0
			}
			return nums[0] + nums[1]
		},
	}
}

// ==============================
// custom interface
// ==============================

func (p Plugin) GetGinEngine() HttpRouter {
	return p.Component("ginengine").(HttpRouter)
}

type BookService struct {
	service any
	plug    Plugin
}

func NewBookService(service any, plug Plugin) *BookService {
	return &BookService{
		service: service,
		plug:    plug,
	}
}

func (s *BookService) ListBooks() []Book {
	resp := s.plug.CallIgnore(s.service, "ListBooks")
	books := make([]Book, 0)
	s.plug.ConvertTo(resp[0], &books)
	return books
}

func (s *BookService) AddBook(b Book) {
	s.plug.CallIgnore(s.service, "AddBook", b)
}

func (s *BookService) UpdateBook(id int, b Book) bool {
	resp := s.plug.CallIgnore(s.service, "UpdateBook", id, b)
	return resp[0].(bool)
}

func (s *BookService) DeleteBook(id int) bool {
	resp := s.plug.CallIgnore(s.service, "DeleteBook", id)
	return resp[0].(bool)
}

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}
