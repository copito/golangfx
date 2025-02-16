package auth

type AuthProvider interface {
	IsValid(token string) bool
	GetUser(token string) (string, error)
}
