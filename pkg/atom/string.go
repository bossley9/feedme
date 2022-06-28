package atom

import (
	"encoding/xml"
	"fmt"
)

func (feed AtomFeed) String() string {
	out, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		fmt.Println("Error: unable to marshal feed")
		return ""
	}

	data := []byte(xml.Header + string(out))
	return string(data)
}
