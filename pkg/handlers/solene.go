package handlers

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"

	"git.sr.ht/~bossley9/feedme/pkg/api"
	"git.sr.ht/~bossley9/feedme/pkg/atom"

	"git.sr.ht/~bossley9/gem"
)

type soleneResponse struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Title string `xml:"title"`
		Items []struct {
			Title       string `xml:"title"`
			Description string `xml:"description"`
			Guid        string `xml:"guid"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

func HandleSolene(w http.ResponseWriter, r *http.Request) {
	url := "https://dataswamp.org/~solene/rss.xml"
	raw, err := api.FetchGet(url)
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	var response soleneResponse
	if err := xml.Unmarshal(raw, &response); err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	data := response.Channel

	feed, err := atom.CreateFeed(url, data.Title, time.Now())
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	feed.AddLink(url, atom.RelSelf)

	for _, item := range data.Items {
		published, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			published = time.Now()
		}

		entry, err := atom.CreateFeedEntry(item.Guid, item.Title, published)
		if err != nil {
			continue
		}

		entry.AddLink(item.Link, atom.RelAlternate)
		entry.SetPublished(published)

		spaceTrimmedContent := strings.TrimSpace(item.Description)
		geminiContent := strings.TrimSuffix(strings.TrimPrefix(spaceTrimmedContent, "<pre>"), "</pre>")
		content := gem.ToHTML(geminiContent)

		entry.SetContent(content, "html")

		feed.AddEntry(entry)
	}

	HandleSuccess(w, r, feed)
}
