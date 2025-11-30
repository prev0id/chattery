package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	"chattery/internal/pb/api/websocketpb"
	"chattery/internal/service/signaling"
	static "chattery/website"
)

const address = "localhost:8080"

func main() {
	appCtx := context.Background()

	http.HandleFunc(static.RootPath, handleRoot)
	http.HandleFunc(static.SettingsPath, handleSettings)
	http.HandleFunc(static.AuthPath, handleAuth)
	// http.HandleFunc(static.NotFoundPath, handle404)
	http.Handle(static.SrcPath, http.FileServer(http.FS(static.Src)))
	http.HandleFunc("/ws", websocketHandler(appCtx))

	slog.Info("starting server", slog.String("address", "http://"+address))

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

// func handle404(w http.ResponseWriter, _ *http.Request) {
// 	w.Write(static.NotFound)
// }

func handleSettings(w http.ResponseWriter, _ *http.Request) {
	w.Write(static.Settings)
}

func handleAuth(w http.ResponseWriter, _ *http.Request) {
	w.Write(static.Auth)
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

func websocketHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signalingService, err := signaling.NewService(w, r)
		if err != nil {
			slog.Error("websocket.New", slog.String("error", err.Error()))
		}

		signalingService.RegisterMessageCallback(websocketpb.Type_ENUM_TYPE_JOIN_CHAT, func(_ context.Context, msg *websocketpb.Message) {
			fmt.Printf("JOIN: %+v\n", msg)
		})

		signalingService.RegisterMessageCallback(websocketpb.Type_ENUM_TYPE_LEAVE_CHAT, func(_ context.Context, msg *websocketpb.Message) {
			fmt.Printf("LEAVE: %+v\n", msg)
		})

		signalingService.ListenAndServe(ctx)
	}
}
