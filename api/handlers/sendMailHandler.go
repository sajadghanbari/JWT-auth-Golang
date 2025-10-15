package handlers

import (
	"JWT-Authentication-go/services"
	"encoding/json"
	"log"
	"net/http"
)

func sendMailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var job services.MailJob
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&job); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	maxQueueSize := 100
	JobQueue := make(chan services.MailJob,maxQueueSize)
	JobQueue <- job

	w.WriteHeader(http.StatusAccepted)
	log.Printf("Mail job queued for %s", job.To)
}