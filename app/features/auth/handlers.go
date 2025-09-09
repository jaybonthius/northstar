package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"northstar/app/features/auth/gen/authdb"
	"northstar/app/features/auth/pages"
	"northstar/app/features/common/utils"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/starfederation/datastar-go/datastar"
	"golang.org/x/crypto/bcrypt"
)

type ValidationErrors struct {
	Username string
	Email    string
	Password string
}

func (v ValidationErrors) HasErrors() bool {
	return v.Username != "" || v.Email != "" || v.Password != ""
}

type authHandlers struct {
	repository *authRepository
	store      sessions.Store
}

func (h *authHandlers) sendGenericError(w http.ResponseWriter, r *http.Request, message string) {
	sse := datastar.NewSSE(w, r)
	errorHTML, _ := utils.RenderTemplToString(r.Context(), pages.GenericAuthError(message))
	if err := sse.PatchElements(errorHTML); err != nil {
		slog.Error("Failed to patch elements", "error", err)
	}
}

func (h *authHandlers) sendSignupErrors(w http.ResponseWriter, r *http.Request, errors ValidationErrors) {
	sse := datastar.NewSSE(w, r)

	var allHTML string
	if errors.Username != "" {
		html, _ := utils.RenderTemplToString(r.Context(), pages.UsernameError(errors.Username))
		allHTML += html
	} else {
		html, _ := utils.RenderTemplToString(r.Context(), pages.UsernameError(""))
		allHTML += html
	}

	if errors.Email != "" {
		html, _ := utils.RenderTemplToString(r.Context(), pages.EmailError(errors.Email))
		allHTML += html
	} else {
		html, _ := utils.RenderTemplToString(r.Context(), pages.EmailError(""))
		allHTML += html
	}

	if errors.Password != "" {
		html, _ := utils.RenderTemplToString(r.Context(), pages.PasswordError(errors.Password))
		allHTML += html
	} else {
		html, _ := utils.RenderTemplToString(r.Context(), pages.PasswordError(""))
		allHTML += html
	}

	if err := sse.PatchElements(allHTML); err != nil {
		slog.Error("Failed to patch elements", "error", err)
	}
}

func (h *authHandlers) sendLoginErrors(w http.ResponseWriter, r *http.Request, errors ValidationErrors) {
	sse := datastar.NewSSE(w, r)

	var allHTML string

	if errors.Email != "" {
		html, _ := utils.RenderTemplToString(r.Context(), pages.EmailError(errors.Email))
		allHTML += html
	} else {
		html, _ := utils.RenderTemplToString(r.Context(), pages.EmailError(""))
		allHTML += html
	}

	if errors.Password != "" {
		html, _ := utils.RenderTemplToString(r.Context(), pages.PasswordError(errors.Password))
		allHTML += html
	} else {
		html, _ := utils.RenderTemplToString(r.Context(), pages.PasswordError(""))
		allHTML += html
	}

	if err := sse.PatchElements(allHTML); err != nil {
		slog.Error("Failed to patch elements", "error", err)
	}
}

func (h *authHandlers) createSession(w http.ResponseWriter, r *http.Request, userID string) error {
	session, err := h.store.Get(r, "auth-session")
	if err != nil {
		return fmt.Errorf("getting session: %w", err)
	}

	session.Values["user_id"] = userID
	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("saving session: %w", err)
	}

	return nil
}

func (h *authHandlers) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	if err := pages.LoginPage().Render(r.Context(), w); err != nil {
		slog.Error("Failed to render login page", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *authHandlers) handleSignupPage(w http.ResponseWriter, r *http.Request) {
	if err := pages.SignupPage().Render(r.Context(), w); err != nil {
		slog.Error("Failed to render signup page", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *authHandlers) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendGenericError(w, r, MsgInvalidMethod)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("Error parsing form data", "error", err)
		h.sendGenericError(w, r, MsgInvalidFormData)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		h.sendGenericError(w, r, MsgMissingCredentials)
		return
	}

	user, validationErr, err := h.validateLogin(r.Context(), email, password)
	if err != nil {
		slog.Error("Error during login validation", "email", email, "error", err)
		h.sendGenericError(w, r, MsgLoginFailed)
		return
	}

	if validationErr.HasErrors() {
		h.sendLoginErrors(w, r, validationErr)
		return
	}

	if err := h.createSession(w, r, user.ID); err != nil {
		slog.Error("Error creating session", "error", err)
		h.sendGenericError(w, r, MsgLoginFailed)
		return
	}

	sse := datastar.NewSSE(w, r)
	if err := sse.ExecuteScript("window.location.href = '/'"); err != nil {
		slog.Error("Failed to execute script", "error", err)
	}
}

func (h *authHandlers) handleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendGenericError(w, r, MsgInvalidMethod)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.sendGenericError(w, r, MsgInvalidFormData)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	validationErr, err := h.validateSignup(r.Context(), username, email, password)
	if err != nil {
		slog.Error("Error during signup validation", "error", err)
		h.sendGenericError(w, r, MsgSignupFailed)
		return
	}

	if validationErr.HasErrors() {
		h.sendSignupErrors(w, r, validationErr)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error hashing password", "error", err)
		h.sendGenericError(w, r, MsgSignupFailed)
		return
	}

	userID := uuid.New().String()
	user, err := h.repository.createUser(r.Context(), authdb.CreateUserParams{
		ID:           userID,
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		slog.Error("Error creating user", "error", err)
		h.sendGenericError(w, r, MsgSignupFailed)
		return
	}

	if err := h.createSession(w, r, user.ID); err != nil {
		slog.Error("Error creating session after signup", "error", err)
		h.sendGenericError(w, r, MsgAccountCreatedLoginFailed)
		return
	}

	sse := datastar.NewSSE(w, r)
	if err := sse.ExecuteScript("window.location.href = '/'"); err != nil {
		slog.Error("Failed to execute script", "error", err)
	}
}

func (h *authHandlers) handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "auth-session")
	slog.Info("Logging out user", "user_id", session.Values["user_id"])
	session.Values["user_id"] = nil
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		slog.Error("Failed to save session", "error", err)
	}

	sse := datastar.NewSSE(w, r)
	if err := sse.Redirect("/"); err != nil {
		slog.Error("Failed to redirect after logout", "error", err)
	}
}

func (h *authHandlers) validateLogin(ctx context.Context, email, password string) (*authdb.User, ValidationErrors, error) {
	var validationErr ValidationErrors

	userExists, err := h.repository.checkIfUserExistsByEmail(ctx, email)
	if err != nil {
		return nil, validationErr, err
	}
	if !userExists {
		validationErr.Email = MsgUserNotFound
		return nil, validationErr, nil
	}

	user, err := h.repository.getUserByEmail(ctx, email)
	if err != nil {
		return nil, validationErr, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		validationErr.Password = MsgInvalidCredentials
		return nil, validationErr, nil
	}

	return &user, validationErr, nil
}

func (h *authHandlers) validateSignup(ctx context.Context, username, email, password string) (ValidationErrors, error) {
	var validationErr ValidationErrors

	userExists, err := h.repository.checkIfUserExistsByUsername(ctx, username)
	if err != nil {
		return validationErr, err
	}
	if userExists {
		validationErr.Username = MsgUsernameAlreadyExists
	}

	emailExists, err := h.repository.checkIfUserExistsByEmail(ctx, email)
	if err != nil {
		return validationErr, err
	}
	if emailExists {
		validationErr.Email = MsgEmailAlreadyExists
	}

	if len(password) < 6 {
		validationErr.Password = MsgPasswordTooShort
	}

	return validationErr, nil
}
