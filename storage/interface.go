package storage

type URLStorage interface {
	Has(id string) bool
	Set(id, url string)
	Get(id string) string
}
