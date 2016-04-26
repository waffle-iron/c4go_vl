package events

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetEventByIDHandlerFunc turns a function with the right signature into a get event by Id handler
type GetEventByIDHandlerFunc func(GetEventByIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetEventByIDHandlerFunc) Handle(params GetEventByIDParams) middleware.Responder {
	return fn(params)
}

// GetEventByIDHandler interface for that can handle valid get event by Id params
type GetEventByIDHandler interface {
	Handle(GetEventByIDParams) middleware.Responder
}

// NewGetEventByID creates a new http.Handler for the get event by Id operation
func NewGetEventByID(ctx *middleware.Context, handler GetEventByIDHandler) *GetEventByID {
	return &GetEventByID{Context: ctx, Handler: handler}
}

/*GetEventByID swagger:route GET /events/{id} events getEventById

Get event by id.

*/
type GetEventByID struct {
	Context *middleware.Context
	Handler GetEventByIDHandler
}

func (o *GetEventByID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewGetEventByIDParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
