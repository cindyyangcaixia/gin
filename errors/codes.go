package errors

import "fmt"

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d, message=%s, error=%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, err error, msg string) *AppError {
	message, ok := ErrorMessages[code]
	if !ok {
		message = "Unknown error"
	}
	return &AppError{
		Code:    code,
		Message: message + msg,
		Err:     err,
	}
}

const (
	// General Errors (90000-90099)
	ErrCodeInvalidRequest = 90001 // Invalid request payload
	ErrCodeInternalServer = 90002 // Unexpected server error

	ErrCodeValidation       = 90003
	ErrCodeInvalidParam     = 90004
	ErrCodeMethodNotAllowed = 90005

	// Authentication Errors (90100-90199)
	ErrCodeInvalidCredentials = 90100 // Invalid username or password
	ErrCodeTokenExpired       = 90101 // JWT token expired
	ErrCodeInvalidToken       = 90102 // Invalid JWT token

	// User Phone Errors (90200-90299)
	ErrCodeInvalidSerial      = 90201 // Invalid serial number
	ErrCodeInvalidPhoneNumber = 90202 // Invalid phone number
	ErrCodeUserNotExist       = 90203

	// Database Errors (90300-90399)
	ErrCodeDatabaseQueryFailed = 90300 // Database query failed
	ErrCodeRecordNotFound      = 90301 // Record not found
)

// ErrorMessages maps error codes to human-readable messages
var ErrorMessages = map[int]string{
	ErrCodeInvalidRequest:      "Invalid request payload",
	ErrCodeInternalServer:      "Internal server error",
	ErrCodeValidation:          "Validation failed",
	ErrCodeInvalidParam:        "Invalid params",
	ErrCodeMethodNotAllowed:    "Method not allowed",
	ErrCodeInvalidCredentials:  "Invalid username or password",
	ErrCodeTokenExpired:        "Token has expired",
	ErrCodeInvalidToken:        "Invalid token",
	ErrCodeInvalidSerial:       "Invalid serial number",
	ErrCodeInvalidPhoneNumber:  "Invalid phone number",
	ErrCodeDatabaseQueryFailed: "Database query failed",
	ErrCodeRecordNotFound:      "Record not found",
	ErrCodeUserNotExist:        "user not exits",
}
