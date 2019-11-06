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
  - title: 
    url: 
    thumbnail: 
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

### _CLI developers prefer Go for portability, performance, and ease of creation_

## **Why use Go for CLIs**

Command line interfaces (CLIs), unlike graphical user interfaces (GUIs), use text-only instructions and syntax to interact with applications and operating systems. While desktop and mobile are primarily GUI-based, cloud and infrastructure are CLI-based due to their easy automation and remote capabilities. CLIs allow users to perform specific computing tasks by typing text commands and receiving system replies via text outputs, and CLIs are easily automated and can be combined to create custom workflows.

{{% gopher gopher=happy align=right %}}

Developers of CLIs find Go to be ideal for designing their applications. Go compiles very quickly into a single binary, works across platforms with a consistent style, and brings a strong development community. “The design of Go lends itself incredibly well to [many] styles of application,” [writes Elliot Forbes](https://tutorialedge.net/golang/building-a-cli-in-go/), software engineer at JP Morgan Chase. “And the ability to cross-compile a binary executable for all major platforms easily is a massive win.”

Specifically, **programs written in Go run on any system without requiring any existing libraries, runtimes, or dependencies**. And **programs written in Go have an immediate startup time**—similar to C or C++ but unobtainable with other programming languages. 

{{% pullquote author="Carolyn Van Slyck, Senior Software Engineer at Microsoft." link="https://www.youtube.com/watch?v=eMz0vni6PAw&list=PL2ntRZ1ySWBdDyspRTNBIKES1Y-P__59_&index=11&t=0s" %}}
CLIs are best designed with predictable, task-oriented commands and you want to use Go.
{{% /pullquote %}}

## **Who uses Go for CLIs**

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

{{% /mediaListBox %}}


## **How to use Go for CLIs**

From a Windows or Mac machine, developers can build a Go program for every one of the dozens of architectures and operating systems Go supports—with no complicated build farms. No other compiled language is so readily deployed.

When developing CLIs in Go, programmers find several tools hugely helpful: Cobra, Viper, and debugger Delve.

[Cobra](https://github.com/spf13/cobra), for example, is both a library for creating powerful modern CLI applications and a program to generate applications and CLI applications in Go. Cobra powers CoreOS, Delve, Docker, Dropbox, Git Lfs, Hugo, Kubernetes, and many other popular apps where handlers/commands can live in separate files or modules. “It also makes documenting each command really simple,” says [Alex Ellis](https://blog.alexellis.io/5-keys-to-a-killer-go-cli/), founder of OpenFaaS.

"Cobra is a great product to write small tools or even large ones," adds Francesc Campoy, VP of product at DGraph Labs and producer of [Just For Func videos](https://www.youtube.com/watch?v=WvWPGVKLvR4). "It's more of a framework than a library, because when you call the binary that would create a skeleton, then you would be adding code in between."

Cobra allows developers to build command-line utilities with commands, subcommands, aliases, configuration files, etc. All Cobra projects follow the [same development cycle](https://www.linode.com/docs/development/go/using-cobra/):  “You first use the Cobra tool to initialize a project, then you create commands and subcommands, and finally you make the desired changes to the generated Go source files in order to support the desired functionality.”

"The framework Cobra provides a generator that adds some boilerplate code for you," [says Nick Geudens](https://ordina-jworks.github.io/development/2018/10/20/make-your-own-cli-with-golang-and-cobra.html), Java consultant at Ordina Belgium. "This is handy because now you can focus more on the logic of your CLI instead of figuring out how to parse flags."

[Viper](https://github.com/spf13/viper) is a complete configuration solution for Go applications, designed to work within an app to handle configuration needs and formats. Cobra and Viper work together.

"If you don’t want to pollute your command line, or if you’re working with sensitive data which you don’t want to show up in the history, it’s a good idea to work with environment variables. To do this, you can use Viper," [suggests Geudens](https://ordina-jworks.github.io/development/2018/10/20/make-your-own-cli-with-golang-and-cobra.html). "Cobra already uses Viper in the generated code."

Viper [supports nested structures](https://scene-si.org/2017/04/20/managing-configuration-with-viper/) in the configuration, allowing CLI developers to manage the configuration for multiple parts of a large application.

"I was tasked with building our CLI tool and found two really great projects, Cobra and Viper, which make building CLI’s easy," [explains Steve Domino](https://medium.com/@skdomino/writing-better-clis-one-snake-at-a-time-d22e50e60056), senior engineer and architect at Strala. "Individually they are very powerful, very flexible and very good at what they do. But together they will help you show your next CLI who is boss!"

For debugging Go code, [Delve](https://github.com/go-delve/delve) is a simple and powerful tool built for programmers used to using a source-level debugger in a compiled language.

## Key Solutions

### Go books on CLIs 

{{% books %}}

{{< headerWithLink header="Frameworks" link="https://pkg.go.dev/search?q=framework" level=3 >}} 

## **Go solutions to legacy challenges**

CLIs face certain challenges in development and distribution.  For example, porting CLIs across operating systems can be difficult with dependencies and often very large binary files. Speed can likewise be a challenge—compiling, loading, or executing—as can creating REST clients for HTTP, XML, and JSON.

Go rises above all these challenges.

Again, programs written in Go run on any system without requiring any existing libraries, runtimes, or dependencies. Ellis summarizes [why Go is best for CLIs](https://blog.alexellis.io/5-keys-to-a-killer-go-cli/): because Go compiles to a single static binary, Go's consistent style is unambiguous and easy for on-boarding, Go loads fast on every platform, and Go makes it easy to create REST clients.

While CLIs do not have graphical user interfaces, the UNIX philosophy that drives them suggests that they be simple, clear, composable, extensible, modular, and small. Go delivers all six elements. With Go, initialization is expressive, automatic, and easy to use. Syntax is clean and light on keywords. Go combines the ease of programming in an interpreted, dynamically typed language with the efficiency and safety of a statically typed, compiled language. And goroutines have little overhead beyond the memory for the stack (which is just a few kilobytes).

{{% pullquote author="Alex Ellis, founder of OpenFaaS." link="https://blog.alexellis.io/5-keys-to-a-killer-go-cli/" %}}
There are many reasons to use Go to build your next killer CLI ... From the speed of compilation and execution, the availability of built-in or high-quality packages, to the ease of automation.
{{% /pullquote %}}

## **Resources for learning more**

*   [Cobra](https://github.com/spf13/cobra) - Commander for modern Go CLI interactions
*   [Viper](https://github.com/spf13/viper) - Go configuration with fangs
*   [Delve](https://github.com/derekparker/delve) - Go debugger
*   [Readline](https://github.com/chzyer/readline) - pure Golang implementation that provides most features in GNU(under MIT license)
*   [wmenu](https://github.com/dixonwille/wmenu) - easy-to-use menu structure for CLI applications that prompts users to make choices
*   [pflag](https://github.com/spf13/pflag) - drop-in replacement for Go's flag package, implementing POSIX/GNU-style flags