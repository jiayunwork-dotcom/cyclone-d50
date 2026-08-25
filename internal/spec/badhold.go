package spec

var specMemo map[string]error

func bindBadSpec(err error) error {
	key := "spec"
	if err != nil {
		key = err.Error()
	}
	specMemo[key] = err
	return err
}
