package kaniko

import (
	"bytes"
	"os/exec"
	shared "github.com/CoreViewInc/CoreNiko/shared"
	environment "github.com/CoreViewInc/CoreNiko/environment"
	"fmt"
	"os"
	"syscall"
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
	Registry		   shared.Registry
	RootDir 		   string
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

	args=append(args,"--ignore-path","/usr/bin/ping")
	args=append(args,"--ignore-path","/usr/bin/newgidmap")
	args=append(args,"--ignore-path","/usr/sbin/arping")
	args=append(args,"--ignore-path","/usr/sbin/clockdiff")
	args=append(args,"--ignore-path","//bin/ping")
	args=append(args,"--ignore-path","//bin/newgidmap")
	args=append(args,"--ignore-path","//sbin/arping")
	args=append(args,"--ignore-path","//sbin/clockdiff")
	return args
}

// Execute runs the executor with the provided configuration
func (ke *KanikoExecutor) Execute() (string, string, error) {
	args := ke.buildArgs()
	//ke.Registry.RecordImage(ke.Destination[0], "/path/to/local/image/or/remote/repository")
	// Change root to the new directory
	if err := syscall.Chroot(ke.RootDir); err != nil {
		fmt.Println("Change root to the new directory failed")
		return "","",err
	}

	// Changing directory to "/"
	if err := os.Chdir("/"); err != nil {
		fmt.Println("Change directory to / failed")
		return "","",err
	}

	fmt.Println(args)
	cmd := exec.Command("/kaniko/executor", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func New(registry shared.Registry,envProvider *environment.EnvProvider) (shared.DockerBuilder, shared.ExecutorInterface) {
	executor := &KanikoExecutor{ // Use the address operator to get a pointer to Executor
		Context:             "/workspace/",
		Dockerfile:          "Dockerfile",
		AdditionalFlags:     make(map[string]string),
		AdditionalMultiArg:  make(map[string][]string),
		Registry: registry,
		Destination: make([]string, 1),
		RootDir: envProvider.Get("RootDir"),
	}
	kanikodocker := &KanikoDocker{Executor: executor} // Use the address operator to get a pointer to KanikoDocker
	return kanikodocker, executor // Return pointers since methods are implemented with pointer receivers
}
