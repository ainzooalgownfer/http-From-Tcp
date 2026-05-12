package response

type Response struct {

}

type StatusCode int 
const (
	Statusok StatusCode = 200
	StatusBadRequest StatusCode = 400
	StatusNotFound StatusCode = 404
	StatusInternalServerError StatusCode = 500
)