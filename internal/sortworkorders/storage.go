package sortworkorders

/*
type OperationStorage struct {
	DB     *sql.DB
	logger logging.Logger
	mu     sync.Mutex
}
*/
type SortWorkOrderRepository interface {
	Getclosedworkorders()
	GetLastJobIdValue1() (string, error)
	GetLastJobIdValue2() (string, error)
	GetLastJobIdValue3() (string, error)

	TestQr() ([]PanelMap, error)
}
