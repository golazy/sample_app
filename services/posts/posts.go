package posts

import (
	"context"
	"embed"
	"fmt"
)

//go:embed content/*.md
var content embed.FS

type Post struct {
	Param string
	Title string
	Body  string
}

type Service struct {
	posts []Post
	index map[string]int
}

type contextKey struct{}

var definitions = []struct {
	param string
	title string
	file  string
}{
	{param: "hello-golazy", title: "Hello, GoLazy", file: "content/hello-golazy.md"},
	{param: "embedded-content", title: "Embedded Content", file: "content/embedded-content.md"},
}

func New() (*Service, error) {
	service := &Service{
		posts: make([]Post, 0, len(definitions)),
		index: make(map[string]int, len(definitions)),
	}
	for _, definition := range definitions {
		body, err := content.ReadFile(definition.file)
		if err != nil {
			return nil, fmt.Errorf("read post %q: %w", definition.param, err)
		}
		service.index[definition.param] = len(service.posts)
		service.posts = append(service.posts, Post{
			Param: definition.param,
			Title: definition.title,
			Body:  string(body),
		})
	}
	return service, nil
}

func WithContext(ctx context.Context, service *Service) context.Context {
	return context.WithValue(ctx, contextKey{}, service)
}

func FromContext(ctx context.Context) (*Service, bool) {
	service, ok := ctx.Value(contextKey{}).(*Service)
	return service, ok
}

func (s *Service) List() []Post {
	return append([]Post(nil), s.posts...)
}

func (s *Service) Get(param string) (Post, bool) {
	position, ok := s.index[param]
	if !ok {
		return Post{}, false
	}
	return s.posts[position], true
}
