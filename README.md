# github-issue-fetcher

## Installation

```
$ make deps
$ make
$ ./github-issue-fetcher
```

## Setup

    $ github-issue-fetcher init

The initialization will be update your `~/.gitconfig`.

1. Add `github-issue-fetcher.user`
2. Add `github-issue-fetcher.token`

## Maintenance

It's possible to release to GitHub using `make` command.

```
$ git checkout master
$ git pull
$ make dist
$ make release
```
