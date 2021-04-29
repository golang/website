---
title: Why Go
---

<section class="Solutions-headline">
  <div class="GoCarousel" id="SolutionsHeroCarousel-carousel">
    <div class="GoCarousel-controlsContainer">
      <div class="GoCarousel-wrapper SolutionsHeroCarousel-wrapper">
        <ul class="js-solutionsHeroCarouselSlides SolutionsHeroCarousel-slides">
          {{range where .Pages "Params.series" "Case Studies" | first 3}}
          <li class="SolutionsHeroCarousel-slide">
            <div class="Solutions-headlineImg">
              <img
                src="/images/{{.Params.carouselImgSrc}}"
                alt="{{.LinkTitle}}"
              />
            </div>
            <div class="Solutions-headlineText">
              <p class="Solutions-headlineNotification">RECENTLY UPDATED</p>
              <h2>
                {{.LinkTitle}}
              </h2>
              <p class="Solutions-headlineBody">
                {{with .Params.quote}}{{.}}{{end}}
                <a href="{{.RelPermalink}}"
                  >Learn more
                  <i class="material-icons Solutions-forwardArrowIcon"
                    >arrow_forward</i
                  >
                </a>
              </p>
            </div>
          </li>
          {{end}}
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
      {{$solutions := where .Pages "Params.series" "Case Studies"}}
      {{range sort $solutions "Params.company" "asc"}}
      <li class="Solutions-card">
        {{if isset .Params "link" }}
        <a
          href="{{.Params.link}}"
          target="_blank"
          rel="noopener"
          class="Solutions-useCaseLink"
        >
          <div
            class="Solutions-useCaseLogo Solutions-useCaseLogo--{{.Params.company}}"
          >
            <img
              loading="lazy"
              alt="{{.Params.company}}"
              src="/images/logos/{{.Params.logoSrc}}"
            />
          </div>
          <div class="Solutions-useCaseBody">
            <h3 class="Solutions-useCaseTitle">{{.LinkTitle}}</h3>
            <p class="Solutions-useCaseDescription">
              {{.Description}}
            </p>
          </div>
          <p class="Solutions-useCaseAction">
            View blog post
            <i class="material-icons Solutions-forwardArrowIcon">open_in_new</i>
          </p>
        </a>
        {{else}}
        <a href="{{.RelPermalink}}" class="Solutions-useCaseLink">
          <div class="Solutions-useCaseLogo">
            <img
              loading="lazy"
              alt="{{.Params.company}}"
              src="/images/logos/{{.Params.logoSrc}}"
            />
          </div>
          <div class="Solutions-useCaseBody">
            <h3 class="Solutions-useCaseTitle">{{.LinkTitle}}</h3>
            <p class="Solutions-useCaseDescription">
              {{with .Params.quote}}{{.}}{{end}}
            </p>
          </div>
          <p class="Solutions-useCaseAction">View case study</p>
        </a>
        {{end}}
      </li>
      {{ end }}
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
      {{range where .Pages "Params.series" "Use Cases"}}
      <li class="Solutions-card">
        <a href="{{.RelPermalink}}" class="Solutions-useCaseLink">
          <div class="Solutions-useCaseLogo">
            {{$icon := .Resources.GetMatch "icon"}} {{if $icon}}
            <img
              loading="lazy"
              alt="{{$icon.Params.alt}}"
              src="{{$icon.RelPermalink}}"
            />
            {{end}}
          </div>
          <div class="Solutions-useCaseBody">
            <h3 class="Solutions-useCaseTitle">{{.LinkTitle}}</h3>
            <p class="Solutions-useCaseDescription">
              {{.Description}}
            </p>
          </div>
          <p class="Solutions-useCaseAction">
            Learn More
          </p>
        </a>
      </li>
      {{end}}
    </ul>
    <div class="Solutions-footer">
      <p>
        Interested in sharing your stories?
        <a
          target="_blank"
          rel="noopener"
          href="https://docs.google.com/forms/d/e/1FAIpQLSdRomKkA2zWQF4UTIYWLVYfjKvOHGA32RjnfavVhqY06yrZTQ/viewform"
        >
          Start here.
        </a>
      </p>
    </div>
  </div>
</section>
