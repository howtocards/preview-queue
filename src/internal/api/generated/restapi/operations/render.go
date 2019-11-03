// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// RenderHandlerFunc turns a function with the right signature into a render handler
type RenderHandlerFunc func(RenderParams) RenderResponder

// Handle executing the request and returning a response
func (fn RenderHandlerFunc) Handle(params RenderParams) RenderResponder {
	return fn(params)
}

// RenderHandler interface for that can handle valid render params
type RenderHandler interface {
	Handle(RenderParams) RenderResponder
}

// NewRender creates a new http.Handler for the render operation
func NewRender(ctx *middleware.Context, handler RenderHandler) *Render {
	return &Render{Context: ctx, Handler: handler}
}

/*Render swagger:route POST /render render

Sending a message to the queue service.

*/
type Render struct {
	Context *middleware.Context
	Handler RenderHandler
}

func (o *Render) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewRenderParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
