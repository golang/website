---
title: "Go for Command-line Interfaces (CLIs)"
linkTitle: "Command-line Interfaces"
description: "With popular open source packages and a robust standard library, use Go to create fast and elegant CLIs."
date: 2019-10-04T15:26:31-04:00
series: Use Cases
books:
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
resources:
- name: icon
  src: globe.png
  params:
    alt: globe
- name: icon-white
  src: globe-white.png
  params:
    alt: globe
---

## CLI developers prefer Go for portability, performance, and ease of creation

Command line interfaces (CLIs), unlike graphical user interfaces (GUIs), are text-only. Cloud and infrastructure applications are primarily CLI-based due to their easy automation and remote capabilities. 

{{% gopher gopher=happy align=right %}}

Developers of CLIs find Go to be ideal for designing their applications. Go compiles very quickly into a single binary, works across platforms with a consistent style, and brings a strong development community.  From a single Windows or Mac laptop, developers can build a Go program for every one of the dozens of architectures and operating systems Go supports in a matter of seconds, no complicated build farms are needed. No other compiled language can be built as portably or quickly. Go applications are built into a single self contained binary making installing Go applications trivial.

{{% pullquote author="Elliot Forbes, Software Engineer at JP Morgan" link="https://tutorialedge.net/golang/building-a-cli-in-go/" %}}
The design of Go lends itself incredibly well to [many] styles of application, and the ability to cross-compile a binary executable for all major platforms easily is a massive win.”
{{% /pullquote %}}

Specifically, **programs written in Go run on any system without requiring any existing libraries, runtimes, or dependencies**. And **programs written in Go have an immediate startup time**—similar to C or C++ but unobtainable with other programming languages. 

{{% pullquote author="Carolyn Van Slyck, Senior Software Engineer at Microsoft." link="https://www.youtube.com/watch?v=eMz0vni6PAw&list=PL2ntRZ1ySWBdDyspRTNBIKES1Y-P__59_&index=11&t=0s" %}}
CLIs are best designed with predictable, task-oriented commands and you want to use Go.
{{% /pullquote %}}

## Featured Go users and projects

{{% mediaList %}}
    {{% mediaListBox img-src="/images/logos/comcast.svg" img-alt="Comcast Logo" img-link="https://xfinity.com" title="" align=top %}}
