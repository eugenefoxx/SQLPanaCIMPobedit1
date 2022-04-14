package panacim

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/filereader"
	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/logging"
	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/utils"
	cp "github.com/otiai10/copy"
	"gopkg.in/ini.v1"
)

const queryDelObj = `
IF OBJECT_ID('dbo.SUMPattern') IS NOT NULL DROP FUNCTION  dbo.SUMPattern;`

const queryGetSUMPattern = `
CREATE FUNCTION dbo.SUMPattern()
RETURNS TABLE
AS
  --BEGIN
  RETURN
SELECT COUNT(DISTINCT PANEL_ID) AS sumPattern
FROM [PanaCIM].[dbo].[panels] where JOB_ID = `

const querySelectSUMPattern = `SELECT * FROM dbo.SUMPattern();`

type SumPattern struct {
	SumPattern string `db:"sumPattern"`
}

func (r PanaCIMStorage) GetSumPattert(jobid string) ([]SumPattern, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	qrDel, err := r.DB.Query(queryDelObj)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qrDel.Close()

	qrFunc, err := r.DB.ExecContext(ctx, queryGetSUMPattern+jobid)
	if err != nil {
		if err.Error() != "sql: function no create" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qrFunc.RowsAffected()

	qr, err := r.DB.QueryContext(ctx, querySelectSUMPattern)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qr.Close()

	var qrs []SumPattern
	for qr.Next() {
		var qrts SumPattern
		if err := qr.Scan(
			&qrts.SumPattern,
		); err != nil {
			return qrs, err
		}
		qrs = append(qrs, qrts)
	}
	if err = qr.Err(); err != nil {
		return qrs, err
	}
	return qrs, nil
}

const queryDelPCB = `
IF OBJECT_ID('dbo.GetQtyPerPanel') IS NOT NULL DROP FUNCTION dbo.GetQtyPerPanel;
`

const queryPatternForPanel = `
CREATE FUNCTION dbo.GetQtyPerPanel()
    RETURNS TABLE
    AS
    RETURN
SELECT [PATTERN_COMBINATIONS_PER_PANEL]
FROM [PanaCIM].[dbo].[product_data]
WHERE [PRODUCT_ID] = (
    SELECT *
    FROM dbo.GetLastProductId()
    )
`

const querySelectPatternForPanel = `SELECT * FROM dbo.GetQtyPerPanel();`

func (r PanaCIMStorage) GetPatternForPanel() ([]ProductData, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	qrDel, err := r.DB.Query(queryDelPCB)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qrDel.Close()

	qrFunc, err := r.DB.ExecContext(ctx, queryPatternForPanel)
	if err != nil {
		if err.Error() != "sql: function no create" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qrFunc.RowsAffected()

	qr, err := r.DB.QueryContext(ctx, querySelectPatternForPanel)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qr.Close()

	var qrs []ProductData
	for qr.Next() {
		var qrts ProductData
		if err := qr.Scan(
			&qrts.PatternPerPanel,
		); err != nil {
			return qrs, err
		}
		qrs = append(qrs, qrts)
	}
	if err = qr.Err(); err != nil {
		return qrs, err
	}
	return qrs, nil
}

const queryDelObjSetupId = `
IF OBJECT_ID('dbo.GetLastSetupId') IS NOT NULL DROP FUNCTION dbo.GetLastSetupId;
`
const queryDelObjProductId = `
IF OBJECT_ID('dbo.GetLastProductId') IS NOT NULL DROP FUNCTION dbo.GetLastProductId;
`

const queryGetSetupId = `
CREATE FUNCTION dbo.GetLastSetupId()
    RETURNS TABLE
    AS
    RETURN
SELECT [SETUP_ID]
FROM [PanaCIM].[dbo].[job_products]
--where SETUP_ID = '9536'
WHERE JOB_ID = `

const queryGetProductId = `
CREATE FUNCTION dbo.GetLastProductId()
    RETURNS TABLE
    AS
    RETURN
SELECT [PRODUCT_ID]
FROM [PanaCIM].[dbo].[product_setup]
WHERE [SETUP_ID] = (
    SELECT *
    FROM dbo.GetLastSetupId()
    )
`

const querySelectProductId = `SELECT * FROM dbo.GetLastProductId();`

func (r PanaCIMStorage) GetProductId(jobid string) ([]ProductSetup, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	qrDel, err := r.DB.Query(queryDelObjSetupId)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qrDel.Close()

	qrFunc, err := r.DB.ExecContext(ctx, queryGetSetupId+jobid)
	if err != nil {
		if err.Error() != "sql: function no create" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qrFunc.RowsAffected()

	qrDelProduct, err := r.DB.Query(queryDelObjProductId)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qrDelProduct.Close()

	qrFuncProductId, err := r.DB.ExecContext(ctx, queryGetProductId)
	if err != nil {
		if err.Error() != "sql: function no create" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qrFuncProductId.RowsAffected()

	qr, err := r.DB.QueryContext(ctx, querySelectProductId)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qr.Close()

	var qrs []ProductSetup
	for qr.Next() {
		var qrts ProductSetup
		if err := qr.Scan(
			&qrts.Product_Id,
		); err != nil {
			return qrs, err
		}
		qrs = append(qrs, qrts)
	}
	if err = qr.Err(); err != nil {
		return qrs, err
	}
	return qrs, nil
}

const querySelectProductName = `
SELECT [PRODUCT_NAME]
FROM [PanaCIM].[dbo].[product_data]
where [PRODUCT_ID] =
`

func (r PanaCIMStorage) GetProductName(productid string) ([]ProductData, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	qr, err := r.DB.QueryContext(ctx, querySelectProductName+productid)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qr.Close()

	var qrs []ProductData
	for qr.Next() {
		var qrts ProductData
		if err := qr.Scan(
			&qrts.ProductName,
		); err != nil {
			return qrs, err
		}
		qrs = append(qrs, qrts)
	}
	if err = qr.Err(); err != nil {
		return qrs, err
	}
	return qrs, nil
}

const querySelectRouteId = `
SELECT [ROUTE_ID]
FROM [PanaCIM].[dbo].[product_setup]
WHERE PRODUCT_ID = `

//@p1
//order by LAST_MODIFIED_TIME desc
//`

func (r PanaCIMStorage) GetRouteId(productid string) ([]ProductSetup, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//productSetup := &ProductSetup{}
	//productSetupList := make(ProductDataLink, 0)
	/*if err := r.DB.QueryRowContext(ctx, querySelectRouteId+productid).Scan(
		&productSetup.Route_Id,
	); err != nil {
		//if err == r.q//sql.ErrNoRows {
		return "", errors.New("record not found") //store.ErrRecordNotFound
		//}
		//	return nil, err
	}*/

	qr, err := r.DB.QueryContext(ctx, querySelectRouteId+productid+`order by LAST_MODIFIED_TIME desc`)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qr.Close()

	var qrs []ProductSetup
	for qr.Next() {
		var qrts ProductSetup
		if err := qr.Scan(
			&qrts.Route_Id,
		); err != nil {
			return qrs, err
		}
		qrs = append(qrs, qrts)
	}
	if err = qr.Err(); err != nil {
		return qrs, err
	}
	return qrs, nil
	//fmt.Printf("productSetup.Route_Id :%v", productSetup.Route_Id)
	//return productSetup.Route_Id, nil
}

const queryDelObjInfoInstallJobId_View = `
IF OBJECT_ID('dbo.InfoInstallLastJobId_View', 'V') IS NOT NULL DROP VIEW dbo.InfoInstallLastJobId_View
`

const queryCreateInfoInstallJobId_View1 = `
CREATE VIEW dbo.InfoInstallLastJobId_View
AS
SELECT
    [PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID],
[PanaCIM].[dbo].[reel_data].PART_NO,
[PanaCIM].[dbo].[reel_data].LOT_NO,
SUM([PanaCIM].[dbo].[Z_CASS_VIEW].PLACE_COUNT) AS PLACE_COUNT,
SUM([PanaCIM].[dbo].[Z_CASS_VIEW].PICKUP_COUNT) AS PICKUP_COUNT,
[PanaCIM].[dbo].[REEL_DATA_VIEW].reel_barcode,
[PanaCIM].[dbo].[reel_data].CURRENT_QUANTITY,
[PanaCIM].[dbo].[reel_data].QUANTITY AS INITIAL_QUANTITY
--(SELECT * FROM dbo.SumInstallComponent([PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID]))
  FROM [PanaCIM].[dbo].[Z_CASS_VIEW]
 -- LEFT JOIN ( SELECT * FROM dbo.SumInstallComponent([PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID]))
  LEFT JOIN [PanaCIM].[dbo].[REEL_DATA_VIEW]
  ON [PanaCIM].[dbo].[REEL_DATA_VIEW].[reel_id] = [PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID]
  LEFT JOIN [PanaCIM].[dbo].[reel_data]
  ON [PanaCIM].[dbo].[REEL_DATA_VIEW].[reel_id] = [PanaCIM].[dbo].[reel_data].REEL_ID
  --where [PanaCIM].[dbo].[Z_CASS_VIEW].JOB_ID = (SELECT * FROM dbo.GetLastJobId()) AND [PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID] IS NOT NULL
  where [PanaCIM].[dbo].[Z_CASS_VIEW].JOB_ID = `
const queryCreateInfoInstallJobId_View2 = `AND [PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID] IS NOT NULL
  group by [PanaCIM].[dbo].[REEL_DATA_VIEW].reel_barcode, 
  [PanaCIM].[dbo].[Z_CASS_VIEW].REEL_ID, 
  [PanaCIM].[dbo].[reel_data].CURRENT_QUANTITY, 
  [PanaCIM].[dbo].[reel_data].QUANTITY,
  [PanaCIM].[dbo].[reel_data].PART_NO,
  [PanaCIM].[dbo].[reel_data].LOT_NO
`
const querySelectInfoInstallJobId_View = `
SELECT
    [PART_NO]
    ,SUM([PLACE_COUNT]) AS SUM_PLACE_COUNT_ALL
FROM [PanaCIM].[dbo].[InfoInstallLastJobId_View]
group by PART_NO;
`

func (r PanaCIMStorage) GetPanacimDataComponentsByJobId(jobid string) ([]InfoInstallLastJobId_View, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	qrDel, err := r.DB.Query(queryDelObjInfoInstallJobId_View)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qrDel.Close()

	qrFunc, err := r.DB.ExecContext(ctx, queryCreateInfoInstallJobId_View1+jobid+queryCreateInfoInstallJobId_View2)
	if err != nil {
		if err.Error() != "sql: function no create" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qrFunc.RowsAffected()

	qr, err := r.DB.QueryContext(ctx, querySelectInfoInstallJobId_View)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qr.Close()

	var qrs []InfoInstallLastJobId_View
	for qr.Next() {
		var qrts InfoInstallLastJobId_View
		if err := qr.Scan(
			&qrts.PartNo,
			&qrts.SumPlaceCount,
		); err != nil {
			return qrs, err
		}
		qrs = append(qrs, qrts)
	}
	if err = qr.Err(); err != nil {
		return qrs, err
	}
	return qrs, nil
}

func (r PanaCIMStorage) WtitePanaCIMDataComponentsToFile(in []InfoInstallLastJobId_View) (err error) {

	panaCIMpath := os.Getenv("panacim")

	panacimFileRemove := panaCIMpath

	if utils.FileExists(panacimFileRemove) {
		os.Remove(panacimFileRemove)
	}

	var partNO string = `PART_NO`
	var sumPlaceCount string = `SUM_PLACE_COUNT`
	panacimFile := panaCIMpath
	if _, err := os.Stat(panacimFile); os.IsNotExist(err) {
		panaFile, err := os.Create(panacimFile)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer panaFile.Close()

		writer := csv.NewWriter(panaFile)
		writer.Write([]string{partNO + `,` + sumPlaceCount})
		writer.Comma = ','
		writer.Flush()
	}

	splitPanaCIM, err := os.OpenFile(panacimFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		r.logger.Errorf(err.Error()) //logger.Errorf(err.Error())
		return nil
	}
	defer splitPanaCIM.Close()

	for _, i := range in {
		var result = []string{i.PartNo + "," + i.SumPlaceCount}

		for _, v := range result {
			_, err = fmt.Fprintln(splitPanaCIM, v)
			if err != nil {
				splitPanaCIM.Close()
				return nil
			}
		}
	}
	return nil

}

const querySelectGetMixName1 = `
SELECT TOP 1 [MIX_NAME]
FROM [PanaCIM].[dbo].[product_setup]
WHERE [PRODUCT_ID] =
`
const querySelectGetMixName2 = `
order by LAST_MODIFIED_TIME desc
`

func (r PanaCIMStorage) GetPanaCIMixName(productid string) ([]ProductSetup, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	qr, err := r.DB.QueryContext(ctx, querySelectGetMixName1+productid+querySelectGetMixName2)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qr.Close()

	var qrs []ProductSetup
	for qr.Next() {
		var qrts ProductSetup
		if err := qr.Scan(
			&qrts.MixName,
		); err != nil {
			return qrs, err
		}
		qrs = append(qrs, qrts)
	}
	if err = qr.Err(); err != nil {
		return qrs, err
	}
	return qrs, nil
}

const querySelecGetParts = `
SELECT [PRIMARY_PN]
,[SUBSTITUTE_PN]
FROM [PanaCIM].[dbo].[substitute_parts]
WHERE [MIX_NAME] = `

func (r PanaCIMStorage) GetPanaCIMParts(mixname string) ([]SubstituteParts, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//qr, err := r.DB.QueryContext(ctx, querySelecGetParts+mixname)
	qr, err := r.DB.QueryContext(ctx, querySelecGetParts+"'"+mixname+"'")
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qr.Close()

	var qrs []SubstituteParts
	for qr.Next() {
		var qrts SubstituteParts
		if err := qr.Scan(
			&qrts.PrimaryPn,
			&qrts.SubstitutePn,
		); err != nil {
			return qrs, err
		}
		qrs = append(qrs, qrts)
	}
	if err = qr.Err(); err != nil {
		return qrs, err
	}
	return qrs, nil
}

func (r PanaCIMStorage) WritePanaCIMPartsToFile(in []SubstituteParts) (err error) {

	substitutepath := os.Getenv("substitute")

	substituteFileRemove := substitutepath

	if utils.FileExists(substituteFileRemove) {
		os.Remove(substituteFileRemove)
	}

	var primaryPN string = "PRIMARY_PN"
	var substitutePN string = "SUBSTITUTE_PN"
	substituteFile := substitutepath
	if _, err := os.Stat(substituteFile); os.IsNotExist(err) {
		substituteFile, err := os.Create(substituteFile)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer substituteFile.Close()

		writer := csv.NewWriter(substituteFile)
		writer.Write([]string{primaryPN + `,` + substitutePN})
		writer.Comma = ','
		writer.Flush()
	}

	splitPanaCIMParts, err := os.OpenFile(substituteFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		r.logger.Errorf(err.Error()) //logger.Errorf(err.Error())
		return nil
	}
	defer splitPanaCIMParts.Close()

	for _, i := range in {
		var result = []string{i.PrimaryPn + "," + i.SubstitutePn}

		for _, v := range result {
			_, err = fmt.Fprintln(splitPanaCIMParts, v)
			if err != nil {
				splitPanaCIMParts.Close()
				return nil
			}
		}
	}
	return nil
}

const querySelectUnixTimeWO = `
SELECT TOP 1 [JOB_ID]
,[EQUIPMENT_ID]
,[SETUP_ID]
,[START_TIME]
,[END_TIME]
,[CLOSING_TYPE]
,[START_OPERATOR_ID]
,[END_OPERATOR_ID]
,[TFR_REASON]
,[LANE_NO]
FROM [PanaCIM].[dbo].[job_history] 
where [JOB_ID] = `

// получаем старт и конец сборки WO в unix-формате
func (r PanaCIMStorage) GetUnixTimeWO(jobid string) ([]Job_History, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	qr, err := r.DB.QueryContext(ctx, querySelectUnixTimeWO+jobid+`order by END_TIME desc`)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return nil, err
		}
	}
	defer qr.Close()

	var qrs []Job_History
	for qr.Next() {
		var qrts Job_History
		if err := qr.Scan(
			&qrts.JOB_ID,
			&qrts.EQUIPMENT_ID,
			&qrts.SETUP_ID,
			&qrts.StartUnixTimeWO,
			&qrts.EndUnixTimeWO,
			&qrts.CLOSING_TYPE,
			&qrts.START_OPERATOR_ID,
			&qrts.END_OPERATOR_ID,
			&qrts.TFR_REASON,
			&qrts.LANE_NO,
		); err != nil {
			return qrs, err
		}
		qrs = append(qrs, qrts)
	}
	if err = qr.Err(); err != nil {
		return qrs, err
	}
	return qrs, nil
}

type U03 struct {
	B            string
	IDNUM        string
	TURN         string
	MS           string
	TS           string
	FAdd         string
	FSAdd        string
	FBLKCode     string
	FBLKSerial   string
	NHAdd        string
	NCAdd        string
	NBLKCode     string
	NBLKSerial   string
	ReelID       string
	F            string
	RCGX         string
	RCGY         string
	RCGA         string
	TCX          string
	TCY          string
	MPosiRecX    string
	MPosiRecY    string
	MPosiRecA    string
	MPosiRecZ    string
	THMAX        string
	THAVE        string
	MNTCX        string
	MNTCY        string
	MNTCA        string
	TLX          string
	TLY          string
	InspectArea  string
	DIDNUM       string
	DS           string
	DispenseID   string
	PARTS        string
	WarpZ        string
	PrePickupLOT string
	PrePickupSTS string
	LoadSV       string
	LoadMV       string
	ReachCZ      string
}

var (
	layoutDate string = "2006/01/02,15:04:05"
)

func (r PanaCIMStorage) GetSumPCBFromU03V2(startUnixTimeWO, finishUnixTimeWO, npm string) (sumstrPCBOrder string) {
	logger := logging.GetLogger()
	npmToUp := strings.ToUpper(npm)
	// конвертация unix времени страта и завершения сборки WO

	tStartWO, err := strconv.ParseInt(startUnixTimeWO, 10, 64)
	if err != nil {
		logger.Panic(err)
		panic(err)
	}
	tmStartWO := time.Unix(tStartWO, 0)
	p_tmStartWO, _ := time.Parse(layoutDate, tmStartWO.Format(layoutDate))
	fmt.Println("p_tmStartWO:", p_tmStartWO)

	tFinishWO, err := strconv.ParseInt(finishUnixTimeWO, 10, 64)
	if err != nil {
		panic(err)
	}
	tmFinishWO := time.Unix(tFinishWO, 0)
	p_tmFinishWO, _ := time.Parse(layoutDate, tmFinishWO.Format(layoutDate))
	fmt.Println("p_tmFinishWO: ", p_tmFinishWO)

	folderFromPanaNPM_1 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resourcePanaCIM/NPM-1/"
	folderFromPanaNPM_2 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resourcePanaCIM/NPM-2/"
	folderFromPanaNPM_3 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resourcePanaCIM/NPM-3/"
	folderFromPanaNPM_4 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resourcePanaCIM/NPM-4/"

	folderToCopyNPM_1 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resource/NPM-1/"
	folderToCopyNPM_2 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resource/NPM-2/"
	folderToCopyNPM_3 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resource/NPM-3/"
	folderToCopyNPM_4 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resource/NPM-4/"

	err = cp.Copy(folderFromPanaNPM_1, folderToCopyNPM_1)
	if err != nil {
		log.Println(err)
	}
	err = cp.Copy(folderFromPanaNPM_2, folderToCopyNPM_2)
	if err != nil {
		log.Println(err)
	}
	err = cp.Copy(folderFromPanaNPM_3, folderToCopyNPM_3)
	if err != nil {
		log.Println(err)
	}
	err = cp.Copy(folderFromPanaNPM_4, folderToCopyNPM_4)
	if err != nil {
		log.Println(err)
	}

	// получить список папок в скопированной директории
	resourcePath := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resource/"
	f_resource, err := ioutil.ReadDir(resourcePath)
	if err != nil {
		log.Panicln(err)
	}
	// Разархивирование архивов, если они есть
	for _, r := range f_resource {
		if r.IsDir() {
			folder, err := ioutil.ReadDir(resourcePath + r.Name())
			if err != nil {
				log.Panicln(err)
			}
			fmt.Println("folder Resource: ", resourcePath+r.Name())
			for _, r2 := range folder {
				if r2.IsDir() {
					processed, err := ioutil.ReadDir(filepath.FromSlash(resourcePath + r.Name() + "/" + r2.Name()))
					if err != nil {
						log.Panicln(err)
					}
					fmt.Println("folder Resource2: ", resourcePath+r.Name()+"/"+r2.Name())
					for _, r3 := range processed {
						dataFolders, err := ioutil.ReadDir(resourcePath + r.Name() + "/" + r2.Name() + "/" + r3.Name())
						if err != nil {
							log.Panicln(err)
						}
						for _, r4 := range dataFolders {
							if !r4.IsDir() {
								if strings.Contains(r4.Name(), ".gz") {
									// Open compressed file
									gzipFile, err := os.Open(resourcePath + r.Name() + "/" + r2.Name() + "/" + r3.Name() + "/" + r4.Name())
									if err != nil {
										log.Fatal(err)
									}
									// Create a gzip reader on top of the file reader
									// Again, it could be any type reader though
									gzipReader, err := gzip.NewReader(gzipFile)
									if err != nil {
										log.Fatal(err)
									}
									//defer gzipReader.Close()
									writeToFile := strings.Trim(r4.Name(), ".gz")
									// Uncompress to a writer. We'll use a file writer
									outfileWriter, err := os.Create(resourcePath + r.Name() + "/" + r2.Name() + "/" + r3.Name() + "/" + writeToFile)
									if err != nil {
										log.Fatal(err)
									}
									//defer outfileWriter.Close()

									// Copy contents of gzipped file to output file
									_, err = io.Copy(outfileWriter, gzipReader)
									if err != nil {
										log.Fatal(err)
									}
									gzipReader.Close()
									outfileWriter.Close()
								}
							}
						}
					}
				}
			}
		}
	}

	sumPCBOrder := 0
	checkDuble := map[string]bool{}
	inputCoreFolder, err := ioutil.ReadDir(resourcePath)
	if err != nil {
		log.Fatal(err)
	}
	for _, npmf := range inputCoreFolder {
		if npmf.IsDir() {
			// расчет кол-ва плат
			if npmf.Name() == "NPM-1" {
				fmt.Println("NPM-1 Great!!!")
				processedf, err := ioutil.ReadDir(resourcePath + npmf.Name())
				if err != nil {
					log.Fatal(err)
				}
				for _, processed := range processedf {
					fmt.Println("2 NPM-1 Great!!!")
					if processed.IsDir() {
						dataf, err := ioutil.ReadDir(resourcePath + npmf.Name() + "/" + processed.Name())
						if err != nil {
							log.Fatal(err)
						}
						for _, data := range dataf {
							fmt.Println("3 NPM-1 Great!!!")
							if data.IsDir() {
								fileu03f, err := ioutil.ReadDir(resourcePath + npmf.Name() + "/" + processed.Name() + "/" + data.Name())
								if err != nil {
									log.Fatal(err)
								}
								for _, fileu03 := range fileu03f {
									if !fileu03.IsDir() {
										//fmt.Println("Finish")
										if strings.Contains(fileu03.Name(), ".u03") && !strings.Contains(fileu03.Name(), ".gz") {
											cfg, err := ini.LoadSources(ini.LoadOptions{
												UnparseableSections: []string{
													//	"Index",
													//	"Information",

													"BRecg",
													"BRecgCalc",
													"ElapseTimeRecog",
													"SBoard",
													"HeightCorrect",
													"MountNormalTrace",
													"MountLatestReel",
													"MountExchangeReel",
													"MountQualityTrace"},
											}, resourcePath+"/"+npmf.Name()+"/"+processed.Name()+"/"+data.Name()+"/"+fileu03.Name())
											if err != nil {
												fmt.Printf("Fail to read file: %v", err)
												os.Exit(1)
											}
											dataFile := cfg.Section("Index").Key("Date").String()
											pdataFile, _ := time.Parse(layoutDate, dataFile)
											if (pdataFile.After(p_tmStartWO) && pdataFile.Before(p_tmFinishWO)) &&
												(strings.EqualFold(strings.ToUpper(cfg.Section("Information").Key("LotName").String()), strings.ToUpper(npmToUp))) {
												// проверка на дубль в файлах по ключу Code checkDuble[cfg.Section("Information").Key("Code").String()] == true
												if checkDuble[cfg.Section("Information").Key("Code").String()] {
													fmt.Println("Code Double: ", cfg.Section("Information").Key("Code").String())
												} else {
													checkDuble[cfg.Section("Information").Key("Code").String()] = true
													// получаем данные по указаннной секции
													readSection := cfg.Section("MountQualityTrace").Body()
													// конвертируем в байты
													writeByte := []byte(readSection)
													// записываем данные в файл
													if err := ioutil.WriteFile("internal/out", writeByte, 0644); err != nil {
														fmt.Printf("%v", err)
													}
													// Пересоздаю данный файл, чтобы его сумму хранить дальше в памяти
													filepcb := "internal/pcb"

													fileRemovePCB := "internal/pcb"

													if _, err := os.Stat(fileRemovePCB); os.IsNotExist(err) {
														pcbFile, err := os.Create(fileRemovePCB)
														if err != nil {
															logger.Errorf(err.Error())
														}
														defer pcbFile.Close()
													}
													e_pcb := os.Remove(fileRemovePCB)
													if e_pcb != nil {
														log.Fatal(e_pcb)
													}
													filepcbRW, err := os.OpenFile(filepcb, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
													if err != nil {
														logger.Errorf(err.Error())
														return
													}
													defer filepcbRW.Close()

													numberPCBs := filereader.Readfileseekerspace("internal/out")
													for _, i := range numberPCBs {
														data := U03{
															B: i[0],
														}
														if data.B != "0" {
															var result = []string{data.B}
															for _, v := range result {
																_, err = fmt.Fprintln(filepcbRW, v)
																if err != nil {
																	filepcbRW.Close()
																	return
																}
															}
														}
													}
													// читаю файл построчно
													getnumberpcb, err := readLines(filepcb)
													if err != nil {
														logger.Errorf(err.Error())
													}
													// убираю дублированные номера
													resnumberpcb := removeDuplicatesinfile(getnumberpcb)
													// считаю уникальные номера плат и записываю в счетчик sumPCBOrder
													for i := 0; i < len(resnumberpcb); i++ {
														//for _, i := range resnumberpcb {
														//fmt.Println("t ", i)
														sumPCBOrder++
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	sumstrPCBOrder = strconv.Itoa(sumPCBOrder)

	return sumstrPCBOrder
}

func (r PanaCIMStorage) GetSumPCBFromU03(startUnixTimeWO, finishUnixTimeWO, npm string) (sumstrPCBOrder string) {
	npmToUp := strings.ToUpper(npm)
	tStartWO, err := strconv.ParseInt(startUnixTimeWO, 10, 64)
	if err != nil {
		r.logger.Fatalf(err.Error()) //panic(err)
	}
	tmStartWO := time.Unix(tStartWO, 0)
	fmt.Println("tmStartWO:", tmStartWO)
	// добавляем +3 ч 15 мин от GMT
	chH_tmStartWO := tmStartWO.Add(time.Hour * 3)
	chM_tmStartWO := chH_tmStartWO.Add(time.Minute * 15)
	fmt.Printf("chM_tmStartWO: %v\n", chM_tmStartWO)

	tFinishWO, err := strconv.ParseInt(finishUnixTimeWO, 10, 64)
	if err != nil {
		r.logger.Fatalf(err.Error()) //panic(err)
	}
	tmFinishWO := time.Unix(tFinishWO, 0)
	fmt.Println("tmFinishWO: ", tmFinishWO)
	// добавляем +3 ч 15 мин от GMT
	chH_tmFinishWO := tmFinishWO.Add(time.Hour * 3)
	chM_tmFinishWO := chH_tmFinishWO.Add(time.Minute * 15)
	fmt.Printf("chM_tmFinishWO: %v\n", chM_tmFinishWO)
	// копирование директории файлов для
	folderFrom := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/u03/processed/"
	folderToCopy := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/u03/processed_copy/"
	err = cp.Copy(folderFrom, folderToCopy)
	if err != nil {
		fmt.Println(err)
	}

	// получить список папок в скопированной директории
	files, err := ioutil.ReadDir(folderToCopy)
	if err != nil {
		r.logger.Fatalf(err.Error()) //log.Fatal(err)
	}
	// создаем копию директории и разархивируем архивы, которые там есть
	for _, f := range files {
		if f.IsDir() {
			folder, err := ioutil.ReadDir(folderToCopy + f.Name())
			if err != nil {
				r.logger.Fatalf(err.Error())
			}
			for _, g := range folder {
				if !g.IsDir() {

					if strings.Contains(g.Name(), ".gz") {

						// fmt.Printf("папка - %v, файл - %v \n", f.Name(), g.Name())
						// Open compressed file
						gzipFile, err := os.Open(folderToCopy + f.Name() + "/" + g.Name())
						if err != nil {
							r.logger.Fatalf(err.Error())
						}

						// Create a gzip reader on top of the file reader
						// Again, it could be any type reader though
						gzipReader, err := gzip.NewReader(gzipFile)
						if err != nil {
							r.logger.Fatalf(err.Error())
						}
						//defer gzipReader.Close()

						writeToFile := strings.Trim(g.Name(), ".gz")
						fmt.Printf("writeToFile: %v\n", writeToFile)
						fmt.Printf("path %v, %v, %v\n", folderToCopy, f.Name(), writeToFile)
						/*_, err = os.Stat(folderToCopy + "/" + f.Name() + "/" + writeToFile)
						if err == nil {
							fmt.Printf("File %s already exists.", folderToCopy+"/"+f.Name()+"/"+writeToFile)
						}*/
						// Uncompress to a writer. We'll use a file writer

						outfileWriter, err := os.Create(folderToCopy + f.Name() + "/" + writeToFile)
						fmt.Println("outfileWriter:", &outfileWriter)

						if err != nil {
							// r.logger.Fatalf(err.Error())
							fmt.Printf("error: %v\n", err)
							r.logger.Printf("error: %v\n", err)
							r.logger.Debugf(err.Error())
							//r.logger.Errorf(err.Error())
						}

						// Copy contents of gzipped file to output file
						_, err = io.Copy(outfileWriter, gzipReader)
						if err != nil {
							r.logger.Fatalf(err.Error())
						}
						gzipReader.Close()
						outfileWriter.Close()
					}
				}
			}
		}
	}

	iniFolders, err := ioutil.ReadDir(folderToCopy)
	if err != nil {
		r.logger.Fatalf(err.Error())
	}
	sumPCBOrder := 0
	for _, ff := range iniFolders {
		if ff.IsDir() {
			iniFolder, err := ioutil.ReadDir(folderToCopy + ff.Name())
			if err != nil {
				r.logger.Fatalf(err.Error())
			}
			for _, gg := range iniFolder {
				if !gg.IsDir() {
					if strings.Contains(gg.Name(), ".u03") && !strings.Contains(gg.Name(), ".gz") {
						//ini.LoadSources
						cfg, err := ini.LoadSources(ini.LoadOptions{
							UnparseableSections: []string{
								//	"Index",
								//	"Information",

								"BRecg",
								"BRecgCalc",
								"ElapseTimeRecog",
								"SBoard",
								"HeightCorrect",
								"MountNormalTrace",
								"MountLatestReel",
								"MountExchangeReel",
								"MountQualityTrace"},
						}, folderToCopy+"/"+ff.Name()+"/"+gg.Name())
						if err != nil {
							fmt.Printf("Fail to read file: %v", err)
							os.Exit(1)
						}
						lotnameToUpper := strings.ToUpper(cfg.Section("Information").Key("LotName").String())
						// if cfg.Section("Information").Key("LotName").String() == "NPM_915-00211_A_S"
						if lotnameToUpper == npmToUp {
							fmt.Printf("NPM_ %v, папка - %v, файл - %v \n", npmToUp, ff.Name(), gg.Name())
							fmt.Println("Date", cfg.Section("Index").Key("Date").Value())
							dataFile := cfg.Section("Index").Key("Date").Value()
							pdataFile, _ := time.Parse(layoutDate, dataFile)
							fmt.Println("pdataFile: ", pdataFile.Format(layoutDate))
							if pdataFile.After(chM_tmStartWO.AddDate(0, 0, -1)) && pdataFile.Before(chM_tmFinishWO.AddDate(0, 0, +1)) {
								// получаем данные по указаннной секции
								readSection := cfg.Section("MountQualityTrace").Body()
								// конвертируем в байты
								writeByte := []byte(readSection)
								// записываем данные в файл
								if err := ioutil.WriteFile("/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/u03/out", writeByte, 0644); err != nil {
									fmt.Printf("%v", err)
								}

								filepcb := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/u03/pcb"

								fileRemovePCB := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/u03/pcb"

								if _, err := os.Stat(fileRemovePCB); os.IsNotExist(err) {
									pcbFile, err := os.Create(fileRemovePCB)
									if err != nil {
										r.logger.Errorf(err.Error())
									}
									defer pcbFile.Close()
								}

								e_pcb := os.Remove(fileRemovePCB)
								if e_pcb != nil {
									r.logger.Errorf(e_pcb.Error()) // log.Fatal(e_pcb)
								}
								filepcbRW, err := os.OpenFile(filepcb, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
								if err != nil {
									r.logger.Errorf(err.Error())
									return
								}
								defer filepcbRW.Close()

								/*outcheck := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/u03/out"
								if _, err := os.Stat(outcheck); os.IsNotExist(err) {
									outFile, err := os.Create(outcheck)
									if err != nil {
										r.logger.Errorf(err.Error())
									}
									defer outFile.Close()
								}*/
								numberPCBs := filereader.Readfileseekerspace("/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/u03/out")

								for _, i := range numberPCBs {
									data := U03{
										B: i[0],
									}
									// убираем строки с нулями
									if data.B != "0" {
										//fmt.Println("f", data.B)
										var result = []string{data.B}
										/*bytePCB := []byte(data.B)
										if err := ioutil.WriteFile("internal/pcb", bytePCB, 0644); err != nil {
											fmt.Printf("%v", err)
										}*/
										//fmt.Println(i[0][1])
										for _, v := range result {
											_, err = fmt.Fprintln(filepcbRW, v)
											if err != nil {
												filepcbRW.Close()
												return
											}
										}
									}
								}

								/*e_out := os.Remove(outcheck)
								if e_out != nil {
									r.logger.Errorf(e_out.Error())
								}*/

								getnumberpcb, err := readLines(filepcb)
								if err != nil {
									r.logger.Errorf(err.Error())
								}

								resnumberpcb := removeDuplicatesinfile(getnumberpcb)
								//fmt.Println("rrr ", resnumberpcb)
								//sum := 0
								//for _, i := range resnumberpcb {
								for i := 0; i < len(resnumberpcb); i++ {
									//fmt.Println("t ", i)

									sumPCBOrder++
								}

							}
						}

					}
				}
			}
		}
	}

	sumstrPCBOrder = strconv.Itoa(sumPCBOrder)

	return sumstrPCBOrder
}

type ReelIdData struct {
	ReelID string
	Qty    string
}

type ReelIdDates struct {
	ReelIdData []ReelIdData
}

const querySelectReelData = `
SELECT [PART_NO]
,[REEL_BARCODE]
,[LOT_NO]
FROM [PanaCIM].[dbo].[reel_data]
where REEL_BARCODE = `

func (r PanaCIMStorage) GetSumComponentFromU03(startUnixTimeWO, finishUnixTimeWO, npm string) {
	logger := logging.GetLogger()
	npmToUp := strings.ToUpper(npm)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// конвертация unix времени страта и завершения сборки WO

	tStartWO, err := strconv.ParseInt(startUnixTimeWO, 10, 64)
	if err != nil {
		logger.Panic(err)
		panic(err)
	}
	tmStartWO := time.Unix(tStartWO, 0)
	p_tmStartWO, _ := time.Parse(layoutDate, tmStartWO.Format(layoutDate))
	fmt.Println("p_tmStartWO:", p_tmStartWO)

	tFinishWO, err := strconv.ParseInt(finishUnixTimeWO, 10, 64)
	if err != nil {
		panic(err)
	}
	tmFinishWO := time.Unix(tFinishWO, 0)
	p_tmFinishWO, _ := time.Parse(layoutDate, tmFinishWO.Format(layoutDate))
	fmt.Println("p_tmFinishWO: ", p_tmFinishWO)
	// получить список папок в скопированной директории
	resourcePath := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resource/"
	// проверяем, не создавался ли ранее файл reel_id
	fileReelId := "internal/reel_id"
	if _, err := os.Stat(fileReelId); os.IsNotExist(err) {
		reelidFile, err := os.Create(fileReelId)
		if err != nil {
			logger.Errorf(err.Error())
		}
		defer reelidFile.Close()
	}
	// удаляем, если ранее был создан
	rm_reedid := os.Remove(fileReelId)
	if rm_reedid != nil {
		logger.Fatalf("%v\n", rm_reedid)
	}

	fileReedIdScrap := "internal/reel_id_scrap"
	if _, err := os.Stat(fileReedIdScrap); os.IsNotExist(err) {
		reedIdScrapFile, err := os.Create(fileReedIdScrap)
		if err != nil {
			logger.Errorf(err.Error())
		}
		defer reedIdScrapFile.Close()
	}

	rm_reedIdScrapFile := os.Remove(fileReedIdScrap)
	if rm_reedIdScrapFile != nil {
		logger.Fatalf("%v\n", rm_reedIdScrapFile)
	}

	//checkDubleComponent := map[string]bool{}
	inputCoreFolder, err := ioutil.ReadDir(resourcePath)
	if err != nil {
		log.Fatal(err)
	}
	for _, npmf := range inputCoreFolder {
		if npmf.IsDir() {
			processedf, err := ioutil.ReadDir(resourcePath + npmf.Name())
			if err != nil {
				log.Fatal(err)
			}
			for _, processed := range processedf {
				fmt.Println("COMPONENT NPM-1 Great!!!")
				if processed.IsDir() {

					dataf, err := ioutil.ReadDir(resourcePath + npmf.Name() + "/" + processed.Name())
					if err != nil {
						log.Fatal(err)
					}
					for _, data := range dataf {
						fmt.Printf("COMPONENT %v Great!!!\n", data.Name())
						if data.IsDir() {
							fileu03f, err := ioutil.ReadDir(resourcePath + npmf.Name() + "/" + processed.Name() + "/" + data.Name())
							if err != nil {
								log.Fatal(err)
							}
							for _, fileu03 := range fileu03f {
								if !fileu03.IsDir() {
									if strings.Contains(fileu03.Name(), ".u03") && !strings.Contains(fileu03.Name(), ".gz") {
										cfg, err := ini.LoadSources(ini.LoadOptions{
											UnparseableSections: []string{
												//	"Index",
												//	"Information",

												"BRecg",
												"BRecgCalc",
												"ElapseTimeRecog",
												"SBoard",
												"HeightCorrect",
												"MountNormalTrace",
												"MountLatestReel",
												"MountExchangeReel",
												"MountQualityTrace"},
										}, resourcePath+"/"+npmf.Name()+"/"+processed.Name()+"/"+data.Name()+"/"+fileu03.Name())
										if err != nil {
											fmt.Printf("Fail to read file: %v", err)
											os.Exit(1)
										}
										dataFile := cfg.Section("Index").Key("Date").String()
										pdataFile, _ := time.Parse(layoutDate, dataFile)
										if (pdataFile.After(p_tmStartWO) && pdataFile.Before(p_tmFinishWO)) &&
											(strings.EqualFold(strings.ToUpper(cfg.Section("Information").Key("LotName").String()), strings.ToUpper(npmToUp))) {
											// проверка на дубль в файлах по ключу Code checkDuble[cfg.Section("Information").Key("Code").String()] == true
											//if checkDubleComponent[cfg.Section("Information").Key("Code").String()] {
											//	fmt.Println("Code Double: ", cfg.Section("Information").Key("Code").String())
											//} else {
											//checkDubleComponent[cfg.Section("Information").Key("Code").String()] = true
											// получаем данные по указаннной секции
											readSection := cfg.Section("MountQualityTrace").Body()
											// конвертируем в байты
											writeByte := []byte(readSection)
											// записываем данные в файл
											if err := ioutil.WriteFile("internal/out", writeByte, 0644); err != nil {
												fmt.Printf("%v", err)
											}

											fileReelIdRW, err := os.OpenFile(fileReelId, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
											if err != nil {
												logger.Errorf(err.Error())
												return
											}
											//defer fileReelIdRW.Close()

											fileReedIdScrapRW, err := os.OpenFile(fileReedIdScrap, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
											if err != nil {
												logger.Errorf(err.Error())
												return
											}
											//defer fileReedIdScrapRW.Close()

											numberPCBs := filereader.Readfileseekerspace("internal/out")
											for _, i := range numberPCBs {
												data := U03{
													B:      i[0],
													ReelID: i[13],
													F:      i[14],
												}
												if data.F == "0" {
													var result = []string{data.ReelID}
													for _, v := range result {
														_, err := fmt.Fprintln(fileReelIdRW, v)
														if err != nil {
															fileReelIdRW.Close()
															return
														}
													}
												}
												// запись reel_id со скрапом
												if data.F == "2" {
													var result = []string{data.ReelID}
													for _, v := range result {
														_, err := fmt.Fprintln(fileReedIdScrapRW, v)
														if err != nil {
															fileReedIdScrapRW.Close()
															return
														}
													}
												}

												//}
											}

											fileReelIdRW.Close()
											fileReedIdScrapRW.Close()
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	file_reelid_unic := "internal/reelid_unic"
	//var reelid_unicFile *os.File
	if _, err := os.Stat(file_reelid_unic); os.IsNotExist(err) {
		reelid_unicFile, err := os.Create(file_reelid_unic)
		if err != nil {
			logger.Errorf(err.Error())
		}
		defer reelid_unicFile.Close()
	}

	rm_reelid_unicFile := os.Remove(file_reelid_unic)
	if rm_reelid_unicFile != nil {
		logger.Fatalf("%v\n", rm_reelid_unicFile)
	}
	// читаю файл построчно
	get_reel_id, err := readLines(fileReelId)
	if err != nil {
		logger.Errorf(err.Error())
	}
	// убираю дублированные номера
	rmdreelid_unic := removeDuplicatesinfile(get_reel_id)

	fileReelIdUnic, err := os.OpenFile(file_reelid_unic, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error creating file %v", err)
	}

	datawriter := bufio.NewWriter(fileReelIdUnic)

	for _, data := range rmdreelid_unic {
		_, err := datawriter.WriteString(data + "\n")
		if err != nil {
			logger.Errorf(err.Error())
		}
	}

	datawriter.Flush()
	fileReelIdUnic.Close()

	arrReelIdALL := filereader.Readfile(fileReelId)
	arrReelIdUnic := filereader.Readfile(file_reelid_unic)

	//var reelIdSumValue []Filds
	reelIdStore := []ReelIdData{}
	for _, i := range arrReelIdUnic {
		sum := 0
		for _, j := range arrReelIdALL {

			if i[0] == j[0] {
				sum += 1
			}
		}
		//reelIdSumValue = append(reelIdSumValue, i[0]+","+strconv.Itoa(sum)+"\n")
		reelIdStore = append(reelIdStore, ReelIdData{ReelID: i[0], Qty: strconv.Itoa(sum)})
		fmt.Printf("reel_id: %v, sum: %v\n", i[0], sum)
	}
	fmt.Printf("reelIdStore: %v\n", reelIdStore)
	// проверяем, создавался ли файл ранее
	fileReelIDScrap_unic := "internal/reelid_scrap_unic"
	if _, err := os.Stat(fileReelIDScrap_unic); os.IsExist(err) {
		reelid_scrap_unicFile, err := os.Create(fileReelIDScrap_unic)
		if err != nil {
			logger.Errorf(err.Error())
		}
		defer reelid_scrap_unicFile.Close()
	}
	/*
		rm_reelid_scrap_unicFile := os.Remove(fileReelIDScrap_unic)
		if rm_reelid_scrap_unicFile != nil {
			logger.Fatalf("%v\n", rm_reelid_scrap_unicFile)
		}*/
	// читаю файл построчно
	get_reel_id_scrap, err := readLines(fileReedIdScrap)
	if err != nil {
		logger.Errorf(err.Error())
	}
	// убираю дублированные номера
	rmreel_id_scrap := removeDuplicatesinfile(get_reel_id_scrap)

	fileReelIdUnicScrap, err := os.OpenFile(fileReelIDScrap_unic, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error creating file %v", err)
	}

	datawriterScrap := bufio.NewWriter(fileReelIdUnicScrap)

	for _, data := range rmreel_id_scrap {
		_, err := datawriterScrap.WriteString(data + "\n")
		if err != nil {
			logger.Errorf(err.Error())
		}
	}

	datawriterScrap.Flush()
	fileReelIdUnicScrap.Close()

	arrReelIdScrapALL := filereader.Readfile(fileReedIdScrap)
	arrReelIdScrapUnic := filereader.Readfile(fileReelIDScrap_unic)
	reelIDScrapStore := []ReelIdData{}
	for _, i := range arrReelIdScrapUnic {
		sum := 0
		for _, j := range arrReelIdScrapALL {
			if i[0] == j[0] {
				sum += 1
			}
		}
		reelIDScrapStore = append(reelIDScrapStore, ReelIdData{ReelID: i[0], Qty: strconv.Itoa(sum)})
		fmt.Printf("Scrap, reel_id: %v, sum: %v\n", i[0], sum)
	}
	fmt.Printf("reelIDScrapStore: %v\n", reelIDScrapStore)
	valuesText := []string{}
	for _, r := range reelIDScrapStore {

		fmt.Printf("r.ReelID - %v\n", r.ReelID)
		valuesText = append(valuesText, "'"+r.ReelID+"'")
	}
	// Join our string slice.
	result := strings.Join(valuesText, " or REEL_BARCODE = ")
	fmt.Println(result)

	qr, err := r.DB.QueryContext(ctx, querySelectReelData+fmt.Sprintln(result))
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			//return nil, err

		}
	}
	defer qr.Close()

	var qrs []Reel_Data
	for qr.Next() {
		var qrts Reel_Data
		if err := qr.Scan(
			&qrts.PART_NO,
			&qrts.REEL_BARCODE,
			&qrts.LOT_NO,
		); err != nil {
			logger.Errorf(err.Error())
		}
		qrs = append(qrs, qrts)
	}
	if err = qr.Err(); err != nil {
		//return qrs, err
		logger.Errorf(err.Error())
	}
	fmt.Println(qrs)
	for _, i := range reelIDScrapStore {
		//fmt.Printf("i.ReelID: %v\n", i.ReelID)
		for _, j := range qrs {
			//	fmt.Printf("j.REEL_BARCODE: %v\n", j.REEL_BARCODE)
			if i.ReelID == j.REEL_BARCODE {
				fmt.Printf("SAP: %v, ID: %v, Qty: %v, Lot: %v\n", j.PART_NO, i.ReelID, i.Qty, j.LOT_NO)
			}
		}
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func removeDuplicatesinfile(elements []string) []string { // change string to int here if required
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{} // change string to int here if required
	result := []string{}             // change string to int here if required

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
