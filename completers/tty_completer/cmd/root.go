package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tty",
	Short: "",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() error {
	return rootCmd.Execute()
}
func init() {
	carapace.Gen(rootCmd).Standalone()

	rootCmd.Flags().Bool("help", false, "display this help and exit")
	rootCmd.Flags().BoolP("silent", "s", false, "print nothing, only return an exit status")
	rootCmd.Flags().Bool("quiet", false, "print nothing, only return an exit status")
	rootCmd.Flags().Bool("version", false, "output version information and exit")
}