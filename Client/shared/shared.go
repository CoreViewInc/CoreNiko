package shared

// Defintion a docker cli must implement
type DockerBuilder interface {
    BuildImage(options BuildOptions, contextPath string, dockerfilePath string)
    TagImage(args []string)
    PushImage(args []string)
    Login(args []string,username string,password string)
}

// Definition for executor, in this case kaniko is the only provided execution provider.
type ExecutorInterface interface {
	Execute() (stdout string, stderr string, err error)
}

// BuildOptions defines options for the docker build command.
type BuildOptions struct {
	Tags []string // To store the tags from the -t flag
}

// variables of a common executor
type Executor struct {
	Context            string
	Dockerfile         string
	Destination        []string
	Cache              bool
	CacheDir           string
	NoPush             bool
	BuildArgs          []string
	AdditionalFlags    map[string]string
	AdditionalMultiArg map[string][]string
	Registry *Registry
}

// DockerImageComponents holds the components of a Docker image tag.
type DockerImageComponents struct {
	Registry string
	Repo     string
	Tag      string
}

// Custom registry definition, This is for the local kaniko builds as they have no concept of a local build registry
type Registry interface {
	Initialize() error
	RecordImage(tag, location string) error
	GetImageLocation(tag string) (string, error)
	Close() error
}

// Reassemble the image destination from the correctly previously parsed tag
func (dic DockerImageComponents) GetFullImageName() string {
	fullImageName := ""
	if dic.Registry != "" {
		fullImageName += dic.Registry + "/"
	}
	fullImageName += dic.Repo
	if dic.Tag != "" {
		fullImageName += ":" + dic.Tag
	}
	return fullImageName
}