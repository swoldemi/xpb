package main

// Version defines the current version of the XPB client
const Version = "1.0.0-rc.1"

// GitSHA will be replaced at build time with
// "-ldflags -X main.GitSHA=xxx", where const is not supported
var GitSHA = ""
