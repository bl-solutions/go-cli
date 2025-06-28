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
package build

import (
    "fmt"
    "time"

    "github.com/briandowns/spinner"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "go-cli/internal/build"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
    Use:   "build [app-name]",
    Short: "Build an application",
    Long:  `Build a specific application.`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        appName := args[0]
        verbose, _ := cmd.Flags().GetBool("verbose")
        
        // Read configuration for the application
        var config build.BuildConfig
        configKey := fmt.Sprintf("apps.%s", appName)
        if err := viper.UnmarshalKey(configKey, &config); err != nil {
            fmt.Printf("Error reading configuration for app '%s': %v\n", appName, err)
            return
        }
        
        // Check if configuration exists
        if config.ProjectPath == "" {
            fmt.Printf("No configuration found for application '%s'\n", appName)
            return
        }
        
        var s *spinner.Spinner
        if !verbose {
            s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
            s.Suffix = fmt.Sprintf(" Building application %s...", appName)
            s.Start()
        }

        err := build.Build(config, verbose)

        if !verbose && s != nil {
            s.Stop()
        }
        
        if err != nil {
            fmt.Printf("Error building application: %v\n", err)
        } else {
            fmt.Printf("Application %s built successfully!\n", appName)
        }
    },
}

func GetCommand() *cobra.Command {
    buildCmd.Flags().Bool("verbose", false, "Show Docker build output")
    return buildCmd
}
