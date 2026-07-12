# Cheatsheet Snippets

The public GoLazy cheatsheet is a generated-app maintenance surface. The latest
catalog page lives in the website repository at:

```text
apps/golazy.dev/app/content/cheatsheets/latest/_index.md
```

Individual latest examples live beside it as
`apps/golazy.dev/app/content/cheatsheets/latest/<slug>.md`, and the grouped
catalog lives in `apps/golazy.dev/app/data/cheatsheets_latest.toml`.
Released snapshots live under
`apps/golazy.dev/app/content/cheatsheets/<version>/<slug>.md` with a matching
versioned catalog. The public catalog route is `/cheatsheets/<version>`, and
each example route is `/cheatsheets/<version>/<slug>`. `/how-tos` redirects
there for search intent, but snippets should use the cheatsheet name.

Use it when a task changes a public API, generated-app convention, or common
sample-app pattern that users are likely to copy.

## Snippet Contract

- Snippets answer "how can I?" for one narrow task.
- Snippets show changes from the generated sample app, not a separate demo app.
- Snippets should be complete enough to copy for that task.
- Code blocks should contain code and comments, not explanatory paragraphs.
- Each snippet gets its own page. Keep the body code-first; the website renders
  the code on the left and uses the frontmatter `description` as the right-side
  commentary.
- Each snippet should point readers toward the fuller guide through the page's
  `See Others` links instead of duplicating guide prose.

## Maintenance Loop

When changing framework behavior, CLI behavior, generated app structure,
controller conventions, forms, SEO, caching, Turbo, assets, mailers, jobs,
services, PWA, MCP, or the sample app skill itself:

1. Update the focused latest guide.
2. Update `cheatsheets/latest/<slug>.md` when a snippet exists or should be
   added, and update `app/data/cheatsheets_latest.toml` when the catalog title,
   grouping, or summary changes.
3. Update this skill or its references when agents need the same convention.
4. Add a changelog note when the public guide surface changes.
