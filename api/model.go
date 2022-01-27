package api

import "github.com/vstdy0/go-project/model"

type shortenURLRequest struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type shortenURLResponse struct {
	Result string `json:"result"`
}

type userURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func userURLsFromCanonical(urls []model.URL, baseURL string) []userURL {
	var userURLs []userURL
	for _, v := range urls {
		userURLs = append(userURLs, userURL{ShortURL: baseURL + "/" + v.ID, OriginalURL: v.URL})
	}

	return userURLs
}
