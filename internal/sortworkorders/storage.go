package sortworkorders

import (
	"database/sql"
	"sync"

	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/logging"
)

type OperationStorage struct {
	DB     *sql.DB
	logger logging.Logger
	mu     sync.Mutex
}
