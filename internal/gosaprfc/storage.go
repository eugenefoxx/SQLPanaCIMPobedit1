package gosaprfc

import "github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/logging"

type goRFCStorage struct {
	logger *logging.Logger
}

func NewGoRFCRepository(logger *logging.Logger) GoRFCRepository {
	return &goRFCStorage{
		logger: logger,
	}
}

type GoRFCRepository interface {
	GetSPP(woname string) (spp interface{}, err error)
}
