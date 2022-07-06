package atom

import (
	// "fmt"
	"errors"
	"time"
)

// s4.1.1

func CreateFeed(id string, title string, updated time.Time) (*AtomFeed, error) {
	if len(id) == 0 {
		return nil, errors.New("atom:id cannot be empty")
	}
	if len(title) == 0 {
		return nil, errors.New("atom:title cannot be empty")
	}
	if updated.IsZero() {
		return nil, errors.New("atom:updated cannot be empty")
	}

	feed := AtomFeed{
		Id: AtomID(id),
		Title: AtomTitle{
			Text: title,
		},
		Updated: AtomDate(updated),
	}

	return &feed, nil
}

// s4.1.2

func CreateFeedEntry(id string, title string, updated time.Time) (*AtomEntry, error) {
	if len(id) == 0 {
		return nil, errors.New("atom:id cannot be empty")
	}
	if len(title) == 0 {
		return nil, errors.New("atom:title cannot be empty")
	}
	if updated.IsZero() {
		return nil, errors.New("atom:updated cannot be empty")
	}

	entry := AtomEntry{
		Id: AtomID(id),
		Title: AtomTitle{
			Text: title,
		},
		Updated: AtomDate(updated),
	}

	return &entry, nil
}

func (feed *AtomFeed) AddEntry(entry *AtomEntry) error {
	feed.Entries = append(feed.Entries, *entry)
	return nil
}

// s4.1.3

func (entry *AtomEntry) SetContent(content string, contentType string) error {
	entry.Content = &AtomContent{
		Text: content,
		Type: contentType,
	}
	return nil
}

// s4.2.1

func (feed *AtomFeed) AddAuthor(name string, uri string, email string) error {
	author := AtomAuthor{
		Name:  AtomName(name),
		Uri:   AtomURI(uri),
		Email: AtomEmail(email),
	}
	feed.Authors = append(feed.Authors, author)

	return nil
}

// s4.2.7

func (feed *AtomFeed) AddLink(href string, rel AtomRelType) error {
	link := AtomLink{
		Href: AtomURI(href),
		Rel:  rel,
	}
	feed.Links = append(feed.Links, link)

	return nil
}

func (entry *AtomEntry) AddLink(href string, rel AtomRelType) error {
	link := AtomLink{
		Href: AtomURI(href),
		Rel:  rel,
	}
	entry.Links = append(entry.Links, link)

	return nil
}

// s4.2.9

func (entry *AtomEntry) SetPublished(published time.Time) error {
	date := AtomDate(published)
	entry.Published = &date
	return nil
}

// s4.2.12

func (feed *AtomFeed) SetSubtitle(text string, textType string) error {
	subtitle := AtomSubtitle{
		Text: text,
		Type: AtomTextType(textType),
	}
	feed.Subtitle = &subtitle

	return nil
}
