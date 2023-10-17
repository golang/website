---
title: "Command-line Interfaces (CLIs)"
linkTitle: "Command-line Interfaces (CLIs)"
description: "With popular open source packages and a robust standard library, use Go to create fast and elegant CLIs."
date: 2019-10-04T15:26:31-04:00
series: Use Cases
icon:
  file: clis-green.svg
  alt: CLI icon
iconDark:
  file: clis-white.svg
  alt: CLI icon
---

## Overview {#overview .sectionHeading}

### CLI developers prefer Go for portability, performance, and ease of creation

Command line interfaces (CLIs), unlike graphical user interfaces (GUIs), are text-only. Cloud and infrastructure applications are primarily CLI-based due to their easy automation and remote capabilities.

## Key benefits {#key-benefits .sectionHeading}

### Leverage fast compile times to build programs that start quickly and run on any system

Developers of CLIs find Go to be ideal for designing their applications. Go compiles very quickly into a single binary, works across platforms with a consistent style, and brings a strong development community. From a single Windows or Mac laptop, developers can build a Go program for every one of the dozens of architectures and operating systems Go supports in a matter of seconds, no complicated build farms are needed. No other compiled language can be built as portably or quickly. Go applications are built into a single self contained binary making installing Go applications trivial.

Specifically, **programs written in Go run on any system without requiring any existing libraries, runtimes, or dependencies**. And **programs written in Go have an immediate startup time**—similar to C or C++ but unobtainable with other programming languages.

## Use Case {#use-case .sectionHeading}

### Use Go for building elegant CLIs

{{backgroundquote `
  author: Steve Domino
  title: senior engineer and architect at Strala
  link: https://medium.com/@skdomino/writing-better-clis-one-snake-at-a-time-d22e50e60056
  quote: |
    I was tasked with building our CLI tool and found two really great projects, Cobra and Viper, which make building CLI’s easy. Individually they are very powerful, very flexible and very good at what they do. But together they will help you show your next CLI who is boss!
`}}

{{backgroundquote `
  author: Francesc Campoy
  title: VP of product at DGraph Labs and producer of Just For Func videos
  link: https://www.youtube.com/watch?v=WvWPGVKLvR4
  quote: |
    Cobra is a great product to write small tools or even large ones. It’s more of a framework than a library, because when you call the binary that would create a skeleton, then you would be adding code in between.”
`}}

When developing CLIs in Go, two tools are widely used: Cobra & Viper.

