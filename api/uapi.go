package api

import (
	"Popkat/constants"
	"Popkat/state"
	"Popkat/types"
	"net/http"

	"github.com/infinitybotlist/eureka/uapi"
)

type DefaultResponder struct{}

func (d DefaultResponder) New(err string, ctx map[string]string) any {
	return types.ApiError{
		Message: err,
		Context: ctx,
	}
}

const (
	TargetTypeUser   = "User"
	TargetTypeServer = "Server"
)

// Authorizes a request
func Authorize(r uapi.Route, req *http.Request) (uapi.AuthData, uapi.HttpResponse, bool) {
	return uapi.AuthData{}, uapi.HttpResponse{}, true
}

func Setup() {
	uapi.SetupState(uapi.UAPIState{
		Logger:  state.Logger,
		Context: state.Context,
		Constants: &uapi.UAPIConstants{
			ResourceNotFound:    constants.ResourceNotFound,
			BadRequest:          constants.BadRequest,
			Forbidden:           constants.Forbidden,
			Unauthorized:        constants.Unauthorized,
			InternalServerError: constants.InternalServerError,
			MethodNotAllowed:    constants.MethodNotAllowed,
			BodyRequired:        constants.BodyRequired,
		},
		DefaultResponder: DefaultResponder{},
		Authorize:        Authorize,
		AuthTypeMap: map[string]string{
			TargetTypeUser:   "user",
			TargetTypeServer: "server",
		},
	})
}
