package cmd

import (
	//embed files
	_ "embed"
	"fmt"
	"snipio/app"

	"github.com/spf13/cobra"
)

//go:embed version
var versionString string

func makeShowCommand(a *app.App) *cobra.Command {

	cmd := cobra.Command{

		Use:   "show",
		Short: "shows ressources",
		Long:  "shows ressources",
	}

	cmd.AddCommand(makeShowLayersCommand(a))

	return &cmd
}

func makeDeleteCommand(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete ressources",
	}
	cmd.AddCommand(makeDeleteLayerCommand(a))
	return cmd
}

func makeDeleteLayerCommand(a *app.App) *cobra.Command {

	opts := app.DeleteLayerOptions{}

	cmd := &cobra.Command{
		Use:   "layer",
		Short: "deletes a layer by name",
		PreRunE: func(cmd *cobra.Command, args []string) error {

			err := opts.Validate()
			if err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			a.DeleteLayer(&opts)

			fmt.Printf("Removed layer: %s from drawing", opts.Layername)
			return nil
		},
	}
	cmd.Flags().StringVarP(&opts.OutputFilename, "output", "o", "", "output file and path name [default export.xml]")
	cmd.Flags().StringVarP(&opts.Filename, "file", "f", "", "draw io model to import")
	cmd.Flags().StringVarP(&opts.Layername, "name", "n", "", "the layer name to delete from file")
	return cmd
}

func makeExtractCommand(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "exports a ressource",
	}
	cmd.AddCommand(makeExtractLayerCommand(a))
	return cmd
}

func makeExtractLayerCommand(a *app.App) *cobra.Command {

	opts := app.ExtractLayerOptions{}
	cmd := &cobra.Command{
		Use:   "layer",
		Short: "exports a layer by name",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := opts.Validate()
			if err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			err := a.ExtractLayer(&opts)
			if err != nil {
				return err
			}

			fmt.Printf("Extracted layer: %s successful into a new file: %s\n", opts.Layername, opts.OutputFile)
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.OutputFile, "output", "o", "", "output file and path name [default export.xml]")
	cmd.Flags().StringVarP(&opts.Filename, "file", "f", "", "draw io model to import")
	cmd.Flags().StringVarP(&opts.Layername, "name", "n", "", "the layer name to extract into a new file")
	return cmd
}

func makeShowLayersCommand(a *app.App) *cobra.Command {
	opts := app.ShowLayersOptions{}
	cmd := cobra.Command{
		Use:   "layers",
		Short: "show layers displays all available layers",
		Long:  "show layers displays all available layers",
		PreRunE: func(cmd *cobra.Command, args []string) error {

			err := opts.Validate()
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			err := a.ShowLayers(&opts)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.Filename, "file", "f", "", "draw io model to import")

	return &cmd
}

func makeClassifyCommand(a *app.App) *cobra.Command {

	cmd := cobra.Command{
		Use:   "classify",
		Short: "classifies a drawing with a watermark",
	}

	cmd.AddCommand(makeClassifyAsDraft(a))
	return &cmd
}

func makeClassifyAsDraft(a *app.App) *cobra.Command {

	var filename string
	var targetFilename string
	cmd := cobra.Command{
		Use:   "draft",
		Short: "classifies the document as draft",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(filename) == 0 {
				return fmt.Errorf("Error please provide a valid draw io file")
			}
			if len(targetFilename) == 0 {
				targetFilename = "export.xml"
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			a.ImportDrawing(filename)
			a.Classify(targetFilename)
			return nil
		},
	}
	cmd.Flags().StringVarP(&filename, "file", "f", "", "draw io model to import")
	cmd.Flags().StringVarP(&targetFilename, "output", "o", "", "output file and path name [default export.xml]")
	return &cmd
}

func makeUnpackCommand(a *app.App) *cobra.Command {
	var filename string
	var targetFilename string
	cmd := cobra.Command{
		Use:   "unpack",
		Short: "unpacks mxfiles and extracts the raw xml",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(filename) == 0 {
				return fmt.Errorf("Error please provide a valid draw io file")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := a.UnpackFile(filename, targetFilename)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&filename, "file", "f", "", "draw io model to import")
	cmd.Flags().StringVarP(&targetFilename, "output", "o", "", "output file and path name. If not provided it will be printed to the console")
	return &cmd

}

func makeMergeCommand(a *app.App) *cobra.Command {

	var inputFilename string
	var outputFilename string
	var mergeObjectFilename string

	cmd := cobra.Command{
		Use:   "merge",
		Short: "merges another drawing ont the imported file",
		PreRunE: func(cmd *cobra.Command, args []string) error {

			if len(inputFilename) == 0 {
				return fmt.Errorf("Please provide a filename to import the drawing")
			}
			if len(mergeObjectFilename) == 0 {
				return fmt.Errorf("Please provide a filename of the file to be merged")
			}

			if len(outputFilename) == 0 {
				outputFilename = "export.xml"
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			err := a.ImportDrawing(inputFilename)
			if err != nil {
				return err
			}

			err = a.Merge(mergeObjectFilename, outputFilename)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputFilename, "file", "f", "", "draw io model to import")
	cmd.Flags().StringVarP(&outputFilename, "output", "o", "", "output file and path name [default export.xml]")
	cmd.Flags().StringVarP(&mergeObjectFilename, "merge-file", "m", "", "the file to be onto the the imported file")
	return &cmd
}

func makeVersionCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "version",
		Short: "shows version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("snipio %s\n", versionString)
		},
	}

	return &cmd
}
