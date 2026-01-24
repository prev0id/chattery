package chat_api

import (
	"net/http"

	"chattery/internal/utils/render"
	"chattery/internal/utils/sliceutil"
)

const searchQueryName = "query"

// Search поиск новых чатов
func (s *Server) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query().Get(searchQueryName)

	chats, err := s.chat.SearchChats(ctx, query)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, SearchChatsResponse{Chats: sliceutil.Map(chats, convertChat)})
}
