package task

import (
	"context"
	"errors"
)

type ContextKey string

var UserIDCtxKey ContextKey = ContextKey("userId")

func getUserIDFromCtx(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDCtxKey).(string)
	if !ok {
		return "", errors.New("user id not found in context")
	}

	return userID, nil
}
