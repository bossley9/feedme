package handlers

import (
	"encoding/xml"
	"html"
	"net/http"
	"strings"
	"sync"
	"time"

	"git.sr.ht/~bossley9/feedme/pkg/api"
	"git.sr.ht/~bossley9/feedme/pkg/atom"

	"github.com/PuerkitoBio/goquery"
)

type jeffGeerlingResponseItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Guid        string `xml:"guid"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
}

type jeffGeerlingResponse struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Title string                     `xml:"title"`
		Link  string                     `xml:"link"`
		Items []jeffGeerlingResponseItem `xml:"item"`
	} `xml:"channel"`
}

func removeAllClasses(s *goquery.Selection) {
	s.RemoveClass()
	s.Children().Each(func(i int, childSel *goquery.Selection) {
		removeAllClasses(childSel)
	})
}

func addJeffGeerlingEntry(feed *atom.AtomFeed, feedUrl string, item jeffGeerlingResponseItem, wg *sync.WaitGroup) {
	defer wg.Done()

	entryUrl := item.Link
	entryTitle := item.Title
	entryUpdated := time.Now()
	entryID := "guid:" + strings.ReplaceAll(item.Guid, " ", ":")

	entry, err := atom.CreateFeedEntry(entryID, entryTitle, entryUpdated)
	if err != nil {
		return
	}

	entry.AddLink(entryUrl, atom.RelAlternate)
	published := getDatetime(item.PubDate, time.RFC1123Z)
	entry.SetPublished(published)
	entry.SetSummary(item.Description, "html")

	doc, err := api.FetchHTML(entryUrl)
	if err != nil {
		feed.AddEntry(entry)
		return
	}

	contentEl := doc.Find("article .node__content")

	// tidy as much as possible
	contentEl.Children().Each(func(i int, s *goquery.Selection) {
		if s.HasClass("node__links") {
			s.Remove()
			return
		}
	})
	removeAllClasses(contentEl)

	contentHtml, err := contentEl.Html()
	if err != nil {
		feed.AddEntry(entry)
		return
	}

	// convert relative links to absolute
	contentHtml = strings.ReplaceAll(contentHtml, `href="/`, `href="`+feedUrl+"/")

	trimmedContentHtml := strings.TrimSpace(contentHtml)

	content := html.EscapeString(trimmedContentHtml)

	if len(content) > 0 {
		entry.SetContent(content, "html")
		// ignore summary if content exists to reduce throughput
		entry.SetSummary("", "html")
	}
	feed.AddEntry(entry)
}

func HandleJeffGeerling(w http.ResponseWriter, r *http.Request) {
	formattedUrl := "https://www.jeffgeerling.com/blog.xml"
	raw, err := api.FetchGet(formattedUrl)
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	var response jeffGeerlingResponse
	if err := xml.Unmarshal(raw, &response); err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	data := response.Channel
	feedUrl := strings.TrimSuffix(data.Link, "/")

	feed, err := atom.CreateFeed(feedUrl, data.Title, time.Now())
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	feed.AddLink(feedUrl, atom.RelSelf)
	feed.AddAuthor("Jeff Geerling", "", "")

	var wg sync.WaitGroup
	for _, item := range data.Items {
		wg.Add(1)
		go addJeffGeerlingEntry(feed, feedUrl, item, &wg)
	}
	wg.Wait()

	HandleSuccess(w, r, feed)
}