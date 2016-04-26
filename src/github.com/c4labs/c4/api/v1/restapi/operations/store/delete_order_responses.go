package store

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit"
)

/*DeleteOrderBadRequest Invalid ID supplied

swagger:response deleteOrderBadRequest
*/
type DeleteOrderBadRequest struct {
}

// NewDeleteOrderBadRequest creates DeleteOrderBadRequest with default headers values
func NewDeleteOrderBadRequest() *DeleteOrderBadRequest {
	return &DeleteOrderBadRequest{}
}

// WriteResponse to the client
func (o *DeleteOrderBadRequest) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(400)
}

/*DeleteOrderNotFound Order not found

swagger:response deleteOrderNotFound
*/
type DeleteOrderNotFound struct {
}

// NewDeleteOrderNotFound creates DeleteOrderNotFound with default headers values
func NewDeleteOrderNotFound() *DeleteOrderNotFound {
	return &DeleteOrderNotFound{}
}

// WriteResponse to the client
func (o *DeleteOrderNotFound) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(404)
}
