<h1>Posts</h1>
<ul>
  {{range .posts}}
    <li><a href="{{path_for "post" .Param}}">{{.Title}}</a></li>
  {{else}}
    <li>No posts are available.</li>
  {{end}}
</ul>