{{pkg "github.com/spf13/cobra" "Cobra"}} is both a library for creating powerful modern CLI applications and a program to generate applications and CLI applications in Go. Cobra powers most of the popular Go applications including CoreOS, Delve, Docker, Dropbox, Git Lfs, Hugo, Kubernetes, and [many more](https://pkg.go.dev/github.com/spf13/cobra?tab=importedby). With integrated command help, autocomplete and documentation “[it] makes documenting each command really simple,” says [Alex Ellis](https://blog.alexellis.io/5-keys-to-a-killer-go-cli/), founder of OpenFaaS.


{{pkg "github.com/spf13/viper" "Viper"}} is a complete configuration solution for Go applications, designed to work within an app to handle configuration needs and formats. Cobra and Viper are designed to work together.

Viper [supports nested structures](https://scene-si.org/2017/04/20/managing-configuration-with-viper/) in the configuration, allowing CLI developers to manage the configuration for multiple parts of a large application. Viper also provides all of the tooling need to easily build twelve factor apps.

"If you don’t want to pollute your command line, or if you’re working with sensitive data which you don’t want to show up in the history, it’s a good idea to work with environment variables. To do this, you can use Viper," [suggests Geudens](https://ordina-jworks.github.io/development/2018/10/20/make-your-own-cli-with-golang-and-cobra.html).

{{projects `
  - company: Comcast
    url: https://xfinity.com/
    logoSrc: comcast.svg
    logoSrcDark: comcast.svg
    desc: Comcast uses Go for a CLI client used to publish and subscribe to its high-traffic sites. The company also supports an open source client library which is written in Go - designed for working with Apache Pulsar.
    ctas:
      - text: Client library for Apache Pulsar
        url: https://github.com/Comcast/pulsar-client-go
      - text: Pulsar CLI Client
        url: https://github.com/Comcast/pulsar-client-go/blob/master/cli/main.go
  - company: GitHub
    url: https://github.com/
    logoSrc: github.svg
    logoSrcDark: github.svg
    desc: GitHub uses Go for a command-line tool that makes it easier to work with GitHub, wrapping git in order to extend it with extra features and commands.
    ctas:
      - text: GitHub command-line tool
        url: https://github.com/github/hub
  - company: Hugo
    url: http://gohugo.io/
    logoSrc: hugo.svg
    logoSrcDark: hugo.svg
    desc: Hugo is one of the most popular Go CLI applications powering thousands of sites, including this one. One reason for its popularity is its ease of install thanks to Go. Hugo author Bjørn Erik Pedersen writes “The single binary takes most of the pain out of installation and upgrades.”
    ctas:
      - text: Hugo Website
        url: https://gohugo.io/
  - company: Kubernetes
    url: https://kubernetes.com/
    logoSrc: kubernetes.svg
    logoSrcDark: kubernetes.svg
    desc: Kubernetes is one of the most popular Go CLI applications. Kubernetes Creator, Joe Beda, said that for writing Kubernetes, “Go was the only logical choice”. Calling Go “the sweet spot” between low level languages like C++ and high level languages like Python.
    ctas:
      - text: Kubernetes + Go
        url: https://blog.gopheracademy.com/birthday-bash-2014/kubernetes-go-crazy-delicious/
  - company: MongoDB
    url: https://mongodb.com/
    logoSrc: mongodb.svg
    logoSrcDark: mongodb.svg
    desc: MongoDB chose to implement their Backup CLI Tool in Go citing Go’s “C-like syntax, strong standard library, the resolution of concurrency problems via goroutines, and painless multi-platform distribution” as reasons.
    ctas:
      - text: MongoDB Backup Service
        url: https://www.mongodb.com/blog/post/go-agent-go
  - company: Netflix
    url: http://netflix.com/
    logoSrc: netflix.svg
    logoSrcDark: netflix.svg
    desc: Netflix uses Go to build the CLI application ChaosMonkey, an application responsible for randomly terminating instances in production to ensure that engineers implement their services to be resilient to instance failures.
    ctas:
      - text: Netflix Techblog Article
        url: https://medium.com/netflix-techblog/application-data-caching-using-ssds-5bf25df851ef
  - company: Stripe
    url: https://stripe.com/
    logoSrc: stripe.svg
    logoSrcDark: stripe.svg
    desc: Stripe uses Go for the Stripe CLI aimed to help build, test, and manage a Stripe integration right from the terminal.
    ctas:
      - text: Stripe CLI
        url: https://github.com/stripe/stripe-cli
  - company: Uber
    url: https://uber.com/
    logoSrc: uber.svg
    logoSrcDark: uber.svg
    desc: Uber uses Go for several CLI tools, including the CLI API for Jaeger, a distributed tracing system used for monitoring microservice distributed systems.
    ctas:
      - text: CLI API for Jaeger
        url: https://www.jaegertracing.io/docs/1.14/cli/
`}}

## Get Started {#get-started .sectionHeading}

### Go books for creating CLIs

{{books `
  - title: Powerful Command-Line Applications in Go
    url: https://www.amazon.com/Powerful-Command-Line-Applications-Go-Maintainable/dp/168050696X
    thumbnail: /images/books/powerful-command-line-applications-in-go.jpg
  - title: Go in Action
    url: https://www.amazon.com/Go-Action-William-Kennedy/dp/1617291781
    thumbnail: /images/books/go-in-action.jpg
  - title: The Go Programming Language
    url: https://www.gopl.io/
    thumbnail: /images/learn/go-programming-language-book.png
  - title: Go Programming Blueprints
    url: https://github.com/matryer/goblueprints
    thumbnail: /images/learn/go-programming-blueprints.png
`}}

{{libraries `
  - title: CLI Libraries
    viewMoreUrl: https://pkg.go.dev/search?q=command%20line%20OR%20CLI
    items:
      - text: spf13/cobra
        url: https://pkg.go.dev/github.com/spf13/cobra?tab=overview
        desc: A library for creating powerful modern CLI applications and a program to generate applications and CLI applications in Go
      - text: spf13/viper
        url: https://pkg.go.dev/github.com/spf13/viper?tab=overview
        desc: A complete configuration solution for Go applications, designed to work within an app to handle configuration needs and formats
      - text: urfave/cli
        url: https://pkg.go.dev/github.com/urfave/cli?tab=overview
        desc: A minimal framework for creating and organizing command line Go applications
      - text: delve
        url: https://pkg.go.dev/github.com/go-delve/delve?tab=overview
        desc: A simple and powerful tool built for programmers used to using a source-level debugger in a compiled language
      - text: chzyer/readline
        url: https://pkg.go.dev/github.com/chzyer/readline?tab=overview
        desc: A pure Golang implementation that provides most features in GNU Readline (under MIT license)
      - text: dixonwille/wmenu
        url: https://pkg.go.dev/github.com/dixonwille/wmenu?tab=overview
        desc: An easy-to-use menu structure for CLI applications that prompts users to make choices
      - text: spf13/pflag
        url: https://pkg.go.dev/github.com/spf13/pflag?tab=overview
        desc: A drop-in replacement for Go’s flag package, implementing POSIX/GNU-style flags
      - text: golang/glog
        url: https://pkg.go.dev/github.com/golang/glog?tab=overview
        desc: Leveled execution logs for Go
      - text: go-prompt
        url: https://pkg.go.dev/github.com/c-bata/go-prompt?tab=overview
        desc: A library for building powerful interactive prompts, making it easier to build cross-platform command line tools using Go.
`}}
