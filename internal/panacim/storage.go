package panacim

import (
	"database/sql"
)

type PanaCIMRepository interface {
	Print()
	GetPanacimDataComponentsByJobIdSAP(string) ([]InfoInstallLastJobId_View, error)
	WritePanacimDataComponentsByJobIdSAPToFile(in []InfoInstallLastJobId_View) (err error)
	WriteDataInfoOrderSAP(wo_name, sum string) error

	GetLastListWO() ([]LastWOData, error)
	GetWOName(name string) ([]LastWOData, error)
	WriteWorkOrderNameToFile(in []LastWOData) (err error)
	WriteListWOToFile(in []LastWOData) (err error)

	GetSumPattert(jobid string) ([]SumPattern, error)
	GetPatternForPanel() ([]ProductData, error)
	GetProductId(jobid string) ([]ProductSetup, error)
	GetProductName(productid string) ([]ProductData, error)
	GetRouteId(productid string) ([]ProductSetup, error)
	GetPanacimDataComponentsByJobId(jobid string) ([]InfoInstallLastJobId_View, error)
	GetPanacimDataComponentsByJobIdAllParamReelid(jobid string) ([]InfoInstallLastJobId_View, error)
	// получить wo_component из db
	GetPanaCIMDBWOComponent(jobid string) ([]InfoInstallLastJobId_View, error)
	// запись результата GetPanaCIMDBWOComponent в файл
	WritePanaCIMDBWOComponentToFile(in []InfoInstallLastJobId_View) (err error)
	// Получить из БД PanaCIM состав списания по ID в скрап
	GetPanaCIMDBScrapID(jobid string) ([]InfoInstallLastJobId_View, error)
	// запись результата GetPanaCIMDBScrapID в файл
	WritePanaCIMDBScrapIDToFile(in []InfoInstallLastJobId_View) (err error)
	// Получить из БД PanaCIM состав парт-номер, сумма, лот на скрап
	GetPanaCIMDBScrap(jobid string) ([]InfoInstallLastJobId_View, error)
	// запись результата GetPanaCIMDBScrap в файл
	WritePanaCIMDBScrapToFile(in []InfoInstallLastJobId_View) (err error)

	WtitePanaCIMDataComponentsToFile(in []InfoInstallLastJobId_View) (err error)
	WtitePanaCIMDataComponentsToFileUnpackId(in []InfoInstallLastJobId_View) (err error)
	GetPanaCIMixName(productid string) ([]ProductSetup, error)
	GetPanaCIMParts(mixname string) ([]SubstituteParts, error)
	WritePanaCIMPartsToFile(in []SubstituteParts) (err error)
	GetUnixTimeWO(jobid string) ([]Job_History, error)
	GetSumPCBFromU03V2(startUnixTimeWO, finishUnixTimeWO, npm string) (sumstrPCBOrder string)
	GetSumPCBFromU03(startUnixTimeWO, finishUnixTimeWO, npm string) (sumstrPCBOrder string)
	GetSumComponentFromU03(startUnixTimeWO, finishUnixTimeWO, npm string) error

	GetPatternTypesPerPanel(product_id string) ([]ProductData, error)
	// сверка данных ЕО на основе SAP и машинных файлов U03
	CompareDataID(spp interface{}) (response bool, err error)
	// сверка данных ЕО на основе SAP и DB PanaCIM
	CompareDataIDFromDB(spp interface{}, jobid string) (response bool, err error)
}

/*type PanaCIMStorage struct {
	DB     *sql.DB
	logger *logging.Logger
	mu     *sync.Mutex
}*/

type LastWOData struct {
	WORKORDERID          string         `db:"WORK_ORDER_ID"`
	WORKORDERNAME        string         `db:"WORK_ORDER_NAME"`
	LOTSIZE              string         `db:"LOT_SIZE"`
	JOBID                string         `db:"JOB_ID"`
	MASTER_WORK_ORDER_ID string         `db:"MASTER_WORK_ORDER_ID"`
	COMMENTS             sql.NullString `db:"COMMENTS"`
}

// [PanaCIM].[dbo].[InfoInstallLastJobId_View]
type InfoInstallLastJobId_View struct {
	ReelID          string `db:"REEL_ID"`
	PartNo          string `db:"PART_NO"`
	Lot             string `db:"LOT_NO"`
	PlaceCount      string `db:"PLACE_COUNT"`
	PickupCount     string `db:"PICKUP_COUNT"`
	ReelBarcode     string `db:"reel_barcode"`
	CurrentQuantity string `db:"CURRENT_QUANTITY"`
	InitialQuantity string `db:"INITIAL_QUANTITY"`
	PlaceCountAll   string `db:"PLACE_COUNT_ALL"`
	PickupCountAll  string `db:"PICKUP_COUNT_ALL"`
	Delta           string `db:"Delta"`
	SumPlaceCount   string `db:"SUM_PLACE_COUNT_ALL"`
}

// [PanaCIM].[dbo].[substitute_parts]
type SubstituteParts struct {
	PrimaryPn    string `db:"PRIMARY_PN"`
	SubstitutePn string `db:"SUBSTITUTE_PN"`
}

// [PanaCIM].[dbo].[product_data]
type ProductData struct {
	ProductName             string `db:"PRODUCT_NAME"`
	PatternPerPanel         string `db:"PATTERN_COMBINATIONS_PER_PANEL"`
	Product_ID              string `db:"PRODUCT_ID"`
	PATTERN_TYPES_PER_PANEL string `db:"PATTERN_TYPES_PER_PANEL"`
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

type Reel_Data struct {
	REEL_ID                   string `db:"REEL_ID"`
	PART_NO                   string `db:"PART_NO"`
	MCID                      string `db:"MCID"`
	VENDOR_NO                 string `db:"VENDOR_NO"`
	LOT_NO                    string `db:"LOT_NO"`
	QUANTITY                  string `db:"QUANTITY"`
	USER_DATA                 string `db:"USER_DATA"`
	REEL_BARCODE              string `db:"REEL_BARCODE"`
	CURRENT_QUANTITY          string `db:"CURRENT_QUANTITY"`
	UPDATE_TIME               string `db:"UPDATE_TIME"`
	MASTER_REEL_ID            string `db:"MASTER_REEL_ID"`
	CREATE_TIME               string `db:"CREATE_TIME"`
	PART_CLASS                string `db:"PART_CLASS"`
	MATERIAL_NAME             string `db:"MATERIAL_NAME"`
	PREV_REEL_ID              string `db:"PREV_REEL_ID"`
	NEXT_REEL_ID              string `db:"NEXT_REEL_ID"`
	ADJUSTED_CURRENT_QUANTITY string `db:"ADJUSTED_CURRENT_QUANTITY"`
	TRAY_QUANTITY             string `db:"TRAY_QUANTITY"`
	BULK_MASTER_ID            string `db:"BULK_MASTER_ID"`
	IS_MSD                    string `db:"IS_MSD"`
	MARKET_USAGE              string `db:"MARKET_USAGE"`
	STICK_QUANTITY            string `db:"STICK_QUANTITY"`
	STICK_COUNT               string `db:"STICK_COUNT"`
	SUPPLY_TYPE               string `db:"SUPPLY_TYPE"`
	MASTER_STICK_REEL_ID      string `db:"MASTER_STICK_REEL_ID"`
	CURRENT_STICK_COUNT       string `db:"CURRENT_STICK_COUNT"`
}
