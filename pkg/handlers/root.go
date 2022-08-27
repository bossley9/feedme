package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

// unique feeds prefixed by "@"
const (
	acastType        = "acast"
	geminiType       = "gemini"
	jeffgeerlingType = "@jeffgeerling"
	odyseeType       = "odysee"
	soundcloudType   = "soundcloud"
)

func getLineType(feedType string) string {
	return "* " + feedType + "\n"
}

func getDefaultUsage() string {
	return `/{type}?{param}={value}

available types are:
` +
		getLineType(acastType) +
		getLineType(geminiType) +
		getLineType(jeffgeerlingType) +
		getLineType(odyseeType) +
		getLineType(soundcloudType)
}

func handleFeed(w http.ResponseWriter, r *http.Request) {
	feedType := mux.Vars(r)["type"]

	switch feedType {
	case acastType:
		HandleAcast(w, r)
	case geminiType:
		HandleGemini(w, r)
	case jeffgeerlingType:
		HandleJeffGeerling(w, r)
	case odyseeType:
		HandleOdysee(w, r)
	case soundcloudType:
		HandleSoundcloud(w, r)
	default:
		HandleNotFound(w, r)
	}
}

func SetupRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true).UseEncodedPath()
	r.HandleFunc("/", HandleDefaultUsage)
	r.HandleFunc("/{type}", handleFeed)
	return r
}
