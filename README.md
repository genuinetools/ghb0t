# ghb0t

[![Travis CI](https://img.shields.io/travis/genuinetools/ghb0t.svg?style=for-the-badge)](https://travis-ci.org/genuinetools/ghb0t)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/genuinetools/ghb0t)
[![Github All Releases](https://img.shields.io/github/downloads/genuinetools/ghb0t/total.svg?style=for-the-badge)](https://github.com/genuinetools/ghb0t/releases)

A GitHub Bot to automatically delete your fork's branches after a pull request
has been merged.

> **NOTE:** This will **never** delete a branch named "master" AND will
**never** delete a branch that is not owned by the current authenticated user.
If the pull request is closed _without_ merging, it will **not** delete it.

<!-- toc -->

<!-- tocstop -->

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/genuinetools/ghb0t/releases).

#### Via Go

```console
$ go get github.com/genuinetools/ghb0t
```

## Usage

```console
$ ghb0t -h
ghb0t -  A GitHub Bot to automatically delete your fork's branches after a pull request has been merged.

Usage: ghb0t <command>

Flags:

  -d         enable debug logging (default: false)
  -interval  check interval (ex. 5ms, 10s, 1m, 3h) (default: 30s)
  -token     GitHub API token (or env var GITHUB_TOKEN) 
  -url       Connect to a specific GitHub server, provide full API URL (ex. https://github.example.com/api/v3/) (default: <none>)

Commands:

  version  Show the version information.
```
