package domain

type Repository interface {
	Save(code, url string)
	Get(code string) (string, bool)
}
