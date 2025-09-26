package dto

type UploadBase64Request struct {
	FileName    string `json:"file_name" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	DataBase64  string `json:"data_base64" binding:"required"`
}

type UploadResponse struct {
	IdImage  string `json:"id_image"`
	FileName string `json:"file_name"`
}
