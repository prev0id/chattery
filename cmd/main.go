package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	static "chattery/website"
)

const address = ":8080"

func main() {
	http.HandleFunc(static.RootPath, handleRoot)
	http.Handle(static.SrcPath, http.FileServer(http.FS(static.Src)))

	slog.Info("starting server", slog.String("address", address))

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("http.ListenAndServe: %s", err.Error())
	}

	slog.Info("gracefully stopped")
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "" && r.URL.Path != static.RootPath {
		http.Redirect(w, r, static.RootPath, http.StatusFound)
		return
	}
	w.Write(static.IndexHTML)
}

func handleStream(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("io.ReadAll", slog.String("error", err.Error()))
		return
	}

	fmt.Println(string(data))

	w.WriteHeader(http.StatusOK)
}
