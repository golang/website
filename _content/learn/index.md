---
title: "Get Started"
breadcrumbTitle: "Learn"
---

<section class="Learn-hero">
  <div class="Container">
    <div class="Learn-heroInner">
      <div class="Learn-heroContent">
        {{breadcrumbs .}}
        <h1>Install the latest version of Go</h1>
        <p>
          Install the latest version of Go. For instructions to download and install
          the Go compilers, tools, and libraries,
          <a href="/doc/install" target="_blank" rel="noopener">
            view the install documentation.
          </a>
        </p>
        <div class="Learn-heroAction">
          <div
            data-version=""
            class="js-latestGoVersion"
          >
            <a
              class="js-downloadBtn"
              href="/dl"
              target="_blank"
              rel="noopener"
            >
              <span class="GoVersionSpan">Download</span>
            </a>
          </div>
        </div>
      </div>
      <div class="Learn-heroGopher">
        <img src="/images/gophers/motorcycle.svg" alt="Go Gopher riding a motorcycle">
      </div>
    </div>
  </div>
</section>

<div class="Learn-columns">
  <aside class="Learn-sidebar">
    <nav class="LearnNav">
      <a class="active" href="#selected-tutorials">
        <svg width="5" height="5" viewBox="0 0 5 5" fill="none" xmlns="http://www.w3.org/2000/svg"><circle cx="2.5" cy="2.5" r="2.5" fill="#007F9F"/></svg>
        <span>Selected tutorials</span>
      </a>
      <a href="#guided-learning-journeys">
      <svg width="5" height="5" viewBox="0 0 5 5" fill="none" xmlns="http://www.w3.org/2000/svg"><circle cx="2.5" cy="2.5" r="2.5" fill="#007F9F"/></svg>
      <span>Guided journeys</span>
      </a>
      <a href="#self-paced-labs">
      <svg width="5" height="5" viewBox="0 0 5 5" fill="none" xmlns="http://www.w3.org/2000/svg"><circle cx="2.5" cy="2.5" r="2.5" fill="#007F9F"/></svg>
      <span>Qwiklabs</span>
      </a>
      <a href="#tutorials">
      <svg width="5" height="5" viewBox="0 0 5 5" fill="none" xmlns="http://www.w3.org/2000/svg"><circle cx="2.5" cy="2.5" r="2.5" fill="#007F9F"/></svg>
      <span>Tutorials</span>
      </a>
      <a href="#training">
      <svg width="5" height="5" viewBox="0 0 5 5" fill="none" xmlns="http://www.w3.org/2000/svg"><circle cx="2.5" cy="2.5" r="2.5" fill="#007F9F"/></svg>
      <span>Training</span>
      </a>
      <a href="#featured-books">
      <svg width="5" height="5" viewBox="0 0 5 5" fill="none" xmlns="http://www.w3.org/2000/svg"><circle cx="2.5" cy="2.5" r="2.5" fill="#007F9F"/></svg>
      <span>Books</span>
      </a>
    </nav>
  </aside>
  <div class="Learn-body">
  <section id="selected-tutorials" class="Learn-tutorials">
    <div class="Container">
      <div class="Learn-learningResourcesHeader">
          <h3>Selected tutorials</h3>
          <p>New to Go and don't know where to start?</p>
      </div>
      <div class="LearnGo-gridContainer">
        <ul class="Learn-cardList">
        {{- range first 3 (data "quickstart.yaml")}}
            <li class="Learn-card">
            {{- template "learn-card" . }}
            </li>
        {{- end}}
        </ul>
      </div>
    </div>
  </section>

  <section id="guided-learning-journeys" class="Learn-guided">
    <div class="Container">
      <div class="Learn-learningResourcesHeader">
        <h3>Guided learning journeys</h3>
        <p>Got the basics and want to learn more?</p>
      </div>
      <div class="LearnGo-gridContainer">
        <ul class="Learn-cardList">
          {{- range first 4 (data "guided.yaml")}}
            <li class="Learn-card">
              {{- template "learn-card" .}}
            </li>
          {{- end}}
        </ul>
      </div>
    </div>
  </section>

  <section id="self-paced-labs" class="Learn-selfPaced">
    <div class="Container">
      <div class="Learn-learningResourcesHeader">
        <h3>Qwiklabs</h3>
        <p>Guided tours of Go programs</p>
      </div>
      <div class="LearnGo-gridContainer">
        <ul class="Learn-cardList">
          {{- range first 3 (data "cloud.yaml")}}
          <li class="Learn-card">
            {{- template "learn-self-paced-card" .}}
          </li>
          </li>
          {{- end}}
        </ul>
      </div>
    </div>
  </section>

  <section id="tutorials" class="Learn-tutorials">
    <div class="Container">
      <div class="Learn-learningResourcesHeader">
        <h3>Tutorials</h3>
        <p></p>
      </div>
      <div class="LearnGo-gridContainer">
        <ul class="Learn-cardList">
          {{- range first 3 (data "tutorials.yaml") }}
            <li class="Learn-card">
              {{- template "learn-card" .}}
            </li>
          {{- end}}
        </ul>
      </div>
    </div>
  </section>

  <section id="training" class="Learn-inPersonTraining">
    <div class="Container">
      <div class="Learn-learningResourcesHeader">
        <h3>Training</h3>
        <p>Guided tours of Go programs</p>
      </div>
      <div class="LearnGo-gridContainer">
        <ul class="Learn-inPersonList">
          {{- range first 4 (data "training.yaml")}}
          <li class="Learn-inPerson">
            <p class="Learn-inPersonTitle">
              <a href="{{.url}}">{{.title}} </a>
            </p>
            <p class="Learn-inPersonBlurb">{{.blurb}}</p>
          </li>
          {{- end}}
        </ul>
      </div>
    </div>
  </section>

  <section id="featured-books" class="Learn-books">
    <div class="Container">
      <div class="Learn-learningResourcesHeader">
        <h3>Books</h3>
        <p></p>
      </div>
      <div class="LearnGo-gridContainer">
        <ul class="Learn-cardList Learn-bookList">
          {{- range first 5 (data "books.yaml")}}
            <li class="Learn-card Learn-book">
              {{template "learn-book" .}}
            </li>
          {{- end}}
        </ul>
      </div>
    </div>
  </section>
  </div>
