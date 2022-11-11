package cloudfunctions

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("RunWatchEmails", runWatchEmails)
}

func runWatchEmails(w http.ResponseWriter, r *http.Request) {
	log.Println("received watch trigger")
}
