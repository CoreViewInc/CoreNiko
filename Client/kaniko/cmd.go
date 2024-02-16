package kaniko

import (
	//cmd "docker/cmd" // import the cmd package that contains our RootCmd
	shared "github.com/CoreViewInc/CoreNiko/shared"
	"fmt"
	"regexp"
	"errors"
	"os"
	io "github.com/CoreViewInc/CoreNiko/io"
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
	fileHandler := io.New()
	fmt.Println("Copying root")
	fileHandler.CopyDirToZip("/", "/kaniko/root.zip")
	stages := []string{}
	if kanikoExecutor, ok := kd.Executor.(*KanikoExecutor); ok {
		if len(contextPath) > 0 {
			kanikoExecutor.Context = contextPath
			fmt.Println("KanikoExecutor context set to:", kanikoExecutor.Context)
		} else {
			fmt.Println("KanikoExecutor context is currently:", kanikoExecutor.Context)
		}

		// Check if dockerfilePath is not empty and the file exists
		if len(dockerfilePath) > 0 {
			if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
				fmt.Printf("Dockerfile '%s' does not exist\n", dockerfilePath)
				return // Exit the function, or handle the error as needed
			} else {
				kanikoExecutor.Dockerfile = dockerfilePath
				fmt.Println("Setting KanikoExecutor Dockerfile to:", kanikoExecutor.Dockerfile)
			}
		} else {
			fmt.Println("No Dockerfile path provided")
			return // Exit the function, or handle the error as needed
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
			stdout, stderr, _ := kanikoExecutor.Execute()
			fmt.Println(stdout)
			fmt.Println(stderr)
		}
	} else {
		fmt.Println("Executor is not of type *KanikoExecutor and does not have a Context field.")
	}
	fmt.Println("Replacing root")
	fileHandler.Unzip("/kaniko/root.zip","/","/kaniko","/azp")
	fmt.Println("Kaniko build complete.")
}

func (kd *KanikoDocker) TagImage(args []string) {
	fmt.Println("Placeholder - tag")
}

func (kd *KanikoDocker) PushImage(args []string) {
	fmt.Println("Placeholder - push")
}

func (kd *KanikoDocker) Login(args []string,username string,password,url string) {
	dockerauth := auth.New()
	if len(username)>0 && len(password)>0{
		dockerauth = auth.NewUserPassAuth(username, password,url)
	}
	dockerauth.CreateDockerConfigJSON()
}

func (kd *KanikoDocker) InspectImage(args []string) (string, error){
	return "",nil
}


func (kd *KanikoDocker) PullImage(imageName string) error {
	return nil
}

func (kd *KanikoDocker) ListImages(args []string) (string, error) {
	return "",nil
}

func (kd *KanikoDocker) ImageHistory(args []string) (string, error) {
	return "placeholder",nil //temporary to provide a debuggable value
}

