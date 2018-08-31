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
	"github.com/spf13/cobra"
	"github.com/wagoodman/stitch/core"
	"fmt"
	"os"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove a stitch project.",
	Long: `Remove a stitch project. Note: this will not delete any files.`,
	Args: cobra.ExactArgs(1),
	Run: deleteProject,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deleteProject(cmd *cobra.Command, args []string) {
	workspace := core.GetWorkspace()
	name := args[0]

	if !workspace.DoesProjectExist(name, "") {
		fmt.Printf("Project '%s' does not exist\n", name)
		os.Exit(1)
	}

	delete(workspace.Projects, name)
	workspace.Save()
	fmt.Printf("Project '%s' deleted\n", name)
}