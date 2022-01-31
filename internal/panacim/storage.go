package panacim

import (
	"database/sql"
	"sync"

	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/logging"
)

type PanaCIMStorage struct {
	DB     *sql.DB
	logger *logging.Logger
	mu     *sync.Mutex
}

// [PanaCIM].[dbo].[InfoInstallLastJobId_View]
type InfoInstallLastJobId_View struct {
	PartNo        string `db:"PART_NO"`
	PlaceCount    string `db:"PLACE_COUNT"`
	SumPlaceCount string `db:"SUM_PLACE_COUNT_ALL"`
}

// [PanaCIM].[dbo].[substitute_parts]
type SubstituteParts struct {
	PrimaryPn    string `db:"PRIMARY_PN"`
	SubstitutePn string `db:"SUBSTITUTE_PN"`
}

// [PanaCIM].[dbo].[product_data]
type ProductData struct {
	ProductName     string `db:"PRODUCT_NAME"`
	PatternPerPanel string `db:"PATTERN_COMBINATIONS_PER_PANEL"`
}

type ProductDataLink []ProductData

// [PanaCIM].[dbo].[job_products]
type JobProducts struct {
	SetupId string `db:"SETUP_ID"`
}

// [PanaCIM].[dbo].[product_setup]
type ProductSetup struct {
	Product_Id string `db:"PRODUCT_ID"`
	Route_Id   string `db:"ROUTE_ID"`
	MixName    string `db:"MIX_NAME"`
}
