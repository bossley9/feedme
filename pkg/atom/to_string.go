package atom

import "encoding/xml"

func (feed AtomFeed) String() (string, error) {
	out, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return "", err
	}

	data := []byte(xml.Header + string(out))
	return string(data), nil
}
