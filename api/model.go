package api

type ShortenURLRequest struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type ShortenURLResponse struct {
	Result string `json:"result"`
}
