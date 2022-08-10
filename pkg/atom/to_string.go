package atom

import "encoding/xml"

func (feed AtomFeed) String() string {
	// silently ignore errors
	out, _ := xml.MarshalIndent(feed, "", "  ")
	data := []byte(xml.Header + string(out))
	return string(data)
}

func (entry AtomEntry) String() string {
	// silently ignore errors
	out, _ := xml.MarshalIndent(entry, "", "  ")
	return string(out)
}
