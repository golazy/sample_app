<article class="mx-auto max-w-3xl space-y-6">
  <a class="text-sm font-medium text-cyan-700 hover:text-cyan-900 dark:text-cyan-300 dark:hover:text-cyan-100" href="{{path_for "posts"}}">Back to posts</a>
  <div class="space-y-3 border-b border-zinc-200 pb-6 dark:border-zinc-800">
    <h1 class="text-3xl font-semibold text-zinc-950 sm:text-4xl dark:text-white">{{.post.Title}}</h1>
    <p class="text-sm text-zinc-500 dark:text-zinc-400">{{word_count .post.Body}} words &middot; {{read_time .post.Body}} min read</p>
  </div>
  <div class="space-y-4 leading-8 text-zinc-700 dark:text-zinc-300 [&_a]:text-cyan-700 dark:[&_a]:text-cyan-300 [&_a]:underline [&_strong]:font-semibold [&_strong]:text-zinc-950 dark:[&_strong]:text-white">{{.body}}</div>
</article>
