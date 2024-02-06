package kaniko

import (
	"bytes"
	"os/exec"
	shared "github.com/CoreViewInc/CoreNiko/shared"
	"fmt"
)

type KanikoExecutor struct {
	Context            string
	Dockerfile         string
	Destination        []string
	Cache              bool
	CacheDir           string
	NoPush             bool
	BuildArgs          []string
	AdditionalFlags    map[string]string
	AdditionalMultiArg map[string][]string
	Registry	shared.Registry
}


func (ke *KanikoExecutor) buildArgs() []string {
	args := []string{}

	if ke.Context != "" {
		args = append(args, "--context", ke.Context)
	}

	if ke.Dockerfile != "" {
		args = append(args, "--dockerfile", ke.Dockerfile)
	}

	for _, destination := range ke.Destination {
		args = append(args, "--destination", destination)
		//disabled, subvert
	}

	if ke.Cache {
		args = append(args, "--cache")
	}

	if ke.CacheDir != "" {
		args = append(args, "--cache-dir", ke.CacheDir)
	}

	if ke.NoPush {
		args = append(args, "--no-push")
	}

	for _, buildArg := range ke.BuildArgs {
		args = append(args, "--build-arg", buildArg)
	}

	for key, value := range ke.AdditionalFlags {
		args = append(args, key, value)
	}

	for key, values := range ke.AdditionalMultiArg {
		for _, value := range values {
			args = append(args, key, value)
		}
	}

	return args
}

// Execute runs the executor with the provided configuration
func (ke *KanikoExecutor) Execute() (string, string, error) {
	args := ke.buildArgs()
	//ke.Registry.RecordImage(ke.Destination[0], "/path/to/local/image/or/remote/repository")
	fmt.Println(args)
	cmd := exec.Command("executor", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func New(registry shared.Registry) (shared.DockerBuilder, shared.ExecutorInterface) {
	executor := &KanikoExecutor{ // Use the address operator to get a pointer to Executor
		Context:             "/workspace/",
		Dockerfile:          "Dockerfile",
		AdditionalFlags:     make(map[string]string),
		AdditionalMultiArg:  make(map[string][]string),
		Registry: registry,
		Destination: make([]string, 1),
	}
	kanikodocker := &KanikoDocker{Executor: executor} // Use the address operator to get a pointer to KanikoDocker
	return kanikodocker, executor // Return pointers since methods are implemented with pointer receivers
}
