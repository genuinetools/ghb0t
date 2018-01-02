# ghb0t

[![Travis CI](https://travis-ci.org/jessfraz/ghb0t.svg?branch=master)](https://travis-ci.org/jessfraz/ghb0t)

A GitHub Bot to automatically delete your fork's branches after a pull request
has been merged.

> **NOTE:** This will **never** delete a branch named "master" AND will
**never** delete a branch that is not owned by the current authenticated user.
If the pull request is closed _without_ merging, it will **not** delete it.

## Installation

#### Binaries

- **darwin** [386](https://github.com/jessfraz/ghb0t/releases/download/v0.1.0/ghb0t-darwin-386) / [amd64](https://github.com/jessfraz/ghb0t/releases/download/v0.1.0/ghb0t-darwin-amd64)
- **freebsd** [386](https://github.com/jessfraz/ghb0t/releases/download/v0.1.0/ghb0t-freebsd-386) / [amd64](https://github.com/jessfraz/ghb0t/releases/download/v0.1.0/ghb0t-freebsd-amd64)
- **linux** [386](https://github.com/jessfraz/ghb0t/releases/download/v0.1.0/ghb0t-linux-386) / [amd64](https://github.com/jessfraz/ghb0t/releases/download/v0.1.0/ghb0t-linux-amd64) / [arm](https://github.com/jessfraz/ghb0t/releases/download/v0.1.0/ghb0t-linux-arm) / [arm64](https://github.com/jessfraz/ghb0t/releases/download/v0.1.0/ghb0t-linux-arm64)
- **solaris** [amd64](https://github.com/jessfraz/ghb0t/releases/download/v0.1.0/ghb0t-solaris-amd64)
- **windows** [386](https://github.com/jessfraz/ghb0t/releases/download/v0.1.0/ghb0t-windows-386) / [amd64](https://github.com/jessfraz/ghb0t/releases/download/v0.1.0/ghb0t-windows-amd64)

#### Via Go

```bash
$ go get github.com/jessfraz/ghb0t
```

## Usage

```
$ ghb0t -h
       _     _      ___  _
  __ _| |__ | |__  / _ \| |_
 / _` | '_ \| '_ \| | | | __|
| (_| | | | | |_) | |_| | |_
 \__, |_| |_|_.__/ \___/ \__|
 |___/

 A GitHub Bot to automatically delete your fork's branches after a pull request has been merged.
 Version: v0.1.0
 Build: be7363d

  -d    run in debug mode
  -interval string
        check interval (ex. 5ms, 10s, 1m, 3h) (default "30s")
  -token string
        GitHub API token (or env var GITHUB_TOKEN)
  -v    print version and exit (shorthand)
  -version
        print version and exit
```
