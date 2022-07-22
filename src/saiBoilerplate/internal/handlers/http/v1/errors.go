package v1

type ServiceErr struct {
	Code    string `json:"code" example:"ERROR_CODE"`
	Message string `json:"message" example:"error description"`
}

var (
	errInternalServerErr = &ServiceErr{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "internal server error",
	}
)
