package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Halo menulis "Halo" ke response ketika endpoint / dipanggil.
// w = tempat kita menulis response ke browser/client.
// r = request dari client (berisi URL, method, header, dll.).
func Halo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Halo")
}

// healthResponse dipakai untuk mengirim response JSON di endpoint /health.
type healthResponse struct {
	Status string `json:"status"`
}

// Health adalah handler untuk GET /health. Mengembalikan JSON {"status":"ok"}
// dan status code 200. Berguna untuk health check (monitoring, load balancer).
func Health(w http.ResponseWriter, r *http.Request) {
	// Hanya terima method GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthResponse{Status: "ok"})
}
