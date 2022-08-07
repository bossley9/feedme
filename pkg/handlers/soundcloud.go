package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"html"
	"net/http"
	"strings"
	"time"

	"git.sr.ht/~bossley9/feedme/pkg/api"
	"git.sr.ht/~bossley9/feedme/pkg/atom"

	"github.com/PuerkitoBio/goquery"
)

type soundcloudResponse struct {
	Collection []struct {
		ArtworkURL       string      `json:"artwork_url"`
		Caption          interface{} `json:"caption"`
		Commentable      bool        `json:"commentable"`
		CommentCount     int         `json:"comment_count"`
		CreatedAt        time.Time   `json:"created_at"`
		Description      string      `json:"description"`
		Downloadable     bool        `json:"downloadable"`
		DownloadCount    int         `json:"download_count"`
		Duration         int         `json:"duration"`
		FullDuration     int         `json:"full_duration"`
		EmbeddableBy     string      `json:"embeddable_by"`
		Genre            string      `json:"genre"`
		HasDownloadsLeft bool        `json:"has_downloads_left"`
		ID               int         `json:"id"`
		Kind             string      `json:"kind"`
		LabelName        interface{} `json:"label_name"`
		LastModified     time.Time   `json:"last_modified"`
		License          string      `json:"license"`
		LikesCount       int         `json:"likes_count"`
		Permalink        string      `json:"permalink"`
		PermalinkURL     string      `json:"permalink_url"`
		PlaybackCount    int         `json:"playback_count"`
		Public           bool        `json:"public"`
		PurchaseTitle    interface{} `json:"purchase_title"`
		PurchaseURL      interface{} `json:"purchase_url"`
		ReleaseDate      interface{} `json:"release_date"`
		RepostsCount     int         `json:"reposts_count"`
		SecretToken      interface{} `json:"secret_token"`
		Sharing          string      `json:"sharing"`
		State            string      `json:"state"`
		Streamable       bool        `json:"streamable"`
		TagList          string      `json:"tag_list"`
		Title            string      `json:"title"`
		TrackFormat      string      `json:"track_format"`
		URI              string      `json:"uri"`
		Urn              string      `json:"urn"`
		UserID           int         `json:"user_id"`
		Visuals          interface{} `json:"visuals"`
		WaveformURL      string      `json:"waveform_url"`
		DisplayDate      time.Time   `json:"display_date"`
		Media            struct {
			Transcodings []struct {
				URL      string `json:"url"`
				Preset   string `json:"preset"`
				Duration int    `json:"duration"`
				Snipped  bool   `json:"snipped"`
				Format   struct {
					Protocol string `json:"protocol"`
					MimeType string `json:"mime_type"`
				} `json:"format"`
				Quality string `json:"quality"`
			} `json:"transcodings"`
		} `json:"media"`
		StationUrn         string `json:"station_urn"`
		StationPermalink   string `json:"station_permalink"`
		TrackAuthorization string `json:"track_authorization"`
		MonetizationModel  string `json:"monetization_model"`
		Policy             string `json:"policy"`
		User               struct {
			AvatarURL      string      `json:"avatar_url"`
			FirstName      string      `json:"first_name"`
			FollowersCount int         `json:"followers_count"`
			FullName       string      `json:"full_name"`
			ID             int         `json:"id"`
			Kind           string      `json:"kind"`
			LastModified   time.Time   `json:"last_modified"`
			LastName       string      `json:"last_name"`
			Permalink      string      `json:"permalink"`
			PermalinkURL   string      `json:"permalink_url"`
			URI            string      `json:"uri"`
			Urn            string      `json:"urn"`
			Username       string      `json:"username"`
			Verified       bool        `json:"verified"`
			City           interface{} `json:"city"`
			CountryCode    interface{} `json:"country_code"`
			Badges         struct {
				Pro          bool `json:"pro"`
				ProUnlimited bool `json:"pro_unlimited"`
				Verified     bool `json:"verified"`
			} `json:"badges"`
			StationUrn       string `json:"station_urn"`
			StationPermalink string `json:"station_permalink"`
		} `json:"user"`
		PublisherMetadata struct {
			ID              int    `json:"id"`
			Urn             string `json:"urn"`
			Artist          string `json:"artist"`
			AlbumTitle      string `json:"album_title"`
			ContainsMusic   bool   `json:"contains_music"`
			UpcOrEan        string `json:"upc_or_ean"`
			Isrc            string `json:"isrc"`
			Explicit        bool   `json:"explicit"`
			CLine           string `json:"c_line"`
			CLineForDisplay string `json:"c_line_for_display"`
			WriterComposer  string `json:"writer_composer"`
			ReleaseTitle    string `json:"release_title"`
		} `json:"publisher_metadata,omitempty"`
	} `json:"collection"`
	NextHref string      `json:"next_href"`
	QueryUrn interface{} `json:"query_urn"`
}

