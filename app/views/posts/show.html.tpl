<article class="mx-auto max-w-3xl space-y-6">
  <a class="text-sm font-medium text-cyan-700 hover:text-cyan-900" href="{{path_for "posts"}}">Back to posts</a>
  <div class="space-y-3 border-b border-zinc-200 pb-6">
    <h1 class="text-3xl font-semibold text-zinc-950 sm:text-4xl">{{.post.Title}}</h1>
    <p class="text-sm text-zinc-500">{{word_count .post.Body}} words &middot; {{read_time .post.Body}} min read</p>
  </div>
  <div class="space-y-4 leading-8 text-zinc-700 [&_a]:text-cyan-700 [&_a]:underline [&_strong]:font-semibold [&_strong]:text-zinc-950">{{.body}}</div>
</article>
