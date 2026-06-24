# GoLazy Sample App

This repository is a small GoLazy application. It demonstrates:

- application dependencies initialized through `context.Context`
- one small service, `helloworldservice`, used by the home controller
- one resource-backed application route at `/`
- embedded production views, local development views, and public files
- fingerprinted asset URLs through `asset_path`, immutable cache policy for
  permanent asset URLs, and asset ETags
- Tailwind utility-class styling compiled through `lazy tailwind`
- Hotwire Turbo and Stimulus controller bundling through `lazy js`
- application-level HTTP integration tests
- single-binary deployment

## Requirements

- Go 1.26 or later
- mise for the standard development environment
- Node.js and npm when regenerating JavaScript assets or Tailwind CSS outside mise

The provided mise toolchain installs Node.js and project CLI helpers for local
development. Go is not installed through mise because Go already bundles
multi-version support through the module `go` directive and toolchain
selection.

When this repository is used inside the GoLazy workspace, the root `go.work`
resolves `golazy.dev` to the sibling framework checkout. The module itself does
not contain a local `replace` directive.

## Run

With mise installed:

```sh
mise trust
mise install
mise exec -- lazy
```

`mise trust` is a one-time local approval for `mise.toml`; mise requires it
because the config loads development environment variables. Project tasks live
as standalone scripts under `.mise/tasks`.

The sample also includes a standalone Go task script:

```sh
mise run hello
```

Mise discovers `hello` from `.mise/tasks/hello.go`. This is an example of using
Go for small project scripts without adding another command to `mise.toml`; it
exists to teach the pattern, not because the sample app needs a helper command.

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

The sample app demonstrates two classes of development environment values:

- Ordinary checked-in values live directly in `mise.toml`, such as
  `SAMPLE_APP_ENV`.
- Secret-shaped checked-in examples live in `.secrets/development.env`, such as
  `SAMPLE_APP_DEVELOPMENT_SECRET`.

Production secrets should come from the deployment environment.

See [.secrets/README.md](.secrets/README.md) for the SOPS and age workflow.
It includes `secrets:new-key`, `secrets:add-key`, `secrets:remove-user`, and
`secrets:users` tasks for managing the public recipients that can decrypt
encrypted development secret files. The task wrappers delegate to
`.mise/tasks/secrets/_lib.sh`, keeping development-secret tooling with the
task namespace and separate from application Go code.

## Routes

Inspect the application route table:

```sh
lazy routes
```

| Method | Path               | Description                   |
|--------|--------------------|-------------------------------|
| `GET`  | `/`                | Home page                     |
| `GET`  | `/styles.css`      | Serve an embedded public asset |
| `GET`  | `/styles-*.css`    | Serve a fingerprinted asset permalink |
| `GET`  | `/assets/importmap.json` | Serve the generated JavaScript importmap |

Other embedded public files are served from the application root. Templates can
use `asset_path` to link the permanent hashed URL for cacheable assets.

## Project Structure

```text
.mise/tasks/         Standalone mise task scripts
app/
  controllers/       Shared base controller and concrete controller packages
  js/                App JavaScript source for lazy js
  public/            Embedded public files and generated JavaScript assets
  styles/            Tailwind input stylesheet
  views/             Layouts and templates
cmd/app/             Application executable
init/                Application composition, dependencies, and routes
js.toml              JavaScript library entrypoints for lazy js
mise.toml            Development toolchain and local env loading
services/            Business services
.secrets/            Checked-in development secret examples and public recipients
test/                Application integration tests
```

Shared dependencies are initialized once in `init/dependencies.go`. Routes are
registered in `init/routes.go`. The application is assembled in `init/app.go`.
Routes construct controller prototypes at app startup. GoLazy borrows pooled
controller instances for each request and resets mutable render state before
reuse.

`init/dependencies.go` registers `helloworldservice`. The home controller reads
that service from context and uses `Hello()` as the page title.

Controller views live at:

```text
app/views/<controller>/<action>.html.tpl
```

Layouts live at:

```text
app/views/layouts/<layout>.html.tpl
```

The framework supplies the shared `app/error` view. Add
`app/views/app/error.html.tpl` only when the sample app should demonstrate an
application-specific error page override.

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
Pass runtime environment variables through the container environment when the
application grows to need them.

## License

This sample application is released under the MIT License. See
[LICENSE](LICENSE).
