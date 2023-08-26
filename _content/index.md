---
title: The Go Programming Language
summary: Go is an open source programming language that makes it simple to build secure, scalable systems.
---

{{$canShare := not googleCN}}

<section class="Hero bluebg">
  <div class="Hero-gridContainer">
    <div class="Hero-blurb">
      <h1>Build simple, secure, scalable systems with Go</h1>
      <ul class="Hero-blurbList">
        <li>
          <svg width="12" height="10" viewBox="0 0 12 10" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M10.8519 0.52594L3.89189 7.10404L1.14811 4.51081L0 5.59592L3.89189 9.27426L12 1.61105L10.8519 0.52594Z" fill="white" fill-opacity="0.87">
          </svg>
          An open-source programming language supported by Google
        </li>
        <li>
          <svg width="12" height="10" viewBox="0 0 12 10" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M10.8519 0.52594L3.89189 7.10404L1.14811 4.51081L0 5.59592L3.89189 9.27426L12 1.61105L10.8519 0.52594Z" fill="white" fill-opacity="0.87">
          </svg>
          Easy to learn and great for teams
        </li>
        <li>
          <svg width="12" height="10" viewBox="0 0 12 10" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M10.8519 0.52594L3.89189 7.10404L1.14811 4.51081L0 5.59592L3.89189 9.27426L12 1.61105L10.8519 0.52594Z" fill="white" fill-opacity="0.87">
          </svg>
          Built-in concurrency and a robust standard library
        </li>
        <li>
          <svg width="12" height="10" viewBox="0 0 12 10" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M10.8519 0.52594L3.89189 7.10404L1.14811 4.51081L0 5.59592L3.89189 9.27426L12 1.61105L10.8519 0.52594Z" fill="white" fill-opacity="0.87">
          </svg>
          Large ecosystem of partners, communities, and tools
        </li>
      </ul>
    </div>
    <div class="Hero-actions">
      <div data-version="" class="js-latestGoVersion">
        <div role="group" class="BtnGroup">
          <div class="BtnGroup-sub">
            <button class="mainBtn go-Button go-Button--accented js-downloadBtn" aria-label="Download" aria-describedby="download-description">
              <span>Download Go <small class="js-goVersion"></small></span>
              <small class="go-Chip go-Chip--accented js-osAndArch"></small>
            </button>
            <button class="go-Button go-Button--accented js-selectBinary" aria-haspopup="menu" aria-label="select binary type to download">
              <img class="go-Icon" height="24" width="24" src="/images/icons/arrow_drop_down_gm_grey_24dp.svg" alt="arrow-dropdown">
            </button>
          </div>
          <div class="binaryMatrix js-binaryMatrix">
            <table class="binaryMatrix-Table">
              <thead>
                <tr>
                  <th scope="col"><p><small>OS</small></p></th>
                  <th scope="col"><p><small>(Arch)</small></p></th>
                  <th scope="col"></th>
                </tr>
              </thead>
              <tbody class="binaryDistribution js-binaryDistribution">
                <tr>
                  <td scope="row"><p>MacOS</p></td>
                  <td>
                    <p><span class="go-Chip go-Chip--accented">ARM64</span></p>
                    <p><span class="go-Chip go-Chip--accented">AMD64</span></p>
                  </td>
                  <td>
                    <p><a href="/doc/install?dc=darwin-arm64"><img class="go-Icon" height="24" width="24" src="/images/icons/insert_drive_file_gm_grey_24dp.svg" alt=""></a></p>
                    <p><a href="/doc/install?dc=darwin-amd64"><img class="go-Icon" height="24" width="24" src="/images/icons/insert_drive_file_gm_grey_24dp.svg" alt=""></a></p>
                  </td>
                </tr>
              </tbody>
              <tbody class="binaryDistribution js-binaryDistribution">
                <tr><td scope="row"><p>Windows</p></td>
                  <td><p><span class="go-Chip go-Chip--accented">AMD64</span></p></td>
                  <td><p><a href="/doc/install?dc=windows-amd64"><img class="go-Icon" height="24" width="24" src="/images/icons/insert_drive_file_gm_grey_24dp.svg" alt=""></a></p>
                  </td>
                </tr>
              </tbody>
              <tbody class="binaryDistribution js-binaryDistribution">
                <tr><td scope="row"><p>Linux</p></td>
                  <td><p><span class="go-Chip go-Chip--accented">AMD64</span></p></td>
                  <td><p>
                    <a href="/doc/install?dc=linux-amd64" data-os="linux" data-arch="amd64"><img class="go-Icon" height="24" width="24" src="/images/icons/insert_drive_file_gm_grey_24dp.svg" alt=""></a></p>
                  </td>
                </tr>
              </tbody>
              <tbody class="otherDistribution">
                <tr><td colspan="3"><a href="/dl"><small>See other <code>go</code> binary releases</small></a></td></tr>
              </tbody>
            </table>
          </div>
        </div>
        <a class="Secondary secondaryBtn" href="/play/" aria-label="Try a Tour of Go" aria-describedby="tryGoTour-description" role="button">
          <span>Not ready to download Go?</span>
          <small>Try a tour of Go</small>
        </a>
        <div class="screen-reader-only" id="download-description" hidden>
          Downloads Go and opens a new window with instructions to install Go.
        </div>
        <div class="screen-reader-only" id="tryGoTour-description" hidden>
          Opens a new window with A Tour of Go.
        </div>
      </div>
      <div class="Hero-footnote">
        <p>
          The <code>go</code> command by default downloads and authenticates
          modules using the Go module mirror and Go checksum database run by
          Google. <a href="/dl" aria-describedby="newwindow-description">Learn more.</a>
        </p>
      </div>
    </div>
    <div class="screen-reader-only" id="newwindow-description" hidden>
      Opens in new window.
    </div>
    <div class="Hero-gopher">
      <img class="Hero-gopherLadder" src="/images/gophers/ladder.svg" alt="Go Gopher climbing a ladder.">
    </div>
  </div>
