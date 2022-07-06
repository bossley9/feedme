package main

import (
	"fmt"
	"time"

	"git.sr.ht/~bossley9/feed-generator/pkg/atom"
)

func main() {
	author := atom.AtomAuthor{
		Name:  "John Doe",
		Uri:   "https://example.com",
		Email: "john.doe@mail.com",
	}

	cat1 := atom.AtomCategory{
		Term:   "cats",
		Scheme: "https://example2.com/scheme.html",
	}

	cat2 := atom.AtomCategory{
		Term:  "dogs",
		Label: "this is a <em>category</em> label, not a <i>random string</i>",
	}

	cats := []atom.AtomCategory{cat1, cat2}

	contrib1 := atom.AtomContributor{
		Name: "John Smith",
	}

	contribs := []atom.AtomContributor{
		contrib1,
	}

	link1 := atom.AtomLink{
		Href: "https://example3.com",
		Rel:  atom.RelAlternate,
	}

	link2 := atom.AtomLink{
		Href:  "http://my-website.org",
		Title: "<h1>How to HTML-escape titles</h1>",
	}

	link3 := atom.AtomLink{
		Title:  "This Feed Link",
		Href:   "https://www.my-website.gov/about.html",
		Rel:    atom.RelSelf,
		Length: 1,
	}

	links := []atom.AtomLink{link1, link2, link3}

	date := atom.AtomDate(time.Now())

	entry1 := atom.AtomEntry{
		Categories: []atom.AtomCategory{
			cat1,
		},
		Content: &atom.AtomContent{
			Text: "<h1>Cats vs Dogs</h1><p>Who will win in this thrilling battle?</p><p>It's still unclear</p>",
		},
		Contributors: []atom.AtomContributor{
			contrib1,
		},
		Id:      "id:my-website-cats-vs-dogs",
		Updated: date,
		Title: atom.AtomTitle{
			Text: "<h1>Cats vs Dogs</h1>",
		},
		Published: &date,
		Rights: &atom.AtomRights{
			Text: "Dog Incorporated. All rights reserved.",
		},
	}

	entries := []atom.AtomEntry{
		entry1,
	}

	feed := &atom.AtomFeed{
		Authors:      []atom.AtomAuthor{author},
		Categories:   cats,
		Contributors: contribs,
		Icon:         "fav.ico",
		Title: atom.AtomTitle{
			Text: "Everything Animals",
		},
		Updated: date,
		Id:      "id:everything-animals",
		Rights: &atom.AtomRights{
			Type: "html",
			Text: "Copyright (c) 2022 <b>Dogs Incorporated</b>.",
		},
		Links:   links,
		Logo:    "dog.jpg",
		Entries: entries,
	}

	err := feed.Validate()
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(*feed)
}
