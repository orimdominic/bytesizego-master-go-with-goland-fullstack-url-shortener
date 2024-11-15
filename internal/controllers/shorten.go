package controllers

import (
	"bytesizego-url-shortener/internal/db"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html/template"
	"net/http"
	"strings"
)

func Shorten(sqlite *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		originalURL := r.FormValue("url")
		if originalURL == "" {
			http.Error(w, "URL not provided", http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
			originalURL = "https://" + originalURL
		}

		hash := generateURLHash(originalURL)
		_, err:= db.CreateUrlRecord(sqlite, originalURL, hash)

		if err!= nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		data := map[string]string{
			"ShortURL": hash,
		}

		tmp, err := template.ParseFiles("internal/views/shorten.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmp.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func generateURLHash(originalURL string) string {
	h := sha256.New()
	h.Write([]byte(originalURL))
	hash := hex.EncodeToString(h.Sum(nil))
	shortHash := hash[:8]
	return shortHash
}

func GetURL(sqlite *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		if hash == "" {
			http.Error(w, "Invalid unique hash id", http.StatusBadRequest)
			return
		}

		url, err := db.GetUrlRecord(sqlite, hash)
		// TODO check if err is err.NoRows() ro return 404
		if err != nil {
			http.Error(w, "Unable to fetch URL", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusPermanentRedirect)
		return
	}
}
