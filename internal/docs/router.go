package docs

import (
    "net/http"

    "github.com/swaggo/swag"
)

// Handler serves the generated swagger JSON at /swagger/doc.json.
func Handler(mux *http.ServeMux) {
    mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
        doc, err := swag.ReadDoc()
        if err != nil {
            http.Error(w, "swagger not available", http.StatusServiceUnavailable)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        _, _ = w.Write([]byte(doc))
    })
}


