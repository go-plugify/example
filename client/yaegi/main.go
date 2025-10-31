package main

import (
	"context"
	"plugify/plugify"
)

func Run(input map[string]any) (any, error) {
	plugify.Logger.Info("Example plugin is running")
	plugify.Ginengine.ReplaceHandler("GET", "/", func(ctx context.Context) {
		plugify.Ginengine.NewHTTPContext(ctx).JSON(200, map[string]string{"message": "Hello from example plugin 2 !!!"})
	})
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

func Methods() map[string]func(any) any {
	return map[string]func(any) any{
		"hello": func(input any) any {
			plugify.Logger.Info("Hello from the 'hello' method!")
			return "Hello, World!"
		},
	}
}

func Destroy(input map[string]any) error {
	plugify.Logger.Info("Example plugin is being destroyed")
	return nil
}
