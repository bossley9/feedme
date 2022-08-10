package handlers

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"git.sr.ht/~bossley9/feedme/pkg/api"
	"git.sr.ht/~bossley9/feedme/pkg/atom"

	"git.sr.ht/~bossley9/gem"
)

const geminiProtocol = "gemini://"
const iSO8601 = "2006-01-02"

const (
	titleState = iota
	subtitleState
	entryState
)

func parseGemlogEntry(feed *atom.AtomFeed, feedUrl string, line string, wg *sync.WaitGroup) {
	defer wg.Done()

	trimmedLine := strings.TrimSpace(strings.TrimPrefix(line, "=>"))
	lineSections := strings.Split(trimmedLine, " ")

	entryUrl := lineSections[0]

	// according to the specification, links with url "entry.gmi"
	// under domain "example.com/sub/path" should resolve as
	// "example.com/sub/path/entry.gmi", not "example.com/entry.gmi".
	baseUrl, _ := url.Parse(feedUrl)
	if !strings.HasPrefix(entryUrl, geminiProtocol) {
		entryUrl = baseUrl.String() + "/" + entryUrl
	}

	updatedText := lineSections[1]
	updated, err := time.Parse(iSO8601, updatedText)
	if err != nil {
		return
	}

	title := strings.Join(lineSections[2:], " ")
	// sanitization as suggested by
	// gemini.circumlunar.space/docs/companion/subscription.gmi
	title = strings.TrimPrefix(title, "- ")

	entry, err := atom.CreateFeedEntry(entryUrl, title, updated)
	if err != nil {
		return
	}

	res, err := api.FetchGemini(entryUrl)
	if err != nil {
		// add alt link entry if unable to fetch
		entry.AddLink(entryUrl, atom.RelAlternate)
		feed.AddEntry(entry)
		return
	}

	entry.AddLink(entryUrl, atom.RelAlternate)

	content := gem.ToHTML(string(res))
	entry.SetContent(content, "html")

	feed.AddEntry(entry)
}

func HandleGemini(w http.ResponseWriter, r *http.Request) {
	encodedUrl := r.FormValue("url")
	if len(encodedUrl) == 0 {
		HandleUsage(w, r, "/gemini?url={ENCODED_URL_WITH_NO_PROTOCOL}")
		return
	}

	decodedUrl, err := url.QueryUnescape(encodedUrl)
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}
	formattedUrl := geminiProtocol + decodedUrl

	res, err := api.FetchGemini(formattedUrl)
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	feed, err := atom.CreateFeed(formattedUrl, "gemlog", time.Now())
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	feed.AddLink(formattedUrl, atom.RelSelf)

	entryMatcher, err := regexp.CompilePOSIX("^=> .* ....-..-..")
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	var wg sync.WaitGroup
	state := titleState
	for _, line := range strings.Split(string(res), "\n") {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		if state == titleState && strings.HasPrefix(line, "#") {
			title := strings.TrimSpace(strings.TrimPrefix(line, "#"))
			feed.SetTitle(title, "text")
			state = subtitleState
		} else if state == subtitleState {
			if strings.HasPrefix(line, "##") {
				subtitle := strings.TrimSpace(strings.TrimPrefix(line, "##"))
				feed.SetSubtitle(subtitle, "text")
			}
			state = entryState
		}

		if state == entryState && entryMatcher.MatchString(line) {
			wg.Add(1)
			go parseGemlogEntry(feed, formattedUrl, line, &wg)
		}
	}

	wg.Wait()

	updated := time.UnixMicro(0) // get unix epoch timestamp 0
	for _, entry := range feed.Entries {
		date := time.Time(entry.Updated)
		if date.Sub(updated) > 0 {
			updated = date
		}
	}
	if !updated.IsZero() {
		feed.SetUpdated(updated)
	}

	HandleSuccess(w, r, feed)
}