</section>
<section class="WhoUses">
  <div class="WhoUses-gridContainer">
    <div class="WhoUses-header">
      <h2 class="WhoUses-headerH2">Companies using Go</h2>
      <p class="WhoUses-subheader">Organizations in every industry use Go to power their software and services
        <a href="/solutions/" class="WhoUsesCaseStudyList-seeAll" aria-describedby="newwindow-description">
        View all stories
       </a>
     </p>
    </div>
  <div class="WhoUsesCaseStudyList">
    <ul class="WhoUsesCaseStudyList-gridContainer">
    {{- range newest (pages "/solutions/*")}}{{if eq .series "Case Studies"}}
      {{- if .link }}
        {{- if .inLandingPageGrid }}
          <li class="WhoUsesCaseStudyList-caseStudy">
            <a href="{{.link}}" aria-label="View CaseStudy of {{.company}}, (opens in new window)" target="_blank" rel="noopener"
              class="WhoUsesCaseStudyList-caseStudyLink">
              <img
                loading="lazy"
                height="48"
                width="30%"
                src="/images/logos/{{.logoSrc}}"
                class="WhoUsesCaseStudyList-logo"
                alt="">
            </a>
          </li>
        {{- end}}
      {{- else}}
        <li class="WhoUsesCaseStudyList-caseStudy">
          <a href="{{.URL}}" aria-label="View CaseStudy of {{.company}}, (opens in new window)" class="WhoUsesCaseStudyList-caseStudyLink">
            <img
              loading="lazy"
              height="48"
              width="30%"
              src="/images/logos/{{.logoSrc}}"
              class="WhoUsesCaseStudyList-logo"
              alt="">
            <p>View case study</p>
          </a>
        </li>
      {{- end}}
    {{- end}}
    {{- end}}
    </ul>
  </div>
