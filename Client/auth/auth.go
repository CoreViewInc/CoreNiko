package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	environment "github.com/CoreViewInc/CoreNiko/environment"
)

type DockerAuth struct {
	Username string
	Password string
}

type DockerAuthConfig struct {
	Auths map[string]map[string]string `json:"auths"`
}
 
func NewEnvAuth(envProvider *environment.EnvProvider) *DockerAuth {
	username := envProvider.Get("DOCKER_USERNAME")
	password := envProvider.Get("DOCKER_PASSWORD")
	return &DockerAuth{
		Username: username,
		Password: password,
	}
}

func NewUserPassAuth(username, password string) *DockerAuth {
	return &DockerAuth{
		Username: username,
		Password: password,
	}
}

func (da *DockerAuth) EncodeCredentials() (string, error) {
	credentials := fmt.Sprintf("%s:%s", da.Username, da.Password)
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	return encodedCredentials, nil
}

func (da *DockerAuth) CreateDockerConfigJSON() error {
	encodedCredentials, err := da.EncodeCredentials()
	if err != nil {
		return fmt.Errorf("error encoding Docker credentials: %v", err)
	}

	dockerConfig := DockerAuthConfig{
		Auths: map[string]map[string]string{
			"https://index.docker.io/v1/": {
				"auth": encodedCredentials,
			},
		},
	}

	jsonData, err := json.Marshal(dockerConfig)
	if err != nil {
		return fmt.Errorf("error marshaling Docker config.json: %v", err)
	}

	err = ioutil.WriteFile("/kaniko/.docker/config.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing Docker config.json: %v", err)
	}

	return nil
}

func New() *DockerAuth {
	return NewEnvAuth(environment.New())
}		