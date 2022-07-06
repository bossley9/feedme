package atom

import (
	"fmt"
	"testing"
	"time"
)

// s4.1.1

func TestAtomFeed_CreateFeed(t *testing.T) {
	feed := makeTestFeed(t)
	ref :=
		fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  %s
  <id>example.com</id>
  <title>My Website</title>
  <updated>2022-07-04T12:34:00Z</updated>
</feed>`, getGenerator())
	test := feed.String()

	assertEqual(t, test, ref)
}

// s4.1.2

func TestAtomEntry_CreateEntry(t *testing.T) {
	entry := makeTestEntry(t)
	ref :=
		`<entry>
  <id>example.com/entry/1</id>
  <title>Entry 1</title>
  <updated>2022-07-04T12:34:00Z</updated>
</entry>`
	test := entry.String()

	assertEqual(t, test, ref)
}

func TestAtomFeed_AddEntry(t *testing.T) {
	feed := makeTestFeed(t)
	entry := makeTestEntry(t)
	ref1 :=
		fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  %s
  <id>example.com</id>
  <title>My Website</title>
  <updated>2022-07-04T12:34:00Z</updated>
</feed>`, getGenerator())
	ref2 :=
		fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  %s
  <id>example.com</id>
  <title>My Website</title>
  <updated>2022-07-04T12:34:00Z</updated>
  <entry>
    <id>example.com/entry/1</id>
    <title>Entry 1</title>
    <updated>2022-07-04T12:34:00Z</updated>
  </entry>
</feed>`, getGenerator())

	assertEqual(t, feed.String(), ref1)

	feed.AddEntry(entry)

	assertEqual(t, feed.String(), ref2)
}

// s4.1.3

func TestAtomEntry_SetContent(t *testing.T) {
	entry := makeTestEntry(t)
	ref :=
		`<entry>
  <content type="html">&lt;p&gt;This is a paragraph&lt;/p&gt;</content>
  <id>example.com/entry/1</id>
  <title>Entry 1</title>
  <updated>2022-07-04T12:34:00Z</updated>
</entry>`

	entry.SetContent("<p>This is a paragraph</p>", "html")

	test := entry.String()
	assertEqual(t, test, ref)
}

// s4.2.1

func TestAtomFeed_AddFirstAuthor(t *testing.T) {
	feed := makeTestFeed(t)
	ref1 :=
		fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  %s
  <id>example.com</id>
  <title>My Website</title>
  <updated>2022-07-04T12:34:00Z</updated>
</feed>`, getGenerator())
	ref2 :=
		fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <author>
    <name>John Doe</name>
  </author>
  %s
  <id>example.com</id>
  <title>My Website</title>
  <updated>2022-07-04T12:34:00Z</updated>
</feed>`, getGenerator())

	assertEqual(t, feed.String(), ref1)

	err := feed.AddAuthor("John Doe", "", "")
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, feed.String(), ref2)
}

func TestAtomFeed_AddDetailedAuthor(t *testing.T) {
	feed := makeTestFeed(t)
	ref :=
		fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <author>
    <name>John Doe</name>
    <uri>johndoe.com</uri>
    <email>john@motors.com</email>
  </author>
  %s
  <id>example.com</id>
  <title>My Website</title>
  <updated>2022-07-04T12:34:00Z</updated>
</feed>`, getGenerator())

	err := feed.AddAuthor("John Doe", "johndoe.com", "john@motors.com")
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, feed.String(), ref)
}

func TestAtomFeed_AddMultipleAuthors(t *testing.T) {
	feed := makeTestFeed(t)
	ref :=
		fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <author>
    <name>John Doe</name>
  </author>
  <author>
    <name>Bob Smith</name>
  </author>
  %s
  <id>example.com</id>
  <title>My Website</title>
  <updated>2022-07-04T12:34:00Z</updated>
</feed>`, getGenerator())

	err1 := feed.AddAuthor("John Doe", "", "")
	if err1 != nil {
		t.Error(err1)
	}
	err2 := feed.AddAuthor("Bob Smith", "", "")
	if err2 != nil {
		t.Error(err2)
	}

	assertEqual(t, feed.String(), ref)
}

// s4.2.7

func TestAtomFeed_AddFirstLink(t *testing.T) {
	feed := makeTestFeed(t)
	ref1 :=
		fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  %s
  <id>example.com</id>
  <title>My Website</title>
  <updated>2022-07-04T12:34:00Z</updated>
</feed>`, getGenerator())
	ref2 :=
		fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  %s
  <id>example.com</id>
  <link href="example.com" rel="self"></link>
  <title>My Website</title>
  <updated>2022-07-04T12:34:00Z</updated>
</feed>`, getGenerator())

	assertEqual(t, feed.String(), ref1)

	err := feed.AddLink("example.com", RelSelf)
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, feed.String(), ref2)
}

func TestAtomFeed_AddMultipleLinks(t *testing.T) {
	feed := makeTestFeed(t)
	err1 := feed.AddLink("example.com", RelSelf)
	if err1 != nil {
		t.Error(err1)
	}
	err2 := feed.AddLink("alt-example.org", RelAlternate)
	if err2 != nil {
		t.Error(err2)
	}
	err3 := feed.AddLink("resource.gov", RelRelated)
	if err3 != nil {
		t.Error(err3)
	}

	ref :=
		fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  %s
  <id>example.com</id>
  <link href="example.com" rel="self"></link>
  <link href="alt-example.org" rel="alternate"></link>
  <link href="resource.gov" rel="related"></link>
  <title>My Website</title>
  <updated>2022-07-04T12:34:00Z</updated>
</feed>`, getGenerator())

	assertEqual(t, feed.String(), ref)
}

func TestAtomEntry_AddFirstLink(t *testing.T) {
	entry := makeTestEntry(t)
	ref1 :=
		`<entry>
  <id>example.com/entry/1</id>
  <title>Entry 1</title>
  <updated>2022-07-04T12:34:00Z</updated>
</entry>`
	ref2 :=
		`<entry>
  <id>example.com/entry/1</id>
  <link href="example.com" rel="self"></link>
  <title>Entry 1</title>
  <updated>2022-07-04T12:34:00Z</updated>
</entry>`

	assertEqual(t, entry.String(), ref1)

	err := entry.AddLink("example.com", RelSelf)
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, entry.String(), ref2)
}

// s4.2.9

func TestAtomEntry_SetPublished(t *testing.T) {
	entry := makeTestEntry(t)
	ref :=
		`<entry>
  <id>example.com/entry/1</id>
  <published>2021-01-01T12:34:00-02:00</published>
  <title>Entry 1</title>
  <updated>2022-07-04T12:34:00Z</updated>
</entry>`

	date, errDate := time.Parse("2006-01-02 15:04 -0700", "2021-01-01 12:34 -0200")
	if errDate != nil {
		t.Error(errDate)
	}

	err := entry.SetPublished(date)
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, entry.String(), ref)
}

// s4.2.12

func TestAtomFeed_SetSubtitle(t *testing.T) {
	feed := makeTestFeed(t)
	ref :=
		fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  %s
  <id>example.com</id>
  <subtitle type="html">&lt;em&gt;Formatted&lt;/em&gt; subtitle</subtitle>
  <title>My Website</title>
  <updated>2022-07-04T12:34:00Z</updated>
</feed>`, getGenerator())

	err := feed.SetSubtitle("<em>Formatted</em> subtitle", "html")
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, feed.String(), ref)
}
