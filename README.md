# GoLazy Sample App

This repository is a small GoLazy application. It demonstrates:

- application dependencies initialized through `context.Context`
- request-local controllers and template rendering
- embedded views, public files, and Markdown posts
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
  init/              Dependency and route initialization
  public/            Embedded public files
  services/          Application services
  views/             Embedded layouts and templates
cmd/app/             Application executable
lib/markdown/        Markdown adapter
test/                Application integration tests
```

Shared dependencies are initialized once in `app/init/context.go`. Routes are
registered in `app/init/routes.go`. Each route constructs a controller for the
current request, so mutable render state is never shared between requests.

Controller views live at:

```text
app/views/<controller>/<action>.html.tpl
```

Layouts live at:

```text
app/views/layouts/<layout>.html.tpl
```

## Verify

```sh
go test ./...
go test -race ./...
go vet ./...
go build -o /tmp/sample-app ./cmd/app
```
