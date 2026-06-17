<section class="grid gap-8 lg:grid-cols-[1.15fr_0.85fr] lg:items-start">
  <div class="space-y-6">
    <div class="space-y-4">
      <p class="text-sm font-medium text-cyan-700">Sample application</p>
      <h1 class="max-w-3xl text-4xl font-semibold text-zinc-950 sm:text-5xl">Hello, world!</h1>
      <p class="max-w-2xl text-lg leading-8 text-zinc-600">This application is running on the GoLazy example framework.</p>
    </div>
    <a class="inline-flex rounded-md bg-zinc-950 px-4 py-2 text-sm font-medium text-white hover:bg-cyan-700" href="{{path_for "posts"}}">Read the embedded posts</a>
  </div>

  <aside class="rounded-md border border-cyan-200 bg-cyan-50 p-5 text-sm leading-6 text-cyan-950">
    <p class="font-medium text-cyan-900">Asset pipeline</p>
    <p class="mt-2">The layout stylesheet is linked through an asset permalink:</p>
    <code class="mt-3 block overflow-x-auto rounded-sm bg-white px-3 py-2 text-cyan-800">{{asset_path "/styles.css"}}</code>
  </aside>
</section>
