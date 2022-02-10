---
title: Why Go
layout: none
---

{{$solutions := pages "/solutions/*"}}
<section class="Solutions-headline">
  <div class="GoCarousel" id="SolutionsHeroCarousel-carousel">
    <div class="GoCarousel-controlsContainer">
      <div class="GoCarousel-wrapper SolutionsHeroCarousel-wrapper">
        <ul class="js-solutionsHeroCarouselSlides SolutionsHeroCarousel-slides">
          {{- $n := 0}}
          {{- range newest $solutions}}
            {{- if eq .series "Case Studies"}}
              {{- $n = add $n 1}}
              {{- if le $n 3}}
              <li class="SolutionsHeroCarousel-slide">
                <div class="Solutions-headlineImg">
                  <img
                    src="/images/{{.carouselImgSrc}}"
                    alt="{{(or .linkTitle .title)}}"
                  />
                </div>
                <div class="Solutions-headlineText">
                  <p class="Solutions-headlineNotification">RECENTLY UPDATED</p>
                  <h2>
                    {{(or .linkTitle .title)}}
                  </h2>
                  <p class="Solutions-headlineBody">
                    {{with .quote}}{{.}}{{end}}
                    <a href="{{.URL}}"
                      >Learn more
                      <i class="material-icons Solutions-forwardArrowIcon"
                        >arrow_forward</i
                      >
                    </a>
                  </p>
                </div>
              </li>
              {{- end}}
            {{- end}}
          {{- end}}
        </ul>
      </div>
      <button
        class="js-solutionsHeroCarouselPrev GoCarousel-controlPrev GoCarousel-controlPrev-solutionsHero"
        hidden
      >
        <i class="GoCarousel-icon material-icons">navigate_before</i>
      </button>
      <button
        class="js-solutionsHeroCarouselNext GoCarousel-controlNext GoCarousel-controlNext-solutionsHero"
      >
        <i class="GoCarousel-icon material-icons">navigate_next</i>
      </button>
    </div>
  </div>
</section>
<section class="Solutions-useCases">
  <div class="Container">
    <div class="SolutionsTabs-tabList js-solutionsTabs" role="tablist">
      <button
        role="tab"
        aria-selected="true"
        class="SolutionsTabs-tab"
        id="btn-companies"
        aria-controls="tab-companies"
      >
        Case studies
      </button>
      <button
        role="tab"
        aria-selected="false"
        class="SolutionsTabs-tab"
        id="btn-tech"
        aria-controls="tab-tech"
      >
        Use cases
      </button>
      <hr />
    </div>
    <ul
      class="js-solutionsList Solutions-cardList"
      aria-expanded="true"
      aria-labelledby="btn-companies"
      id="tab-companies"
      role="tabpanel"
      tabindex="0"
    >
      {{- range $solutions}}
      {{- if eq .series "Case Studies"}}
      <li class="Solutions-card">
        {{- if .link}}
        <a
          href="{{.link}}"
          target="_blank"
          rel="noopener"
          class="Solutions-useCaseLink"
        >
          <div
            class="Solutions-useCaseLogo Solutions-useCaseLogo--{{.company}}"
          >
            <img
              class="DarkMode-img"
              loading="lazy"
              alt="{{.company}}"
              src="/images/logos/{{.logoSrcDark}}"
            />
            <img
              class="LightMode-img"
              loading="lazy"
              alt="{{.company}}"
              src="/images/logos/{{.logoSrc}}"
            />
          </div>
          <div class="Solutions-useCaseBody">
            <h3 class="Solutions-useCaseTitle">{{or .linkTitle .title}}</h3>
            <p class="Solutions-useCaseDescription">
              {{.description}}
            </p>
          </div>
          <p class="Solutions-useCaseAction">
            View blog post
            <i class="material-icons Solutions-forwardArrowIcon">open_in_new</i>
          </p>
        </a>
        {{- else}}
        <a href="{{.URL}}" class="Solutions-useCaseLink">
          <div class="Solutions-useCaseLogo">
            <img
              class="DarkMode-img"
              loading="lazy"
              alt="{{.company}}"
              src="/images/logos/{{.logoSrcDark}}"
            />
            <img
              class="LightMode-img"
              loading="lazy"
              alt="{{.company}}"
              src="/images/logos/{{.logoSrc}}"
            />
          </div>
          <div class="Solutions-useCaseBody">
            <h3 class="Solutions-useCaseTitle">{{or .linkTitle .title}}</h3>
            <p class="Solutions-useCaseDescription">
              {{with .quote}}{{.}}{{end}}
            </p>
          </div>
          <p class="Solutions-useCaseAction">View case study</p>
        </a>
        {{- end}}
      </li>
      {{- end}}
      {{- end}}
    </ul>
    <ul
      class="js-solutionsList Solutions-cardList"
      aria-expanded="false"
      aria-labelledby="btn-tech"
      id="tab-tech"
      role="tabpanel"
      tabindex="0"
      hidden
    >
      {{- range newest $solutions}}{{if eq .series "Use Cases"}}
      <li class="Solutions-card">
        <a href="{{.URL}}" class="Solutions-useCaseLink">
          <div class="Solutions-useCaseLogo">
            {{- $icon := .icon}}
            {{- $iconDark := .iconDark}}
            {{- if $icon}}
            <img
              class="LightMode-img"
              loading="lazy"
              alt="{{$icon.alt}}"
              src="{{path.Dir .URL}}/{{$icon.file}}"
            />
            {{- end}}
            {{- if $iconDark}}
            <img
              class="DarkMode-img"
              loading="lazy"
              alt="{{$iconDark.alt}}"
              src="{{path.Dir .URL}}/{{$iconDark.file}}"
            />
            {{- end}}
          </div>
          <div class="Solutions-useCaseBody">
            <h3 class="Solutions-useCaseTitle">{{or .linkTitle .title}}</h3>
            <p class="Solutions-useCaseDescription">
              {{.description}}
            </p>
          </div>
          <p class="Solutions-useCaseAction">
            Learn More
          </p>
        </a>
      </li>
      {{- end}}
      {{- end}}
    </ul>
  </div>
</section>
