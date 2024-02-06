
package cmd

import (
	"github.com/spf13/cobra"
	shared "github.com/CoreViewInc/CoreNiko/shared"
	"fmt"
)

var (
	 DockerCLI *DockerCLIBase
	 buildOptions shared.BuildOptions // global variable to hold build options
	 username string
	 password string
	 labels []string
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

var pullCmd = &cobra.Command{
	Use:   "pull NAME[:TAG]",
	Short: "Pull an image or a repository from a registry",
	Long:  `This command is used to pull an image or a repository from a Docker registry.`,
	Args:  cobra.ExactArgs(1), // Expect exactly one argument, the name of the image
	Run: func(cmd *cobra.Command, args []string) {
		imageName := args[0]
		DockerCLI.Service.PullImage(imageName)
	},
}

var inspectCmd = &cobra.Command{
    Use:   "inspect [OPTIONS] NAME|ID [NAME|ID...]",
    Short: "Return low-level information on Docker objects",
    Long: `Return low-level information on Docker objects, including containers, images, volumes, nodes, networks, services, and more.
By default, docker inspect will render all results in a JSON array. This command is capable of inspecting multiple targets at a time.`,
    // Allowing for a variable number of arguments 
    Args:  cobra.MinimumNArgs(1), 
	RunE: func(cmd *cobra.Command, args []string) error {
	    // Retrieve flags. This example assumes necessary flags are added to this Cobra command elsewhere in the code.
	    format, _ := cmd.Flags().GetString("format")
	    size, _ := cmd.Flags().GetBool("size")
	    targetType, _ := cmd.Flags().GetString("type")
	 
	    infoArgs := []string{}

	    // Adding type flag to arguments if specified
	    if targetType != "" {
	        infoArgs = append(infoArgs, "--type="+targetType)
	    }

	    // Adding format flag to arguments if specified
	    if format != "" {
	        infoArgs = append(infoArgs, "--format="+format)
	    }

	    // Adding size flag to arguments if specified
	    if size {
	        infoArgs = append(infoArgs, "--size")
	    }

	    // Adding the actual object names or IDs to inspect
	    infoArgs = append(infoArgs, args...)

	    // Executing inspection with the collected arguments
	    result, err := DockerCLI.Service.InspectImage(infoArgs)
	    if err != nil {
	        return err
	    }

	    // Printing result
	    fmt.Println(result)
	    return nil
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	RootCmd.AddCommand(tagCmd)
	RootCmd.AddCommand(pushCmd)
	RootCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(pullCmd)
	RootCmd.AddCommand(inspectCmd)
	buildCmd.Flags().StringArrayVarP(&buildOptions.Tags, "tag", "t", []string{}, "Name and optionally a tag in the 'name:tag' format")
	buildCmd.Flags().StringArrayVar(&labels, "label", []string{}, "Set metadata for an image")
	buildCmd.Flags().StringP("file", "f", "Dockerfile", "Name of the Dockerfile")
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Username for registry authentication")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Password for registry authentication")
	inspectCmd.Flags().StringP("type", "", "", "Specify the type of object to inspect (container, image, etc.)")
	inspectCmd.Flags().StringP("format", "f", "", "Format the output using the given Go template")
	inspectCmd.Flags().BoolP("size", "s", false, "Display total file sizes if the type is container")	
}