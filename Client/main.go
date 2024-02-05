package main

import (
	"os"
	"fmt"
	docker "docker/cmd" // import the cmd package that contains our RootCmd
	kanikocmd "docker/kaniko" // import the cmd package that contains our RootCmd
	registry "docker/registry"
)

func main() {
	registry, err := registry.New("registry.db")
	if err!=nil{
		panic(err)
	}
	KanikoDockerCLI,KanikoCLI := kanikocmd.New(registry)
	docker.NewDockerCLI(KanikoDockerCLI,KanikoCLI)
	// Execute the root command from cmd package
	if err := docker.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}