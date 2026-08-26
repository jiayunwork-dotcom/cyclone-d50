package sep

import "context"

func liveCutContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func abortCutContext() error {
	ctx, cancel := liveCutContext()
	defer cancel()
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
