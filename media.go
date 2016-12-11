package studio

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var (
	MediaKind = "Media"
)

// Media contains the generic information for any form of
// media.
type Media struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Uploader    int       `json:"uploader"`
	Date        time.Time `json:"date"`
	Type        int       `json:"type"`
	Views       int       `json:"views"`
	Articles    []int64   `json:"articles"`
	Thumbnail   string    `json:"thumbnail"`
	Content     string    `json:"content"`
	Key         int64     `json:"key"`
}

// MediaHandler handles MediaRequests
type MediaHandler struct {
	router *mux.Router
}

// GetMediaFile returns the media file location
func (audio *MediaAudio) GetMediaFile() string {
	return audio.MediaFile
}

// GetMediaFile returns the media file location
func (video *MediaVideo) GetMediaFile() string {
	return video.MediaFile
}

func InitMediaHandler(router *mux.Router) *MediaHandler {
	sub := router.PathPrefix("/media").Subrouter()
	handler := &MediaHandler{sub}
	handler.router.HandleFunc("/get", MediaGet)
}

func MediaGet(w http.ResponseWriter, r *http.Request) {
	context := appengine.NewContext(r)

	if id := r.URL.Query().Get("key"); id != "" {
		key, keyErr := strconv.ParseInt(id, 10, 64)
		if keyErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(keyErr)
			return
		}
		queryKey := datastore.NewKey(context, MediaKey, "", key, nil)
		media := &Media{}
		queryErr := datastore.Get(context, queryKey, media)
		if queryErr != nil {
			if queryErr == datastore.ErrNoSuchEntity {
				http.NotFound(w, r)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(queryErr)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(article)
	} else if article := r.URL.Query().Get("article"); article != "" {
		articleID, articleErr := strconv.ParseInt(article, 10, 64)
		if articleErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(keyErr)
			return
		}
		query := datastore.NewQuery(MediaKind).Filter("Articles =", articleID).Limit(25)
		var media []Media
		queryErr := query.GetAll(context, &media)
		if queryErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		json.NewEncoder(w).Encode(media)
	} else {
		query := datastore.NewQuery(MediaKind).Limit(25)
		for key, value := range r.URL.Query() {
			switch key {
			case "title":
				query = query.Filter("Title =", value)
			case "description":
				query = query.Filter("Description =", value)
			case "url":
				query = query.Filter("URL =", value)
			case "limit":
				limitVal, limitErr := strconv.Atoi(value)
				if limitErr != nil {
					log.Println(limitErr)
					continue
				}
				query = query.Limit(limitVal)
			}
		}
		var article []Article
		_, err := query.GetAll(context, &article)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		json.NewEncoder(w).Encode(article)
	}
}
