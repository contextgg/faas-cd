package cmd

import (
	"fmt"
	"context"

	"github.com/contextgg/faas-cd/templater"

	"github.com/openfaas/faas-cli/stack"
	"github.com/spf13/cobra"
)

// packCmd represents the pack command
var packCmd = &cobra.Command{
	Use:   `pack`,
	Short: "pack makes docker buildables",
	Long:  `pack converts functions into buildable docker specs`,
	Example: `
  faas-cd pack -f https://domain/path/service.yml
  faas-cd pack -f ./service.yml`,
	RunE: runPack,
}

func init() {
	rootCmd.AddCommand(packCmd)
}

func runPack(cmd *cobra.Command, args []string) error {
	parsedServices, err := stack.ParseYAMLFile(yamlFile, regex, filter, envsubst)
	if err != nil {
		return err
	}

	var opts []templater.TemplaterOption
	for _, templateSource := range parsedServices.StackConfiguration.TemplateConfigs {
		opts = append(opts, templater.AddLocationOption(templateSource.Name, templateSource.Source))
	}

	t := templater.NewTemplater(opts...)
	for name, fn := range parsedServices.Functions {
		// Need to fetch templates.
		t.AddFunction(name, fn.Language)
	}

	downloaded, err := t.Download(context.Background())
	if err != nil {
		return err
	}
	for _, path := range downloaded {
		fmt.Printf("Fetched %s", path)
	}

	packed, err := t.Pack(context.Background())
	if err != nil {
		return err
	}
	for _, path := range packed {
		fmt.Printf("Packed %s", path)
	}

	return nil
}