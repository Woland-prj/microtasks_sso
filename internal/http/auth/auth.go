package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/Woland-prj/microtasks_sso/internal/domain/cerrors"
	"github.com/Woland-prj/microtasks_sso/internal/domain/dtos"
	"github.com/Woland-prj/microtasks_sso/internal/domain/entities"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type AuthService interface {
	Login(
		ctx context.Context,
		dto dtos.LoginDto,
	) (*entities.JwtTokenPair, error)
	Register(
		ctx context.Context,
		dto dtos.RegisterDto,
	) (int64, error)
	Refresh(
		ctx context.Context,
		dto dtos.RefreshDto,
	) (*entities.JwtTokenPair, error)
}

type serverAPI struct {
	authService AuthService
	validate    *validator.Validate
}

func Register(
	router *chi.Mux, 
	service AuthService, 
	validate *validator.Validate,
) {
	api := serverAPI{authService: service, validate: validate}
	router.Post("/login", api.Login())
	router.Post("/register", api.Register())
	router.Get("/refresh", api.Refresh())
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	AppId    int64  `json:"app_id" validate:"required"`
}

type LoginResponse struct {
	AuthToken    string `json:"auth_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Error        string `json:"error,omitempty"`
}

func (api *serverAPI) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, LoginResponse{Error: "Invalid request"})
			return
		}

		err = api.validate.Struct(req)
		if err != nil {
			render.JSON(w, r, LoginResponse{Error: "Invalid credentials"})
			return
		}

		tokens, err := api.authService.Login(r.Context(), dtos.LoginDto{
			Email:    req.Email,
			Password: req.Password,
			AppId:    req.AppId,
		})

		if err != nil {
			if errors.Is(err, &cerrors.InvalidCredentialsError{}) {
				render.JSON(w, r, LoginResponse{Error: "Invalid credentials"})
				return
			}
			render.JSON(w, r, LoginResponse{Error: "Internal error"})
			return
		}

		render.JSON(w, r, LoginResponse{
			AuthToken:    tokens.AuthToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	Uid int64 `json:"uid,omitempty"`
	Error string `json:"error,omitempty"`
}

func (api *serverAPI) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, RegisterResponse{Error: "Invalid request"})
			return
		}

		err = api.validate.Struct(req)
		if err != nil {
			render.JSON(w, r, RegisterResponse{Error: "Invalid credentials"})
			return
		}

		uid, err := api.authService.Register(r.Context(), dtos.RegisterDto{
			Email:    req.Email,
			Password: req.Password,
		})

		if err != nil {
			if errors.Is(err, &cerrors.AlreadyExistsError{}) {
				render.JSON(w, r, RegisterResponse{Error: "User already exists"})
				return
			}
			render.JSON(w, r, RegisterResponse{Error: "Internal error"})
			return
		}

		render.JSON(w, r, RegisterResponse{Uid: uid})
	}
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,jwt"`
	AppId        int64  `json:"app_id" validate:"required"`
}

func (api *serverAPI) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RefreshRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, LoginResponse{Error: "Invalid request"})
			return
		}

		err = api.validate.Struct(req)
		if err != nil {
			render.JSON(w, r, LoginResponse{Error: "Bad format"})
			return
		}

		tokens, err := api.authService.Refresh(r.Context(), dtos.RefreshDto{
			RefreshToken: req.RefreshToken,
			AppId:        req.AppId,
		})

		if err != nil {
			var cErr cerrors.InvalidTokenError
			if errors.As(err, &cErr) {
				switch(cErr.Subject()) {
					case cerrors.TokenExpired:
						render.JSON(w, r, LoginResponse{Error: "Token expired"})
						return
					case cerrors.TokenBadFormat:
						render.JSON(w, r, LoginResponse{Error: "Fake token"})
						return
				}
			}
			render.JSON(w, r, LoginResponse{Error: "Internal error"})
			return
		}

		render.JSON(w, r, LoginResponse{
			AuthToken:    tokens.AuthToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}