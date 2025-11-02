package plugify

import (
	"example.com/server/service"

	goplugify "github.com/go-plugify/go-plugify"
	ginadapter "github.com/go-plugify/webadapters/gin"
)

var BookService = new(service.BookService)

var Logger = new(goplugify.DefaultLogger)

var Util = new(goplugify.Util)

var Ginengine = new(ginadapter.HttpRouter)

type ServiceBook = service.Book

type ServiceFictionBook = service.FictionBook
