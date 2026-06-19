<section class="grid gap-8 lg:grid-cols-[1.15fr_0.85fr] lg:items-start">
  <div class="space-y-6">
    <div class="space-y-4">
      <p class="text-sm font-medium text-cyan-700 dark:text-cyan-300">Sample application</p>
      <h1 class="max-w-3xl text-4xl font-semibold text-zinc-950 sm:text-5xl dark:text-white">Hello, world!</h1>
      <p class="max-w-2xl text-lg leading-8 text-zinc-600 dark:text-zinc-300">This application is running on the GoLazy example framework.</p>
    </div>
    <a class="inline-flex rounded-md bg-zinc-950 px-4 py-2 text-sm font-medium text-white hover:bg-cyan-700 dark:bg-cyan-400 dark:text-zinc-950 dark:hover:bg-cyan-300" href="{{path_for "posts"}}">Read the embedded posts</a>
    <a class="inline-flex rounded-md border border-cyan-200 px-4 py-2 text-sm font-medium text-cyan-700 hover:bg-cyan-50 dark:border-cyan-900 dark:text-cyan-200 dark:hover:bg-cyan-900/40" href="/flash">Trigger a sample flash</a>
  </div>

  <aside class="rounded-md border border-cyan-200 bg-cyan-50 p-5 text-sm leading-6 text-cyan-950 dark:border-cyan-900/60 dark:bg-cyan-950/30 dark:text-cyan-100">
    <p class="font-medium text-cyan-900 dark:text-cyan-200">Asset pipeline</p>
    <p class="mt-2">The layout stylesheet is linked through an asset permalink:</p>
    <code class="mt-3 block overflow-x-auto rounded-sm bg-white px-3 py-2 text-cyan-800 dark:bg-zinc-900 dark:text-cyan-200">{{asset_path "/styles.css"}}</code>
  </aside>
</section>
