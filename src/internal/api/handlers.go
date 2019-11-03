package api

import (
	"net/http"

	"github.com/go-openapi/swag"
	"github.com/howtocards/preview-queue/src/internal/api/generated/models"
	"github.com/howtocards/preview-queue/src/internal/api/generated/restapi/operations"
	"github.com/howtocards/preview-queue/src/internal/queue"
)

func (api *service) render(params operations.RenderParams) operations.RenderResponder {
	_, log := fromRequest(params.HTTPRequest)

	event := queue.Event{
		URL:      string(params.Args.URL),
		Selector: string(params.Args.Selector),
		Callback: string(params.Args.Callback),
		AppID:    string(params.Args.AppID),
		UserID:   string(params.Args.UserID),
	}
	log.Info("params", event)

	err := api.queue.Send(event)
	if err != nil {
		return operations.NewRenderDefault(http.StatusInternalServerError).
			WithPayload(&models.Error{
				Code:    swag.Int32(http.StatusInternalServerError),
				Message: swag.String(err.Error()),
			})
	}

	return operations.NewRenderOK()
}
