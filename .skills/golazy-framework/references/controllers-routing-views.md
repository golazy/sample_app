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

## Generated Action Arguments

Use generator arguments when an action should receive a typed request-derived
value instead of reading untyped state from the controller:

```go
type User struct {
	ID string
}

func (c *BaseController) GenUser(
	ctx context.Context,
	r *http.Request,
) (User, error) {
	user, err := userFromSession(ctx, r)
	if err != nil {
		return User{}, err
	}
	if user.ID == "" {
		return User{}, lazycontroller.Error(
			http.StatusUnauthorized,
			fmt.Errorf("login required"),
		)
	}
	return user, nil
}

func (c *AdminController) Index(user User) error {
	c.Set("user", user)
	return nil
}
```

`GenX` methods may receive `context.Context`, `*http.Request`, route
parameters, and other generated values. They return `T` or `(T, error)`.
Generated values are cached by type for the current request.

When a generator returns a non-nil error, GoLazy does not call the action. It
passes the error through the normal controller error path, including
`HandleError(http.ResponseWriter, *http.Request, error)` when the concrete
controller or embedded base controller implements it. For protected areas,
prefer a typed `User` generator plus a `HandleError` redirect/status policy over
setting a `CurrentUser` string in `BeforeAction`.

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
