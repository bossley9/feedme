package atom

import "encoding/xml"

func marshalItem(v any) []byte {
	out, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		// silently ignore errors
		return []byte{}
	}
	return out
}

func (feed AtomFeed) String() string {
	out := marshalItem(feed)
	data := []byte(xml.Header + string(out))
	return string(data)
}

func (entry AtomEntry) String() string {
	out := marshalItem(entry)
	return string(out)
}
