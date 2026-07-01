# Controllers, Routing, And Views

Controllers are request-local. Routes point at controller constructors, not
shared controller instances.

## Controller Shape

```go
type PostsController struct {
	controllers.BaseController
	posts postservice.Service
}

func New(ctx context.Context) (*PostsController, error) {
	base, err := controllers.NewBaseController(ctx)
	if err != nil {
		return nil, err
	}
	posts, ok := postservice.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("posts service missing")
	}
	return &PostsController{BaseController: base, posts: posts}, nil
}
```

Constructor rules:

- Accept only `context.Context`.
- Return `(*Controller, error)`.
- Build the shared base controller first.
- Resolve services through typed context helpers.
- Fail early when a required service is missing.

## Action Shape

The standard action signature uses `net/http` values and returns `error`:

```go
func (c *PostsController) Index(_ http.ResponseWriter, _ *http.Request) error {
	c.Set("title", "Posts")
	c.Set("posts", c.posts.List())
	return nil
}
```

If the action returns `nil` without writing a response, GoLazy renders the
default view for the matched controller/action. Use `c.Render("other")` only
when an action intentionally renders a different view.

Use `c.Status`, `c.Header`, and `c.ContentType` to set metadata before the
automatic render. Use the raw `http.ResponseWriter` only when the action owns
the whole response body.

## Expected Errors

Return `lazycontroller.Error(status, err)` for expected HTTP failures such as
missing records or forbidden actions. Unexpected errors become `500`.

Add `app/views/app/error.html.tpl` only when the app needs to override the
framework default error page.

## Routes

Routes live in `init/routes.go`:

```go
func Draw(router *lazyroutes.Scope) {
	router.Get("/posts", postcontroller.New, (*postcontroller.PostsController).Index)
	router.Get("/posts/{post_id}", postcontroller.New, (*postcontroller.PostsController).Show)
}
```

Path parameters use the standard library `http.ServeMux` syntax and are read
with `r.PathValue("post_id")`.

Use `router.Resources(controller.New)` for conventional CRUD actions:

```text
GET    /posts                 Index
GET    /posts/new             New
POST   /posts                 Create
GET    /posts/{post_id}       Show
GET    /posts/{post_id}/edit  Edit
PATCH  /posts/{post_id}       Update
PUT    /posts/{post_id}       Update
DELETE /posts/{post_id}       Delete
```

`Resources` registers only actions that exist on the controller. Customize
resource path, singular, plural, parameter name, model mapping, collection
routes, member routes, or nested resources in the resource callback.

## Views And Layouts

Default controller views live at:

```text
app/views/<controller>/<action>.html.tpl
```

The default layout is:

```text
app/views/layouts/app.html.tpl
```

Template values come from controller `Set` calls:

```gotemplate
<h1>{{.title}}</h1>
```

Layouts receive the rendered controller view as `.content`.

Use `{{path_for "route_name"}}`, `{{stylesheet "/styles.css"}}`, and
`{{asset_path "/path"}}` instead of hard-coding generated paths.

Go templates escape values by default. Convert only trusted renderer output to
`template.HTML`.
