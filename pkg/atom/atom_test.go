package atom

import (
	"encoding/xml"
	"fmt"
	"testing"
	"time"
)

// s3.1

func TestAtomTextConstruct_MarshalText(t *testing.T) {
	text := AtomTextConstruct{
		Text: "This is xml data. <p> tags should be escaped.",
	}
	ref := "<AtomTextConstruct>This is xml data. &lt;p&gt; tags should be escaped.</AtomTextConstruct>"

	out, err := xml.MarshalIndent(text, "", "  ")
	if err != nil {
		t.Error("Error: unable to marshal xml")
	}
	test := string(out)

	assertEqual(t, test, ref)
}

func TestAtomTextConstruct_MarshalType(t *testing.T) {
	text := AtomTextConstruct{
		Type: "html",
		Text: "Hello world",
	}
	ref := "<AtomTextConstruct type=\"html\">Hello world</AtomTextConstruct>"

	out, err := xml.MarshalIndent(text, "", "  ")
	if err != nil {
		t.Error("Error: unable to marshal xml")
	}
	test := string(out)

	assertEqual(t, test, ref)
}

// s3.2

func TestAtomPersonConstruct_MarshalElements(t *testing.T) {
	person := AtomPersonConstruct{
		Name:  "John Doe",
		Uri:   "johndoe.com",
		Email: "me@johndoe.com",
	}
	ref :=
		`<AtomPersonConstruct>
  <name>John Doe</name>
  <uri>johndoe.com</uri>
  <email>me@johndoe.com</email>
</AtomPersonConstruct>`

	out, err := xml.MarshalIndent(person, "", "  ")
	if err != nil {
		t.Error("Error: unable to marshal xml")
	}
	test := string(out)

	assertEqual(t, test, ref)
}

func TestAtomPersonConstruct_HideOptionalElements(t *testing.T) {
	person := AtomPersonConstruct{
		Name: "John Doe",
	}
	ref :=
		`<AtomPersonConstruct>
  <name>John Doe</name>
</AtomPersonConstruct>`

	out, err := xml.MarshalIndent(person, "", "  ")
	if err != nil {
		t.Error("Error: unable to marshal xml")
	}
	test := string(out)

	assertEqual(t, test, ref)
}

// s3.3

func TestAtomDateConstruct_Format(t *testing.T) {
	date, err := time.Parse("2006-01-02 15:04", "2022-07-04 12:34")
	if err != nil {
		t.Error("Error: unable parse datetime")
	}
	ref := "<Time>2022-07-04T12:34:00Z</Time>"

	out, err := xml.MarshalIndent(date, "", "  ")
	if err != nil {
		t.Error("Error: unable to marshal xml")
	}
	test := string(out)

	assertEqual(t, test, ref)
}

// s4.1.1

func TestAtomFeed_Format(t *testing.T) {
	date, err := time.Parse("2006-01-02 15:04", "2022-07-04 12:34")
	if err != nil {
		t.Error("Error: unable parse datetime")
	}
	feed := AtomFeed{
		Id: "example.com",
		Title: AtomTitle{
			Text: "How I Generate Atom Feeds",
		},
		Updated: AtomDate(date),
	}
	ref :=
		fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  %s
  <id>example.com</id>
  <title>How I Generate Atom Feeds</title>
  <updated>2022-07-04T12:34:00Z</updated>
</feed>`, getGenerator())

	test := feed.String()

	assertEqual(t, test, ref)
}

// s4.1.2

func TestAtomEntry_Format(t *testing.T) {
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
	ref :=
		`<entry>
  <id>example.com/entry1</id>
  <title>Entry 1</title>
  <updated>2022-07-04T12:34:00Z</updated>
</entry>`

	out, err := xml.MarshalIndent(entry, "", "  ")
	if err != nil {
		t.Error("Error: unable to marshal xml")
	}
	test := string(out)

	assertEqual(t, test, ref)
}

// s4.2.2

func TestAtomCategory_Format(t *testing.T) {
	cat := AtomCategory{
		Term: "cats",
	}
	ref := `<category term="cats"></category>`

	out, err := xml.MarshalIndent(cat, "", "  ")
	if err != nil {
		t.Error("Error: unable to marshal xml")
	}
	test := string(out)

	assertEqual(t, test, ref)
}

func TestAtomCategory_EscapeLabel(t *testing.T) {
	cat := AtomCategory{
		Term:  "cats",
		Label: "<em>only</em> pets",
	}
	ref := `<category term="cats" label="&lt;em&gt;only&lt;/em&gt; pets"></category>`

	out, err := xml.MarshalIndent(cat, "", "  ")
	if err != nil {
		t.Error("Error: unable to marshal xml")
	}
	test := string(out)

	assertEqual(t, test, ref)
}

// s4.2.7

func TestAtomLink_Format(t *testing.T) {
	link := AtomLink{
		Href:     "example.com",
		Rel:      RelAlternate,
		Type:     "text/html",
		HrefLang: "en",
		Title:    "My website link",
		Length:   4,
	}
	ref := `<link href="example.com" rel="alternate" type="text/html" hreflang="en" title="My website link" length="4"></link>`

	out, err := xml.MarshalIndent(link, "", "  ")
	if err != nil {
		t.Error("Error: unable to marshal xml")
	}
	test := string(out)

	assertEqual(t, test, ref)
}

func TestAtomLink_EscapeTitle(t *testing.T) {
	link := AtomLink{
		Href:  "example.com",
		Title: "<em>This should be escaped</em>",
	}
	ref := `<link href="example.com" title="&lt;em&gt;This should be escaped&lt;/em&gt;"></link>`

	out, err := xml.MarshalIndent(link, "", "  ")
	if err != nil {
		t.Error("Error: unable to marshal xml")
	}
	test := string(out)

	assertEqual(t, test, ref)
}
