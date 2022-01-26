package sortworkorders

import (
	"database/sql"

	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/logging"
)

type OperationStorage struct {
	DB     *sql.DB
	logger *logging.Logger
}
