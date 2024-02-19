package main

import (
	"os"
	"fmt"
	docker "github.com/CoreViewInc/CoreNiko/cmd" // import the cmd package that contains our RootCmd
	kanikocmd "github.com/CoreViewInc/CoreNiko/kaniko" // import the cmd package that contains our RootCmd
	registry "github.com/CoreViewInc/CoreNiko/registry"
	environment "github.com/CoreViewInc/CoreNiko/environment"
)

func main() {
	registry, err := registry.New("registry.db")
	if err!=nil{
		panic(err)
	}
	KanikoDockerCLI,KanikoCLI := kanikocmd.New(registry,environment.New())
	docker.NewDockerCLI(KanikoDockerCLI,KanikoCLI)
	// Execute the root command from cmd package
	if err := docker.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}