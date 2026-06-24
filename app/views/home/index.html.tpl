<main class="min-h-screen overflow-hidden">
  <section class="grid min-h-screen grid-rows-[auto_1fr_auto]">
    <header class="flex items-center justify-between border-b-2 border-[#080808] px-5 py-5 sm:px-8 lg:px-12">
      <a class="text-xl font-black tracking-tight" href="{{path_for "home"}}">GoLazy</a>
      <span class="rounded-full border-2 border-[#080808] bg-[#fddd00] px-4 py-2 text-xs font-black uppercase tracking-[0.16em]">sample</span>
    </header>

    <div class="grid items-center gap-10 px-5 py-12 sm:px-8 lg:grid-cols-[minmax(0,1.05fr)_minmax(20rem,0.95fr)] lg:px-12">
      <div class="max-w-4xl">
        <p class="mb-5 inline-flex border-2 border-[#080808] bg-white px-3 py-2 text-xs font-black uppercase tracking-[0.18em] shadow-[6px_6px_0_#00add8]">service driven</p>
        <h1 class="text-[clamp(5rem,19vw,15rem)] font-black uppercase leading-[0.78] tracking-normal">
          <span
            data-controller="hello"
            data-hello-target="output"
            data-hello-text-value="{{.title}}"
            data-hello-speed-value="95"
          >{{.title}}</span><span class="text-[#ce3262]" aria-hidden="true">.</span>
        </h1>
        <p class="mt-8 max-w-2xl font-sans text-lg font-semibold leading-8 sm:text-xl">
          A compact GoLazy application with one service, one controller, one
          resource-backed route, Tailwind styles, and a small Stimulus effect.
        </p>
      </div>

      <div class="border-2 border-[#080808] bg-[#080808] p-3 shadow-[12px_12px_0_#ce3262]">
        <div class="flex items-center gap-2 border-b border-[#f7f4eb]/30 px-2 pb-3">
          <span class="h-3 w-3 rounded-full bg-[#ce3262]"></span>
          <span class="h-3 w-3 rounded-full bg-[#fddd00]"></span>
          <span class="h-3 w-3 rounded-full bg-[#00add8]"></span>
          <span class="ml-auto font-mono text-xs text-[#f7f4eb]/70">home_controller.go</span>
        </div>
        <pre class="overflow-x-auto p-4 font-mono text-sm leading-7 text-[#f7f4eb]"><code><span class="text-[#ce3262]">func</span> (c *HomeController) <span class="text-[#00add8]">Index</span>() error {
  c.<span class="text-[#00add8]">Set</span>(<span class="text-[#fddd00]">"title"</span>, c.helloService.<span class="text-[#00add8]">Hello</span>())
  <span class="text-[#ce3262]">return</span> nil
}</code></pre>
      </div>
    </div>

    <footer class="grid gap-3 border-t-2 border-[#080808] px-5 py-5 text-sm font-black uppercase tracking-[0.14em] sm:grid-cols-3 sm:px-8 lg:px-12">
      <span>service: helloworldservice</span>
      <span>controller: home</span>
      <span>route: GET /</span>
    </footer>
  </section>
</main>
