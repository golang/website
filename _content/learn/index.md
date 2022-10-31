---
title: "Get Started"
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
              Download
            </a>
          </div>
        </div>
        <p>
          Download packages for
          <a class="js-downloadWin">Windows 64-bit</a>,
          <a class="js-downloadMac">macOS</a>,
          <a class="js-downloadLinux">Linux</a>, and
          <a href="/dl/">more</a>.
        </p>
      </div>
      <div class="Learn-heroGopher">
        <img src="/images/gophers/motorcycle.svg" alt="Go Gopher riding a motorcycle">
      </div>
    </div>
    <div class="LearnGo-gridContainer">
      <ul class="Learn-quickstarts Learn-cardList">
        {{- range first 3 (data "quickstart.yaml")}}
          <li class="Learn-quickstart Learn-card">
            {{- template "learn-card" .}}
          </li>
        {{- end}}
      </ul>
    </div>
  </div>
</section>

<section class="Learn-learningResources">
  <h2>Learning Resources</h2>
</section>

<section id="guided-learning-journeys" class="Learn-guided">
  <div class="Container">
    <div class="Learn-learningResourcesHeader">
      <h3>Guided learning journeys</h3>
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

<section id="online-learning" class="Learn-online">
  <div class="Container">
    <div class="Learn-learningResourcesHeader">
      <h3>Online learning</h3>
    </div>
    <div class="LearnGo-gridContainer">
      <ul class="Learn-cardList">
        {{- range first 4 (data "courses.yaml") }}
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
      <h3>Google Cloud Self-Paced Labs</h3>
    </div>
    <div class="LearnGo-gridContainer">
      <ul class="Learn-cardList">
        {{- range first 4 (data "cloud.yaml")}}
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
    </div>
    <div class="LearnGo-gridContainer">
      <ul class="Learn-cardList">
        {{- range first 4 (data "tutorials.yaml") }}
          <li class="Learn-card">
            {{- template "learn-card" .}}
          </li>
        {{- end}}
      </ul>
    </div>
  </div>
</section>

<section id="featured-books" class="Learn-books">
  <div class="Container">
    <div class="Learn-learningResourcesHeader">
      <h3>Featured books</h3>
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

<section class="Learn-inPersonTraining">
  <div class="Container">
    <div class="Learn-learningResourcesHeader">
      <h3>In-person training</h3>
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
      <p class="Card-contentBody">{{raw .content}}</p>
      <div class="Card-contentCta">
        <a href="{{.url}}" target="_blank">
          <span>{{.cta}}</span>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              fill="#007D9C"
              fill-rule="evenodd"
              d="M5 5v14h14v-7h2v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5c0-1.1.9-2 2-2h7v2H5zm9 0V3h7v7h-2V6.4l-9.8 9.8-1.4-1.4L17.6 5H14z"
              clip-rule="evenodd"
            />
          </svg>
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
        <p class="Book-title">{{.title}}</p>
        <p class="Book-description">{{.description}}</p>
        <div class="Book-cta">
          <span>{{.cta}}</span>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              fill="#007D9C"
              fill-rule="evenodd"
              d="M5 5v14h14v-7h2v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5c0-1.1.9-2 2-2h7v2H5zm9 0V3h7v7h-2V6.4l-9.8 9.8-1.4-1.4L17.6 5H14z"
              clip-rule="evenodd"
            />
          </svg>
        </div>
      </div>
    </div>
  </a>
</div>
{{- end}}
