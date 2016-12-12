package studio

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"
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
type ArticleHandler struct{}

// InitArticleHandler returns a new initialized instance of an
// article handler.
func InitArticleHandler() *ArticleHandler {
	handler := &ArticleHandler{}
	http.HandleFunc("/article/get", ArticleGet)
	http.HandleFunc("/article/add", ArticleAdd)
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(article)
	}
}

func ArticleAdd(w http.ResponseWriter, r *http.Request) {
	context := appengine.NewContext(r)
	apiReq := &APIRequest{r}

	log.Println(r.URL.Query())
	missingFields := &MissingFields{Message: missingFieldsMsg, Fields: []string{}}
	article := &Article{}
	article.Title = apiReq.GetParameter("title", missingFields)
	article.Description = apiReq.GetParameter("description", missingFields)
	article.URL = apiReq.GetParameter("url", missingFields)

	timestamp := apiReq.GetParameter("timestamp", missingFields)
	log.Println(article)
	intTime, intErr := strconv.ParseInt(timestamp, 10, 64)
	if intErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(intErr)
		return
	}
	article.Time = time.Unix(intTime, 0)

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
	id, _, idErr := datastore.AllocateIDs(context, ArticleKind, nil, 1)
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
	log.Println(key)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}
