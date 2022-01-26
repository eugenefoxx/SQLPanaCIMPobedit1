package panacim

import (
	"database/sql"

	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/logging"
)

type PanaCIMStorage struct {
	DB     *sql.DB
	logger *logging.Logger
}
