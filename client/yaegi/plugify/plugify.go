package plugify

import (
	"example.com/server/service"

	goplugify "github.com/go-plugify/go-plugify"
)

var BookService = new(service.BookService)

var Logger = new(goplugify.DefaultLogger)

type ServiceBook = service.Book
