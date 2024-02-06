
package cmd

import (
	"github.com/spf13/cobra"
	shared "github.com/CoreViewInc/CoreNiko/shared" 
)

var (
	 DockerCLI *DockerCLIBase
	 buildOptions shared.BuildOptions // global variable to hold build options
	 username string
	 password string
)

// DockerCLI provides methods which execute docker commands.
type DockerCLIBase struct {
	Service shared.DockerBuilder
}

// NewDockerCLI creates a new DockerCLI with the provided service.
func NewDockerCLI(service shared.DockerBuilder,kaniko shared.ExecutorInterface) *DockerCLIBase {
	DockerCLI = &DockerCLIBase{Service: service}
	return DockerCLI
}


var RootCmd = &cobra.Command {
	Use:   "docker",
	Short: "CoreNiko Docker CLI",
	Long:  `CoreNiko is a Docker CLI Client implementation that uses Kaniko as a backend.`,
}

var buildCmd = &cobra.Command {
	Use:   "build [OPTIONS] PATH",
	Short: "Build an image from a Dockerfile",
	Long:  `This command is used to build an image from a Dockerfile and can accept extra options like tags.`,
	Args:  cobra.ExactArgs(1), // Make sure there is at least one argument - the context path
	Run: func(cmd *cobra.Command, args []string) {
		contextPath := args[0] // The context path is the first argument
		// Retrieve the dockerfile path from the --file flag, or use the default "Dockerfile"
		dockerfilePath, _ := cmd.Flags().GetString("file")
		// Call the build service with the options, context path, and dockerfile path
		DockerCLI.Service.BuildImage(buildOptions, contextPath, dockerfilePath)
	},
}

var tagCmd = &cobra.Command {
	Use:   "tag SOURCE_IMAGE[:TAG] TARGET_IMAGE[:TAG]",
	Short: "Create a tag TARGET_IMAGE that refers to SOURCE_IMAGE",
	Long:  `This command is used to create a tag for a source image reference.`,
	Run:   func(cmd *cobra.Command, args []string) { DockerCLI.Service.TagImage(args) },
}

var pushCmd = &cobra.Command {
	Use:   "push NAME[:TAG]",
	Short: "Push an image or a repository to a registry",
	Long:  `This command is used to push an image or a repository to a Docker registry.`,
	Run:   func(cmd *cobra.Command, args []string) { DockerCLI.Service.PushImage(args) },
}

var loginCmd = &cobra.Command {
	Use:   "login [OPTIONS] [SERVER]",
	Short: "Log in to a Docker registry",
	Long:  `This command is used to log in to a Docker registry. If no server is specified, the default is to log in to the registry at Docker Hub.`,
	Run:   func(cmd *cobra.Command, args []string) { DockerCLI.Service.Login(args, username, password) },
}


func init() {
	RootCmd.AddCommand(buildCmd)
	RootCmd.AddCommand(tagCmd)
	RootCmd.AddCommand(pushCmd)
	RootCmd.AddCommand(loginCmd)
	buildCmd.Flags().StringArrayVarP(&buildOptions.Tags, "tag", "t", []string{}, "Name and optionally a tag in the 'name:tag' format")
	buildCmd.Flags().StringP("file", "f", "Dockerfile", "Name of the Dockerfile")
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Username for registry authentication")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Password for registry authentication")
	
}