<section class="space-y-6">
  <div class="space-y-2">
    <p class="text-sm font-medium text-cyan-700">Embedded content</p>
    <h1 class="text-3xl font-semibold text-zinc-950 sm:text-4xl">Posts</h1>
    <p class="max-w-2xl text-zinc-600">These posts are Markdown files embedded into the sample application.</p>
  </div>

  <ul class="divide-y divide-zinc-200 overflow-hidden rounded-md border border-zinc-200 bg-white">
    {{range .posts}}
      <li>
        <a class="block px-4 py-4 font-medium text-zinc-800 hover:bg-cyan-50 hover:text-cyan-800" href="{{path_for "post" .Param}}">{{.Title}}</a>
      </li>
    {{else}}
      <li class="px-4 py-4 text-zinc-500">No posts are available.</li>
    {{end}}
  </ul>
</section>
