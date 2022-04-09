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

type Job_History struct {
	JOB_ID            string         `db:"JOB_ID"`
	EQUIPMENT_ID      string         `db:"EQUIPMENT_ID"`
	SETUP_ID          string         `db:"SETUP_ID"`
	StartUnixTimeWO   string         `db:"START_TIME"`
	EndUnixTimeWO     string         `db:"END_TIME"`
	CLOSING_TYPE      string         `db:"CLOSING_TYPE"`
	START_OPERATOR_ID sql.NullString `db:"START_OPERATOR_ID"`
	END_OPERATOR_ID   sql.NullString `db:"END_OPERATOR_ID"`
	TFR_REASON        sql.NullString `db:"TFR_REASON"`
	LANE_NO           string         `db:"LANE_NO"`
}
