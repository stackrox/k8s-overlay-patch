/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/stackrox/k8s-overlay-patch/pkg/patch"
	"github.com/stackrox/k8s-overlay-patch/pkg/types"
	"io"
	"os"
	"sigs.k8s.io/yaml"

	"github.com/spf13/cobra"
)

var patchFilePath string
var manifestFilePath string
var namespace string
var outFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-overlay-patch",
	Short: "Applies overlays to rendered k8s manifests",
	RunE: func(cmd *cobra.Command, args []string) error {
		patchFile, err := os.Open(patchFilePath)
		if err != nil {
			return err
		}
		defer patchFile.Close()

		var manifestFile *os.File
		if manifestFilePath != "" {
			manifestFile, err = os.Open(manifestFilePath)
			if err != nil {
				return err
			}
			defer manifestFile.Close()
		} else {
			manifestFile = os.Stdin
		}

		patchBytes, err := io.ReadAll(patchFile)
		if err != nil {
			return err
		}

		manifestBytes, err := io.ReadAll(manifestFile)
		if err != nil {
			return err
		}

		var overlayObj types.OverlayObject
		if err := yaml.Unmarshal(patchBytes, &overlayObj); err != nil {
			return err
		}

		result, err := patch.YAMLManifestPatch(string(manifestBytes), namespace, overlayObj.Overlays)
		if err != nil {
			return err
		}

		var out = cmd.OutOrStdout()
		if outFile != "" {
			outF, err := os.Create(outFile)
			if err != nil {
				return err
			}
			defer outF.Close()
			out = outF
		}

		_, err = out.Write([]byte(result))
		return err

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
	rootCmd.Flags().StringVarP(&patchFilePath, "patch-file", "p", "", "File containing the patch to apply")
	rootCmd.Flags().StringVarP(&manifestFilePath, "manifest-file", "m", "", "File containing the rendered manifests to patch")
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace to use when patching the manifests")
	rootCmd.Flags().StringVarP(&outFile, "out", "o", "", "File to write the patched manifests to")
}
