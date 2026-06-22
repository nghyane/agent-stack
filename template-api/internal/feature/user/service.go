package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nghiahoang/template-api/internal/feature/user/data"
	"github.com/nghiahoang/template-api/internal/platform/session"
)

// Sentinel errors. The handler maps these to HTTP status codes; the service
// stays transport-agnostic and never imports huma or net/http.
var (
	ErrEmailTaken   = errors.New("email already registered")
	ErrInvalidLogin = errors.New("invalid credentials")
	ErrNotAuthed    = errors.New("not authenticated")
	ErrUserNotFound = errors.New("user not found")
)

// Service holds the user feature's business logic and its dependencies.
type Service struct {
	q        *data.Queries
	sessions *session.Manager
}

// NewService builds the user service.
func NewService(q *data.Queries, sessions *session.Manager) *Service {
	return &Service{q: q, sessions: sessions}
}

// Register creates a new user and starts a session.
func (s *Service) Register(ctx context.Context, email, password, name string) (UserResponse, error) {
	if _, err := s.q.GetUserByEmail(ctx, email); err == nil {
		return UserResponse{}, ErrEmailTaken
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return UserResponse{}, fmt.Errorf("lookup email: %w", err)
	}

	hash, err := hashPassword(password)
	if err != nil {
		return UserResponse{}, fmt.Errorf("hash password: %w", err)
	}

	row, err := s.q.CreateUser(ctx, data.CreateUserParams{
		Email:    email,
		Password: hash,
		Name:     name,
	})
	if err != nil {
		return UserResponse{}, fmt.Errorf("create user: %w", err)
	}

	if err := s.sessions.Login(ctx, row.ID.String()); err != nil {
		return UserResponse{}, fmt.Errorf("start session: %w", err)
	}
	return responseFromCreate(row), nil
}

// Login authenticates a user and starts a session.
func (s *Service) Login(ctx context.Context, email, password string) (UserResponse, error) {
	u, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserResponse{}, ErrInvalidLogin
		}
		return UserResponse{}, fmt.Errorf("lookup email: %w", err)
	}

	if !checkPassword(u.Password, password) {
		return UserResponse{}, ErrInvalidLogin
	}

	if err := s.sessions.Login(ctx, u.ID.String()); err != nil {
		return UserResponse{}, fmt.Errorf("start session: %w", err)
	}
	return responseFromUser(u), nil
}

// Logout ends the current session.
func (s *Service) Logout(ctx context.Context) error {
	return s.sessions.Logout(ctx)
}

// Current returns the authenticated user, or ErrNotAuthed if none.
func (s *Service) Current(ctx context.Context) (UserResponse, error) {
	id := s.sessions.UserID(ctx)
	if id == "" {
		return UserResponse{}, ErrNotAuthed
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return UserResponse{}, ErrNotAuthed
	}

	u, err := s.q.GetUserByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserResponse{}, ErrUserNotFound
		}
		return UserResponse{}, fmt.Errorf("lookup user: %w", err)
	}
	return responseFromByID(u), nil
}
