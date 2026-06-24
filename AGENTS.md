# GoLazy Application Instructions

This repository is a GoLazy application. Keep this file focused on durable
project conventions that should apply to coding agents and automation.

## Project Shape

- Application code lives under `app`.
- Shared controller behavior lives in `app/controllers/base_controller.go`.
- Concrete controllers live in package directories under `app/controllers`;
  the sample route is handled by `app/controllers/home_controller/homecontroller.go`.
- Business services live in top-level `services`, outside the web-facing
  `app` tree.
- Views live in `app/views`; layouts live in `app/views/layouts`.
- Public assets live in `app/public` and are embedded into production builds.
- The executable entrypoint is `cmd/app`.
- Application composition lives in `init`: dependencies, routes, and app config.
- HTTP integration tests belong in `test`.

## GoLazy Conventions

- Initialize shared dependencies once through `init.Dependencies`.
- Register routes only through `func Draw(router *lazyroutes.Scope)` in
  `init/routes.go`.
- Controller constructors should receive only `context.Context` and return
  `(*Controller, error)`.
- Concrete controllers should embed `controllers.BaseController` so shared
  behavior stays consistent. Add `app/views/app/error.html.tpl` only when this
  app should override the framework default error page.
- The home route sets `title` from `c.helloService.Hello()` and renders the
  sample home view.
- Controllers are request-local. Do not share mutable render state between
  requests.
- Controller actions should return `error`; use the standard
  `http.ResponseWriter` and `*http.Request` signature unless route arguments
  materially simplify the action.
- Use `Set("name", value)` for template data. Template data is escaped by
  default; only trusted framework-generated HTML should become `template.HTML`.
- Keep production secrets out of source. Put ordinary checked-in development
  values in `mise.toml` and secret-shaped checked-in examples in
  `.secrets/development.env`.

## Assets

- App-owned browser JavaScript source lives in `app/js`. The main entry is
  `app/js/app.js`; Stimulus controllers live under `app/js/controllers`.
  The default entry imports Hotwire Turbo and Stimulus through
  `// golazy:turbo` and `// golazy:stimulus`, with matching entrypoints in
  `js.toml`.
- `lazy js` manages JavaScript assets from `js.toml`, `package.json`,
  and the lockfile, expands GoLazy directives in `app/js/app.js`, and bundles
  app JavaScript into `app/public/assets/lazyshaft`. Commit generated
  importmaps and `app/public/assets` outputs, but do not commit `node_modules`.
- Tailwind source lives in `app/styles/application.css`.
- `lazy tailwind` compiles the public stylesheet at `app/public/styles.css`.
- Keep sample app styling in Tailwind utility classes rather than custom CSS;
  `app/styles/application.css` should normally stay as the Tailwind import.
- Do not edit generated importmaps, lazyshaft bundles, or compiled CSS by hand
  when the source manifest, package files, or Tailwind source should be changed
  instead.

## Commands

Keep project-specific mise tasks as standalone scripts under `.mise/tasks`;
`mise.toml` should stay focused on tool and environment configuration.
Mise is the standard development environment for this app template, but do not
add Go to `[tools]`; Go already bundles multi-version support through the
module `go` directive and toolchain selection.
Secret-recipient tasks live under `.mise/tasks/secrets`. Shared shell helpers
for that task namespace may live beside them as hidden support files such as
`.mise/tasks/secrets/_lib.sh`; do not add a separate `.mise/scripts`
convention. Keep public age recipients in `.secrets/recipients.txt`, keep
generated `.sops.yaml` recipient rules committed, and keep private age
identities under ignored `.secrets/keys`.

Start the development server:

```sh
lazy
```

Run the app without the GoLazy CLI:

```sh
go run ./cmd/app
```

Inspect routes:

```sh
lazy routes
```

Regenerate assets when inputs change:

```sh
lazy tailwind
lazy js
```

Run verification:

```sh
go test ./...
go test -race ./...
go vet ./...
```

When building an executable, choose an explicit output path so it does not
collide with the `app/` directory.

## Editing Expectations

- Prefer standard-library HTTP, templates, and small explicit dependencies.
- Keep framework-generic behavior in GoLazy packages, not in application code.
- Update application-level tests when routes, rendering, assets, sessions, or
  services change.
- Keep `README.md` aligned with user-visible app structure and commands.
