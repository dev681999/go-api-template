package user

import (
	"errors"
	apperr "go-api-template/internal/error"
	"go-api-template/internal/openapi"
	"go-api-template/pkg/util"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// TransportUser handles transport for service
type TransportUser struct {
	logger zerolog.Logger
	srv    Service
}

// NewTransport creates a new transport
func NewTransport(
	logger zerolog.Logger,
	srv Service,
) TransportUser {
	return TransportUser{
		logger: logger,
		srv:    srv,
	}
}

// RegisterUser registers a user
func (t TransportUser) RegisterUser(c echo.Context) error {
	ctx := c.Request().Context()

	req := &openapi.UserRegistrationRequest{}
	{
		err := c.Bind(req)
		if err != nil {
			t.logger.Err(err).Msg("")
			return apperr.NewWithCode("transport", err.Error(), http.StatusBadRequest, err)
		}
	}

	_, err := t.srv.Create(ctx, &User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Mobile:    req.Mobile,
		Address:   req.Address,
		Active:    true,
	})
	if err != nil {
		t.logger.Err(err).Msg("")
		if errors.Is(err, ErrUserAlreadyExists) {
			return apperr.NewWithCode("transport", ErrUserAlreadyExists.Reason, http.StatusBadRequest, err)
		}

		return apperr.NewWithCode("transport", ErrInternalService.Reason, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, openapi.Status{
		Message: "user registered",
	})
}

// LoginUser authenicates users
func (t TransportUser) LoginUser(c echo.Context) error {
	ctx := c.Request().Context()

	req := &openapi.UserLoginRequest{}
	if err := c.Bind(req); err != nil {
		t.logger.Err(err).Msg("")
		return apperr.NewWithCode("transport", "bad request", http.StatusBadRequest, err)
	}

	u, err := t.srv.ValidateByEmailAndPassword(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return apperr.NewWithCode("transport", ErrUserNotFound.Reason, http.StatusForbidden, err)
		} else if errors.Is(err, ErrWrongCredentials) {
			return apperr.NewWithCode("transport", ErrWrongCredentials.Reason, http.StatusForbidden, err)
		}

		return apperr.NewWithCode("transport", ErrInternalService.Reason, http.StatusInternalServerError, err)
	}

	token, err := t.srv.GenerateToken(ctx, u, ClaimTypeNormal)
	if err != nil {
		t.logger.Err(err).Msg("")
		return apperr.NewWithCode("transport", ErrInternalService.Reason, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, openapi.UserLoginResponse{
		Token:     token,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	})
}

// UserProfile gets users profile
func (t TransportUser) UserProfile(c echo.Context) error {
	claims := GetClaimFromEchoContext(c)

	u, err := t.srv.FindByID(c.Request().Context(), claims.UserID)
	if err != nil {
		t.logger.Err(err).Msg("")
		if errors.Is(err, ErrUserNotFound) {
			return apperr.NewWithCode("transport", ErrUserNotFound.Reason, http.StatusBadRequest, err)
		}

		return apperr.NewWithCode("transport", ErrInternalService.Reason, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, openapi.UserProfile{
		Address:   u.Address,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Mobile:    u.Mobile,
		ImageUrl:  u.ImageURL,
	})
}

// UpdateUserProfile updates user profile
func (t TransportUser) UpdateUserProfile(c echo.Context) error {
	ctx := c.Request().Context()
	claims := GetClaimFromEchoContext(c)

	req := &openapi.UserProfileUpdateRequest{}
	if err := c.Bind(req); err != nil {
		t.logger.Err(err).Msg("")
		return apperr.NewWithCode("transport", "bad request", http.StatusBadRequest, err)
	}

	u, err := t.srv.Update(ctx, claims.UserID, &Update{
		Email:     req.Email,
		Mobile:    req.Mobile,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		ImageURL:  req.ImageUrl,
		Address:   req.Address,
	})
	if err != nil {
		t.logger.Err(err).Msg("")
		if errors.Is(err, ErrUserNotFound) {
			return apperr.NewWithCode("transport", ErrUserNotFound.Reason, http.StatusForbidden, err)
		}

		return apperr.NewWithCode("transport", ErrInternalService.Reason, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, openapi.UserProfile{
		Address:   u.Address,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Mobile:    u.Mobile,
		ImageUrl:  u.ImageURL,
	})
}

// ActivateUser activates a user
func (t TransportUser) ActivateUser(c echo.Context, params openapi.ActivateUserParams) error {
	ctx := c.Request().Context()
	claims := GetClaimFromEchoContext(c)

	if claims.Type != ClaimTypeActivation {
		err := errors.New("invalid token type")
		t.logger.Err(err).Msg("")
		return apperr.NewWithCode("transport", "invalid token", http.StatusForbidden, err)
	}

	_, err := t.srv.Update(ctx, claims.UserID, &Update{
		Active: util.BoolPtr(true),
	})
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return apperr.NewWithCode("transport", ErrUserNotFound.Reason, http.StatusForbidden, err)
		}

		return apperr.NewWithCode("transport", ErrInternalService.Reason, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, openapi.Status{
		Message: "account activated",
	})
}
