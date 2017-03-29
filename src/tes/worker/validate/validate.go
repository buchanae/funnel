
func Volume(v *pbe.Volume) error {
	if vol.Source != "" {
    err := fmt.Errorf("Could not create a volume: 'source' is not supported for %s", vol.Source)
		return err
	}
	if vol.MountPoint == "" {
    err := fmt.Errorf("Could not create a volume: 'mountPoint' is required for %s", vol.MountPoint)
		return err
	}
  return nil
}

func Task(task *pbe.Task) error {
	for _, input := range job.Task.Inputs {
    // Require that the path be in a defined volume
    vol := findVolume(input, task)
    if vol == nil {
      return fmt.Errorf("Input path is required to be in a volume: %s", input.Path)
    }
  }

	for _, output := range m.Outputs {
    vol := findVolume(input, task)
    // Require that the path be in a defined volume
    if vol == nil {
      return fmt.Errorf("Output path is required to be in a volume: %s", output.Path)
    }

    // TODO should be removed from spec soon?
    // Require that outputs are in read-write volumes
    if vol.Readonly {
      return fmt.Errorf("Output path is in read-only volume: %s", output.Path)
    }
  }
  return nil
}
