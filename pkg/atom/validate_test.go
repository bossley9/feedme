package atom

import (
	"testing"
	"time"
)

// s3.1

func TestValidateAtomTextConstruct_MissingText(t *testing.T) {
	text := AtomTextConstruct{
		Text: "",
	}

	err := text.Validate()
	if err == nil {
		t.Error("Expected atomTextConstruct with no text attribute to throw error")
	}
}

// s3.1

func TestValidateAtomPersonConstruct_MissingName(t *testing.T) {
	person := AtomPersonConstruct{
		Name: "",
	}

	err := person.Validate()
	if err == nil {
		t.Error("Expected atomPersonConstruct with no atom:name element to throw error")
	}
}

// s4.1.1

func TestValidateAtomFeed_MissingAuthor(t *testing.T) {
	date, dateErr := time.Parse("2006-01-02 15:04", "2022-07-04 12:34")
	if dateErr != nil {
		t.Error("Error: unable parse datetime")
	}
	feed := AtomFeed{
		Id: "example.com",
		Title: AtomTitle{
			Text: "How I Generate Atom Feeds",
		},
		Updated: AtomDate(date),
	}
	err := feed.Validate()
	if err == nil {
		t.Error("Expected atom:feed with no authors to throw error (simplified)")
	}
}

func TestValidateAtomFeed_MissingId(t *testing.T) {
	date, dateErr := time.Parse("2006-01-02 15:04", "2022-07-04 12:34")
	if dateErr != nil {
		t.Error("Error: unable parse datetime")
	}
	author := AtomAuthor{
		Name: "John Doe",
	}
	feed := AtomFeed{
		Authors: []AtomAuthor{author},
		Title: AtomTitle{
			Text: "How I Generate Atom Feeds",
		},
		Updated: AtomDate(date),
	}
	err := feed.Validate()
	if err == nil {
		t.Error("Expected atom:feed with no id element to throw error")
	}
}

func TestValidateAtomFeed_NoLinks(t *testing.T) {
	date, dateErr := time.Parse("2006-01-02 15:04", "2022-07-04 12:34")
	if dateErr != nil {
		t.Error("Error: unable parse datetime")
	}
	author := AtomAuthor{
		Name: "John Doe",
	}
	feed := AtomFeed{
		Authors: []AtomAuthor{author},
		Id:      "example.com",
		Title: AtomTitle{
			Text: "How I Generate Atom Feeds",
		},
		Updated: AtomDate(date),
	}
	err := feed.Validate()
	if err == nil {
		t.Error("Expected atom:feed with no links to throw error")
	}
}

func TestValidateAtomFeed_MissingLinkRelSelf(t *testing.T) {
	date, dateErr := time.Parse("2006-01-02 15:04", "2022-07-04 12:34")
	if dateErr != nil {
		t.Error("Error: unable parse datetime")
	}
	author := AtomAuthor{
		Name: "John Doe",
	}
	links := []AtomLink{
		AtomLink{
			Href: "example.com",
		},
	}
	feed := AtomFeed{
		Authors: []AtomAuthor{author},
		Id:      "example.com",
		Links:   links,
		Title: AtomTitle{
			Text: "How I Generate Atom Feeds",
		},
		Updated: AtomDate(date),
	}
	err := feed.Validate()
	if err == nil {
		t.Error("Expected atom:feed with no link with rel=\"self\" to throw error")
	}
}

func TestValidateAtomFeed_MissingTitle(t *testing.T) {
	date, dateErr := time.Parse("2006-01-02 15:04", "2022-07-04 12:34")
	if dateErr != nil {
		t.Error("Error: unable parse datetime")
	}
	author := AtomAuthor{
		Name: "John Doe",
	}
	links := []AtomLink{
		AtomLink{
			Href: "example.com",
			Rel:  RelSelf,
		},
	}
	feed := AtomFeed{
		Authors: []AtomAuthor{author},
		Id:      "example.com",
		Links:   links,
		Updated: AtomDate(date),
	}
	err := feed.Validate()
	if err == nil {
		t.Error("Expected atom:feed with no title to throw error")
	}
}

