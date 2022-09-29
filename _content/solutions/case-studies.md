---
title: Case Studies
layout: none
---

{{$solutions := pages "/solutions/*"}}
<section class="Solutions-useCases">
  <div class="Container">
    <ul class="MarketingCardList">
      {{- range $solutions}}
        {{- if eq .series "Case Studies"}}
          <li class="MarketingCard">
            {{- if .link}}
            <a
            href="{{.link}}"
            target="_blank"
            rel="noopener"
            >
            {{- else}}
            <a href="{{.URL}}">
            {{- end}}
              <div class="MarketingCard-section">
                <img
                  height="36"
                  class="DarkMode-img"
                  loading="lazy"
                  alt="{{.company}}"
                  src="/images/logos/{{.logoSrcDark}}"
                />
                <img
                  height="36"
                  class="LightMode-img"
                  loading="lazy"
                  alt="{{.company}}"
                  src="/images/logos/{{.logoSrc}}"
                />
              </div>
              <div class="MarketingCard-section MarketingCard-section__spacer">
                <h2 class="MarketingCard-title">{{or .linkTitle .title}}</h2>
                <p class="MarketingCard-body">
                {{- if .link}}
                  {{.description}}
                {{- else}}
                  {{with .quote}}{{.}}{{end}}
                {{- end}}
                </p>
              </div>
              <div class="MarketingCard-section__bottom">
                <p class="MarketingCard-action">View Case Study</p>
              </div>
            </a>
          </li>
        {{- end}}
      {{- end}}
    <ul>
  </div>
</section>
