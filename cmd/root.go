/*
Copyright © 2025 Efe Aslan Söyler efeaslan1703@gmail.com

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
	"strings"

	"github.com/efeaslansoyler/go-wordlistgen/internal/generator"
	"github.com/efeaslansoyler/go-wordlistgen/internal/tui"
	"github.com/spf13/cobra"
)

var (
	// CLI mode flags
	cliMode        bool
	firstName      string
	lastName       string
	birthday       string
	relatedWords   string
	minLength      string
	maxLength      string
	outputFilePath string
	enableLeet     bool
	enableCap      bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-wordlistgen",
	Short: "Generate wordlists for password cracking based on personal information",
	Long: `Go-Wordlistgen is a tool that creates customized wordlists for password cracking.
It generates potential passwords by combining personal information such as names,
birth dates, and related words with common variations like capitalization and 
leet (1337) speak substitutions.

You can use either the interactive TUI mode (default) or the CLI mode with flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cliMode {
			runCLIMode()
		} else {
			tui.Start()
		}
	},
}

func runCLIMode() {
	if firstName == "" || lastName == "" {
		fmt.Println("Error: both first name and last name are required")
		fmt.Println("Use --help for more information")
		os.Exit(1)
	}

	firstNames := strings.Fields(firstName)
	lastNames := strings.Fields(lastName)
	birthdaySlice := []string{}

	if birthday != "" {
		birthdaySlice = strings.Split(birthday, "/")
	}

	var relatedWordsSlice []string
	if relatedWords != "" {
		for _, word := range strings.Split(relatedWords, ",") {
			if trimmed := strings.TrimSpace(word); trimmed != "" {
				relatedWordsSlice = append(relatedWordsSlice, trimmed)
			}
		}
	}

	opts := generator.Options{
		InputFirstName:    firstNames,
		InputLastName:     lastNames,
		InputBirthday:     birthdaySlice,
		InputRelatedWords: relatedWordsSlice,
		InputMinLength:    minLength,
		InputMaxLength:    maxLength,
		EnableLeet:        enableLeet,
		EnableCapitalize:  enableCap,
		OutputFilePath:    outputFilePath,
	}

	fmt.Println("Generating wordlist...")
	err := generator.Run(opts)
	if err != nil {
		fmt.Printf("Error generating wordlist: %v\n", err)
		os.Exit(1)
	}

	outputPath := outputFilePath
	if outputPath == "" {
		outputPath = "wordlist.txt"
	}
	fmt.Printf("Wordlist successfully generated at: %s\n", outputPath)
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
	// Define flags for CLI mode
	rootCmd.Flags().BoolVarP(&cliMode, "cli", "c", false, "Run in CLI mode instead of TUI mode")

	// Input flags
	rootCmd.Flags().StringVarP(&firstName, "firstname", "f", "", "First name (and middle name if needed)")
	rootCmd.Flags().StringVarP(&lastName, "lastname", "l", "", "Last name")
	rootCmd.Flags().StringVarP(&birthday, "birthday", "b", "", "Birthday in format DD/MM/YYYY (or similar, use / to separate)")
	rootCmd.Flags().StringVarP(&relatedWords, "words", "w", "", "Related words separated by commas")
	rootCmd.Flags().StringVar(&minLength, "min", "", "Minimum password length (default 6)")
	rootCmd.Flags().StringVar(&maxLength, "max", "", "Maximum password length (default 12)")
	rootCmd.Flags().StringVarP(&outputFilePath, "output", "o", "", "Output file path (default wordlist.txt)")

	// Options flags
	rootCmd.Flags().BoolVar(&enableLeet, "leet", false, "Enable leet speak variations (1337)")
	rootCmd.Flags().BoolVar(&enableCap, "caps", false, "Enable capitalization variations")
}
