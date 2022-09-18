package errors

const (
	ErrDBConnect       = "Failed to open connection to the database!"
	ErrInvalidLogLevel = "Unsupported log level!"
	ErrConfigFileRead  = "Error reading config file!"
	ErrContextRun      = "Error executing context!"
)

type GenericError struct {
	Code string
}

func (ge *GenericError) Error() string {
	return ge.Code
}
