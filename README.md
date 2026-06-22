# GoLazy Sample App

This repository is a small GoLazy application. It demonstrates:

- application dependencies initialized through `context.Context`
- pooled controller instances with request-local render state
- application helpers registered with the view renderer
- secure cookie session configuration through `lazyapp.Config.Sessions`
- embedded production views, local development views, public files, and
  Markdown posts
- fingerprinted asset URLs through `asset_path`, immutable cache policy for
  permanent asset URLs, and asset ETags
- Tailwind stylesheet compilation through `lazy tailwind`
- JavaScript library and app-module bundling through `lazy js`
- application-level HTTP integration tests
- single-binary deployment

## Requirements

- mise when using the provided development toolchain
- Go 1.26 or later
- Node.js and npm when regenerating JavaScript assets or Tailwind CSS

The provided mise toolchain installs Go and Node.js for local development.

When this repository is used inside the GoLazy workspace, the root `go.work`
resolves `golazy.dev` to the sibling framework checkout. The module itself does
not contain a local `replace` directive.

## Run

With mise installed:

```sh
mise trust
mise install
mise run dev
```

`mise trust` is a one-time local approval for `mise.toml`; mise requires it
because the config loads development environment variables. Project tasks live
as standalone scripts under `.mise/tasks`.

The sample also includes a standalone Go task script:

```sh
mise run hello
```

Mise discovers `hello` from `.mise/tasks/hello.go`. This is an example of using
Go for small project scripts without adding another command to `mise.toml`.

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

## Development Secrets

The sample app includes checked-in development values under `.secrets/`.
`mise.toml` installs Go, Node.js, and the `age`, `sops`, and `usage` tools,
then loads `.secrets/development.env` for commands run through mise.

`SECURE_COOKIE_KEY` configures the session cookie signing key. In development
the checked-in value is intentionally low ceremony. In production, the
deployment environment is responsible for providing `SECURE_COOKIE_KEY` and any
other application environment variables.

See [.secrets/README.md](.secrets/README.md) for the SOPS and age workflow.

## Routes

Inspect the application route table:

```sh
lazy routes
```

| Method | Path               | Description                   |
|--------|--------------------|-------------------------------|
| `GET`  | `/`                | Home page                     |
| `GET`  | `/posts`           | List embedded posts as HTML or Markdown |
| `GET`  | `/posts/{post_id}` | Render an embedded post as HTML or Markdown |
| `GET`  | `/styles.css`      | Serve an embedded public asset |
| `GET`  | `/styles-*.css`    | Serve a fingerprinted asset permalink |
| `GET`  | `/assets/importmap.json` | Serve the generated JavaScript importmap |

Other embedded public files are served from the application root. Templates can
use `asset_path` to link the permanent hashed URL for cacheable assets.
Send `Accept: text/markdown`, request `/posts.md`, or request
`/posts/{post_id}.md` to receive the raw embedded Markdown instead of the HTML
page. The `.html` suffix keeps the normal HTML rendering.

## Project Structure

```text
.mise/tasks/         Standalone mise task scripts
app/
  controllers/       Controllers and request-local render hooks
  helpers/           Template helpers registered by the app
  js/                App JavaScript source for lazy js
  public/            Embedded public files and generated JavaScript assets
  styles/            Tailwind input stylesheets when lazy tailwind is enabled
  views/             Layouts and templates
cmd/app/             Application executable
init/                Application composition, dependencies, and routes
js.toml              JavaScript library entrypoints for lazy js
lib/markdown/        Markdown adapter
mise.toml            Development toolchain and local env loading
services/            Business services
.secrets/            Checked-in development secret examples
test/                Application integration tests
```

Shared dependencies are initialized once in `init/context.go`. Routes are
registered in `init/routes.go`. The application is assembled in `init/app.go`.
Routes construct controller prototypes at app startup. GoLazy borrows pooled
controller instances for each request and resets mutable render state before
reuse.

`init/app.go` also configures cookie-backed sessions. The app reads
`SECURE_COOKIE_KEY` from the environment. In development, mise loads the
checked-in example value from `.secrets/development.env`; production deployments
provide their own value through the runtime environment. The framework expands
short keys deterministically before signing cookies.

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
lazy tailwind
lazy js
go test ./...
go test -race ./...
go vet ./...
go build -o /tmp/sample-app ./cmd/app
```

## Docker

Refresh generated assets before building the image:

```sh
lazy tailwind
lazy js
docker build -t sample-app .
docker run --rm -p 127.0.0.1:3000:3000 sample-app
```

The image runs the compiled application with `ADDR=0.0.0.0:3000`.
Provide production secrets such as `SECURE_COOKIE_KEY` through the container
environment.

## License

This sample application is released under the MIT License. See
[LICENSE](LICENSE).
