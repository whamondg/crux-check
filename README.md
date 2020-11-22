## Overview

Crux-check is a simple utility for checking how real-world Chrome users experience a particular website. It uses the [crux-api](https://web.dev/chrome-ux-report-api) to fetch the current field data for a site and displays an overview of the [Core Web Vitals](https://web.dev/vitals/#core-web-vitals).  For each metric it indicates if the 75th percentile achieves a good score.

## Installation
```
  > brew tap whamondg/homebrew-crux-check
  > brew install crux-check
  > crux-check -h
```
## Usage

`crux-check -u https://www.example.com`

Also supports a "," separated list of URLs:

`crux-check -u https://www.example.com,https://www.example.com/foo`

## Releasing

[Goreleaser](https://goreleaser.com) is used to build, release and publish crux-check. It can be installed on a mac via Homebrew.

To verify the build configuration locally and build a binary for testing a dry-run can be executed:

`goreleaser --snapshot --skip-publish --rm-dist`

Using [GitHub Actions](https://docs.github.com/en/free-pro-team@latest/actions) crux-check has been configured to build automatically when code is pushed to the main branch.

In order for a new version of the homebrew formula to be published the repo needs to have been tagged.  For example:

`git tag -a v0.1.0 -m "First working version" && git push origin v0.0.1`

