package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/lukinairina90/banking_backend/internal/repository"
	"github.com/lukinairina90/banking_backend/internal/service"
	"github.com/lukinairina90/banking_backend/internal/transport/rest"
	"github.com/lukinairina90/banking_backend/pkg/config"
	"github.com/lukinairina90/banking_backend/pkg/database"
	"github.com/lukinairina90/banking_backend/pkg/hash"
	"github.com/lukinairina90/banking_backend/pkg/lib/generator"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		logrus.WithError(err).Fatalf("error parsing config from env variables: %s", err.Error())
	}

	// init db
	db, err := database.CreateConn(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.SSLMode)
	if err != nil {
		logrus.WithError(err).Fatalf("failed to connection db: %s", err.Error())
	}

	defer db.Close()

	// Init casbin for RBAC middleware.
	enforser, err := casbin.NewEnforcer(cfg.RBACConfig.ModelFilePath, cfg.RBACConfig.PolicyFilePath)
	if err != nil {
		logrus.WithError(err).Fatal("error initialization casbin enforcer")
	}

	// init deps
	hasher := hash.NewSHA1Hasher(cfg.UserPasswordSalt)
	randomGenerator := generator.NewGenerator("UA", "123456")

	//initialize repositories
	usersRepository := repository.NewUsers(db)
	tokensRepository := repository.NewTokens(db)
	rolesRepository := repository.NewRoles(db)
	accountRepository := repository.NewAccount(db)
	transactionRepository := repository.NewTransactions(db)
	cardRepository := repository.NewCard(db)
	eventRepository := repository.NewEvent(db)

	//initialize services
	usersService := service.NewUsers(usersRepository, tokensRepository, rolesRepository, eventRepository, hasher, []byte(cfg.TokenSecret), cfg.TokenTTL)
	accountService := service.NewAccount(accountRepository, transactionRepository, eventRepository, randomGenerator)
	transactionService := service.NewTransaction(transactionRepository, accountRepository)
	cardService := service.NewCard(cardRepository, usersRepository, accountRepository, eventRepository, randomGenerator)
	eventService := service.NewEvent(eventRepository)

	//initialize transports
	authTransport := rest.NewAuth(usersService)
	accountTransport := rest.NewAccount(accountService)
	transactionTransport := rest.NewTransaction(transactionService)
	cardTransport := rest.NewCard(cardService)
	eventTransport := rest.NewEvent(eventService)

	rbacMiddleware := rest.RBACMiddleware(enforser, rolesRepository)

	// init routes
	g := gin.New()
	//g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	g.Use(rest.LoggingMiddleware())
	authTransport.InjectRoutes(g, rbacMiddleware)
	accountTransport.InjectRoutes(g, authTransport.AuthMiddleware(), rbacMiddleware)
	transactionTransport.InjectRoutes(g, authTransport.AuthMiddleware(), rbacMiddleware)
	cardTransport.InjectRoutes(g, authTransport.AuthMiddleware(), rbacMiddleware)
	eventTransport.InjectRoutes(g, authTransport.AuthMiddleware(), rbacMiddleware)

	fmt.Println("Server run...")
	if err := g.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		logrus.Fatalf("error occured while running http server %s", err.Error())
	}
}
