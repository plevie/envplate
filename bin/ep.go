package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/plevie/envplate"
	"github.com/spf13/cobra"
	"github.com/yawn/doubledash" 
)

var (
	build   string
	version string
)

func init() {
	os.Args = doubledash.Args
}

func main() {

	var ( // flags
		prefix  *string
		backup  *bool
		dryRun  *bool
		strict  *bool
		verbose *bool
	)

	root := &cobra.Command{
		Use:   "ep",
		Short: fmt.Sprintf("envplate %s (%s) provides trivial templating for configuration files using environment keys", version, build),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			envplate.Logger.Verbose = *verbose
		},
		Run: func(cmd *cobra.Command, args []string) {

			var h = envplate.Handler{
				Prefix: *prefix
				Backup: *backup,
				DryRun: *dryRun,
				Strict: *strict,
			}

			if err := h.Apply(args); err != nil {
				os.Exit(1)
			}

			if h.DryRun {
				os.Exit(0)
			}

			if len(doubledash.Xtra) > 0 {

				if err := syscall.Exec(doubledash.Xtra[0], doubledash.Xtra, os.Environ()); err != nil {
					log.Fatalf("Cannot exec '%v': %v", doubledash.Xtra, err)
				}

			}

		},
	}

	// flag parsing
	prefix = root.Flags().StringP("prefix", "p", false, "Only preplace env vars starting with prefix")
	backup = root.Flags().BoolP("backup", "b", false, "Create a backup file when using inline mode")
	dryRun = root.Flags().BoolP("dry-run", "d", false, "Dry-run - output templates to stdout instead of inline replacement")
	strict = root.Flags().BoolP("strict", "s", false, "Strict-mode - fail when falling back on defaults")
	verbose = root.Flags().BoolP("verbose", "v", false, "Verbose logging")

	if err := root.Execute(); err != nil {
		log.Fatalf("Failed to start the application: %v", err)
	}

}
