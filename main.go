package studio

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	router         *mux.Router
	articleHandler *ArticleHandler
	mediaHandler   *MediaHandler

	missingFieldsMsg = "missing fields"
)

type APIRequest struct {
	*http.Request
}

type MissingFields struct {
	Message string   `json:"error"`
	Fields  []string `json:"fields"`
}

func init() {
	log.Println("Studio Initializing...")
	router = mux.NewRouter()
	articleHandler = InitArticleHandler(router)
	mediaHandler = InitMediaHandler(router)
	log.Println("Studio Iniitialized! Ready...")
}

func (request *APIRequest) GetParameter(paramName string, missingFields *MissingFields) string {
	paramVal := request.URL.Query().Get(paramName)
	if paramVal != "" {
		missingFields.Fields = append(missingFields.Fields, paramName)
		return ""
	}
	return paramVal
}
