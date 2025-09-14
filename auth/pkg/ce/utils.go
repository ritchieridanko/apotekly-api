package ce

import "net/http"

func MapToExternalErrorCode(internalCode int) (externalCode int) {
	switch internalCode {
	case ErrCodeInvalidPayload, ErrCodeInvalidParams:
		return http.StatusBadRequest
	case ErrCodeInvalidAction, ErrCodeDBNoChange:
		return http.StatusUnauthorized
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeLocked:
		return http.StatusLocked
	default:
		return http.StatusInternalServerError
	}
}
