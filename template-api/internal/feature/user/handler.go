package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nghiahoang/template-api/internal/feature/user/data"
	"github.com/nghiahoang/template-api/internal/platform/session"
)

// Handler is the transport layer for the user feature. It owns Huma I/O types,
// validation (via struct tags), and mapping service errors to HTTP responses.
type Handler struct {
	svc *Service
}

// NewHandler wires the feature: queries + service + handler.
func NewHandler(pool *pgxpool.Pool, sessions *session.Manager) *Handler {
	svc := NewService(data.New(pool), sessions)
	return &Handler{svc: svc}
}

// RegisterRoutes mounts every operation of this feature onto the API.
// Adding a feature = add a folder like this and call its RegisterRoutes once.
func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "register",
		Method:        http.MethodPost,
		Path:          "/auth/register",
		Summary:       "Register a new user",
		Tags:          []string{"auth"},
		DefaultStatus: http.StatusCreated,
	}, h.register)

	huma.Register(api, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/auth/login",
		Summary:     "Log in",
		Tags:        []string{"auth"},
	}, h.login)

	huma.Register(api, huma.Operation{
		OperationID: "logout",
		Method:      http.MethodPost,
		Path:        "/auth/logout",
		Summary:     "Log out",
		Tags:        []string{"auth"},
	}, h.logout)

	huma.Register(api, huma.Operation{
		OperationID: "me",
		Method:      http.MethodGet,
		Path:        "/auth/me",
		Summary:     "Get current user",
		Tags:        []string{"auth"},
	}, h.me)
}

// --- I/O types (validation lives in struct tags) ---

type registerInput struct {
	Body struct {
		Email    string `json:"email" format:"email" doc:"User email"`
		Password string `json:"password" minLength:"8" maxLength:"72" doc:"Password (8-72 chars)"`
		Name     string `json:"name" maxLength:"100" doc:"Display name"`
	}
}

type loginInput struct {
	Body struct {
		Email    string `json:"email" format:"email"`
		Password string `json:"password" minLength:"1"`
	}
}

type userOutput struct {
	Body UserResponse
}

type logoutOutput struct {
	Body struct {
		OK bool `json:"ok"`
	}
}

// --- Handlers (thin: parse, call service, map errors) ---

func (h *Handler) register(ctx context.Context, in *registerInput) (*userOutput, error) {
	dto, err := h.svc.Register(ctx, in.Body.Email, in.Body.Password, in.Body.Name)
	if err != nil {
		return nil, toHTTP(err)
	}
	return &userOutput{Body: dto}, nil
}

func (h *Handler) login(ctx context.Context, in *loginInput) (*userOutput, error) {
	dto, err := h.svc.Login(ctx, in.Body.Email, in.Body.Password)
	if err != nil {
		return nil, toHTTP(err)
	}
	return &userOutput{Body: dto}, nil
}

func (h *Handler) logout(ctx context.Context, _ *struct{}) (*logoutOutput, error) {
	if err := h.svc.Logout(ctx); err != nil {
		return nil, toHTTP(err)
	}
	out := &logoutOutput{}
	out.Body.OK = true
	return out, nil
}

func (h *Handler) me(ctx context.Context, _ *struct{}) (*userOutput, error) {
	dto, err := h.svc.Current(ctx)
	if err != nil {
		return nil, toHTTP(err)
	}
	return &userOutput{Body: dto}, nil
}

// toHTTP maps the feature's sentinel errors to HTTP responses. Unknown errors
// become 500 with the original cause attached for logging.
func toHTTP(err error) error {
	switch {
	case errors.Is(err, ErrEmailTaken):
		return huma.Error409Conflict(err.Error())
	case errors.Is(err, ErrInvalidLogin),
		errors.Is(err, ErrNotAuthed),
		errors.Is(err, ErrUserNotFound):
		return huma.Error401Unauthorized(err.Error())
	default:
		return huma.Error500InternalServerError("internal error", err)
	}
}
