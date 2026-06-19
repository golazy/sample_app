# GoLazy Application Instructions

This repository is a GoLazy application. Keep this file focused on durable
project conventions that should apply to coding agents and automation.

## Project Shape

- Application code lives under `app`.
- Controllers live in `app/controllers`.
- Services live in `app/services`.
- Helpers live in `app/helpers` and are registered from `init/app.go`.
- Views live in `app/views`; layouts live in `app/views/layouts`.
- Public assets live in `app/public` and are embedded into production builds.
- The executable entrypoint is `cmd/app`.
- Application composition lives in `init`: context, routes, and app config.
- HTTP integration tests belong in `test`.

## GoLazy Conventions

- Initialize shared dependencies once through `init.Context`.
- Register routes only through `func Draw(router *lazyroutes.Scope)` in
  `init/routes.go`.
- Controller constructors should receive only `context.Context` and return
  `(*Controller, error)`.
- Controllers are request-local. Do not share mutable render state between
  requests.
- Controller actions should return `error`; use the standard
  `http.ResponseWriter` and `*http.Request` signature unless route arguments
  materially simplify the action.
- Use `Set("name", value)` for template data. Template data is escaped by
  default; only trusted framework-generated HTML should become `template.HTML`.
- Keep production secrets out of source. Session keys should come from the
  environment once the app is configured for production.

## Assets

- App-owned browser JavaScript source lives in `app/js`. The main entry is
  `app/js/app.js`; Stimulus controllers live under `app/js/controllers`.
- `lazy js` manages JavaScript assets from `js.toml`, `package.json`,
  and the lockfile, expands GoLazy directives in `app/js/app.js`, and bundles
  app JavaScript into `app/public/assets/lazyshaft`. Commit generated
  importmaps and `app/public/assets` outputs, but do not commit `node_modules`.
- Tailwind source lives in `app/styles/application.css`.
- `lazy tailwind` compiles the public stylesheet at `app/public/styles.css`.
- Do not edit generated importmaps, lazyshaft bundles, or compiled CSS by hand
  when the source manifest, package files, or Tailwind source should be changed
  instead.

## Commands

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
