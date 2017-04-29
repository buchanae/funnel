package util

type CallList []func() error

func (dl *CallList) Add(f func() error) {
	*dl = append(dl, f)
}
func (dl *CallList) AddUnchecked(f func()) {
  dl.Add(func() error {
    f()
    return nil
  })
}

func (dl *CallList) Run() error {
	for _, f := range dl {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}
