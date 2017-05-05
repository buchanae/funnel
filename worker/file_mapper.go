package worker

import (
	"fmt"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"path"
)


func TaskVolumes(task *tes.Task) []Volume {
  var volumes []Volume

	for _, input := range task.Inputs {
    volumes = append(volumes, Volume{
      ContainerPath: input.Path,
      Readonly:      true,
    }
	}

	for _, vol := range task.Volumes {
    volumes = append(volumes, Volume{
      ContainerPath: vol,
      Readonly:      false,
    }
	}

	for _, output := range task.Outputs {
    p := output.Path
    if output.Type == tes.FileType_FILE {
      p = path.Dir(p)
    }

    volumes = append(volumes, Volume{
      ContainerPath: p,
      Readonly: false,
    }
	}

  return volumes
}
