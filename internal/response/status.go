package response

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func getStatusReason(statusCode StatusCode) string {
	var reason string
	switch statusCode {
	case StatusOK:
		reason = "OK"
	case StatusBadRequest:
		reason = "Bad Request"
	case StatusInternalServerError:
		reason = "Internal Server Error"
	}

	return reason
}
