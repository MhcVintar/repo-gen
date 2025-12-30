package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var (
		source         string
		repository     string
		destination    string
		packageName    string
		implementation string
	)

	cmd := &cobra.Command{
		Use:   "repogen",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().StringVarP(&source, "source", "s", "", "source file with repository interface definition")
	cmd.Flags().StringVarP(&repository, "repository", "r", "", "repository interface symbol")
	cmd.Flags().StringVarP(&destination, "destination", "d", "", "destination for generated file")
	cmd.Flags().StringVarP(&packageName, "package", "p", "", "package name for generated file")
	cmd.Flags().StringVarP(&implementation, "implementation", "i", "", "implementation symbol for generated repository")

	requiredFlags := []string{"source", "repository", "destination", "package", "implementation"}
	for _, f := range requiredFlags {
		if err := cmd.MarkFlagRequired(f); err != nil {
			log.Fatalf("could not mark %q as required: %v", f, err)
		}
	}

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
