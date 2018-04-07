package auth

type Config struct {
	AuthType AuthType
	Scopes   []string
}

type AuthType int

const (
	NONE AuthType = iota
	LOGGEDIN
	SCOPE
)

func AnyScope(scopes ...string) *Config {
	return &Config{
		AuthType: SCOPE,
		Scopes:   scopes,
	}
}

var LoggedIn = &Config{
	AuthType: LOGGEDIN,
	Scopes:   []string{},
}

var None = &Config{
	AuthType: NONE,
}
