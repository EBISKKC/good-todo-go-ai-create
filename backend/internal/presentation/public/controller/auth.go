package controller

import (
	"net/http"

	"good-todo-go/internal/domain/repository"
	"good-todo-go/internal/presentation/public/api"
	"good-todo-go/internal/presentation/public/presenter"
	"good-todo-go/internal/usecase"
	"good-todo-go/internal/usecase/input"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authUsecase   usecase.IAuthInteractor
	authPresenter presenter.IAuthPresenter
	tenantRepo    repository.ITenantRepository
}

func NewAuthController(
	authUsecase usecase.IAuthInteractor,
	authPresenter presenter.IAuthPresenter,
	tenantRepo repository.ITenantRepository,
) *AuthController {
	return &AuthController{
		authUsecase:   authUsecase,
		authPresenter: authPresenter,
		tenantRepo:    tenantRepo,
	}
}

func (ctrl *AuthController) Register(c echo.Context, req api.RegisterRequest) error {
	out, err := ctrl.authUsecase.Register(c.Request().Context(), &input.RegisterInput{
		Email:    string(req.Email),
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		if err == usecase.ErrUserAlreadyExists {
			return echo.NewHTTPError(http.StatusBadRequest, "user already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctrl.authPresenter.Register(c, out)
}

func (ctrl *AuthController) Login(c echo.Context, req api.LoginRequest) error {
	// First, find tenant by slug
	tenant, err := ctrl.tenantRepo.FindBySlug(c.Request().Context(), req.TenantSlug)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if tenant == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	// Cast to access the method with tenant ID
	authInteractor, ok := ctrl.authUsecase.(*usecase.AuthInteractor)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	out, err := authInteractor.LoginWithTenant(c.Request().Context(), tenant.ID, &input.LoginInput{
		Email:    string(req.Email),
		Password: req.Password,
	})
	if err != nil {
		if err == usecase.ErrInvalidCredentials {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}
		if err == usecase.ErrEmailNotVerified {
			return echo.NewHTTPError(http.StatusUnauthorized, "email not verified")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctrl.authPresenter.Login(c, out)
}

func (ctrl *AuthController) VerifyEmail(c echo.Context, req api.VerifyEmailRequest) error {
	out, err := ctrl.authUsecase.VerifyEmail(c.Request().Context(), &input.VerifyEmailInput{
		Token: req.Token,
	})
	if err != nil {
		if err == usecase.ErrInvalidToken {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid token")
		}
		if err == usecase.ErrTokenExpired {
			return echo.NewHTTPError(http.StatusBadRequest, "token expired")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctrl.authPresenter.VerifyEmail(c, out)
}

func (ctrl *AuthController) RefreshToken(c echo.Context, req api.RefreshTokenRequest) error {
	out, err := ctrl.authUsecase.RefreshToken(c.Request().Context(), &input.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		if err == usecase.ErrInvalidToken {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctrl.authPresenter.RefreshToken(c, out)
}
