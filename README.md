# tealib

`tealib` is a proof-of-concept Go library for creating programs that follow the [Elm architecture](https://guide.elm-lang.org/architecture/) (TEA, hence the name.) The primary goal is that library users need not write impure code and can offload all the messy side effects to the library.

To get started, have a look at the samples.

 * [exit](examples/exit/main.go) - a program that exits immediately.
 * [echo](examples/echo/) - a simple line-based echo server.

Given the experimental nature of this library, expect breaking changes any moment without warning. Make sure to vendor it.
