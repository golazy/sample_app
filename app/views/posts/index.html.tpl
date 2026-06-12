<h1>Posts</h1>
<ul>
  {{range .posts}}
    <li><a href="/posts/{{.Param}}">{{.Title}}</a></li>
  {{else}}
    <li>No posts are available.</li>
  {{end}}
</ul>
