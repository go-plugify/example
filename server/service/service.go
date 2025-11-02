package service

import "sync"

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

type FictionBook struct{ Book }
type ScienceBook struct{ Book }
type HistoryBook struct{ Book }

type AllKindBook struct {
	FictionBook FictionBook
	ScienceBook ScienceBook
	HistoryBook HistoryBook
}
