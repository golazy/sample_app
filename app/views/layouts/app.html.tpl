<!doctype html>
<html lang="{{seo_lang}}" class="dark scheme-dark">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    {{seo}}
    {{stylesheet "/styles.css"}}
    {{importmap "/assets/importmap.json"}}
    <script type="module">import "/js/app.js"</script>
  </head>
  <body class="min-h-screen bg-white text-zinc-950 antialiased dark:bg-zinc-950 dark:text-zinc-100">
    <div class="mx-auto flex min-h-screen w-full max-w-5xl flex-col px-5 py-6 sm:px-8 lg:px-10">
      <header class="flex flex-col gap-4 border-b border-zinc-200 pb-5 sm:flex-row sm:items-center sm:justify-between dark:border-zinc-800">
        <a class="inline-flex items-center gap-2 text-base font-semibold text-zinc-950 dark:text-white" href="{{path_for "root"}}">
          <span class="h-3 w-3 rounded-sm bg-cyan-500 dark:bg-cyan-400"></span>
          GoLazy
        </a>
        <nav class="flex items-center gap-2 text-sm">
          <a class="rounded-md border border-transparent px-3 py-2 font-medium text-zinc-600 hover:border-cyan-200 hover:bg-cyan-50 hover:text-cyan-800 dark:text-zinc-300 dark:hover:border-cyan-800 dark:hover:bg-cyan-950/40 dark:hover:text-cyan-200" href="{{path_for "root"}}">Home</a>
          <a class="rounded-md border border-transparent px-3 py-2 font-medium text-zinc-600 hover:border-cyan-200 hover:bg-cyan-50 hover:text-cyan-800 dark:text-zinc-300 dark:hover:border-cyan-800 dark:hover:bg-cyan-950/40 dark:hover:text-cyan-200" href="{{path_for "posts"}}">Posts</a>
        </nav>
      </header>
      <main class="flex-1 py-10">
        {{if .visitCount}}
          <p class="mb-4 text-sm text-zinc-500 dark:text-zinc-400">Visit count: {{.visitCount}}</p>
        {{end}}
        {{if .flashMessages}}
          <div class="mb-4 space-y-2">
            {{range .flashMessages}}
              <p class="rounded-md border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-900 dark:border-emerald-900/60 dark:bg-emerald-950/30 dark:text-emerald-100">
                {{.}}
              </p>
            {{end}}
          </div>
        {{end}}
        {{.content}}
      </main>
      <footer class="border-t border-zinc-200 py-6 text-sm text-zinc-500 dark:border-zinc-800 dark:text-zinc-400">Rendered at {{.currentTime}}</footer>
    </div>
  </body>
</html>
