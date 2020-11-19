# Go CloudEvents function detection

This repo has some playthings that I used for learning buildpacks as well as how
to use Go tools for doing inspection of source code in Go. Basic idea is to be
able to use these tools to detect if a given Go source code is compatible with
[CloudEvents go-sdk](https://github.com/cloud-events/sdk-go/). Rough idea being
that we'll inspect the files and determine if we can autogen enough scaffolding
to just inject the function code in, and have the user just be able to write the
function code alone without having to spin up the web server for
it. But... Wait, there's more, because the function specified has nothing that
specifies it as an HTTP handler, we should be able to build a buildpack that
actually is protocol agnostic, and the function could be built targeting
different protocols.

# Building

go-buildpack/bin/detect is a binary that has been built like this:

```shell
GOOS=linux GOARCH=amd64 go build ./cmd/detect/main.go && cp main ~/buildpack-go/go-buildpack/bin/detect
```

If you want to target a different arch / OS, you'd have to modify the GO* params
as appropriate.

# Supported function signatures

While the CloudEvents SDK supports multiple signatures, I thought the most
useful ones actually involve events, so the supported function signatures are
these:

```
func(event.Event)
func(event.Event) protocol.Result
func(context.Context, event.Event)
func(context.Context, event.Event) protocol.Result
func(event.Event) *event.Event
func(event.Event) (*event.Event, protocol.Result)
func(context.Context, event.Event) *event.Event
func(context.Context, event.Event) (*event.Event, protocol.Result)
```
