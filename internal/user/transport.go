package user

import (
	"go-api-template/internal/openapi"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Transport handles transport for service
type Transport struct {
	logger        zerolog.Logger
	srv           Service
	tokenGenrator func(u *User) (string, error)
}

// NewTransport creates a new transport
func NewTransport(
	logger zerolog.Logger,
	srv Service,
	tokenGenrator func(u *User) (string, error),
) Transport {
	return Transport{
		logger:        logger,
		srv:           srv,
		tokenGenrator: tokenGenrator,
	}
}

// RegisterUser registers a user
func (h Transport) RegisterUser(c echo.Context) error {
	ctx := c.Request().Context()

	req := &openapi.UserRegistrationRequest{}
	err := c.Bind(req)
	if err != nil {
		h.logger.Err(err).Msg("")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err = h.srv.Create(ctx, &User{
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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, openapi.Status{
		Message: "user registered",
	})
}

// LoginUser authenicates users
func (h Transport) LoginUser(c echo.Context) error {
	ctx := c.Request().Context()

	req := &openapi.UserLoginRequest{}
	err := c.Bind(req)
	if err != nil {
		h.logger.Err(err).Msg("")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	u, err := h.srv.FindByEmail(ctx, req.Email)
	if err != nil {
		h.logger.Err(err).Msg("")
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	token, err := h.tokenGenrator(u)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, openapi.UserLoginResponse{
		Token:     token,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	})
}
