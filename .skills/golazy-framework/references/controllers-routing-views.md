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

## Generated Action And Hook Arguments

Use generator arguments when an action or `BeforeAction` hook should receive a
typed request-derived value instead of reading untyped state from the
controller:

```go
type User string

func (c *BaseController) GenUser(
	r *http.Request,
) (*User, error) {
	userID, err := userIDFromSession(r)
	if err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, lazycontroller.Error(
			http.StatusUnauthorized,
			fmt.Errorf("login required"),
		)
	}
	user := User(userID)
	return &user, nil
}

func (c *BaseController) BeforeAction(user *User) error {
	c.Layout("admin")
	c.Set("user", user)
	return nil
}

func (c *AdminController) Index() error {
	return nil
}
```

`GenX` methods may receive `context.Context`, `*http.Request`, route
parameters, and other generated values. They return `T` or `(T, error)`.
Generated values are cached by type inside the current action or hook call.

When a generator returns a non-nil error, GoLazy does not call the action or
hook that requested it. It passes the error through the normal controller error
path, including `HandleError(http.ResponseWriter, *http.Request, error)` when
the concrete controller or embedded base controller implements it. For
protected areas, prefer a typed `User` generator, generated-argument
`BeforeAction`, and `HandleError` redirect/status policy over setting a
`CurrentUser` string or adding auth-only user parameters to every action.

## Expected Errors

Return `lazycontroller.Error(status, err)` for expected HTTP failures such as
missing records or forbidden actions. Unexpected errors become `500`.

When the condition comes from a helper or service, return a typed or sentinel
error and let the application base controller decide the response:

```go
var ErrNotFound = errors.New("not found")

func (c *BaseController) HandleError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) error {
	if errors.Is(err, ErrNotFound) {
		err = lazycontroller.Error(http.StatusNotFound, err)
	}
	return c.Base.HandleError(w, r, err)
}
```

Helpers should not render content or write `404` responses directly. They
should return errors such as `ErrNotFound`; `HandleError` owns the mapping to
status codes, redirects, or custom error views.

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
