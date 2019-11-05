---
title: "Go for Command-line Interfaces (CLIs)"
linkTitle: "Command-line Interfaces"
description: "With popular open source packages and a robust standard library, use Go to create fast and elegant CLIs."
date: 2019-10-04T15:26:31-04:00
series: Use Cases
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
    {{% mediaListBox img-src="/images/logos/stripe.svg" img-alt="Stripe Logo"  img-link="https://stripe.com" title="" align=top  %}}
Stripe [uses Go in one of their CLIs](https://github.com/stripe/stripe-cli) aimed to help build, test, and manage a Stripe integration right from the terminal.
{{% /mediaListBox %}}

{{% /mediaListBox %}}

## **Key Solutions**

{{< headerWithLink header="CLI Frameworks" link="https://pkg.go.dev/search?q=command+line" level=3 >}}

*   [Cobra](https://github.com/spf13/cobra), a library for creating powerful modern CLI applications and a program to generate applications and CLI applications in Go.
*   [Viper](https://github.com/spf13/viper), a complete configuration solution for Go applications, designed to work within an app to handle configuration needs and formats. Cobra and Viper work together.
*   [Delve](https://github.com/go-delve/delve), a simple and powerful tool built for programmers used to using a source-level debugger in a compiled language.
*   [urfave](https://pkg.go.dev/github.com/urfave/cli), a minimal framework for creating and organizing command line Go applications.

{{< headerWithLink header="Tools" link="https://pkg.go.dev/search?q=command+line+tool" level=3 >}}

*   [Readline](https://github.com/chzyer/readline), a pure Golang implementation that provides most features in GNU(under MIT license)
*   [wmenu](https://github.com/dixonwille/wmenu), an easy-to-use menu structure for CLI applications that prompts users to make choices
*   [pflag](https://github.com/spf13/pflag), a drop-in replacement for Go's flag package, implementing POSIX/GNU-style flags