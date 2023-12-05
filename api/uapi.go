package api

import (
	"Popkat/constants"
	"Popkat/state"
	"Popkat/types"
	"net/http"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/infinitybotlist/eureka/uapi"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type DefaultResponder struct{}

func (d DefaultResponder) New(err string, ctx map[string]string) any {
	return types.ApiError{
		Message: err,
		Context: ctx,
	}
}

const (
	TargetTypeUser    = "User"
	TargetTypeService = "Service"
)

// Authorizes a request
func Authorize(r uapi.Route, req *http.Request) (uapi.AuthData, uapi.HttpResponse, bool) {
	authHeader := req.Header.Get("Authorization")

	if len(r.Auth) > 0 && authHeader == "" && !r.AuthOptional {
		return uapi.AuthData{}, uapi.DefaultResponse(http.StatusUnauthorized), false
	}

	authData := uapi.AuthData{}

	for _, auth := range r.Auth {
		if authData.Authorized {
			break
		}

		if authHeader == "" {
			continue
		}

		var urlIds []string

		switch auth.Type {
		case TargetTypeUser:
			var id pgtype.Text
			var banned bool

			err := state.Pool.QueryRow(state.Context, "SELECT user_id, banned FROM users WHERE token = $1", strings.Replace(authHeader, "User ", "", 1)).Scan(&id, &banned)

			if err != nil {
				continue
			}

			if !id.Valid {
				continue
			}

			authData = uapi.AuthData{
				TargetType: TargetTypeUser,
				ID:         id.String,
				Authorized: true,
				Banned:     banned,
			}
			urlIds = []string{id.String}
		case TargetTypeService:
			var id pgtype.Text

			err := state.Pool.QueryRow(state.Context, "SELECT id FROM services WHERE token = $1", strings.Replace(authHeader, "Service ", "", 1)).Scan(&id)

			if err != nil {
				continue
			}

			if !id.Valid {
				continue
			}

			authData = uapi.AuthData{
				TargetType: TargetTypeService,
				ID:         id.String,
				Authorized: true,
			}
			urlIds = []string{id.String}
		}

		if auth.URLVar != "" {
			state.Logger.Info("Checking URL variable against user ID from auth token", zap.String("URLVar", auth.URLVar))
			gotUserId := chi.URLParam(req, auth.URLVar)
			if !slices.Contains(urlIds, gotUserId) {
				authData = uapi.AuthData{}
			}
		}

		if authData.Banned && auth.AllowedScope != "ban_exempt" {
			return uapi.AuthData{}, uapi.HttpResponse{
				Status: http.StatusForbidden,
				Json:   types.ApiError{Message: "Hello, there. The Popkat CDN Account associated with this token has been terminated due to being in violation with our Terms of Service and/or Community Guidelines. If you think this was a mistake, please contact support by emailing 'contact@select-list.xyz'."},
			}, false
		}
	}

	if len(r.Auth) > 0 && !authData.Authorized && !r.AuthOptional {
		return uapi.AuthData{}, uapi.DefaultResponse(http.StatusUnauthorized), false
	}

	return authData, uapi.HttpResponse{}, true
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
			TargetTypeUser:    "user",
			TargetTypeService: "service",
		},
		Redis: state.Redis,
	})
}
