<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{if .title}}{{.title}} · {{end}}GoLazy</title>
    <link rel="stylesheet" href="/styles.css">
  </head>
  <body>
    <nav>
      <a href="/">Home</a>
      <a href="/posts">Posts</a>
    </nav>
    <main>
      {{$content := .content}}
      {{$content}}
    </main>
    <footer>Rendered at {{.currentTime}}</footer>
  </body>
</html>
