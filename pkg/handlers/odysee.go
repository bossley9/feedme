package handlers

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"

	"git.sr.ht/~bossley9/feedme/pkg/api"
	"git.sr.ht/~bossley9/feedme/pkg/atom"
)

type odyseeResponse struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Title         string `xml:"title"`
		Link          string `xml:"link"`
		LastBuildDate string `xml:"lastBuildDate"`
		Description   string `xml:"description"`
		Items         []struct {
			Title       string `xml:"title"`
			Description string `xml:"description"`
			Guid        string `xml:"guid"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

func HandleOdysee(w http.ResponseWriter, r *http.Request) {
	channelID := r.FormValue("channelid")
	if len(channelID) == 0 {
		HandleUsage(w, r, "/odysee?channelid={CHANNEL_ID}")
		return
	}

	formattedUrl := "https://lbryfeed.melroy.org/channel/odysee/" + channelID
	raw, err := api.FetchGet(formattedUrl)
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	var response odyseeResponse
	if err := xml.Unmarshal(raw, &response); err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	data := response.Channel
	feedTitle := strings.TrimSuffix(data.Title, " on Odysee")

	feed, err := atom.CreateFeed(data.Link, feedTitle, time.Now())
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	feed.AddLink(data.Link, atom.RelSelf)
	feed.SetSubtitle(data.Description, "text")

	for _, item := range data.Items {
		id := item.Link
		title := strings.TrimSuffix(item.Title, " (video)")

		entry, err := atom.CreateFeedEntry(id, title, time.Now())
		if err != nil {
			continue
		}

		published := getDatetime(item.PubDate, time.RFC1123)
		entry.SetPublished(published)

		entry.SetContent(item.Description, "html")

		feed.AddEntry(entry)
	}

	HandleSuccess(w, r, feed)
}