</section>
<section class="TestimonialsGo">
  <div class="GoCarousel">
    <div class="GoCarousel-controlsContainer">
      <div class="GoCarousel-wrapper">
        <ul class="js-testimonialsGoQuotes TestimonialsGo-quotes">
          {{- range $index, $element := data "/testimonials.yaml"}}
            <li class="TestimonialsGo-quoteGroup GoCarousel-slide" id="quote_slide{{$index}}">
              <div class="TestimonialsGo-quoteSingleItem">
                <div class="TestimonialsGo-quoteSection">
                  <p class="TestimonialsGo-quote">{{raw .quote}}</p>
                  <div class="TestimonialsGo-author">— {{.name}},
                    <span class="NoWrapSpan">{{.title}}</span>
                    <span class="NoWrapSpan"> at {{.company}}</span>
                  </div>
                </div>
              </div>
            </li>
          {{- end}}
        </ul>
      </div>
    <button class="js-testimonialsPrev GoCarousel-controlPrev" hidden>
      <i class="GoCarousel-icon material-icons">navigate_before</i>
    </button>
    <button class="js-testimonialsNext GoCarousel-controlNext">
      <i class="GoCarousel-icon material-icons">navigate_next</i>
    </button>
  </div>
  </div>
</section>
<section class="Playground">
  <div class="Playground-gridContainer">
    <div class="Playground-headerContainer">
      <h2 class="HomeSection-header">Try Go</h2>
    </div>
    <div class="Playground-inputContainer">
      <div class="Playground-preContainer">
        Press Esc to move out of the editor.
      </div>
      <textarea class="Playground-input js-playgroundCodeEl" spellcheck="false" aria-label="Try Go" aria-describedby="editor-description" id="code">
// You can edit this code!
// Click here and start typing.
package main
import "fmt"
func main() {
  fmt.Println("Hello, 世界")
}</textarea>
    </div>
    <div class="screen-reader-only" id="editor-description" hidden>
      Press Esc to move out of the editor.
    </div>
    <div class="Playground-outputContainer js-playgroundOutputEl">
      <pre class="Playground-output"><noscript>Hello, 世界</noscript></pre>
    </div>
    <div class="Playground-controls">
      <select class="Playground-selectExample js-playgroundToysEl" aria-label="Code examples">
      <option value="hello.go">Hello, World!</option>
      <option value="life.go">Conway's Game of Life</option>
      <option value="fib.go">Fibonacci Closure</option>
      <option value="peano.go">Peano Integers</option>
      <option value="pi.go">Concurrent pi</option>
      <option value="sieve.go">Concurrent Prime Sieve</option>
      <option value="solitaire.go">Peg Solitaire Solver</option>
      <option value="tree.go">Tree Comparison</option>
      </select>
      <div class="Playground-buttons">
      <button class="Button Button--primary js-playgroundRunEl Playground-runButton" title="Run this code [shift-enter]">Run</button>
      <div class="Playground-secondaryButtons">
        {{- if $canShare}}
        <button class="Button js-playgroundShareEl Playground-button" title="Share in Go Playground">Share</button>
        {{- end}}
        <a class="Button tour Playground-button" href="/tour/" title="Tour Go from your browser">Tour</a>
      </div>
      </div>
    </div>
  </div>
</section>
<section class="WhyGo">
  <div class="WhyGo-gridContainer">
    <div class="WhyGo-header">
      <h2 class="WhyGo-headerH2">What’s possible with Go</h2>
      <p class="WhyGo-subheader">
        Use Go for a variety of software development purposes
      </p>
    </div>
    <ul class="WhyGo-reasons">
      {{- range first 4 (data "/resources.yaml")}}
        <li class="WhyGo-reason">
          <div class="WhyGo-reasonDetails">
            <div class="WhyGo-reasonIcon" role="presentation">
              <img class="DarkMode-img" src="{{.iconDark}}" alt="{{.iconName}}">
              <img class="LightMode-img" src="{{.icon}}" alt="{{.iconName}}">
            </div>
            <div class="WhyGo-reasonText">
              <h3 class="WhyGo-reasonTitle">{{.title}}</h3>
              <p>
                {{.description}}
              </p>
            </div>
          </div>
          <div class="WhyGo-reasonFooter">
            <div class="WhyGo-reasonPackages">
              <div class="WhyGo-reasonPackagesHeader">
                <img src="/images/icons/package.svg" alt="Packages.">
                Popular Packages:
              </div>
              <ul class="WhyGo-reasonPackagesList">
                {{- range .packages }}
                  <li class="WhyGo-reasonPackage">
                    <a class="WhyGo-reasonLink" href="{{.url}}" target="_blank" rel="noopener">
                      {{.title}}
                    </a>
                  </li>
                  {{- end}}
              </ul>
            </div>
            <div class="WhyGo-reasonLearnMoreLink">
              <a href="{{.link}}" aria-describedby="newwindow-description">Learn More 
              <i class="material-icons WhyGo-forwardArrowIcon" aria-hidden="true">arrow_forward</i></a>
            </div>
          </div>
        </li>
      {{- end}}
      {{- if gt (len (data "resources.yaml")) 3}}
        <li class="WhyGo-reason">
          <div class="WhyGo-reasonShowMore">
            <div class="WhyGo-reasonShowMoreImgWrapper">
              <img
                class="WhyGo-reasonShowMoreImg"
                loading="lazy"
                height="148"
                width="229"
                src="/images/gophers/biplane.svg"
                alt="Go Gopher is skateboarding.">
            </div>
            <div class="WhyGo-reasonShowMoreLink">
              <a href="/solutions/use-cases" aria-describedby="newwindow-description">More use cases 
              <i class="material-icons
              WhyGo-forwardArrowIcon" aria-hidden="true">arrow_forward</i></a>
            </div>
          </div>
        </li>
      {{- end}}
    </ul>
  </div>
