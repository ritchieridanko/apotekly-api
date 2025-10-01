package ce

import (
	"database/sql"
	"errors"
)

type internalErrorCode string

// internal error codes (for logs/debugging)
const (
	CodeContextValueNotFound internalErrorCode = "CONTEXT_VALUE_NOT_FOUND_ERROR"
	CodeDBDuplicateData      internalErrorCode = "DB_DUPLICATE_DATA_ERROR"
	CodeDBQueryExecution     internalErrorCode = "DB_QUERY_EXECUTION_ERROR"
	CodeDBTransaction        internalErrorCode = "DB_TRANSACTION_ERROR"
	CodeFileBuffer           internalErrorCode = "FILE_BUFFER_ERROR"
	CodeFileUploadFailed     internalErrorCode = "FILE_UPLOAD_FAILED_ERROR"
	CodeInvalidParams        internalErrorCode = "INVALID_PARAMS_ERROR"
	CodeInvalidPayload       internalErrorCode = "INVALID_PAYLOAD_ERROR"
	CodeRequestFile          internalErrorCode = "REQUEST_FILE_ERROR"
)

// external error messages (for end-users)
const (
	MsgInternalServer  string = "Internal server error."
	MsgInvalidParams   string = "Invalid params."
	MsgInvalidPayload  string = "Invalid payload."
	MsgUnauthenticated string = "Unauthenticated."
)

// internal error logs
var (
	ErrDBAffectNoRows error = errors.New("query execution affected no rows")
	ErrDBQueryNoRows  error = sql.ErrNoRows
)
