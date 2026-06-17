<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{if .title}}{{.title}} · {{end}}GoLazy</title>
    {{stylesheet "/styles.css"}}
    {{importmap "/assets/importmap.json"}}
    <script type="module" src="/javascript/application.js"></script>
  </head>
  <body class="min-h-screen bg-white text-zinc-950 antialiased">
    <div class="mx-auto flex min-h-screen w-full max-w-5xl flex-col px-5 py-6 sm:px-8 lg:px-10">
      <header class="flex flex-col gap-4 border-b border-zinc-200 pb-5 sm:flex-row sm:items-center sm:justify-between">
        <a class="inline-flex items-center gap-2 text-base font-semibold text-zinc-950" href="{{path_for "root"}}">
          <span class="h-3 w-3 rounded-sm bg-cyan-500"></span>
          GoLazy
        </a>
        <nav class="flex items-center gap-2 text-sm">
          <a class="rounded-md border border-transparent px-3 py-2 font-medium text-zinc-600 hover:border-cyan-200 hover:bg-cyan-50 hover:text-cyan-800" href="{{path_for "root"}}">Home</a>
          <a class="rounded-md border border-transparent px-3 py-2 font-medium text-zinc-600 hover:border-cyan-200 hover:bg-cyan-50 hover:text-cyan-800" href="{{path_for "posts"}}">Posts</a>
        </nav>
      </header>
      <main class="flex-1 py-10">
        {{.content}}
      </main>
      <footer class="border-t border-zinc-200 py-6 text-sm text-zinc-500">Rendered at {{.currentTime}}</footer>
    </div>
  </body>
</html>
