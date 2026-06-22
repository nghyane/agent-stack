// Package settings (feature) exposes the runtime operational settings over HTTP.
// It is a thin transport over the platform settings.Manager: GET reads the live
// snapshot, PUT mutates it (and the change propagates to all instances without
// a restart). Mutating settings requires an authenticated session.
package settings

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nghiahoang/template-api/internal/platform/session"
	platsettings "github.com/nghiahoang/template-api/internal/platform/settings"
)

// Handler is the transport layer for reading/updating operational settings.
type Handler struct {
	mgr      *platsettings.Manager
	sessions *session.Manager
}

// NewHandler wires the feature over the platform settings manager.
func NewHandler(mgr *platsettings.Manager, sessions *session.Manager) *Handler {
	return &Handler{mgr: mgr, sessions: sessions}
}

// RegisterRoutes mounts the settings endpoints.
func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "getSettings",
		Method:      http.MethodGet,
		Path:        "/settings",
		Summary:     "Get current operational settings",
		Tags:        []string{"settings"},
	}, h.get)

	huma.Register(api, huma.Operation{
		OperationID: "updateSettings",
		Method:      http.MethodPut,
		Path:        "/settings",
		Summary:     "Update operational settings (no restart)",
		Tags:        []string{"settings"},
	}, h.update)
}

// SettingsResponse is the public shape of the operational settings.
type SettingsResponse struct {
	AIProvider string `json:"aiProvider"`
	AIModel    string `json:"aiModel"`
	AIEnabled  bool   `json:"aiEnabled"`
}

func responseFrom(s platsettings.Settings) SettingsResponse {
	return SettingsResponse{
		AIProvider: s.AIProvider,
		AIModel:    s.AIModel,
		AIEnabled:  s.AIEnabled,
	}
}

type settingsOutput struct {
	Body SettingsResponse
}

type updateInput struct {
	Body struct {
		AIProvider string `json:"aiProvider" minLength:"1" maxLength:"50"`
		AIModel    string `json:"aiModel" minLength:"1" maxLength:"100"`
		AIEnabled  bool   `json:"aiEnabled"`
	}
}

func (h *Handler) get(ctx context.Context, _ *struct{}) (*settingsOutput, error) {
	return &settingsOutput{Body: responseFrom(h.mgr.Current())}, nil
}

func (h *Handler) update(ctx context.Context, in *updateInput) (*settingsOutput, error) {
	if h.sessions.UserID(ctx) == "" {
		return nil, huma.Error401Unauthorized("not authenticated")
	}

	err := h.mgr.Update(ctx, func(s *platsettings.Settings) {
		s.AIProvider = in.Body.AIProvider
		s.AIModel = in.Body.AIModel
		s.AIEnabled = in.Body.AIEnabled
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("update settings", err)
	}
	return &settingsOutput{Body: responseFrom(h.mgr.Current())}, nil
}
