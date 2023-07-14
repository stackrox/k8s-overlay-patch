package cmd

import (
	"bytes"
	"github.com/stackrox/k8s-overlay-patch/pkg/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRoot(t *testing.T) {

	rootCmd.SetArgs([]string{
		"-n",
		"test-namespace",
		"-p",
		"../pkg/testdata/patch.yaml",
		"-m",
		"../pkg/testdata/manifest.yaml",
	})

	var wr = bytes.NewBufferString("")
	rootCmd.SetOut(wr)
	err := rootCmd.Execute()
	require.NoError(t, err)

	// parse the output
	outManifest := wr.String()
	objs, err := object.ParseK8sObjectsFromYAMLManifest(outManifest)
	require.NoError(t, err)

	objMap := objs.ToNameKindMap()

	t.Log(objMap)

	deploymentU, ok := objMap["Deployment:test-deployment"]
	require.True(t, ok)
	assert.Equal(t, "annotation", deploymentU.UnstructuredObject().GetAnnotations()["my"])

	serviceU, ok := objMap["Service:test-service"]
	require.True(t, ok)
	assert.Equal(t, "annotation", serviceU.UnstructuredObject().GetAnnotations()["my"])

}
