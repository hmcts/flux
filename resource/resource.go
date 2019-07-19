package resource

import (
	"github.com/weaveworks/flux/image"
	"github.com/weaveworks/flux/policy"
)

// For the minute we just care about
type Resource interface {
	ResourceID() ID       // name, to correlate with what's in the cluster
	Policies() policy.Set // policy for this resource; e.g., whether it is locked, automated, ignored
	Source() string       // where did this come from (informational)
	Bytes() []byte        // the definition, for sending to cluster.Sync
}

type ImagePath struct {
	Registry   string
	Repository string
	Tag        string
}

func (m ImagePath) Map(i image.Ref) map[string]string {
	var im map[string]string
	if m.Repository == ""  {
		return im
	}
	switch {
	case m.Registry != "" && m.Tag != "":
		im[m.Registry] = i.Domain
		im[m.Repository] = i.Image
		im[m.Tag] = i.Tag
	case m.Registry != "":
		im[m.Registry] = i.Domain
		im[m.Repository] = i.Image + ":" + i.Tag
	case m.Tag != "":
		im[m.Repository] = i.Name.String()
		im[m.Tag] = i.Tag
	default:
		im[m.Repository] = i.String()
	}
	return im
}

type Container struct {
	Name    string
	Image   image.Ref
	Mapping ImagePath
}

type Workload interface {
	Resource
	Containers() []Container
	// SetContainerImage mutates this workload so that the container
	// named has the image given. This is not expected to have an
	// effect on any underlying file or cluster resource.
	SetContainerImage(container string, ref image.Ref) error
}
