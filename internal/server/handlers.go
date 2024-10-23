package server

import "net/http"

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST request allowed!", http.StatusMethodNotAllowed)
		return
	}
}
