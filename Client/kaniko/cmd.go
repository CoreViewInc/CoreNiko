package kaniko

import (
	//cmd "docker/cmd" // import the cmd package that contains our RootCmd
	shared "github.com/CoreViewInc/CoreNiko/shared"
	"fmt"
	"regexp"
	"errors"
	auth "github.com/CoreViewInc/CoreNiko/auth"
)

type KanikoDocker struct{
	Executor shared.ExecutorInterface
}

// ParseDockerImageTag takes a Docker image tag string and extracts its components.
func (kd *KanikoDocker) ParseDockerImageTag(imageTag string) (shared.DockerImageComponents,error) {
	var components shared.DockerImageComponents
	tagRegex := regexp.MustCompile(`^(?:(?P<Registry>.+?)/)?(?P<Repo>[^/:]+)(?::(?P<Tag>.+))?$`)
	matches := tagRegex.FindStringSubmatch(imageTag)
	for i, name := range tagRegex.SubexpNames() {
		if i != 0 && name != "" && i <= len(matches) {
			switch name {
			case "Registry":
				components.Registry = matches[i]
			case "Repo":
				components.Repo = matches[i]
			case "Tag":
				components.Tag = matches[i]
			}
		}
	}
	if components.Tag == "" { // assign a default tag if none were found
		components.Tag = "latest"
	}

	// Check if we have valid components to construct a registry
	if components.Registry == "" && components.Repo == "" {
		return shared.DockerImageComponents{}, errors.New("invalid image tag: neither registry nor repository specified")
	}
	if components.Repo == "" {
		return shared.DockerImageComponents{}, errors.New("invalid image tag: repository cannot be empty")
	}
	return components, nil
}

func (kd *KanikoDocker) BuildImage(options shared.BuildOptions, contextPath string, dockerfilePath string) {
	stages := []string{}
	if kanikoExecutor, ok := kd.Executor.(*KanikoExecutor); ok {
		if len(contextPath) > 0 {
			kanikoExecutor.Context = contextPath
			fmt.Println("KanikoExecutor context set to:", kanikoExecutor.Context)
		} else {
			fmt.Println("KanikoExecutor context is currently:", kanikoExecutor.Context)
		}
		if len(dockerfilePath) > 0{
			kanikoExecutor.Dockerfile = dockerfilePath
		}else{
			fmt.Println("KanikoExecutor Dockerfile is currently:", kanikoExecutor.Dockerfile)
		}

		//build must have atleast a tag otherwise it should generated random uuid
		if len(options.Tags)>0{
			for _,tag := range options.Tags {
				parsed_tag,err := kd.ParseDockerImageTag(tag)
				if err!=nil{
					panic(err)
				}
				stages = append(stages,parsed_tag.GetFullImageName())
			}
		}else{
			//no tag provided
			fmt.Println("No tag providing error!")
		}
		for _,stage := range stages{
			kanikoExecutor.Destination[0] = stage
			kanikoExecutor.Execute()
		}
	} else {
		fmt.Println("Executor is not of type *KanikoExecutor and does not have a Context field.")
	}
}

func (kd *KanikoDocker) TagImage(args []string) {
	fmt.Println("Placeholder - tag")
}

func (kd *KanikoDocker) PushImage(args []string) {
	fmt.Println("Placeholder - push")
}

func (kd *KanikoDocker) Login(args []string,username string,password string) {
	dockerauth := auth.New()
	if len(username)>0 && len(password)>0{
		dockerauth = auth.NewUserPassAuth(username, password)
	}
	dockerauth.CreateDockerConfigJSON()
}

func (kd *KanikoDocker) InspectImage(args []string) (string, error){
	return "",nil
}


func (kd *KanikoDocker) PullImage(imageName string) error {
	return nil
}

