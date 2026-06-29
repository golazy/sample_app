<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{if .title}}{{.title}} - {{end}}GoLazy Sample App</title>
    {{stylesheet "/styles.css"}}
    {{importmap "/assets/importmap.json"}}
    <script type="module">import "app.js"</script>
  </head>
  <body class="min-h-screen bg-[#f7f4eb] text-[#080808] antialiased">
    {{.content}}
  </body>
</html>
