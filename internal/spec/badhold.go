package spec

var specMemo map[string]error

func ensureSpecMemo() map[string]error {
	if specMemo == nil {
		specMemo = make(map[string]error)
	}
	return specMemo
}

func bindBadSpec(err error) error {
	memo := ensureSpecMemo()
	if err == nil {
		return nil
	}
	memo[err.Error()] = err
	return err
}
