# Go HTTP handler function detection

This repo has some playthings that I used for learning buildpacks as well as how
to use Go tools for doing inspection of source code in Go. Basic idea is to be
able to use these tools to detect if a given Go source code is compatible with
[net.http.HandlerFunc](https://godoc.org/net/http#HandlerFunc). Rough idea being
that we'll inspect the files and determine if we can autogen enough scaffolding
to just inject the HTTP handler code in, and have the user just be able to write the
handler code alone without having to spin up the web server for
it. 

# Building

go-buildpack/bin/detect is a binary that has been built like this:

```shell
GOOS=linux GOARCH=amd64 go build ./cmd/detect/main.go && cp main ~/buildpack-go/go-buildpack/bin/detect
```

If you want to target a different arch / OS, you'd have to modify the GO* params
as appropriate.

# Supported function signatures

Basically, you just need to have a file that imports the net/http and
has an externally visible (capitalized) function that implements the HTTP handler.
