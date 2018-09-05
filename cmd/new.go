// Copyright © 2018 Alex Goodman
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wagoodman/stitch/core"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [project] [url]",
	Short: "Create a new project",
	Long: `Create a new project that will be made of one or more repositories.

todo: example usage with git url, local url, plain http url, and gist url
`,
	Args: cobra.RangeArgs(1,2),
	Run:  newProject,
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func newProject(cmd *cobra.Command, args []string) {
	workspace := core.GetWorkspace()
	url := args[0]

	if workspace.DoesProjectExist("", url) {
		fmt.Printf("Unable to add project: project already exists\n")
		os.Exit(1)
	}

	project, err := core.ReadConfig(url)
	core.Check(err, "unable to load project")

	// add stitch-project to the workspace store
	if err := workspace.AddProject(project); err != nil {
		fmt.Printf("unable to add project: %s\n", err)
		os.Exit(1)
	}

	for _, repository := range project.Repositories {
		if !core.PathExists(repository.Path) {
			fmt.Printf("Cloning '%s' repository...\n", repository.Name)

			err := repository.Clone()
			core.Check(err, "unable to clone repository")

		} else {
			// Note: we may not want to pull the latest...

			// fmt.Printf("Repository '%s' already exists. pulling latest...\n", repository.Name)
			// err := repository.Pull()
			// if err == git.NoErrAlreadyUpToDate {
			// 	fmt.Println("...already up to date")
			// } else {
			// 	Check(err, "unable to pull latest from repository")
			// }

			fmt.Printf("Repository '%s' already exists. \n", repository.Name)
		}
	}

	// test loading this project
	project.Load()

	fmt.Println("Added project!")
	workspace.CurrentProjectUrl = project.Repository.GitURL
	workspace.Save()
}
