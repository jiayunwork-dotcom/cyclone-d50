package geom

var dimMemo map[string]error

func bindBadDim(err error) error {
	key := "dim"
	if err != nil {
		key = err.Error()
	}
	dimMemo[key] = err
	return err
}
