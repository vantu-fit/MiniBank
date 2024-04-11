package gapi

import (
	"github.com/vantu-fit/master-go-be/pb"
	"github.com/vantu-fit/master-go-be/token"
	"github.com/vantu-fit/master-go-be/utils"
	"github.com/vantu-fit/master-go-be/worker"

	db "github.com/vantu-fit/master-go-be/db/sqlc"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store          db.Store
	maker          token.Maker
	config         utils.Config
	taskDitributor worker.TaskDistributor
}

func NewServer(store db.Store, taskDitributor worker.TaskDistributor , config utils.Config) (*Server, error) {

	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		store:  store,
		maker:  maker,
		config: config,
		taskDitributor: taskDitributor,
	}

	return server, nil

}
