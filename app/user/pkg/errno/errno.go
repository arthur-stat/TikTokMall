package errno

type ErrNo struct {
	Code    int32
	Message string
}

func (e ErrNo) Error() string {
	return e.Message
}

func NewErrNo(code int32, msg string) *ErrNo {
	return &ErrNo{
		Code:    code,
		Message: msg,
	}
}

var (
	Success = NewErrNo(0, "Success")
	
	// Common errors
	ServiceErr = NewErrNo(10001, "Service is unable to start successfully")
	ParamErr   = NewErrNo(10002, "Wrong Parameter has been given")
	
	// Auth errors
	AuthorizationFailedErr = NewErrNo(20001, "Authorization failed")
	
	// User errors
	UserAlreadyExistErr     = NewErrNo(30001, "User already exists")
	UserNotExistErr         = NewErrNo(30002, "User does not exist")
	PasswordIncorrectErr    = NewErrNo(30003, "Password is incorrect")
	PasswordsDoNotMatchErr  = NewErrNo(30004, "Passwords do not match")
)