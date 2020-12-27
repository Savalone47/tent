/*
Copyright © 2020 FARHAN HASIN CHOWDHURY <MAIL@FARHAN.INFO>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/fhsinchy/tent/utils"

	"github.com/spf13/cobra"

	"github.com/containers/podman/v2/pkg/specgen"
)

var isDefault bool

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		connText := utils.GetContext()

		service := args[0]

		switch service {
		case "mysql":
			tag := "latest"
			password := "secret"
			containerName := "tent-mysql"
			var hostPort uint16 = 3306

			if !isDefault {
				var tagInput string
				var passwordInput string
				var portInput uint16

				fmt.Print("Which tag you want to use? (default: latest): ")
				fmt.Scanln(&tagInput)

				fmt.Print("Password for the root user? (default: secret): ")
				fmt.Scanln(&passwordInput)

				fmt.Print("Host system port? (default: 3306): ")
				fmt.Scanln(&portInput)

				if tagInput != "" {
					tag = tagInput
				}

				if passwordInput != "" {
					password = passwordInput
				}

				if portInput != 0 {
					hostPort = portInput
				}
			}

			portMapping := specgen.PortMapping{
				ContainerPort: 3306,
				HostPort:      hostPort,
			}

			rawImage := "docker.io/mysql:" + tag
			env := make(map[string]string)
			env["MYSQL_ROOT_PASSWORD"] = password

			utils.PullImage(connText, rawImage)
			utils.CreateContainer(connText, rawImage, env, containerName, portMapping)
			utils.StartContainer(connText, containerName)
		default:
			fmt.Println("invalid service name given")
		}
	},
}

func init() {
	startCmd.Flags().BoolVarP(&isDefault, "default", "d", false, "starts the service with default options")

	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
