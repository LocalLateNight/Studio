package studio

import (
	"log"

	"github.com/gorilla/mux"
)

var (
	router *mux.Router
)

func init() {
	log.Println("Studio Initializing...")
	router = mux.NewRouter()
	articleHandler := InitArticleHandler(router)
	mediaHandler := InitMediaHandler(router)
	log.Println("Studio Iniitialized! Ready...")
}
