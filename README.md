# feedme

An Atom feed generator adhering to RFC 4287 standards

## Table of Contents

1. [Introduction](#introduction)
2. [Usage](#usage)

## Introduction

The purpose of this service is to generate Atom feeds for sites, services, or protocols which do not natively generate feeds (or do not generate useful feeds).

Atom is defined according to [RFC 4287](https://datatracker.ietf.org/doc/html/rfc4287). Note that not all possible definitions of an Atom feed can be generated with this repository (see [Appendix B](https://datatracker.ietf.org/doc/html/rfc4287#appendix-B)); however, each generated Atom feed can be validated according to this standard.

## Usage

Feedme's feed generating capabilities can be accessed via its atom package in golang:

```go
import (
	"fmt"
	"time"

	"git.sr.ht/~bossley9/feedme/pkg/atom"
)

func main() {
	date, errDate := time.Parse("2006-01-02 15:04", "2022-07-08 08:26")
	if errDate != nil {
		t.Error(errDate)
	}

	feed, errCreateFeed := atom.CreateFeed("example.com", "My Website", date)
	if errCreateFeed != nil {
		t.Error(errCreateFeed)
	}

	fmt.Println(feed.String())

	// outputs:
  //
	// <?xml version="1.0" encoding="UTF-8"?>
	// <feed xmlns="http://www.w3.org/2005/Atom">
	//   <id>example.com</id>
	//   <title>My Website</title>
	//   <updated>2022-07-08T08:26:00Z</updated>
	// </feed>
}
```
