/*
Copyright © 2024 Ashish Alex ashish.alex10@gmail.com

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
package cmd

import (
	"fmt"
	"os"

	"github.com/ashish10alex/formatdataform/internal/version"
	"github.com/spf13/cobra"
)

var getVersionInfo bool
var sqlfluffConfigPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "formatdataform",
	Short: "Format .sqlx files",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//TODO: might add version information here
	Run: func(cmd *cobra.Command, args []string) {

		if getVersionInfo {
			versionInfo := version.Get()
			fmt.Println("formatdataform is a command line too to format .sqlx files in Dataform project")
			fmt.Println("")
			fmt.Println("Git Version: ", versionInfo.GitVersion)
			fmt.Println("Git Commit:  ", versionInfo.GitCommit)
			fmt.Println("Build Date:  ", versionInfo.BuildDate)
			fmt.Println("Go Version:  ", versionInfo.GoVersion)
			fmt.Println("Platform:    ", versionInfo.Platform)
		} else {
			cmd.Help()
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.formatdataform.yaml)")

	rootCmd.Flags().BoolVarP(&getVersionInfo, "version", "v", false, "Returns the version of the binary")
	rootCmd.PersistentFlags().StringP("sqlfluff_config_path", "c", "", "Path to the sqlfluff config file")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
