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
  <body>
    <nav>
      <a href="{{path_for "root"}}">Home</a>
      <a href="{{path_for "posts"}}">Posts</a>
    </nav>
    <main>
      {{.content}}
    </main>
    <footer>Rendered at {{.currentTime}}</footer>
  </body>
</html>
