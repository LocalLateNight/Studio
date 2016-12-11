package studio

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"

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

func InitMediaHandler(router *mux.Router) *MediaHandler {
	sub := router.PathPrefix("/media").Subrouter()
	handler := &MediaHandler{sub}
	handler.router.HandleFunc("/get", MediaGet)
	handler.router.HandleFunc("/add", MediaAdd)
	return handler
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
		queryKey := datastore.NewKey(context, MediaKind, "", key, nil)
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
		json.NewEncoder(w).Encode(media)
	} else if article := r.URL.Query().Get("article"); article != "" {
		articleID, articleErr := strconv.ParseInt(article, 10, 64)
		if articleErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(articleErr)
			return
		}
		query := datastore.NewQuery(MediaKind).Filter("Articles =", articleID).Limit(25)
		var media []Media
		_, queryErr := query.GetAll(context, &media)
		if queryErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(queryErr)
			return
		}
		json.NewEncoder(w).Encode(media)
	} else {
		query := datastore.NewQuery(MediaKind).Limit(25)
		for key, value := range r.URL.Query() {
			switch key {
			case "title":
				query = query.Filter("Title =", value[0])
			case "description":
				query = query.Filter("Description =", value[0])
			case "url":
				query = query.Filter("URL =", value[0])
			case "limit":
				limitVal, limitErr := strconv.Atoi(value[0])
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

func MediaAdd(w http.ResponseWriter, r *http.Request) {
	context := appengine.NewContext(r)
	apiReq := &APIRequest{r}

	missingFields := &MissingFields{Message: missingFieldsMsg, Fields: []string{}}
	media := &Media{}
	media.Title = apiReq.GetParameter("title", missingFields)
	media.Description = apiReq.GetParameter("description", missingFields)
	media.Uploader = apiReq.GetParameter("uploader", missingFields)
	media.Content = apiReq.GetParameter("content", missingFields)
	media.Thumbnail = apiReq.GetParameter("thumbnail", missingFields)

	date := apiReq.GetParameter("date", missingFields)
	intTime, intErr := strconv.ParseInt(date, 10, 64)
	if intErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(intErr)
		return
	}

	articles := []int64{}
	for _, articleStr := range apiReq.URL.Query()["article"] {
		articleInt, articleErr := strconv.ParseInt(articleStr, 10, 64)
		if articleErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(intErr)
			return
		}
		articles = append(articles, articleInt)
	}
	media.Articles = articles

	media.Date = time.Unix(intTime, 0)
	media.Views = 0

	if len(missingFields.Fields) > 0 {
		jsonMsg, jsonErr := json.Marshal(missingFields)
		if jsonErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(intErr)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(jsonMsg), http.StatusBadRequest)
	}
	id, _, idErr := datastore.AllocateIDs(context, ArticleKind, nil, 0)
	if idErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(idErr)
		return
	}
	key := datastore.NewKey(context, ArticleKind, "", id, nil)
	article.Key = key.IntID()
	_, putErr := datastore.Put(context, key, article)
	if putErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(putErr)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}
