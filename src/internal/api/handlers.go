package api

import (
	"net/http"

	"github.com/go-openapi/swag"
	"github.com/howtocards/preview-queue/src/internal/api/generated/models"
	"github.com/howtocards/preview-queue/src/internal/api/generated/restapi/operations"
	"github.com/howtocards/preview-queue/src/internal/queue"
)

func (api *service) renderCard(params operations.RenderCardParams) operations.RenderCardResponder {
	_, log := fromRequest(params.HTTPRequest)

	event := queue.Card{
		Card: string(params.Body.Card),
		Event: queue.Event{
			Extra:    params.Body.Extra,
			Callback: string(params.Body.Callback),
		},
	}
	log.Info("params", event)

	err := api.queue.Send(event, swag.StringValue(params.AppName))
	if err != nil {
		return operations.NewRenderCardDefault(http.StatusInternalServerError).
			WithPayload(&models.Error{
				Code:    swag.Int32(http.StatusInternalServerError),
				Message: swag.String(err.Error()),
			})
	}

	return operations.NewRenderCardOK()
}

func (api *service) renderUser(params operations.RenderUserParams) operations.RenderUserResponder {
	_, log := fromRequest(params.HTTPRequest)

	event := queue.User{
		User: string(params.Body.User),
		Event: queue.Event{
			Extra:    params.Body.Extra,
			Callback: string(params.Body.Callback),
		},
	}
	log.Info("params", event)

	err := api.queue.Send(event, swag.StringValue(params.AppName))
	if err != nil {
		return operations.NewRenderUserDefault(http.StatusInternalServerError).
			WithPayload(&models.Error{
				Code:    swag.Int32(http.StatusInternalServerError),
				Message: swag.String(err.Error()),
			})
	}

	return operations.NewRenderUserOK()
}
