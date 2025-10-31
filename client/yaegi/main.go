package main

import (
	"plugify/plugify"
)

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

var Methods = map[string]func(any) any{}

func Destroy(any) error {
	return nil
}