func fetchSoundcloudClientID(htmlDoc *goquery.Document) (string, error) {
	// reliant on the fact that the last crossorigin script contains the client id
	clientIDUrl, exists := htmlDoc.Find("script[crossorigin]").Last().Attr("src")
	if !exists {
		return "", errors.New("unable to find Soundcloud client id script source in document")
	}

	clientJSRaw, err := api.FetchGet(clientIDUrl)
	if err != nil {
		return "", err
	}

	client_js := string(clientJSRaw)
	client_id_key := "client_id"
	index_client_id_key := strings.Index(client_js, client_id_key)
	client_id_raw := client_js[index_client_id_key+len(client_id_key):]

	quote := "\""
	client_id_raw_1 := client_id_raw[strings.Index(client_id_raw, quote)+1:]
	client_id_raw_2 := client_id_raw_1[:strings.Index(client_id_raw_1, quote)]
	return client_id_raw_2, nil
}

func HandleSoundcloud(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	if len(user) == 0 {
		HandleUsage(w, r, "/soundcloud?user={USERNAME_FROM_URL}")
		return
	}

	formattedUrl := "https://soundcloud.com/" + user + "/tracks"

	res, err := api.FetchGet(formattedUrl)
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	reader := bytes.NewReader(res)
	htmlDoc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	// display name might be different that url username
	username := user
	displayName, exists := htmlDoc.Find("meta[property='og:title']").Attr("content")
	if exists {
		username = displayName
	}

	feed, err := atom.CreateFeed(formattedUrl, username, time.Now())
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	feed.AddLink(formattedUrl, atom.RelSelf)

	feed.AddAuthor(username, "", "")
	feed.SetSubtitle(username+"'s Soundcloud tracks", "text")

	image, exists := htmlDoc.Find("meta[property='og:image']").Attr("content")
	if exists {
		feed.SetLogo(image)
	}

	// get userID
	userIDUrl, exists := htmlDoc.Find("meta[property='al:ios:url']").Attr("content")
	if !exists {
		HandleBadRequest(w, r, errors.New("unable to find Soundcloud user id in document"))
		return
	}
	userIDSegments := strings.Split(userIDUrl, ":")
	userID := userIDSegments[len(userIDSegments)-1]

	clientID, err := fetchSoundcloudClientID(htmlDoc)
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	// fetch data
	data_url := "https://api-v2.soundcloud.com/users/" + userID + "/tracks?representation=&offset=&limit=30&client_id=" + clientID
	data, err := api.FetchGet(data_url)
	if err != nil {
		HandleBadRequest(w, r, err)
		return
	}

	var sc_json soundcloudResponse
	json.Unmarshal(data, &sc_json)

	for _, track := range sc_json.Collection {
		title := html.EscapeString(track.Title)

		entry, err := atom.CreateFeedEntry(track.PermalinkURL, track.Title, track.LastModified)
		if err != nil {
			continue
		}

		entry.AddLink(track.PermalinkURL, atom.RelSelf)
		entry.SetPublished(track.CreatedAt)

		var content strings.Builder
		content.WriteString("<h2>" + title + " by " + html.EscapeString(track.User.Username+"</h2>"))
		content.WriteString(`<img src="` + track.ArtworkURL + `" alt="` + title + `" />`)
		content.WriteString("<p>" + html.EscapeString(track.Description) + "</p>")
		entry.SetContent(content.String(), "html")

		if len(track.Genre) > 0 {
			entry.AddCategory(track.Genre, "", "")
		}

		feed.AddEntry(entry)
	}

	HandleSuccess(w, r, feed)
}
