/*
Copyright © 2024 Mathieu DE SOUSA <m.desousa@bl-solutions.co>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cluster

import (
    "fmt"
    "time"

    "github.com/briandowns/spinner"
    "github.com/spf13/cobra"
    "go-cli/internal/cluster"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
    Use:   "cluster",
    Short: "Manage cluster operations",
    Long:  `Manage cluster operations including create, delete, start, and stop.`,
    Run: func(cmd *cobra.Command, args []string) {
        cmd.Help()
    },
}

// createCmd represents the create subcommand
var createCmd = &cobra.Command{
    Use:   "create",
    Short: "Create a new cluster",
    Long:  `Create a new cluster with the specified configuration.`,
    Run: func(cmd *cobra.Command, args []string) {
        verbose, _ := cmd.Flags().GetBool("verbose")
        
        var s *spinner.Spinner
        if !verbose {
            s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
            s.Suffix = " Creating cluster..."
            s.Start()
        }

        err := cluster.Create(verbose)

        if !verbose && s != nil {
            s.Stop()
        }
        
        if err != nil {
            fmt.Printf("Error creating cluster: %v\n", err)
        } else {
            fmt.Println("Cluster created successfully!")
        }
    },
}

// deleteCmd represents the delete subcommand
var deleteCmd = &cobra.Command{
    Use:   "delete",
    Short: "Delete a cluster",
    Long:  `Delete an existing cluster.`,
    Run: func(cmd *cobra.Command, args []string) {
        verbose, _ := cmd.Flags().GetBool("verbose")
        
        // Ask for confirmation
        fmt.Print("Are you sure you want to delete the cluster? (y/N): ")
        var response string
        fmt.Scanln(&response)
        
        if response != "y" && response != "Y" {
            fmt.Println("Cluster deletion cancelled.")
            return
        }
        
        var s *spinner.Spinner
        if !verbose {
            s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
            s.Suffix = " Deleting cluster..."
            s.Start()
        }

        err := cluster.Delete(verbose)

        if !verbose && s != nil {
            s.Stop()
        }
        
        if err != nil {
            fmt.Printf("Error deleting cluster: %v\n", err)
        } else {
            fmt.Println("Cluster deleted successfully!")
        }
    },
}


func GetCommand() *cobra.Command {
    createCmd.Flags().Bool("verbose", false, "Show k3d output")
    deleteCmd.Flags().Bool("verbose", false, "Show k3d output")
    clusterCmd.AddCommand(createCmd)
    clusterCmd.AddCommand(deleteCmd)
    return clusterCmd
}
