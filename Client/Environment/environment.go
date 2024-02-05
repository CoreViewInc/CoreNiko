package environment

import (
	"os"
)

// EnvProvider holds the environment variables map
type EnvProvider struct {
	EnvVars map[string]string
}

func (ep *EnvProvider) Get(key string) string {
	return ep.EnvVars[key]
}

func New() *EnvProvider {
	ep := &EnvProvider{EnvVars: make(map[string]string)}
	for _, env := range os.Environ() {
		pair := splitPair(env, '=')
		ep.EnvVars[pair[0]] = pair[1]
	}
	return ep
}

func splitPair(s string, sep byte) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s} // If separator not found, return the whole string as the first element.
}

type Environment interface {
	Get() string
	New() *EnvProvider
}
