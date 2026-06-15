# GoLazy Sample App

This repository is a small GoLazy application. It demonstrates:

- application dependencies initialized through `context.Context`
- request-local controllers and template rendering
- embedded production views, local development views, public files, and
  Markdown posts
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

The application listens on <http://localhost:8080>. Set `ADDR` to use another
port or address:

```sh
ADDR=3000 go run ./cmd/app
ADDR=127.0.0.1:3000 go run ./cmd/app
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
| `GET`  | `/styles.css`      | Serve an embedded public file |

Other embedded public files are served from the application root.

## Project Structure

```text
app/
  controllers/       Request-local controllers
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
