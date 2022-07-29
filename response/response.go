package request

type FileResponse struct {
	Filename string
	Content  []byte
}

type Response struct {
	Message    string
	File       *FileResponse
	ClientName string
	ClientIp   string
}
