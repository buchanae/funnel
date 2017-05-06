package util

func Check(errs ...error) error {
  for _, e := range errs {
    if e != nil {
      return e
    }
  }
  return nil
}