</section>
<section class="GettingStartedGo">
  <div class="GettingStartedGo-gridContainer">
    <div class="GettingStartedGo-header">
      <h2 class="GettingStartedGo-headerH2">Get started with Go</h2>
      <p class="GettingStartedGo-headerDesc">
        Explore a wealth of learning resources, including guided journeys, courses, books, and more.
      </p>
      <div class="GettingStartedGo-ctas">
        <a class="GettingStartedGo-primaryCta" href="/learn/"aria-describedby="newwindow-description">Get Started</a>
        <a href="/doc/install/" aria-describedby="newwindow-description">Download Go</a>
      </div>
    </div>
    <div class="GettingStartedGo-resourcesSection">
      <ul class="GettingStartedGo-resourcesList">
        <li class="GettingStartedGo-resourcesHeader">
          Resources to start on your own
        </li>
        <li class="GettingStartedGo-resourceItem">
          <a href="/learn#guided-learning-journeys" class="GettingStartedGo-resourceItemTitle" aria-describedby="newwindow-description">
            Guided learning journeys
          </a>
          <div class="GettingStartedGo-resourceItemDescription">
            Step-by-step tutorials to get your feet wet
          </div>
        </li>
        <li class="GettingStartedGo-resourceItem">
          <a href="/learn#online-learning" class="GettingStartedGo-resourceItemTitle" aria-describedby="newwindow-description">
            Online learning
          </a>
          <div class="GettingStartedGo-resourceItemDescription">
            Browse resources and learn at your own pace
          </div>
        </li>
        <li class="GettingStartedGo-resourceItem">
          <a href="/learn#featured-books" class="GettingStartedGo-resourceItemTitle" aria-describedby="newwindow-description">
            Featured books
          </a>
          <div class="GettingStartedGo-resourceItemDescription">
            Read through structured chapters and theories
          </div>
        </li>
        <li class="GettingStartedGo-resourceItem">
          <a href="/learn#self-paced-labs" class="GettingStartedGo-resourceItemTitle" aria-describedby="newwindow-description">
            Cloud Self-paced labs
          </a>
          <div class="GettingStartedGo-resourceItemDescription">
            Jump in to deploying Go apps on GCP
          </div>
        </li>
      </ul>
      <ul class="GettingStartedGo-resourcesList">
        <li class="GettingStartedGo-resourcesHeader">
          In-Person Trainings
        </li>
        {{- range first 4 (data "/learn/training.yaml")}}
          <li class="GettingStartedGo-resourceItem">
            <a href="{{.url}}" class="GettingStartedGo-resourceItemTitle" aria-describedby="newwindow-description">
              {{.title}}
            </a>
            <div class="GettingStartedGo-resourceItemDescription">
              {{.blurb}}
            </div>
          </li>
        {{- end}}
      </ul>
    </div>
  </div>
</section>
<script src="/js/index.js" defer></script>
