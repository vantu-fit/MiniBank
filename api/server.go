package api

import (
	"github.com/vantu-fit/master-go-be/token"
	"github.com/vantu-fit/master-go-be/utils"

	"github.com/gin-gonic/gin"
	db "github.com/vantu-fit/master-go-be/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
	maker  token.Maker
	config utils.Config
}

func NewServer(store db.Store) (*Server, error) {
	config, err := utils.LoadConfig("..")
	if err != nil {
		return nil, err
	}

	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{store: store}
	router := gin.Default()

	server.maker = maker
	server.config = config

	authRoutes := router.Group("/").Use(authMiddleWare(server.maker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)

	authRoutes.POST("/transfer", server.createTransfer)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	router.POST("/tokens/renew_access", server.renewAccessToken)

	server.router = router

	return server, nil

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)

}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}

}
