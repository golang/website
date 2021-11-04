---
title: "Security"
---

<section class="Security-hero">
  <div class="Container">
    <div class="Security-heroInner">
      <div class="Security-heroContent">
        {{breadcrumbs .}}
        <h1>Build secure applications with Go</h1>
        <p>
          Gopherum quis custodiet?
        </p>
      </div>
      <div class="Security-heroGopher">
        <img src="/images/gophers/motorcycle.svg" alt="Go Gopher riding a motorcycle">
      </div>
    </div>
  </div>
</section>

<section class="Security-foundations">
  <div class="Container">
    <div class="Security-gridContainer">
      <ul class="Security-cardList">
        {{- range (data "foundations.yaml") }}
          <li class="Security-card">
            {{- template "security-card" .}}
          </li>
        {{- end }}
      </ul>
    </div>
    <hr />
    <div class="Security-comingSoon">
        <div class="Security-comingSoonTitle">
          <h3>Coming Soon</h3>
        </div>
        <div class="Security-comingSoonContent">
          <ul>
            <li>Native support for fuzz testing, maybe OSS-Fuzz integration</li>
            <li>Vulnerabilities database curated by the Go team, with low-noise auditing tools</li>
          </ul>
        </div>
        <div class="Security-comingSoonImage">
          <img src="/images/gophers/motorcycle.svg" alt="Go Gopher riding a motorcycle">
        </div>
    </div>
  </div>
</section>

<section class="Security-getStarted">
  <div class="Container">
    <div class="Security-sectionHeader">
      <h2>Get Started</h2>
    </div>
    <div class="Security-gridContainer">
      <ul class="Security-cardList">
        {{- range (data "getstarted.yaml") }}
          <li class="Security-card">
            {{- template "security-card" .}}
          </li>
        {{- end }}
      </ul>
    </div>
  </div>
</section>

<section class="Security-recentupdates">
  <div class="Container">
    <div class="Security-sectionHeader">
      <h2>Recent Updates</h2>
    </div>
    <div class="Security-gridContainer">
      <ul class="Security-cardList">
        {{- range (data "recentUpdates.yaml") }}
          <li class="Security-card">
            {{- template "security-card" . }}
          </li>
        {{- end }}
      </ul>
    </div>
  </div>
</section>

<section class="Security-secondary-cta">
  <div class="Container">
    <div class="Security-secondary-cta-body">
      <h2>Start building software efficiently and securely with Go</h2>
      <a href="/" rel="noopener"><span>Get Started</span></a>
    </div>
    <div class="Security-secondary-cta-image">
      <img src="/images/gophers/newscaster.svg" alt="Go Gophers surrounding scientific machine">
    </div>
  </div>
</section>

{{define "security-card"}}
<div class="Card">
  <div class="Card-inner">
    {{- if .icon}}
    <div class="Card-icon">
      <img src="{{.icon}}"/>
    </div>
    {{- end}}
    <div class="Card-content">
      <div class="Card-contentTitle">{{.title}}</div>
      <div class="Card-contentBody">
        {{- if .content}}
          {{.content}}
        {{- end}}

        {{- if .contentList}}
          <ul>
            {{- range $index, $element := .contentList}}
              <li>
              {{ $element.title }}
              {{- if $element.url}}
              <a href="{{$element.url}}" target="_blank" rel="noopener" >{{$element.url}}</a>
              {{- end}}
              </li>
            {{- end}}
          </ul>
        {{- end}}
      </div>
      <div class="Card-contentCta">
        <a href="{{.url}}" target="_blank" rel="noopener">
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