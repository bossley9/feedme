package atom

import (
	"fmt"
	"testing"
	"time"
)

func getGenerator() string {
	return fmt.Sprintf(`<generator uri="%s" version="%s">Generated via %s</generator>`, PKG, VERSION, NAME)
}

func makeTestFeed(t *testing.T) *AtomFeed {
	date, errDate := time.Parse("2006-01-02 15:04", "2022-07-04 12:34")
	if errDate != nil {
		t.Error(errDate)
	}
	feed, errCreateFeed := CreateFeed("example.com", "My Website", date)
	if errCreateFeed != nil {
		t.Error(errCreateFeed)
	}

	return feed
}

func makeTestEntry(t *testing.T) *AtomEntry {
	date, errDate := time.Parse("2006-01-02 15:04", "2022-07-04 12:34")
	if errDate != nil {
		t.Error(errDate)
	}
	entry, errCreateEntry := CreateFeedEntry("example.com/entry/1", "Entry 1", date)
	if errCreateEntry != nil {
		t.Error(errCreateEntry)
	}
	return entry
}
