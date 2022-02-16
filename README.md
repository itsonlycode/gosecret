# gosecret

Gosecret is a program for storing structured information encrypted in a git
repository based on the model used in gopass

-- WARNING: NOT READY FOR USE ... might destroy your password store !!!

# Aknowledgement

All credit goes to gopasswd for all the gopass code in this project.
This is basically gopass with password related content ripped out of it.


# gopass



<p align="center">
    <img src="docs/logo.png" height="250" alt="gopass Gopher by Vincent Leinweber, remixed from the Renée French original Gopher" title="gopass Gopher by Vincent Leinweber, remixed from the Renée French original Gopher" />
</p>

# gopass

[![Build Status](https://img.shields.io/github/workflow/status/gopasspw/gopass/Build%20gopass/master)](https://github.com/gopasspw/gopass/actions/workflows/build.yml?query=branch%3Amaster)
[![Packaging status](https://repology.org/badge/tiny-repos/gopass.svg)](https://repology.org/project/gopass/versions)
[![Go Report Card](https://goreportcard.com/badge/github.com/gopasspw/gopass)](https://goreportcard.com/report/github.com/gopasspw/gopass)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/gopasspw/gopass/blob/master/LICENSE)
[![Github All Releases](https://img.shields.io/github/downloads/gopasspw/gopass/total.svg)](https://github.com/gopasspw/gopass/releases)
[![codecov](https://codecov.io/gh/gopasspw/gopass/branch/master/graph/badge.svg)](https://codecov.io/gh/gopasspw/gopass)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/1899/badge)](https://bestpractices.coreinfrastructure.org/projects/1899)
[![Gopass Slack](https://img.shields.io/badge/%23gopass-Slack-brightgreen)](https://docs.google.com/forms/d/e/1FAIpQLScxOPX_OLDaG5ak2E1kNdcFw9fJvPCr8xUaPGLyW8cyNUEnJw/viewform?usp=sf_link)

## Introduction

gopass is a password manager for the command line written in Go. It supports all major operating systems (Linux, MacOS, BSD) as well as Windows.

For detailed usage and installation instructions please check out our [documentation](docs/).

## Features

Please see [docs/features.md](https://github.com/gopasspw/gopass/blob/master/docs/features.md) for an extensive list of all features along with several usage examples. Some examples are available in our
[example password store](https://github.com/gopasspw/password-store-example).

| **Feature**                 | **State**     | **Description**                                                   |
| --------------------------- | ------------- | ----------------------------------------------------------------- |
| Secure secret storage       | *stable*      | Securely storing encrypted secrets                                |
| Recipient management        | *beta*        | Easily manage multiple users of each store                        |
| Multiple stores             | *stable*      | Mount multiple stores in your root store, like file systems       |
| PAGER support               | *stable*      | Automatically invoke a pager on long output                       |
| JSON API                    | *integration* | Allow gopass to be used as a native extension for browser plugins |
| Automatic fuzzy search      | *stable*      | Automatically search for matching store entries if a literal entry was not found |
| gopass sync                 | *stable*      | Easy to use syncing of remote repos and GPG keys                  |
| Desktop Notifications       | *stable*      | Display desktop notifications and completing long running operations |
| REPL                        | *beta*        | Integrated Read-Eval-Print-Loop shell with autocompletion. |
| Extensions                  |               | Extend gopass with custom commands using our API                  |

## Installation

Please see [docs/setup.md](https://github.com/gopasspw/gopass/blob/master/docs/setup.md).

If you have [Go](https://golang.org/) 1.16 (or greater) installed:

```bash
go get github.com/gopasspw/gopass
```

WARNING: Please prefer releases, unless you want to contribute to the
development of gopass. The master branch might not be very well tested and
can contain breaking changes without further notice.

## Getting Started

Either initialize a new git repository or clone an existing one.

### New information store

```
$ gopass setup

Hint: `gopass setup` will use `gpg` encryption and `git` storage by default.

### Existing password store

```
$ gopass clone git@gitlab.example.org:john/passwords.git

Your password store is ready to use! Have a look around: `gopass list`
```

## Upgrade

To use the self-updater run:
```bash
gopass update
```

or to upgrade with Go installed, run:
```bash
go get -u github.com/itsonlycode/gosecret
```

Otherwise, use the setup docs mentioned in the installation section to reinstall the latest version.

## Development

This project uses [GitHub Flow](https://guides.github.com/introduction/flow/). In other words, create feature branches from master, open an PR against master, and rebase onto master if necessary.

We aim for compatibility with the [latest stable Go Release](https://golang.org/dl/) only.

While this project is maintained by volunteers in their free time we aim to triage issues weekly and release a new version at least every quarter.

## Credit & License

gosecret is licensed under the terms of the MIT license. You can find the complete text in `LICENSE`.

Please refer to the Git commit log for a complete list of contributors.

## Community

gosecret is developed in the open. Here are some of the channels we use to communicate and contribute:

* Issue tracker: Use the [GitHub issue tracker](https://github.com/itsonlycode/gosecret/issues) to file bugs and feature requests.

## Contributing

We welcome any contributions.

<!---
Please see the [CONTRIBUTING.md](https://github.com/gopasspw/gopass/blob/master/CONTRIBUTING.md) file for instructions on how to submit changes.
--->
