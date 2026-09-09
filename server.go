package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func startMockServer() {
	manifest := manifest{
		Version:     "v1.1.0",
		DownloadURL: "http://127.0.0.1:8080/nametag.bin",
		SHA256:      "",
		Description: "New release",
	}

	_, err := json.Marshal(manifest)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(manifest)
	})

	http.HandleFunc("/nametag.bin", func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile("./nametag")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(content)
	})

	fmt.Println("serving mock update server on :8080")
	_ = http.ListenAndServe(":8080", nil)
}
