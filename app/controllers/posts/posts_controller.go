package posts

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"golazy.dev/lazycontroller"
	"golazy.dev/lazyseo"
	"sample_app/app/controllers"
	"sample_app/lib/markdown"
	postservice "sample_app/services/posts"
)

var Markdown = lazycontroller.NewFormat(
	"text/markdown",
	lazycontroller.As("markdown"),
	lazycontroller.Suffix("md", "markdown"),
)

type PostsController struct {
	controllers.BaseController
	posts *postservice.Service
}

type postsIndexMetadata struct{}

func (postsIndexMetadata) Title() string {
	return "Posts"
}

func (postsIndexMetadata) Description() string {
	return "Read sample GoLazy posts served from embedded Markdown content."
}

func (postsIndexMetadata) Canonical() string {
	return "/posts"
}

type postMetadata struct {
	post postservice.Post
}

func (m postMetadata) Title() string {
	return m.post.Title
}

func (m postMetadata) Description() string {
	return description(m.post.Body)
}

func (m postMetadata) Canonical() string {
	return "/posts/" + m.post.Param
}

func (m postMetadata) Kind() lazyseo.PageKind {
	return lazyseo.Article
}

func New(ctx context.Context) (*PostsController, error) {
	base, err := controllers.NewBaseController(ctx)
	if err != nil {
		return nil, err
	}
	posts, ok := postservice.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("posts service is missing from application context")
	}
	return &PostsController{BaseController: base, posts: posts}, nil
}

func (c *PostsController) Index(w http.ResponseWriter, _ *http.Request) error {
	return c.Wants(lazycontroller.Formats{
		lazycontroller.HTML: func() error {
			c.Metadata(postsIndexMetadata{})
			c.Set("posts", c.posts.List())
			return nil
		},
		Markdown: func() error {
			return c.renderMarkdownIndex(w)
		},
	})
}

func (c *PostsController) Show(w http.ResponseWriter, r *http.Request) error {
	slug := r.PathValue("post_id")
	post, ok := c.posts.Get(slug)
	if !ok {
		return lazycontroller.Error(
			http.StatusNotFound,
			fmt.Errorf("post %q not found", slug),
		)
	}

	return c.Wants(lazycontroller.Formats{
		lazycontroller.HTML: func() error {
			return c.renderHTMLPost(post)
		},
		Markdown: func() error {
			return c.renderMarkdownPost(w, post)
		},
	})
}

func (c *PostsController) renderHTMLPost(post postservice.Post) error {
	body, err := markdown.Convert(post.Body)
	if err != nil {
		return fmt.Errorf("render post markdown: %w", err)
	}

	c.Metadata(postMetadata{post: post})
	c.Set("post", post)
	c.Set("body", template.HTML(body))
	return nil
}

func description(body string) string {
	const limit = 160
	words := strings.Fields(body)
	var out strings.Builder
	for _, word := range words {
		if out.Len() > 0 && out.Len()+1+len(word) > limit {
			break
		}
		if out.Len() > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(word)
	}
	return out.String()
}

func (c *PostsController) renderMarkdownIndex(w http.ResponseWriter) error {
	c.ContentType("text/markdown; charset=utf-8")
	if _, err := fmt.Fprintln(w, "# Posts"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	for _, post := range c.posts.List() {
		if _, err := fmt.Fprintf(w, "- [%s](/posts/%s)\n", post.Title, post.Param); err != nil {
			return err
		}
	}
	return nil
}

func (c *PostsController) renderMarkdownPost(w http.ResponseWriter, post postservice.Post) error {
	c.ContentType("text/markdown; charset=utf-8")
	_, err := w.Write([]byte(post.Body))
	return err
}
