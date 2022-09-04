package handlers

import (
	"encoding/xml"
	"net/http"
	"strings"
	"sync"
	"time"

	"git.sr.ht/~bossley9/feedme/pkg/api"
	"git.sr.ht/~bossley9/feedme/pkg/atom"
)

type tkdodoResponseItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Guid        string `xml:"guid"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
}

type tkdodoResponse struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Title         string               `xml:"title"`
		Link          string               `xml:"link"`
		LastBuildDate string               `xml:"lastBuildDate"`
		Items         []tkdodoResponseItem `xml:"item"`
	} `xml:"channel"`
}

func addTkdodoEntry(feed *atom.AtomFeed, feedUrl string, item tkdodoResponseItem, wg *sync.WaitGroup) {
	defer wg.Done()

	entryUrl := item.Link

	entry, err := atom.CreateFeedEntry(item.Guid, item.Title, time.Now())
	if err != nil {
		return
	}

	entry.AddLink(entryUrl, atom.RelAlternate)
	published := getDatetime(item.PubDate, time.RFC1123)
	entry.SetPublished(published)

	doc, err := api.FetchHTML(entryUrl)
	if err != nil {
		feed.AddEntry(entry)
		return
	}

	contentEl := doc.Find("main section")

	// remove splash image
	first := contentEl.Children().First()
	if first.HasClass("gatsby-resp-image-wrapper") {
		first.Remove()
	}

	contentHtml, err := contentEl.Html()
	if err != nil {
		feed.AddEntry(entry)
		return
	}

	// convert relative links to absolute
	contentHtml = strings.ReplaceAll(contentHtml, `href="/`, `href="`+feedUrl+"/")

	content := strings.TrimSpace(contentHtml)

	if len(content) > 0 {
		entry.SetContent(content, "html")
	} else {
		// should ideally never occur
		entry.SetSummary(item.Description, "html")
	}
	feed.AddEntry(entry)
}

func HandleTkdodo(w http.ResponseWriter, r *http.Request) {
	formattedUrl := "https://tkdodo.eu/blog/rss.xml"
	raw, err := api.FetchGet(formattedUrl)
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	var response tkdodoResponse
	if err := xml.Unmarshal(raw, &response); err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	data := response.Channel
	feedUrl := data.Link
	feedTitle := strings.TrimSuffix(data.Title, "'s blog")
	updated := getDatetime(data.LastBuildDate, time.RFC1123)

	feed, err := atom.CreateFeed(feedUrl, feedTitle, updated)
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	feed.AddLink(feedUrl, atom.RelSelf)
	feed.AddAuthor("Dominik Dorfmeister", "", "")

	var wg sync.WaitGroup
	for _, item := range data.Items {
		wg.Add(1)
		go addTkdodoEntry(feed, feedUrl, item, &wg)
	}
	wg.Wait()

	HandleSuccess(w, r, feed)
}
