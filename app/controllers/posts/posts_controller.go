package posts

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"golazy.dev/lazycontroller"
	"sample_app/app/controllers"
	postservice "sample_app/app/services/posts"
	"sample_app/lib/markdown"
)

type PostsController struct {
	controllers.BaseController
	posts *postservice.Service
}

func New(ctx context.Context) (*PostsController, error) {
	base, err := controllers.NewBaseController(ctx, "posts")
	if err != nil {
		return nil, err
	}
	posts, ok := postservice.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("posts service is missing from application context")
	}
	return &PostsController{BaseController: base, posts: posts}, nil
}

func (c *PostsController) Index(
	_ http.ResponseWriter,
	_ *http.Request,
) error {
	c.Set("title", "Posts")
	c.Set("posts", c.posts.List())
	return c.Render("index")
}

func (c *PostsController) Show(
	_ http.ResponseWriter,
	r *http.Request,
) error {
	slug := r.PathValue("post_id")
	post, ok := c.posts.Get(slug)
	if !ok {
		return lazycontroller.Error(
			http.StatusNotFound,
			fmt.Errorf("post %q not found", slug),
		)
	}

	body, err := markdown.Convert(post.Body)
	if err != nil {
		return fmt.Errorf("render post markdown: %w", err)
	}

	c.Set("title", post.Title)
	c.Set("post", post)
	c.Set("body", template.HTML(body))
	return c.Render("show")
}
