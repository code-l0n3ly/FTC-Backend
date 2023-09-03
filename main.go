package FTC_App

import (
	"log"
	"net/http"
)

func main() {
	router := initializeRouter()
	log.Fatal(http.ListenAndServe(":8000", router))
}
