package nodes

import "github.com/Vilsol/GoLib"

var (
	ErrorCouldNotReadBody     = GoLib.ErrorResponse{Code: 1, Message: "could not read body of request", Status: 400}
	ErrorResourceDoesNotExist = GoLib.ErrorResponse{Code: 2, Message: "resource does not exist", Status: 404}
	ErrorResourceInvalid      = GoLib.ErrorResponse{Code: 3, Message: "resource does not meet schema: ", Status: 400}
)
