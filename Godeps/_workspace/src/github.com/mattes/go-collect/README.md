# go-collect

[![Build Status](https://travis-ci.org/mattes/go-collect.svg?branch=v0)](https://travis-ci.org/mattes/go-collect)
[![GoDoc](https://godoc.org/gopkg.in/mattes/go-collect.v0?status.svg)](https://godoc.org/gopkg.in/mattes/go-collect.v0)

``go-collect`` fetches ``key=value`` content from different [sources](https://github.com/mattes/go-collect/tree/v0/source)
and from command line [flags](https://github.com/mattes/go-collect/tree/v0/flags)
and merges everything into [data](https://github.com/mattes/go-collect/tree/v0/data) objects.

Available sources:

  * [Files](https://github.com/mattes/go-collect/tree/v0/source/file)
  * [URL queries via flags](https://github.com/mattes/go-collect/tree/v0/source/urlquery)
  * Please feel free to add more sources, just implement the [Source interface](https://godoc.org/gopkg.in/mattes/go-collect.v0#Source)

# Usage

```
go get http://gopkg.in/mattes/go-collect.v0
```

# Tools built with go-collect

* [fugu](https://github.com/mattes/fugu) - Swiss Army knife for Docker

