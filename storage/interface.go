package storage

type URLStorage interface {
	Has(id string) bool
	Set(id, url string) (string, error)
	Get(id string) string
}
