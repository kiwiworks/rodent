package props

import (
	"fmt"

	"go.uber.org/zap"
)

func AuthUserId(id fmt.Stringer) zap.Field {
	return zap.String("auth.user.id", id.String())
}
