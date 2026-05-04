package image

type PostUploadImageResponse struct {
	Avatar string `json:"avatar"`
}

type DeleteImageResponse struct {
	Avatar string `json:"avatar"`
}

func convertPostUploadImageResponse(avatar string) *PostUploadImageResponse {
	return &PostUploadImageResponse{
		Avatar: avatar,
	}
}

func convertDeleteImageResponse(avatar string) *DeleteImageResponse {
	return &DeleteImageResponse{
		Avatar: avatar,
	}
}
