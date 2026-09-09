package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func startMockServer() {
	release := githubRelease{
		TagName: "v1.1.0",
		Assets: []releaseAsset{
			{Name: "nametag-linux-amd64.tar.gz", BrowserDownloadURL: "http://127.0.0.1:8080/nametag.bin"},
		},
	}

	_, err := json.Marshal(release)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(release)
	})

	http.HandleFunc("/nametag.bin", func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile("./nametag")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(content)
	})

	fmt.Println("serving mock release server on :8080")
	_ = http.ListenAndServe(":8080", nil)
}
