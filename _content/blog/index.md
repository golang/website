---
title: The Go Blog
---

<div id="blogindex">

{{range first 10 (newest (pages "/blog/*.md")) -}}
{{if .date}}
<p class="blogtitle">
  <a href="{{.URL}}" aria-describedby="blog-description">{{.title}}</a>, <span class="date">{{.date.Format "2 January 2006"}}</span><br>
  <span class="author">{{with .by}}{{by .}}<br>{{end}}</span>
  {{with .Tags}}<span class="tags">{{range .}}{{.}} {{end}}</span>{{end}}
</p>
<p class="blogsummary">
  {{.summary}}
</p>
{{end}}
{{end}}

<p class="blogtitle">
<a href="/blog/all" aria-label="More articles" aria-describedby="blog-description">More articles...</a>
</p>

<div class="screen-reader-only" id="blog-description" hidden>
    Opens in new window.
</div>
