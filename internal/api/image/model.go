package image

type PostUploadImageResponse struct {
	Avatar string `json:"avatar"`
}

func convertPostUploadImageResponse(avatar string) *PostUploadImageResponse {
	return &PostUploadImageResponse{
		Avatar: avatar,
	}
}
