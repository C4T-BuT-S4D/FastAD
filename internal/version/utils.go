package version

import versionpb "github.com/c4t-but-s4d/fastad/pkg/proto/data/version"

func NewVersionProto(version int) *versionpb.Version {
	return &versionpb.Version{
		Version: int64(version),
	}
}
