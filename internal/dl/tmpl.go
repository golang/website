// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package dl

// TODO(adg): refactor this to use the tools/godoc/static template.

const templateHTML = `
{{define "root"}}
<!DOCTYPE html>
<html lang="en">
<meta charset="utf-8">
<meta name="description" content="Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#00ADD8">
<title>Downloads - The Go Programming Language</title>
<link href="https://fonts.googleapis.com/css?family=Work+Sans:600|Roboto:400,700" rel="stylesheet">
<link href="https://fonts.googleapis.com/css?family=Product+Sans&text=Supported%20by%20Google&display=swap" rel="stylesheet">
<link type="text/css" rel="stylesheet" href="/lib/godoc/style.css">
<script>window.initFuncs = [];</script>
<style>
  table.codetable {
    margin-left: 20px;
    margin-right: 20px;
    border-collapse: collapse;
  }
  table.codetable tr {
    background-color: #f0f0f0;
  }
  table.codetable tr:nth-child(2n), table.codetable tr.first {
    background-color: white;
  }
  table.codetable td, table.codetable th {
    white-space: nowrap;
    padding: 6px 10px;
  }
  table.codetable tt {
    font-size: xx-small;
  }
  table.codetable tr.highlight td {
    font-weight: bold;
  }
  a.downloadBox {
    display: block;
    color: #222;
    border: 1px solid #375EAB;
    border-radius: 5px;
    background: #E0EBF5;
    width: 280px;
    float: left;
    margin-left: 10px;
    margin-bottom: 10px;
    padding: 10px;
  }
  a.downloadBox:hover {
    text-decoration: none;
  }
  .downloadBox .platform {
    font-size: large;
  }
  .downloadBox .filename {
    color: #007d9c;
    font-weight: bold;
    line-height: 1.5em;
  }
  a.downloadBox:hover .filename {
    text-decoration: underline;
  }
  .downloadBox .size {
    font-size: small;
    font-weight: normal;
  }
  .downloadBox .reqs {
    font-size: small;
    font-style: italic;
  }
  .downloadBox .checksum {
    font-size: 5pt;
  }
</style>
<body class="Site">
<header class="Header js-header">
  <nav class="Header-nav">
    <a href="/"><img class="Header-logo" src="/lib/godoc/images/go-logo-blue.svg" alt="Go"></a>
    <button class="Header-menuButton js-headerMenuButton" aria-label="Main menu" aria-expanded="false">
      <div class="Header-menuButtonInner"></div>
    </button>
    <ul class="Header-menu">
      <li class="Header-menuItem"><a href="/doc/">Documents</a></li>
      <li class="Header-menuItem"><a href="/pkg/">Packages</a></li>
      <li class="Header-menuItem"><a href="/project/">The Project</a></li>
      <li class="Header-menuItem"><a href="/help/">Help</a></li>
      {{if not .GoogleCN}}
        <li class="Header-menuItem"><a href="/blog/">Blog</a></li>
        <li class="Header-menuItem"><a href="https://play.golang.org/">Play</a></li>
      {{end}}
      <li class="Header-menuItem Header-menuItem--search">
        <form class="HeaderSearch" role="search" action="/search">
          <input class="HeaderSearch-input"
                type="search"
                name="q"
                placeholder="Search"
                aria-label="Search"
                autocapitalize="off"
                autocomplete="off"
                autocorrect="off"
                spellcheck="false"
                required>
          <button class="HeaderSearch-submit">
            <!-- magnifying glass: --><svg class="HeaderSearch-icon" width="24" height="24" viewBox="0 0 24 24"><title>Search</title><path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/><path d="M0 0h24v24H0z" fill="none"/></svg>
          </button>
        </form>
      </li>
    </ul>
  </nav>
</header>

<main id="page" class="Site-content">
<div class="container">

<h1>Downloads</h1>

<p>
After downloading a binary release suitable for your system,
please follow the <a href="/doc/install">installation instructions</a>.
</p>

<p>
If you are building from source,
follow the <a href="/doc/install/source">source installation instructions</a>.
</p>

<p>
See the <a href="/doc/devel/release.html">release history</a> for more
information about Go releases.
</p>

<p>
  As of Go 1.13, the go command by default downloads and authenticates
  modules using the Go module mirror and Go checksum database run by Google. See
  <a href="https://proxy.golang.org/privacy">https://proxy.golang.org/privacy</a>
  for privacy information about these services and the
  <a href="/cmd/go/#hdr-Module_downloading_and_verification">go command documentation</a>
  for configuration details including how to disable the use of these servers or use
  different ones.
</p>

{{with .Featured}}
<h3 id="featured">Featured downloads</h3>
{{range .}}
{{template "download" .}}
{{end}}
{{end}}

<div style="clear: both;"></div>

{{with .Stable}}
<h3 id="stable">Stable versions</h3>
{{template "releases" .}}
{{end}}

{{with .Unstable}}
<h3 id="unstable">Unstable version</h3>
{{template "releases" .}}
{{end}}

{{with .Archive}}
<div class="toggle" id="archive">
  <div class="collapsed">
    <h3 class="toggleButton" title="Click to show versions">Archived versions ▹</h3>
  </div>
  <div class="expanded">
    <h3 class="toggleButton" title="Click to hide versions">Archived versions ▾</h3>
    {{template "releases" .}}
  </div>
</div>
{{end}}

</div><!-- .container -->
</main><!-- #page -->
<footer>
  <div class="Footer">
    <img class="Footer-gopher" src="/lib/godoc/images/footer-gopher.jpg" alt="The Go Gopher">
    <ul class="Footer-links">
      <li class="Footer-link"><a href="/doc/copyright.html">Copyright</a></li>
      <li class="Footer-link"><a href="/doc/tos.html">Terms of Service</a></li>
      <li class="Footer-link"><a href="http://www.google.com/intl/en/policies/privacy/">Privacy Policy</a></li>
      <li class="Footer-link"><a href="http://golang.org/issues/new?title=x/website:" target="_blank" rel="noopener">Report a website issue</a></li>
    </ul>
    <a class="Footer-supportedBy" href="https://google.com">Supported by Google</a>
  </div>
</footer>
<script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-11222381-2', 'auto');
  ga('send', 'pageview');

</script>
</body>
<script src="/lib/godoc/jquery.js"></script>
<script src="/lib/godoc/godocs.js"></script>
<script>
$(document).ready(function() {
  $('a.download').click(function(e) {
    // Try using the link text as the file name,
    // unless there's a child element of class 'filename'.
    var filename = $(this).text();
    var child = $(this).find('.filename');
    if (child.length > 0) {
      filename = child.text();
    }

    // This must be kept in sync with the filenameRE in godocs.js.
    var filenameRE = /^go1\.\d+(\.\d+)?([a-z0-9]+)?\.([a-z0-9]+)(-[a-z0-9]+)?(-osx10\.[68])?\.([a-z.]+)$/;
    var m = filenameRE.exec(filename);
    if (!m) {
      // Don't redirect to the download page if it won't recognize this file.
      // (Should not happen.)
      return;
    }

    var dest = "/doc/install";
    if (filename.indexOf(".src.") != -1) {
      dest += "/source";
    }
    dest += "?download=" + filename;

    e.preventDefault();
    e.stopPropagation();
    window.location = dest;
  });
});
</script>
{{end}}

{{define "releases"}}
{{range .}}
<div class="toggle{{if .Visible}}Visible{{end}}" id="{{.Version}}">
	<div class="collapsed">
		<h2 class="toggleButton" title="Click to show downloads for this version">{{.Version}} ▹</h2>
	</div>
	<div class="expanded">
		<h2 class="toggleButton" title="Click to hide downloads for this version">{{.Version}} ▾</h2>
		{{if .Stable}}{{else}}
			<p>This is an <b>unstable</b> version of Go. Use with caution.</p>
			<p>If you already have Go installed, you can install this version by running:</p>
<pre>
go get golang.org/dl/{{.Version}}
</pre>
			<p>Then, use the <code>{{.Version}}</code> command instead of the <code>go</code> command to use {{.Version}}.</p>
		{{end}}
		{{template "files" .}}
	</div>
</div>
{{end}}
{{end}}

{{define "files"}}
<table class="codetable">
<thead>
<tr class="first">
  <th>File name</th>
  <th>Kind</th>
  <th>OS</th>
  <th>Arch</th>
  <th>Size</th>
  {{/* Use the checksum type of the first file for the column heading. */}}
  <th>{{(index .Files 0).ChecksumType}} Checksum</th>
</tr>
</thead>
{{if .SplitPortTable}}
  {{range .Files}}{{if .PrimaryPort}}{{template "file" .}}{{end}}{{end}}

  {{/* TODO(cbro): add a link to an explanatory doc page */}}
  <tr class="first"><th colspan="6" class="first">Other Ports</th></tr>
  {{range .Files}}{{if not .PrimaryPort}}{{template "file" .}}{{end}}{{end}}
{{else}}
  {{range .Files}}{{template "file" .}}{{end}}
{{end}}
</table>
{{end}}

{{define "file"}}
<tr{{if .Highlight}} class="highlight"{{end}}>
  <td class="filename"><a class="download" href="{{.URL}}">{{.Filename}}</a></td>
  <td>{{pretty .Kind}}</td>
  <td>{{.PrettyOS}}</td>
  <td>{{pretty .Arch}}</td>
  <td>{{.PrettySize}}</td>
  <td><tt>{{.PrettyChecksum}}</tt></td>
</tr>
{{end}}

{{define "download"}}
<a class="download downloadBox" href="{{.URL}}">
<div class="platform">{{.Platform}}</div>
{{with .Requirements}}<div class="reqs">{{.}}</div>{{end}}
<div>
  <span class="filename">{{.Filename}}</span>
  {{if .Size}}<span class="size">({{.PrettySize}})</span>{{end}}
</div>
</a>
{{end}}
`