Comcast [uses Go for a CLI client](https://github.com/Comcast/pulsar-client-go/blob/master/cli/main.go) used to publish and subscribe to it's high-traffic sites. The company also supports an [open source client library](https://github.com/Comcast/pulsar-client-go) which is written in Go - designed for working with Apache Pulsar.
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/github.svg" img-alt="GitHub Logo"  img-link="https://github.com" title="" align=top  %}}
GitHub [uses Go for a command-line tool](https://github.com/github/hub) that makes it easier to work with GitHub, wrapping git in order to extend it with extra features and commands.
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/uber.svg" img-alt="Uber Logo"  img-link="https://uber.com" title="" align=top  %}}
Uber uses Go for several CLI tools, including the [CLI API for Jaeger](https://www.jaegertracing.io/docs/1.14/cli/), a distributed tracing system used for monitoring micro-service distributed systems.
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/stripe.svg" img-alt="Stripe Logo"  img-link="https://stripe.com" title="" align=top  %}}
Stripe [uses Go for the Stripe CLI](https://github.com/stripe/stripe-cli) aimed to help build, test, and manage a Stripe integration right from the terminal.
{{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/netflix.svg" img-alt="Netflix Logo" title="" align=top %}}
Netflix uses Go to build the CLI application [ChaosMonkey](https://medium.com/netflix-techblog/application-data-caching-using-ssds-5bf25df851ef), an application responsible for randomly terminating instances in production to ensure that engineers implement their services to be resilient to instance failures. 
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/mongodb.svg" img-alt="MongoDB Logo"  img-link="https://mongodb.com" title="" align=top  %}}
MongoDB choose to [implement their Backup CLI Tool in Go](https://www.mongodb.com/blog/post/go-agent-go) citing Go's "C-like syntax, strong standard library, the resolution of concurrency problems via goroutines, and painless multi-platform distribution" as reasons. 
{{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/kubernetes.svg" img-alt="Kubernetes Logo"  img-link="https://kubernetes.com" title="" align=top  %}}
Kubernetes is one of the most popular Go CLI applications. [Kubernetes Creator, Joe Beda, said that for writing Kubernetes](https://blog.gopheracademy.com/birthday-bash-2014/kubernetes-go-crazy-delicious/), "Go was the only logical choice". Calling Go "the sweet spot" between low level languages like C++ and high level languages like Python. 
{{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/hugo.svg" img-alt="Hugo Logo"  img-link="https://gohugo.io" title="" align=top  %}}
Hugo is one of the most popular Go CLI applications powering thousands of sites including this one. One reason for it's popularity is it's ease of install thanks to Go. Hugo author [Bjørn Erik Pedersen writes](https://gohugo.io/news/lets-celebrate-hugos-5th-birthday/) "The single binary takes most of the pain out of installation and upgrades."
{{% /mediaListBox %}}

{{% /mediaListBox %}}

## How to use Go for CLIs

When developing CLIs in Go, programmers find several tools hugely helpful: Cobra, Viper, and debugger Delve.

 {{< pkg "github.com/spf13/cobra" Cobra >}} is both a library for creating powerful modern CLI applications and a program to generate applications and CLI applications in Go. Cobra powers CoreOS, Delve, Docker, Dropbox, Git Lfs, Hugo, Kubernetes, and many other popular apps where handlers/commands can live in separate files or modules. “It also makes documenting each command really simple,” says [Alex Ellis](https://blog.alexellis.io/5-keys-to-a-killer-go-cli/), founder of OpenFaaS.

{{% pullquote author="Francesc Campoy, VP of product at DGraph Labs and producer of Just For Func videos" link="https://www.youtube.com/watch?v=WvWPGVKLvR4" %}}
Cobra is a great product to write small tools or even large ones. It's more of a framework than a library, because when you call the binary that would create a skeleton, then you would be adding code in between."
{{% /pullquote %}}

Cobra allows developers to build command-line utilities with commands, subcommands, aliases, configuration files, etc. All Cobra projects follow the [same development cycle](https://www.linode.com/docs/development/go/using-cobra/):  “You first use the Cobra tool to initialize a project, then you create commands and subcommands, and finally you make the desired changes to the generated Go source files in order to support the desired functionality.”

{{< gopher gopher=peach  >}}
 {{< pkg "github.com/spf13/viper" Viper >}} is a complete configuration solution for Go applications, designed to work within an app to handle configuration needs and formats. Cobra and Viper work together.

"If you don’t want to pollute your command line, or if you’re working with sensitive data which you don’t want to show up in the history, it’s a good idea to work with environment variables. To do this, you can use Viper," [suggests Geudens](https://ordina-jworks.github.io/development/2018/10/20/make-your-own-cli-with-golang-and-cobra.html). "Cobra already uses Viper in the generated code."

Viper [supports nested structures](https://scene-si.org/2017/04/20/managing-configuration-with-viper/) in the configuration, allowing CLI developers to manage the configuration for multiple parts of a large application.

{{< pullquote author="Steve Domino, senior engineer and architect at Strala" link="https://medium.com/@skdomino/writing-better-clis-one-snake-at-a-time-d22e50e60056" >}}
I was tasked with building our CLI tool and found two really great projects, Cobra and Viper, which make building CLI’s easy. Individually they are very powerful, very flexible and very good at what they do. But together they will help you show your next CLI who is boss!
{{< /pullquote >}}

For debugging Go code, {{< pkg "github.com/go-delve/delve" delve >}} is a simple and powerful tool built for programmers used to using a source-level debugger in a compiled language.

## Key Solutions

### Go books for creating CLIs 

{{% books %}}

{{< headerWithLink header="CLI libraries" search="command line OR CLI" level=3 >}} 

*   {{< pkg "github.com/spf13/cobra" >}}, a library for creating powerful modern CLI applications and a program to generate applications and CLI applications in Go
*   {{< pkg "github.com/spf13/viper" >}}, a complete configuration solution for Go applications, designed to work within an app to handle configuration needs and formats
*   {{< pkg "github.com/urfave/cli" >}}, a minimal framework for creating and organizing command line Go applications
*   {{< pkg "github.com/go-delve/delve" delve >}}, a simple and powerful tool built for programmers used to using a source-level debugger in a compiled language
*   {{< pkg "https://github.com/chzyer/readline" >}}, a pure Golang implementation that provides most features in GNU Readline (under MIT license)
*   {{< pkg "https://github.com/dixonwille/wmenu" >}}, an easy-to-use menu structure for CLI applications that prompts users to make choices
*   {{< pkg "https://github.com/spf13/pflag" >}}, a drop-in replacement for Go's flag package, implementing POSIX/GNU-style flags
*   {{< pkg "https://github.com/golang/glog" >}}, leveled execution logs for Go