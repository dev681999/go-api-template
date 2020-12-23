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

// Transport handles transport for service
type Transport struct {
	logger zerolog.Logger
	srv    Service
}

// NewTransport creates a new transport
func NewTransport(
	logger zerolog.Logger,
	srv Service,
) Transport {
	return Transport{
		logger: logger,
		srv:    srv,
	}
}

// RegisterUser registers a user
func (h Transport) RegisterUser(c echo.Context) error {
	ctx := c.Request().Context()

	req := &openapi.UserRegistrationRequest{}
	{
		err := c.Bind(req)
		if err != nil {
			h.logger.Err(err).Msg("")
			return apperr.New("transport", err.Error(), http.StatusBadRequest, err)
		}
	}

	_, err := h.srv.Create(ctx, &User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Mobile:    req.Mobile,
		Address:   req.Address,
		Active:    true,
	})
	if err != nil {
		h.logger.Err(err).Msg("")
		if err.IsMatchesCode(ErrUserAlreadyExists) {
			return apperr.New("transport", err.Reason, http.StatusBadRequest, err)
		}

		return apperr.New("transport", err.Reason, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, openapi.Status{
		Message: "user registered",
	})
}

// LoginUser authenicates users
func (h Transport) LoginUser(c echo.Context) error {
	ctx := c.Request().Context()

	req := &openapi.UserLoginRequest{}
	if err := c.Bind(req); err != nil {
		h.logger.Err(err).Msg("")
		return apperr.New("transport", "bad request", http.StatusBadRequest, err)
	}

	u, err := h.srv.ValidateByEmailAndPassword(ctx, req.Email, req.Password)
	if err != nil {
		if err.IsMatchesCode(errRepoUserNotFound) {
			return apperr.New("transport", err.Reason, http.StatusForbidden, err)
		}

		return apperr.New("transport", err.Reason, http.StatusInternalServerError, err)
	}

	token, err := h.srv.GenerateToken(ctx, u, ClaimTypeNormal)
	if err != nil {
		h.logger.Err(err).Msg("")
		return apperr.New("transport", err.Reason, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, openapi.UserLoginResponse{
		Token:     token,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	})
}

// UserProfile gets users profile
func (h Transport) UserProfile(c echo.Context) error {
	claims := GetClaimFromEchoContext(c)

	u, err := h.srv.FindByID(c.Request().Context(), claims.UserID)
	if err != nil {
		h.logger.Err(err).Msg("")
		return apperr.New("transport", err.Reason, http.StatusInternalServerError, err)
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
func (h Transport) UpdateUserProfile(c echo.Context) error {
	ctx := c.Request().Context()
	claims := GetClaimFromEchoContext(c)

	req := &openapi.UserProfileUpdateRequest{}
	if err := c.Bind(req); err != nil {
		h.logger.Err(err).Msg("")
		return apperr.New("transport", "bad request", http.StatusBadRequest, err)
	}

	u, err := h.srv.Update(ctx, claims.UserID, &Update{
		Email:     req.Email,
		Mobile:    req.Mobile,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		ImageURL:  req.ImageUrl,
		Address:   req.Address,
	})
	if err != nil {
		h.logger.Err(err).Msg("")
		if err.IsMatchesCode(errRepoUserNotFound) {
			return apperr.New("transport", err.Reason, http.StatusForbidden, err)
		}

		return apperr.New("transport", err.Reason, http.StatusInternalServerError, err)
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
func (h Transport) ActivateUser(c echo.Context, params openapi.ActivateUserParams) error {
	ctx := c.Request().Context()
	claims := GetClaimFromEchoContext(c)

	if claims.Type != ClaimTypeActivation {
		err := errors.New("invalid token type")
		h.logger.Err(err).Msg("")
		return apperr.New("transport", "invalid token", http.StatusForbidden, err)
	}

	_, err := h.srv.Update(ctx, claims.UserID, &Update{
		Active: util.BoolPtr(true),
	})
	if err != nil {
		if err.IsMatchesCode(errRepoUserNotFound) {
			return apperr.New("transport", err.Reason, http.StatusForbidden, err)
		}

		return apperr.New("transport", err.Reason, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, openapi.Status{
		Message: "account activated",
	})
}
