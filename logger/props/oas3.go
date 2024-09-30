package props

import "go.uber.org/zap"

func Oas3OperationId(operationId string) zap.Field {
	return zap.String("oas3.operation.id", operationId)
}
