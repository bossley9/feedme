package atom

import (
	"errors"
	"fmt"
	"time"
)

func FeedError(message string, rfcSection string) error {
	str := message

	if len(rfcSection) > 0 {
		str = str + fmt.Sprintf(". See https://datatracker.ietf.org/doc/html/rfc4287#section-%s for more details.", rfcSection)
	}

	return errors.New(str)
}

// s3.1

func (text AtomTextConstruct) Validate() error {
	if len(text.Text) == 0 {
		return FeedError("text constructs must have textual content", "3.1")
	}
	return nil
}

// s3.2

func (person AtomPersonConstruct) Validate() error {
	if len(person.Name) == 0 {
		return FeedError("person constructs must contain exactly one atom:name element", "3.2.1")
	}
	return nil
}

// s4.1.1

func (feed AtomFeed) Validate() error {
	if len(feed.Authors) == 0 {
		return FeedError("atom:feed elements must contain at least one author (simplified)", "4.1.1")
	}
	for _, author := range feed.Authors {
		err := AtomPersonConstruct(author).Validate()
		if err != nil {
			return err
		}
	}

	for _, cat := range feed.Categories {
		err := cat.Validate()
		if err != nil {
			return err
		}
	}

	for _, contrib := range feed.Contributors {
		err := AtomPersonConstruct(contrib).Validate()
		if err != nil {
			return err
		}
	}

	if len(feed.Id) == 0 {
		return FeedError("atom:feed elements must contain exactly one id", "4.1.1")
	}

	// ensure at least one link exists with rel="self"
	hasSelfLink := false
	for _, link := range feed.Links {
		err := link.Validate()
		if err != nil {
			return err
		}

		if link.Rel == RelSelf {
			if hasSelfLink {
				return FeedError("atom:feed elements must not contain more than one link with rel=\"self\" (simplified)", "4.1.1")
			}
			hasSelfLink = true
		}
	}
	if !hasSelfLink {
		return FeedError("atom:feed elements should contain one link with rel=\"self\" (simplified)", "4.1.1")
	}

	if len(feed.Title.Text) == 0 {
		return FeedError("atom:feed elements must contain exactly one title element", "4.1.1")
	}

	if time.Time(feed.Updated).IsZero() {
		return FeedError("atom:feed elements must contain exactly one valid updated element", "4.1.1")
	}

	for _, entry := range feed.Entries {
		err := entry.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// s4.1.2

func (entry AtomEntry) Validate() error {
	for _, author := range entry.Authors {
		err := AtomPersonConstruct(author).Validate()
		if err != nil {
			return err
		}
	}

	for _, cat := range entry.Categories {
		err := cat.Validate()
		if err != nil {
			return err
		}
	}

	for _, contrib := range entry.Contributors {
		err := AtomPersonConstruct(contrib).Validate()
		if err != nil {
			return err
		}
	}

	if len(entry.Id) == 0 {
		return FeedError("atom:entry elements must contain exactly one id", "4.1.2")
	}

	if len(entry.Title.Text) == 0 {
		return FeedError("atom:feed elements must contain exactly one title element", "4.1.1")
	}

	if time.Time(entry.Updated).IsZero() {
		return FeedError("atom:feed elements must contain exactly one valid updated element", "4.1.1")
	}

	return nil
}

// s4.2.2

func (category AtomCategory) Validate() error {
	if len(category.Term) == 0 {
		return FeedError("atom:category elements must have a term attribute", "4.2.2.1")
	}
	return nil
}

// s4.2.7

func (link AtomLink) Validate() error {
	if len(link.Href) == 0 {
		return FeedError("atom:link elements must have an href attribute", "4.2.7.1")
	}
	return nil
}
