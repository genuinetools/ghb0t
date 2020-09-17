# ghb0t

[![make-all](https://github.com/genuinetools/ghb0t/workflows/make%20all/badge.svg)](https://github.com/genuinetools/ghb0t/actions?query=workflow%3A%22make+all%22)
[![make-image](https://github.com/genuinetools/ghb0t/workflows/make%20image/badge.svg)](https://github.com/genuinetools/ghb0t/actions?query=workflow%3A%22make+image%22)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/genuinetools/ghb0t)
[![Github All Releases](https://img.shields.io/github/downloads/genuinetools/ghb0t/total.svg?style=for-the-badge)](https://github.com/genuinetools/ghb0t/releases)

A GitHub Bot to automatically delete your fork's branches after a pull request
has been merged.

> **NOTE:** This will **never** delete a branch named "master" AND will
**never** delete a branch that is not owned by the current authenticated user.
If the pull request is closed _without_ merging, it will **not** delete it.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Installation](#installation)
    - [Binaries](#binaries)
    - [Via Go](#via-go)
- [Usage](#usage)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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
