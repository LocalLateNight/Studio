package studio

import (
	"log"
	"net/http"
)

var (
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
	articleHandler = InitArticleHandler()
	mediaHandler = InitMediaHandler()
	log.Println("Studio Iniitialized! Ready...")
}

func (request *APIRequest) GetParameter(paramName string, missingFields *MissingFields) string {
	paramVal := request.URL.Query().Get(paramName)
	if paramVal == "" {
		missingFields.Fields = append(missingFields.Fields, paramName)
		return ""
	}
	return paramVal
}
