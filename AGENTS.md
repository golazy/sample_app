# GoLazy Application Instructions

This repository is a GoLazy application. Keep this file focused on durable
project conventions that should apply to coding agents and automation.

## Ownership

- Business and domain logic lives in top-level `services`. Services own use
  cases, domain models, business rules, data access, and domain-facing
  integrations.
- `app` is the presentation and communication layer. It owns HTTP controllers,
  views, forms, browser behavior, mailers, job adapters, and MCP adapters.
- Authorization may live at the controller boundary. Mailers may format and
  deliver messages. Jobs and MCP tools may translate their inputs. These
  adapters must call services for business behavior instead of duplicating it.
- `init` is the composition root for dependencies, routes, and app config. It
  must not contain business logic.
- Shared controller behavior lives in `app/controllers/base_controller.go`.
- Concrete controllers live in package directories under `app/controllers`;
  the sample route is handled by `app/controllers/home_controller/homecontroller.go`.
- Views live in `app/views`; layouts live in `app/views/layouts`.
- Public assets live in `app/public` and are embedded into production builds.
- The executable entrypoint is `cmd/app`.
- Application composition lives in `init`: dependencies, routes, and app config.
- HTTP integration tests belong in `test`.
- Optional local datasets live under `datasets/<name>/<service>.dump` when this
  app has development services with dump/load tasks.

## GoLazy Conventions

- Initialize shared dependencies once through `init.Dependencies`.
- Give each service a `New` constructor. Prefer the `lazydeps.Func[T]` lifecycle
  shape so `lazydeps.Service` can register it directly, install it in context,
  and run cleanup during dependency-ordered shutdown.
- Register routes only through `func Draw(router *lazyroutes.Scope)` in
  `init/routes.go`.
- Start with `router.Resources` and conventional `Index`, `New`, `Create`,
  `Show`, `Edit`, `Update`, and `Delete` actions. Add collection, member, and
  nested routes through the resource callback. Use top-level verb routes only
  for endpoints that are not meaningfully resources.
- Model multiple content pages as a resource with `Show`; select an explicit
  view with `Render` or `Variants` instead of adding one `Get` route and custom
  action per page.
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
- For submitted forms, prefer a typed `GenX` form generator such as
  `GenPasswordForm() (*PasswordForm, error)` that allocates the form and calls
  `c.Decode(form)`, then let the action receive `*PasswordForm`. Decode errors
  continue through `HandleError`.
- Use controller `SessionGet`, `SessionSet`, `SessionDelete`, `FlashSet`, and
  `FlashGet` for browser session and flash behavior. Keep auth wrapper types
  such as `AuthenticatedUser` in application code. Form validation should
  return normal joined errors from `lazyerrors.Validator(form)` and optional
  form-owned `Validate() error` methods. On validation failure, set
  `http.StatusUnprocessableEntity`, restore form/error template data, and
  render the form view; on success, flash if needed and redirect with a named
  route and `http.StatusSeeOther`.
- Use `Set("name", value)` for template data. Template data is escaped by
  default; only trusted framework-generated HTML should become `template.HTML`.
- Keep production secrets out of source. Put ordinary checked-in development
  values in `mise.toml` and secret-shaped checked-in examples in
  `.secrets/development.env`.

## Guide Contract

This repository is the template for `lazy new`. Its embedded skill must contain
the app-development conventions an agent needs without requiring files from a
sibling framework, website, workspace, or documentation repository.

Keep `AGENTS.md` as the concise shared source for generated-app guidance. If a
workflow grows too large for this file, split it into the current skill layout:
one `.skills/<name>/SKILL.md` entrypoint per workflow, with optional
`references/`, `scripts/`, `templates/`, or `examples/` folders beside it.
When `.skills/` exists, inspect `.skills/*/SKILL.md` before task-specific work
and use the matching skill instead of duplicating the workflow here.
Use `.skills/golazy-framework/SKILL.md` before adding framework-shaped
application behavior such as routes, controllers, services, views, assets,
sessions, jobs, mailers, storage, or tests.

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
Use committed `mise.toml` values for shared development defaults. Use
`mise.local.toml` or ignored env files for machine-specific ports, paths, and
experiments.
Secret-recipient tasks live under `.mise/tasks/secrets`. Shared shell helpers
for that task namespace may live beside them as hidden support files such as
`.mise/tasks/secrets/_lib.sh`; do not add a separate `.mise/scripts`
convention. Keep public age recipients in `.secrets/recipients.txt`, keep
generated `.sops.yaml` recipient rules committed, and keep private age
identities under ignored `.secrets/keys`.

When this app adds a local development service such as PostgreSQL or MinIO,
put lifecycle tasks under `.mise/tasks/<service>/`:

```text
<service>:start    # run in the foreground; Ctrl-C stops it
<service>:check    # exit 0 only when the service is ready
<service>:create   # create the local database, bucket, or schema if missing
<service>:migrate  # apply pending migrations
<service>:dump     # receive one output path
<service>:load     # receive one input path
```

Use `lazy.toml` to list services only when the app needs explicit service order
or selection. Otherwise `lazy` discovers services from `:start` task files.
When services exist, `lazy` starts them as managed subprocesses, waits for
`check`, runs `create` and `migrate`, and then starts the app. Ctrl-C stops the
app first and then the managed services.

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

Create or restore a local dataset when service dump/load tasks exist:

```sh
lazy dump baseline
lazy load baseline
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
