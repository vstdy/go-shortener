package shortener

type URLService interface {
	AddURL(id string) string
	GetURL(id string) string
}
