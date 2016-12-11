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
	ArticleKind = "Article"
)

// Article represents a page that relates or explains information
// in a video.
type Article struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Time        time.Time `json:"timestamp"`
	Key         int64     `json:"key"`
}

// ArticleHandler handles article requests
type ArticleHandler struct {
	router *mux.Router
}

// InitArticleHandler returns a new initialized instance of an
// article handler.
func InitArticleHandler(router *mux.Router) *ArticleHandler {
	sub := router.PathPrefix("/article").Subrouter()
	handler := &ArticleHandler{sub}
	handler.router.HandleFunc("/get", ArticleGet)
	handler.router.HandleFunc("/add", ArticleAdd)
	handler.router.HandleFunc("/edit", ArticleEdit)
	return handler
}

func ArticleGet(w http.ResponseWriter, r *http.Request) {
	context := appengine.NewContext(r)

	if id := r.URL.Query().Get("key"); id != "" {
		key, keyErr := strconv.ParseInt(id, 10, 64)
		if keyErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(keyErr)
			return
		}
		queryKey := datastore.NewKey(context, ArticleKind, "", key, nil)
		article := &Article{}
		queryErr := datastore.Get(context, queryKey, article)
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
	} else {
		query := datastore.NewQuery(ArticleKind).Limit(25)
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
