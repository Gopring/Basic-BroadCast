package request

import "context"

type requestKey struct{}

func With(ctx context.Context, r *Request) context.Context {
	return context.WithValue(ctx, requestKey{}, r)
}

func From(ctx context.Context) *Request {
	if ctx == nil {
		return nil
	}

	r, ok := ctx.Value(requestKey{}).(*Request)
	if !ok {
		return nil
	}
	return r
}