</div>

<script async src="/js/jumplinks.js"></script>

{{define "learn-card"}}
<div class="Card">
  <div class="Card-inner">
    {{- if .thumbnailDark}}
    <div
      class="Card-thumbnail DarkMode-img"
      style="background-image: url('{{.thumbnailDark}}')"
    ></div>
    {{- else if .thumbnail}}
    <div
      class="Card-thumbnail DarkMode-img"
      style="background-image: url('{{.thumbnail}}')"
    ></div>
    {{- end}}
    {{- if .thumbnail}}
    <div
      class="Card-thumbnail LightMode-img"
      style="background-image: url('{{.thumbnail}}')"
    ></div>
    {{- end}}
    <div class="Card-content">
      <div class="Card-contentTitle">{{.title}}</div>
      <p class="Card-contentBody Card-lineClamp">{{raw .content}}</p>
      <div class="Card-contentCta">
        <a href="{{.url}}" target="_blank">
          <span>{{.cta}}</span>
        </a>
      </div>
    </div>
  </div>
</div>
{{- end}}

{{define "learn-self-paced-card"}}
<div class="Card">
  <a href="{{.url}}" target="_blank" rel="noopener">
    <div class="Card-inner">
      {{- if .thumbnail}}
      <div
        class="Card-thumbnail"
        style="background-image: url('{{.thumbnail}}')"
      ></div>
      {{- end}}
      <div class="Card-content">
        <div class="Card-contentTitle">{{.title}}</div>
        <div class="Card-selfPacedFooter">
          <div class="Card-selfPacedCredits">
            <span>{{ .length }}</span> •
            <span>{{.credits}} Credits</span>
          </div>
          <div class="Card-selfPacedRating">
            <div class="Card-starRating" style="width: {{ .rating }}rem"></div>
          </div>
        </div>
      </div>
    </div>
  </a>
</div>
{{- end}}

{{define "learn-book"}}
<div class="Book">
  <a href="{{.url}}" target="_blank" rel="noopener">
    <div class="Book-inner">
      {{- if .thumbnail}}
      <div class="Book-thumbnail">
        <img alt="{{.title}} thumbnail." src="{{.thumbnail}}" />
      </div>
      {{- end}}
      <div class="Book-content">
        <p class="Book-eyebrow">{{.eyebrow}}</p>
        <p class="Book-title">{{.title}}</p>
        <p class="Book-description">{{.description}}</p>
        <div class="Book-cta">
          <span>view book</span>
        </div>
      </div>
    </div>
  </a>
</div>
{{- end}}
