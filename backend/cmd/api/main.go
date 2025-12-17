package main

import (
	"fmt"
	"log"

	"good-todo-go/internal/domain/repository"
	"good-todo-go/internal/ent/generated"
	"good-todo-go/internal/infrastructure/environment"
	infrarepo "good-todo-go/internal/infrastructure/repository"
	"good-todo-go/internal/pkg"
	"good-todo-go/internal/presentation/public/api"
	"good-todo-go/internal/presentation/public/controller"
	"good-todo-go/internal/presentation/public/presenter"
	"good-todo-go/internal/presentation/public/router"
	"good-todo-go/internal/presentation/public/router/middleware"
	"good-todo-go/internal/usecase"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"

	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

func main() {
	container := dig.New()

	// Environment
	if err := container.Provide(environment.NewEnvironment); err != nil {
		log.Fatal(err)
	}

	// Database
	if err := container.Provide(func(env *environment.Environment) (*generated.Client, error) {
		db, err := sql.Open("postgres", env.GetAppDSN())
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}
		drv := entsql.OpenDB(dialect.Postgres, db)
		return generated.NewClient(generated.Driver(drv)), nil
	}); err != nil {
		log.Fatal(err)
	}

	// PKG Services
	if err := container.Provide(func(env *environment.Environment) *pkg.JWTService {
		return pkg.NewJWTService(env.JWTSecret)
	}); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(pkg.NewPasswordService); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(func() pkg.IUUIDGenerator {
		return pkg.NewUUIDGenerator()
	}); err != nil {
		log.Fatal(err)
	}

	// Repositories
	if err := container.Provide(func(client *generated.Client) repository.ITenantRepository {
		return infrarepo.NewTenantRepository(client)
	}); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(func(client *generated.Client) repository.IUserRepository {
		return infrarepo.NewUserRepository(client)
	}); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(func(client *generated.Client) repository.ITodoRepository {
		return infrarepo.NewTodoRepository(client)
	}); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(func(env *environment.Environment) repository.IAuthRepository {
		return infrarepo.NewAuthRepository(env)
	}); err != nil {
		log.Fatal(err)
	}

	// Usecases
	if err := container.Provide(func(
		tenantRepo repository.ITenantRepository,
		userRepo repository.IUserRepository,
		authRepo repository.IAuthRepository,
		jwtService *pkg.JWTService,
		passwordService *pkg.PasswordService,
		uuidGenerator pkg.IUUIDGenerator,
	) usecase.IAuthInteractor {
		return usecase.NewAuthInteractor(tenantRepo, userRepo, authRepo, jwtService, passwordService, uuidGenerator)
	}); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(func(userRepo repository.IUserRepository) usecase.IUserInteractor {
		return usecase.NewUserInteractor(userRepo)
	}); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(func(todoRepo repository.ITodoRepository, uuidGenerator pkg.IUUIDGenerator) usecase.ITodoInteractor {
		return usecase.NewTodoInteractor(todoRepo, uuidGenerator)
	}); err != nil {
		log.Fatal(err)
	}

	// Presenters
	if err := container.Provide(func() presenter.IAuthPresenter {
		return presenter.NewAuthPresenter()
	}); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(func() presenter.IUserPresenter {
		return presenter.NewUserPresenter()
	}); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(func() presenter.ITodoPresenter {
		return presenter.NewTodoPresenter()
	}); err != nil {
		log.Fatal(err)
	}

	// Controllers
	if err := container.Provide(func(
		authUsecase usecase.IAuthInteractor,
		authPresenter presenter.IAuthPresenter,
		tenantRepo repository.ITenantRepository,
	) *controller.AuthController {
		return controller.NewAuthController(authUsecase, authPresenter, tenantRepo)
	}); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(func(
		userUsecase usecase.IUserInteractor,
		userPresenter presenter.IUserPresenter,
	) *controller.UserController {
		return controller.NewUserController(userUsecase, userPresenter)
	}); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(func(
		todoUsecase usecase.ITodoInteractor,
		todoPresenter presenter.ITodoPresenter,
	) *controller.TodoController {
		return controller.NewTodoController(todoUsecase, todoPresenter)
	}); err != nil {
		log.Fatal(err)
	}

	// Middleware
	if err := container.Provide(func(jwtService *pkg.JWTService) *middleware.JWTAuthMiddleware {
		return middleware.NewJWTAuthMiddleware(jwtService)
	}); err != nil {
		log.Fatal(err)
	}

	// Server
	if err := container.Provide(router.NewServer); err != nil {
		log.Fatal(err)
	}

	// Start server
	err := container.Invoke(func(
		env *environment.Environment,
		server *router.Server,
		jwtAuth *middleware.JWTAuthMiddleware,
	) error {
		e := echo.New()
		e.Use(echomiddleware.Logger())
		e.Use(echomiddleware.Recover())
		e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
			AllowOrigins: []string{"http://localhost:3000"},
			AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		}))

		// Register routes with authentication middleware for protected endpoints
		apiGroup := e.Group("/api/v1")

		// Health check (no auth)
		apiGroup.GET("/health", func(c echo.Context) error {
			return server.GetHealth(c)
		})

		// Auth routes (no auth)
		apiGroup.POST("/auth/register", func(c echo.Context) error {
			return server.Register(c)
		})
		apiGroup.POST("/auth/login", func(c echo.Context) error {
			return server.Login(c)
		})
		apiGroup.POST("/auth/verify-email", func(c echo.Context) error {
			return server.VerifyEmail(c)
		})
		apiGroup.POST("/auth/refresh", func(c echo.Context) error {
			return server.RefreshToken(c)
		})

		// Protected routes
		protected := apiGroup.Group("", jwtAuth.Authenticate)
		protected.GET("/me", func(c echo.Context) error {
			return server.GetMe(c)
		})
		protected.PUT("/me", func(c echo.Context) error {
			return server.UpdateMe(c)
		})
		protected.GET("/todos", func(c echo.Context) error {
			return server.ListTodos(c)
		})
		protected.GET("/todos-public", func(c echo.Context) error {
			return server.ListPublicTodos(c)
		})
		protected.POST("/todos", func(c echo.Context) error {
			return server.CreateTodo(c)
		})
		protected.PUT("/todos/:id", func(c echo.Context) error {
			return server.UpdateTodo(c, c.Param("id"))
		})
		protected.DELETE("/todos/:id", func(c echo.Context) error {
			return server.DeleteTodo(c, c.Param("id"))
		})

		// Verify ServerInterface implementation
		var _ api.ServerInterface = server

		log.Printf("Starting server on port %s", env.PublicAPIPort)
		return e.Start(":" + env.PublicAPIPort)
	})

	if err != nil {
		log.Fatal(err)
	}
}
