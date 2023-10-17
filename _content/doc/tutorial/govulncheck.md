<!--{
  "Title": "Tutorial: Find and fix vulnerable dependencies with govulncheck",
  "HideTOC": true,
  "Breadcrumb": true
}-->

Govulncheck is a low-noise tool that helps you find and fix vulnerable
dependencies in your Go projects. It does this by scanning your project's
dependencies for known vulnerabilities and then identifying any direct or
indirect calls to those vulnerabilities in your code.

In this tutorial, you will learn how to use govulncheck to scan a simple
program for vulnerabilities. You will also learn how to prioritize and
evaluate vulnerabilities so that you can focus on fixing the most important
ones first.

To learn more about govulncheck, see the
[govulncheck documentation](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck),
and this [blog post on vulnerability management](https://go.dev/blog/vuln) for Go.
We'd also love to [hear your feedback](https://go.dev/s/govulncheck-feedback).

## Prerequisites

- **Go 1.18 or later.** Govulncheck is designed to work with Go 1.18 onwards.
  (For installation instructions, see [Installing Go](https://go.dev/doc/install).)
  We recommend using the latest version of Go to follow this tutorial.
- **A code editor.** Any editor you have will work fine.
- **A command terminal.** Go works well using any terminal on Linux and Mac, and on PowerShell or cmd in Windows.

The tutorial will take you through the following steps:

1. Create a sample Go module with a vulnerable dependency
2. Install and run govulncheck
3. Evaluate vulnerabilities
4. Upgrade vulnerable dependencies

## Create a sample Go module with a vulnerable dependency

**Step 1.** To begin, create a new folder called `vuln-tutorial` and initialize a Go module.
(If you are new to Go modules, check out https://go.dev/doc/tutorial/create-module).

For example, from your home directory, run the following:

```
$ mkdir vuln-tutorial
$ cd vuln-tutorial
$ go mod init vuln.tutorial
```

**Step 2.** Create a file called `main.go` within the `vuln-tutorial` folder, and copy
the following code into it:

```
package main

import (
        "fmt"
        "os"

        "golang.org/x/text/language"
)

func main() {
        for _, arg := range os.Args[1:] {
                tag, err := language.Parse(arg)
                if err != nil {
                        fmt.Printf("%s: error: %v\n", arg, err)
                } else if tag == language.Und {
                        fmt.Printf("%s: undefined\n", arg)
                } else {
                        fmt.Printf("%s: tag %s\n", arg, tag)
                }
        }
}
```

This sample program takes a list of language tags as command line arguments
and prints a message for each tag indicating if it was parsed successfully,
the tag is undefined, or whether there was an error while parsing the tag.

**Step 3.** Run `go mod tidy`, which will populate the `go.mod` file with all the
dependencies required by the code you added to `main.go` in the previous step.

From the `vuln-tutorial` folder, run:

```
$ go mod tidy
```

You should see this output:

```
go: finding module for package golang.org/x/text/language
go: downloading golang.org/x/text v0.9.0
go: found golang.org/x/text/language in golang.org/x/text v0.9.0
```

**Step 4.** Open your `go.mod` file to verify that it looks like this:

```
module vuln.tutorial

go 1.20

require golang.org/x/text v0.9.0
```

**Step 5.** Downgrade the version of `golang.org/x/text` to v0.3.5, which contains known
vulnerabilities. Run:

```
$ go get golang.org/x/text@v0.3.5
```

You should see this output:

```
go: downgraded golang.org/x/text v0.9.0 => v0.3.5
```

The `go.mod` file should now read:

```
module vuln.tutorial

go 1.20

require golang.org/x/text v0.3.5
```

Now, let’s see govulncheck in action.


## Install and run govulncheck

**Step 6.** Install govulncheck with the `go install` command:

```
$ go install golang.org/x/vuln/cmd/govulncheck@latest
```

**Step 7.** From the folder you want to analyze (in this case, `vuln-tutorial`). Run:

```
$ govulncheck ./...
```

You should see this output:

```
govulncheck is an experimental tool. Share feedback at https://go.dev/s/govulncheck-feedback.

Using go1.20.3 and govulncheck@v0.0.0 with
vulnerability data from https://vuln.go.dev (last modified 2023-04-18 21:32:26 +0000 UTC).

Scanning your code and 46 packages across 1 dependent module for known vulnerabilities...
Your code is affected by 1 vulnerability from 1 module.

Vulnerability #1: GO-2021-0113
  Due to improper index calculation, an incorrectly formatted
  language tag can cause Parse to panic via an out of bounds read.
  If Parse is used to process untrusted user inputs, this may be
  used as a vector for a denial of service attack.

  More info: https://pkg.go.dev/vuln/GO-2021-0113

  Module: golang.org/x/text
    Found in: golang.org/x/text@v0.3.5
    Fixed in: golang.org/x/text@v0.3.7

    Call stacks in your code:
      main.go:12:29: vuln.tutorial.main calls golang.org/x/text/language.Parse

=== Informational ===

Found 1 vulnerability in packages that you import, but there are no call
stacks leading to the use of this vulnerability. You may not need to
take any action. See https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck
for details.

Vulnerability #1: GO-2022-1059
  An attacker may cause a denial of service by crafting an
  Accept-Language header which ParseAcceptLanguage will take
  significant time to parse.
  More info: https://pkg.go.dev/vuln/GO-2022-1059
  Found in: golang.org/x/text@v0.3.5
  Fixed in: golang.org/x/text@v0.3.8

```

### Interpreting the output

<font size="2">  *Note: If you are not using the latest version of Go,
you may see additional vulnerabilities from the standard library. </font>

Our code is affected by one vulnerability,
[GO-2021-0113](https://pkg.go.dev/vuln/GO-2021-0113), because it directly calls
the `Parse` function of `golang.org/x/text/language` at a vulnerable version
(v0.3.5).

Another vulnerability, [GO-2022-1059](https://pkg.go.dev/vuln/GO-2022-1059),
exists in the `golang.org/x/text` module at v0.3.5.  However, it is reported as
"Informational" because our code never (directly or indirectly) calls any of
its vulnerable functions.

Now, let's evaluate the vulnerabilities and determine an action to take.

### Evaluate vulnerabilities

a. Evaluate vulnerabilities.

First, read the description of the vulnerability and determine if it actually
applies to your code and your use case. If you need more information, visit
the "More info" link.

Based on the description, vulnerability GO-2021-0113 can cause a panic when
`Parse` is used to process untrusted user inputs. Let's suppose that we intend
our program to withstand untrusted inputs, and we are concerned about denial of
service, so the vulnerability likely applies.

GO-2022-1059 likely does not affect our code, because our code does not call
any vulnerable functions from that report.

b. Decide on an action.

To mitigate GO-2021-0113, we have a few options:
- **Option 1: Upgrade to a fixed version.** If there is a fix available,
  we can remove a vulnerable dependency by upgrading to a fixed version of the module.
- **Option 2: Stop using the vulnerable symbol(s).** We could choose to
  remove all calls to the vulnerable function in our code.
  We would need to find an alternative or implement it ourselves.

In this case, a fix is available, and the `Parse` function is integral to our
program. Let's upgrade our dependency to the "fixed in" version, v0.3.7.

We decided to deprioritize fixing the informational vulnerability,
GO-2022-1059, but because it is in the same module as GO-2021-0113, and because the fixed in version for it is v0.3.8, we can
easily remove both at the same time by upgrading to v0.3.8.

## Upgrade vulnerable dependencies

Luckily, upgrading vulnerable dependencies is quite simple.

**Step 8.** Upgrade `golang.org/x/text` to v0.3.8:

```
$ go get golang.org/x/text@v0.3.8
```

You should see this output:

```
go: upgraded golang.org/x/text v0.3.5 => v0.3.8
```

(Note that we could have also chosen to upgrade to `latest`, or any other version after v0.3.8).

**Step 9.** Now run govulncheck again:

```
$ govulncheck ./...
```

You will now see this output:

```
govulncheck is an experimental tool. Share feedback at https://go.dev/s/govulncheck-feedback.

Using go1.20.3 and govulncheck@v0.0.0 with
vulnerability data from https://vuln.go.dev (last modified 2023-04-06 19:19:26 +0000 UTC).

Scanning your code and 46 packages across 1 dependent module for known vulnerabilities...
No vulnerabilities found.
```

Finally, govulncheck confirms that there are no vulnerabilities found.

By regularly scanning your dependencies with command govulncheck, you can
safeguard your codebase by identifying, prioritizing, and addressing
vulnerabilities.
