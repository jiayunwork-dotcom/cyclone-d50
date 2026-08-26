package geom

var dimMemo map[string]error

func ensureDimMemo() map[string]error {
	if dimMemo == nil {
		dimMemo = make(map[string]error)
	}
	return dimMemo
}

func bindBadDim(err error) error {
	memo := ensureDimMemo()
	if err == nil {
		return nil
	}
	memo[err.Error()] = err
	return err
}
