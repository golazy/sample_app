# GoLazy Sample App

This repository is a small GoLazy application. It demonstrates:

- application dependencies initialized through `context.Context`
- request-local controllers and template rendering
- application helpers registered with the view renderer
- secure cookie session configuration through `lazyapp.Config.Sessions`
- embedded production views, local development views, public files, and
  Markdown posts
- fingerprinted asset URLs through `asset_path`, immutable cache policy for
  permanent asset URLs, and asset ETags
- application-level HTTP integration tests
- single-binary deployment

## Requirements

- Go 1.26 or later

When this repository is used inside the GoLazy workspace, the root `go.work`
resolves `golazy.dev` to the sibling framework checkout. The module itself does
not contain a local `replace` directive.

## Run

With the GoLazy CLI installed:

```sh
lazy
```

Or run the application directly:

```sh
go run ./cmd/app
```

The application listens on <http://localhost:3000>. Set `ADDR` or `PORT` to use
another port or address:

```sh
PORT=4000 go run ./cmd/app
ADDR=127.0.0.1:4000 go run ./cmd/app
```

## Routes

Inspect the application route table:

```sh
lazy routes
```

| Method | Path               | Description                   |
|--------|--------------------|-------------------------------|
| `GET`  | `/`                | Home page                     |
| `GET`  | `/posts`           | List embedded posts           |
| `GET`  | `/posts/{post_id}` | Render an embedded post       |
| `GET`  | `/styles.css`      | Serve an embedded public asset |
| `GET`  | `/styles-*.css`    | Serve a fingerprinted asset permalink |

Other embedded public files are served from the application root. Templates can
use `asset_path` to link the permanent hashed URL for cacheable assets.

## Project Structure

```text
app/
  controllers/       Request-local controllers
  helpers/           Template helpers registered by the app
  public/            Embedded public files
  services/          Application services
  views/             Layouts and templates
cmd/app/             Application executable
init/                Application composition, dependencies, and routes
lib/markdown/        Markdown adapter
test/                Application integration tests
```

Shared dependencies are initialized once in `init/context.go`. Routes are
registered in `init/routes.go`. The application is assembled in `init/app.go`.
Each route constructs a controller for the current request, so mutable render
state is never shared between requests.

`init/app.go` also configures cookie-backed sessions. The template keeps a
short development `lazysession.Config.Key` in source with a TODO showing where
production apps should load `SECURE_COOKIE_KEY`; `lazy new` replaces that
template key with fresh random key material for every generated app. The
framework expands this short key deterministically before signing cookies.

Application helpers live in `app/helpers`. `init/app.go` registers them through
`lazyapp.Config.Helpers`, and the post view uses `word_count` and `read_time`
to render Markdown metadata.

Controller views live at:

```text
app/views/<controller>/<action>.html.tpl
```

Layouts live at:

```text
app/views/layouts/<layout>.html.tpl
```

Production builds embed views into the binary. The GoLazy CLI runs the app with
the `lazydev` build tag so templates are read from disk during local
development.

## Verify

```sh
go test ./...
go test -race ./...
go vet ./...
go build -o /tmp/sample-app ./cmd/app
```

## License

This sample application is released under the MIT License. See
[LICENSE](LICENSE).
