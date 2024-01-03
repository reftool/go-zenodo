gozenodo
===

A Go library for interacting with the [Zenodo](https://zenodo.org/) API.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Language](https://img.shields.io/badge/language-go-blue.svg)](https://golang.org/)

**Please Note: This is NOT ready for production.**

# Installation

```shell
go get github.com/reftool/go-zenodo
```

# Usage

For more information, see `examples/usage.go`. I will be adding documentation soon. :-)

## Setup

```go
gozenodo.SetAccessToken("h4kmNNh9QL2ZD1YqiVW8MSqZapbAmvriroH2Bh57sjRk3SWkVjgzKO7Gn2c1")
gozenodo.SetSandboxMode(true)
```

## Create Deposition

```go
// Create blank deposition
deposition, err := gozenodo.CreateDeposition()
if err != nil {
    panic(err)
}
fmt.Println("Created Deposition with ID: ", deposition.ID)
```

## Add a file to a deposition

```go
// Add a file to the deposition bucket
newFileUpload, err := gozenodo.UploadFile(deposition.Links.Bucket, "test.txt", "examples/test.txt")
if err != nil {
    panic(err)
}
fmt.Println("Uploaded file: Key=", newFileUpload.Key)
```

## List a deposition

```go
// List all of your depositions
all, err := gozenodo.ListDepositions()
if err != nil {
    panic(err)
}

for _, d := range all {
    fmt.Println("Found Deposition: ", d.ID)
}
```

# Contributing

## Submit a ticket/issue

If you find a bug, have a feature suggestion or would like to ask for help, please create an issue in GitHubs issue tracker and try to label the issue correctly. I'll respond as soon as possible.

## Code

1. Create or pick an issue
2. Fork
3. Create a branch with its name being the first then last initial of your name, forward slash then the title using dashes. (Example: `jh/add-better-contributing-guide`)
4. Please try to add test coverage and update references/documentation for all of your changes. Update the .gitignore if it changes. Don't push secrets. etc... <3
5. Submit a PR and wait for review. :-)

**Please note that all changes you submit are licensed under the same license we use, which is the MIT license.**