func TestValidateAtomFeed_MissingUpdated(t *testing.T) {
	author := AtomAuthor{
		Name: "John Doe",
	}
	links := []AtomLink{
		AtomLink{
			Href: "example.com",
			Rel:  RelSelf,
		},
	}
	feed := AtomFeed{
		Authors: []AtomAuthor{author},
		Id:      "example.com",
		Links:   links,
		Title: AtomTitle{
			Text: "How I Generate Atom Feeds",
		},
	}
	err := feed.Validate()
	if err == nil {
		t.Error("Expected atom:feed with no updated date to throw error")
	}
}

func TestValidateAtomFeed_ValidFeed(t *testing.T) {
	date, dateErr := time.Parse("2006-01-02 15:04", "2022-07-04 12:34")
	if dateErr != nil {
		t.Error("Error: unable parse datetime")
	}
	author := AtomAuthor{
		Name: "John Doe",
	}
	links := []AtomLink{
		AtomLink{
			Href: "example.com",
			Rel:  RelSelf,
		},
	}
	feed := AtomFeed{
		Authors: []AtomAuthor{author},
		Id:      "example.com",
		Links:   links,
		Title: AtomTitle{
			Text: "How I Generate Atom Feeds",
		},
		Updated: AtomDate(date),
	}
	err := feed.Validate()
	if err != nil {
		t.Error(err)
	}
}

// s4.1.2

func TestValidateAtomEntry_MissingId(t *testing.T) {
	date, dateErr := time.Parse("2006-01-02 15:04", "2022-07-04 12:34")
	if dateErr != nil {
		t.Error("Error: unable parse datetime")
	}
	entry := AtomEntry{
		Title: AtomTitle{
			Text: "Entry 1",
		},
		Updated: AtomDate(date),
	}
	err := entry.Validate()
	if err == nil {
		t.Error("Expected atom:entry with no id element to throw error")
	}
}

func TestValidateAtomEntry_MissingTitle(t *testing.T) {
	date, dateErr := time.Parse("2006-01-02 15:04", "2022-07-04 12:34")
	if dateErr != nil {
		t.Error("Error: unable parse datetime")
	}
	entry := AtomEntry{
		Id:      "example.com/entry1",
		Updated: AtomDate(date),
	}
	err := entry.Validate()
	if err == nil {
		t.Error("Expected atom:entry with no title element to throw error")
	}
}

func TestValidateAtomEntry_MissingUpdated(t *testing.T) {
	entry := AtomEntry{
		Id: "example.com/entry1",
		Title: AtomTitle{
			Text: "Entry 1",
		},
	}
	err := entry.Validate()
	if err == nil {
		t.Error("Expected atom:entry with no updated element to throw error")
	}
}

func TestValidateAtomEntry_ValidEntry(t *testing.T) {
	date, dateErr := time.Parse("2006-01-02 15:04", "2022-07-04 12:34")
	if dateErr != nil {
		t.Error("Error: unable parse datetime")
	}
	entry := AtomEntry{
		Id: "example.com/entry1",
		Title: AtomTitle{
			Text: "Entry 1",
		},
		Updated: AtomDate(date),
	}
	err := entry.Validate()
	if err != nil {
		t.Error(err)
	}
}

// s4.2.2

func TestValidateAtomCategory_MissingTerm(t *testing.T) {
	cat := AtomCategory{}
	err := cat.Validate()
	if err == nil {
		t.Error("Expected atom:category to have a term attribute")
	}
}

func TestValidateAtomCategory_ValidCategory(t *testing.T) {
	cat := AtomCategory{
		Term:   "cats",
		Scheme: "example.com/tags/cats",
		Label:  "a type of pet",
	}
	err := cat.Validate()
	if err != nil {
		t.Error(err)
	}
}

// s4.2.7

func TestValidateAtomLink_MissingHref(t *testing.T) {
	link := AtomLink{}
	err := link.Validate()
	if err == nil {
		t.Error("Expected atom:link to have an href attribute")
	}
}
