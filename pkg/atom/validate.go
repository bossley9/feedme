package atom

import (
	"errors"
	"fmt"
	"time"
)

func feedError(message string, rfcSection string) error {
	str := message

	if len(rfcSection) > 0 {
		str = str + fmt.Sprintf(". See https://datatracker.ietf.org/doc/html/rfc4287#section-%s for more details.", rfcSection)
	}

	return errors.New(str)
}

func (feed AtomFeed) Validate() error {
	if len(feed.Id) == 0 {
		return feedError("atom:feed elements must contain exactly one id", "4.1.1")
	}

	if time.Time(feed.Updated).IsZero() {
		return feedError("atom:feed elements must contain exactly one valid updated element", "4.1.1")
	}

	if len(feed.Title.Text) == 0 {
		return feedError("atom:feed elements must contain exactly one title element", "4.1.1")
	}

	if len(feed.Authors) == 0 {
		return feedError("atom:feed elements must contain at least one author (simplified)", "4.1.1")
	}
	for _, author := range feed.Authors {
		err := AtomPerson(author).Validate()
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
		err := AtomPerson(contrib).Validate()
		if err != nil {
			return err
		}
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
				return feedError("atom:feed elements must not contain more than one link with rel=\"self\" (simplified)", "4.1.1")
			}
			hasSelfLink = true
		}
	}
	if !hasSelfLink {
		return feedError("atom:feed elements should contain one link with rel=\"self\" (simplified)", "4.1.1")
	}

	for _, entry := range feed.Entries {
		err := entry.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (entry AtomEntry) Validate() error {
	if len(entry.Id) == 0 {
		return feedError("atom:entry elements must contain exactly one id", "4.1.2")
	}

	if time.Time(entry.Updated).IsZero() {
		return feedError("atom:feed elements must contain exactly one valid updated element", "4.1.1")
	}

	if len(entry.Title.Text) == 0 {
		return feedError("atom:feed elements must contain exactly one title element", "4.1.1")
	}

	for _, cat := range entry.Categories {
		err := cat.Validate()
		if err != nil {
			return err
		}
	}

	for _, contrib := range entry.Contributors {
		err := AtomPerson(contrib).Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (link AtomLink) Validate() error {
	if len(link.Href) == 0 {
		return feedError("atom:link elements must have an href attribute", "4.2.7.1")
	}
	return nil
}

func (person AtomPerson) Validate() error {
	if len(person.Name) == 0 {
		return feedError("atom:person constructs must contain a name element", "3.2.1")
	}
	return nil
}

func (category AtomCategory) Validate() error {
	if len(category.Term) == 0 {
		return feedError("atom:category elements must have a term attribute", "4.2.2.1")
	}
	return nil
}
