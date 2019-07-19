package kubernetes

import (
	"fmt"
	"strings"

	"github.com/weaveworks/flux/image"
	"github.com/weaveworks/flux/resource"
)

// updateContainer takes a YAML document stream (one or more YAML
// docs, as bytes), a resource ID referring to a controller, a
// container name, and the name of the new image that should be used
// for the container. It returns a new YAML stream where the image for
// the container has been replaced with the imageRef supplied.
func updateContainer(in []byte, resource resource.ID, container string, newImageID image.Ref) ([]byte, error) {
	namespace, kind, name := resource.Components()
	if _, ok := resourceKinds[strings.ToLower(kind)]; !ok {
		return nil, UpdateNotSupportedError(kind)
	}
	return (KubeYAML{}).Image(in, namespace, kind, name, container, newImageID.String())
}

// updatePaths takes a YAML document stream (one or more YAML
// docs, as bytes), as resource ID referring to a controller,
// image paths, and the new image that should be applied to
// those paths. It returns a new YAML stream where the values
// of the paths have been replaced with the imageRef supplied.
func updatePaths(in []byte, resource resource.ID, paths resource.ImagePath, newImageID image.Ref) ([]byte, error) {
	namespace, kind, name := resource.Components()
	if _, ok := resourceKinds[strings.ToLower(kind)]; !ok {
		return nil, UpdateNotSupportedError(kind)
	}
	var args []string
	for k, v := range paths.Map(newImageID) {
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}
	return (KubeYAML{}).Set(in, namespace, kind, name, args...)
}
