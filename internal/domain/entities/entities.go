package entities

type User struct {
	UID      uint64
	Email    string
	PassHash string
}

type App struct {
	ID            int64
	Name          string
	AuthSecret    string
	RefreshSecret string
}

type JwtTokenPair struct {
	AuthToken    string
	RefreshToken string
}
