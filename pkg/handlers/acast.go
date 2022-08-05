package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"git.sr.ht/~bossley9/feedme/pkg/api"
	"git.sr.ht/~bossley9/feedme/pkg/atom"
)

type acastResponse struct {
	XMLName    xml.Name `xml:"rss"`
	Text       string   `xml:",chardata"`
	Version    string   `xml:"version,attr"`
	Atom       string   `xml:"atom,attr"`
	Googleplay string   `xml:"googleplay,attr"`
	Itunes     string   `xml:"itunes,attr"`
	Media      string   `xml:"media,attr"`
	Podaccess  string   `xml:"podaccess,attr"`
	Acast      string   `xml:"acast,attr"`
	Channel    struct {
		Text      string `xml:",chardata"`
		Ttl       string `xml:"ttl"`
		Generator string `xml:"generator"`
		Title     string `xml:"title"`
		Link      struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Language    string `xml:"language"`
		Copyright   string `xml:"copyright"`
		Keywords    string `xml:"keywords"`
		Author      string `xml:"author"`
		Subtitle    string `xml:"subtitle"`
		Summary     string `xml:"summary"`
		Description string `xml:"description"`
		Explicit    string `xml:"explicit"`
		Owner       struct {
			Text  string `xml:",chardata"`
			Name  string `xml:"name"`
			Email string `xml:"email"`
		} `xml:"owner"`
		ShowId    string `xml:"showId"`
		ShowUrl   string `xml:"showUrl"`
		Signature struct {
			Text      string `xml:",chardata"`
			Key       string `xml:"key,attr"`
			Algorithm string `xml:"algorithm,attr"`
		} `xml:"signature"`
		Settings string `xml:"settings"`
		Network  struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Slug string `xml:"slug,attr"`
		} `xml:"network"`
		ImportedFeed string `xml:"importedFeed"`
		Type         string `xml:"type"`
		Image        struct {
			Text  string `xml:",chardata"`
			Href  string `xml:"href,attr"`
			URL   string `xml:"url"`
			Link  string `xml:"link"`
			Title string `xml:"title"`
		} `xml:"image"`
		NewFeedURL string `xml:"new-feed-url"`
		Item       []struct {
			Text      string `xml:",chardata"`
			Title     string `xml:"title"`
			PubDate   string `xml:"pubDate"` // adheres to RFC 1123
			Duration  string `xml:"duration"`
			Enclosure struct {
				Text   string `xml:",chardata"`
				URL    string `xml:"url,attr"`
				Length string `xml:"length,attr"`
				Type   string `xml:"type,attr"`
			} `xml:"enclosure"`
			Guid struct {
				Text        string `xml:",chardata"`
				IsPermaLink string `xml:"isPermaLink,attr"`
			} `xml:"guid"`
			Explicit    string `xml:"explicit"`
			Link        string `xml:"link"`
			EpisodeId   string `xml:"episodeId"`
			EpisodeUrl  string `xml:"episodeUrl"`
			Settings    string `xml:"settings"`
			EpisodeType string `xml:"episodeType"`
			Season      string `xml:"season"`
			Episode     string `xml:"episode"`
			Image       struct {
				Text string `xml:",chardata"`
				Href string `xml:"href,attr"`
			} `xml:"image"`
			Description string `xml:"description"`
			Summary     string `xml:"summary"`
			Subtitle    string `xml:"subtitle"`
		} `xml:"item"`
		Category []struct {
			Text     string `xml:",chardata"`
			AttrText string `xml:"text,attr"`
			Category struct {
				Text     string `xml:",chardata"`
				AttrText string `xml:"text,attr"`
			} `xml:"category"`
		} `xml:"category"`
	} `xml:"channel"`
}

func handleAcastUsage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "usage: /acast?show={SHOW_ID}")
}

func HandleAcast(w http.ResponseWriter, r *http.Request) {
	showID := r.FormValue("show")
	if len(showID) == 0 {
		handleAcastUsage(w, r)
		return
	}

	url := "https://feeds.acast.com/public/shows/" + showID
	raw, err := api.FetchGet(url)
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	var response acastResponse
	if err := xml.Unmarshal(raw, &response); err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	data := response.Channel

	// Acast provides no update time so we use the current time
	feed, err := atom.CreateFeed(url, data.Title, time.Now())
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	feed.AddAuthor(data.Owner.Name, "", data.Owner.Email)
	feed.AddLink(url, atom.RelSelf)
	feed.SetSubtitle(data.Subtitle, "text")
	feed.SetCopyright(data.Copyright, "text")
	feed.SetLogo(data.Image.URL)

	for _, cat := range data.Category {
		// Categories can be (single-level) nested within categories
		category := cat.AttrText
		if len(cat.Category.AttrText) > 0 {
			category = cat.Category.AttrText
		}
		feed.AddCategory(category, "", "")
	}

	for _, item := range data.Item {
		published, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			published = time.Now()
		}

		entry, err := atom.CreateFeedEntry(item.EpisodeId, item.Title, published)
		if err != nil {
			continue
		}

		entry.AddLink(item.Enclosure.URL, atom.RelSelf)
		entry.SetPublished(published)
		entry.SetSummary(item.Description, "html")

		feed.AddEntry(entry)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, feed.String())
}
