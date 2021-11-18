---
title: Blog Index
---

<div id="blogindex">

{{range newest (pages "/blog/*.md") -}}
{{if .date}}
<p class="blogtitle">
  <a href="{{.URL}}">{{.title}}</a>, <span class="date">{{.date.Format "2 January 2006"}}</span><br>
  <span class="author">{{with .by}}{{by .}}<br>{{end}}</span>
  {{with .Tags}}<span class="tags">{{range .}}{{.}} {{end}}</span>{{end}}
</p>
<p class="blogsummary">
  {{.summary}}
</p>
{{end}}
{{end}}

</div>
