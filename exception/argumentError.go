package exception

func BuildArgumentError() *ArgumentError{
	return &ArgumentError{BaseCodeError: BaseCodeError{ErrorCode: error_code_argument_error, OutPutMessage: "参数错误",}}
}

type ArgumentError struct {
	BaseCodeError
	ParamError  error
	ParamErrorMsg string
}

func (runtime *ArgumentError) SetParamError(err error) *ArgumentError {
	runtime.ParamError = err
	return runtime
}

func (runtime *ArgumentError) SetParamErrorMsg(msg string) *ArgumentError {
	runtime.ParamErrorMsg = msg
	return runtime
}