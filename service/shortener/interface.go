package shortener

type URLService interface {
	AddURL(id string) (string, error)
	GetURL(id string) string
}
