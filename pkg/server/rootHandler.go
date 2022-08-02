package server

import (
	"net/http"

	h "git.sr.ht/~bossley9/feedme/pkg/handlers"

	"github.com/gorilla/mux"
)

func handleFeed(w http.ResponseWriter, r *http.Request) {
	feedType := mux.Vars(r)["type"]

	// unique feeds are prefixed by "@"
	switch feedType {
	case "acast":
		h.HandleAcast(w, r)
	case "gemini":
		h.HandleGemini(w, r)
	case "soundcloud":
		h.HandleSoundcloud(w, r)
	case "@solene":
		h.HandleSolene(w, r)
	case "@jeffgeerling":
		h.HandleJeffGeerling(w, r)
	default:
		h.HandleNotFound(w, r)
	}
}
