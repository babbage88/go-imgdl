package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate shell completion scripts",
	Long: `To load completions:

Zsh:

  $ source <(image-downloader completion zsh)

  $ image-downloader completion zsh > ~/.zsh/completions/_image-downloader
  $ fpath=(~/.zsh/completions $fpath)
  $ autoload -U compinit && compinit

Bash:

  $ source <(image-downloader completion bash)
  $ image-downloader completion bash > /etc/bash_completion.d/image-downloader

Fish:

  $ image-downloader completion fish | source
  $ image-downloader completion fish > ~/.config/fish/completions/image-downloader.fish
`,
}

var bashCompletionCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate bash completion script",
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

var zshCompletionCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate zsh completion script",
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenZshCompletion(os.Stdout)
	},
}

var fishCompletionCmd = &cobra.Command{
	Use:   "fish",
	Short: "Generate fish completion script",
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenFishCompletion(os.Stdout, true)
	},
}

func init() {
	completionCmd.AddCommand(bashCompletionCmd)
	completionCmd.AddCommand(zshCompletionCmd)
	completionCmd.AddCommand(fishCompletionCmd)
	rootCmd.AddCommand(completionCmd)
}
