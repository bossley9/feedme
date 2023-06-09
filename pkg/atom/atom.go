// see RFC 4287

package atom

import (
	"encoding/xml"
	"time"
)

// s3.1

type AtomTextConstruct struct {
	Type AtomTextType `xml:"type,attr,omitempty"`
	Text string       `xml:",chardata"` // required
}

// s3.1.1

type AtomTextType string // one of "text", "html", "xhtml" - default "text"

// s3.2

type AtomPersonConstruct struct {
	Name  AtomName  `xml:"name"` // required
	Uri   AtomURI   `xml:"uri,omitempty"`
	Email AtomEmail `xml:"email,omitempty"`
}

// s3.2.1

type AtomName string

// s3.2.2

type AtomURI string // IRI reference as defined in RFC 3987

// s3.2.3

type AtomEmail string // email address as defined in RFC 2822

// s3.3

type AtomDate time.Time // datetime reference as defined in RFC 3339

func (date AtomDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	str := time.Time(date).Format(time.RFC3339)
	e.EncodeElement(str, start)
	return nil
}

// s4.1.1

type AtomFeed struct {
	XMLName      xml.Name          `xml:"http://www.w3.org/2005/Atom feed"`
	Authors      []AtomAuthor      `xml:"author"` // >1 required
	Categories   []AtomCategory    `xml:"category"`
	Contributors []AtomContributor `xml:"contributor"`
	Generator    AtomGenerator     `xml:",omitempty"`
	Icon         AtomIcon          `xml:"icon,omitempty"` // should use aspect ratio 1(h):1(v)
	Id           AtomID            `xml:"id"`             // required
	Links        []AtomLink        `xml:"link"`
	Logo         AtomLogo          `xml:"logo,omitempty"` // should use aspect ratio 2(h):1(v)
	Rights       *AtomRights       `xml:"rights,omitempty"`
	Subtitle     *AtomSubtitle     `xml:"subtitle,omitempty"`
	Title        AtomTitle         `xml:"title"`   // required
	Updated      AtomDate          `xml:"updated"` // required
	Entries      []AtomEntry       `xml:"entry"`
}

// s4.1.2

type AtomEntry struct {
	XMLName      xml.Name          `xml:"entry"`
	Authors      []AtomAuthor      `xml:"author"`
	Categories   []AtomCategory    `xml:"category"`
	Content      *AtomContent      `xml:"content"`
	Contributors []AtomContributor `xml:"contributor"`
	Id           AtomID            `xml:"id"` // required
	Links        []AtomLink        `xml:"link"`
	Published    *AtomDate         `xml:"published,omitempty"`
	Rights       *AtomRights       `xml:"rights,omitempty"`
	Source       *AtomSource       `xml:"source,omitempty"`
	Summary      *AtomSummary      `xml:"summary,omitempty"`
	Title        AtomTitle         `xml:"title"`   // required
	Updated      AtomDate          `xml:"updated"` // required
}

// s4.1.2

type AtomContent struct {
	Type string  `xml:"type,attr,omitempty"` // one of "text", "html", "xhtml", or a MIME media type - default "text"
	Src  AtomURI `xml:"src,attr,omitempty"`
	Text string  `xml:",chardata"`
}

// s4.2.1

type AtomAuthor AtomPersonConstruct

// s4.2.2

// self-closing element
type AtomCategory struct {
	XMLName xml.Name `xml:"category"`
	Term    string   `xml:"term,attr"` // required
	Scheme  AtomURI  `xml:"scheme,attr,omitempty"`
	Label   string   `xml:"label,attr,omitempty"` // HTML-escaped
}

// s4.2.3

type AtomContributor AtomPersonConstruct

// s4.2.4

type AtomGenerator struct {
	Uri     AtomURI `xml:"uri,attr"`
	Version string  `xml:"version,attr"`
	Text    string  `xml:",innerxml"`
}

func (generator AtomGenerator) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	uri := xml.Attr{
		Name:  xml.Name{Local: "uri"},
		Value: PKG,
	}

	version := xml.Attr{
		Name:  xml.Name{Local: "version"},
		Value: VERSION,
	}

	element := xml.StartElement{
		Name: xml.Name{Local: "generator"},
		Attr: []xml.Attr{uri, version},
	}

	e.EncodeElement("Generated via "+NAME, element)

	return nil
}

// s4.2.5

type AtomIcon AtomURI

// s4.2.6

// atom:id elements must ALWAYS assure uniqueness.
// See https://datatracker.ietf.org/doc/html/rfc4287#section-4.2.6 for
// strategies on generating unique id elements.
type AtomID AtomURI

// s4.2.7

// self-closing element
type AtomLink struct {
	XMLName  xml.Name        `xml:"link"`
	Href     AtomURI         `xml:"href,attr"` // required
	Rel      AtomRelType     `xml:"rel,attr,omitempty"`
	Type     AtomMediaType   `xml:"type,attr,omitempty"`
	HrefLang AtomLanguageTag `xml:"hreflang,attr,omitempty"`
	Title    string          `xml:"title,attr,omitempty"` // HTML-escaped
	Length   uint            `xml:"length,attr,omitempty"`
}

// s4.2.7.2

type AtomRelType int

const (
	RelUnknown AtomRelType = iota
	RelAlternate
	RelRelated
	RelSelf
	RelEnclosure
	RelVia
)

func (rel AtomRelType) MarshalText() ([]byte, error) {
	var str string
	switch rel {
	case RelAlternate:
		str = "alternate"
	case RelRelated:
		str = "related"
	case RelSelf:
		str = "self"
	case RelEnclosure:
		str = "enclosure"
	case RelVia:
		str = "via"
	default:
		str = ""
	}
	return []byte(str), nil
}

// s4.2.7.3

type AtomMediaType string // MIME media type

// s4.2.7.4

type AtomLanguageTag string // language type as defined in RFC 3066

// s4.2.8

type AtomLogo AtomURI

// s4.2.10

type AtomRights AtomTextConstruct

// s4.2.11

// unimplemented due to language constraints
type AtomSource struct{} // extension of atom:feed without entries

// s4.2.12

type AtomSubtitle AtomTextConstruct

// s4.2.13

type AtomSummary AtomTextConstruct

// s4.2.14

type AtomTitle AtomTextConstruct
