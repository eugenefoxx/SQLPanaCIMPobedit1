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
	"sync"
	"time"

	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/filereader"
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

func (r *panaCIMStorage) GetSumPattert(jobid string) ([]SumPattern, error) {
	// logger := logging.GetLogger()
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

func (r *panaCIMStorage) GetPatternForPanel() ([]ProductData, error) {
	// logger := logging.GetLogger()
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

const querySelectPatternTypesPerPanel = `
SELECT [PATTERN_TYPES_PER_PANEL]
FROM [PanaCIM].[dbo].[product_data]
WHERE PRODUCT_ID = `

// ???????????????? ???? ???? PatternTypesPerPanel
func (r *panaCIMStorage) GetPatternTypesPerPanel(product_id string) ([]ProductData, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	qr, err := r.DB.QueryContext(ctx, querySelectPatternTypesPerPanel+product_id)
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
			&qrts.PATTERN_TYPES_PER_PANEL,
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

func (r *panaCIMStorage) GetProductId(jobid string) ([]ProductSetup, error) {
	// logger := logging.GetLogger()
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

func (r *panaCIMStorage) GetProductName(productid string) ([]ProductData, error) {
	// logger := logging.GetLogger()
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
// ???????????????? ?????????? ?????????? ????????????
func (r *panaCIMStorage) GetRouteId(productid string) ([]ProductSetup, error) {
	// logger := logging.GetLogger()
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
const querySelecInfoInstallJobId_ViewbyParameter = `
SELECT
    [REEL_ID]
        ,[PART_NO]
        ,[LOT_NO]
        , [PLACE_COUNT]
        , [PICKUP_COUNT]
        , [reel_barcode]
        , [CURRENT_QUANTITY]
        , [INITIAL_QUANTITY]
        , (SELECT *
    FROM dbo.SumPLACE_COUNT_ALL_REEL_ID([PanaCIM].[dbo].[InfoInstallLastJobId_View].[REEL_ID])) AS PLACE_COUNT_ALL
        , (SELECT *
    FROM dbo.SumPICKUP_COUNT_ALL_REEL_ID([PanaCIM].[dbo].[InfoInstallLastJobId_View].[REEL_ID])) AS PICKUP_COUNT_ALL
        , ([PICKUP_COUNT] - [PLACE_COUNT]) AS Delta
FROM [PanaCIM].[dbo].[InfoInstallLastJobId_View]
order by PART_NO;`

func (r *panaCIMStorage) GetPanacimDataComponentsByJobId(jobid string) ([]InfoInstallLastJobId_View, error) {
	// logger := logging.GetLogger()
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

// ???????????????? ???????????? ???? id ???????????????????? (??????????????????????, ??????????????, ????????????) ???? ???????????? ?????????????????????? InfoInstallJobId_View
func (r *panaCIMStorage) GetPanacimDataComponentsByJobIdAllParamReelid(jobid string) ([]InfoInstallLastJobId_View, error) {
	// logger := logging.GetLogger()
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

	qr, err := r.DB.QueryContext(ctx, querySelecInfoInstallJobId_ViewbyParameter)
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
			&qrts.ReelID,
			&qrts.PartNo,
			&qrts.Lot,
			&qrts.PlaceCount,
			&qrts.PickupCount,
			&qrts.ReelBarcode,
			&qrts.CurrentQuantity,
			&qrts.InitialQuantity,
			&qrts.PlaceCountAll,
			&qrts.PickupCountAll,
			&qrts.Delta,
			//&qrts.SumPlaceCount,
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

const querySelectGetWOComponent = `
SELECT 
[PART_NO]
, SUM([PLACE_COUNT]) AS SUM_PLACE_COUNT
,LOT_NO
FROM [PanaCIM].dbo.InfoInstallLastJobId_View
group by PART_NO, LOT_NO;`

// ???????????????? ???? ???? PanaCIM ???????????? ???????????????? ????????-??????????, ??????????, ??????
func (r *panaCIMStorage) GetPanaCIMDBWOComponent(jobid string) ([]InfoInstallLastJobId_View, error) {
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

	qr, err := r.DB.QueryContext(ctx, querySelectGetWOComponent)
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
			&qrts.Lot,
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

func (r *panaCIMStorage) WritePanaCIMDBWOComponentToFile(in []InfoInstallLastJobId_View) (err error) {
	wo_componentPath := os.Getenv("wo_component")

	if utils.FileExists(wo_componentPath) {
		os.Remove(wo_componentPath)
	}

	var partNO string = `PART_NO`
	var sum string = `SUM`
	var lot string = `Lot`

	if _, err := os.Stat(wo_componentPath); os.IsNotExist(err) {
		wo_comp_c, err := os.Create(wo_componentPath)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer wo_comp_c.Close()

		writer := csv.NewWriter(wo_comp_c)
		writer.Write([]string{partNO, sum, lot})
		writer.Comma = ','
		writer.Flush()
	}

	splitWOComponent, err := os.OpenFile(wo_componentPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		r.logger.Errorf(err.Error())
		return nil
	}
	defer splitWOComponent.Close()

	for _, i := range in {
		var result = []string{i.PartNo + "," + i.SumPlaceCount + "," + i.Lot}
		for _, v := range result {
			_, err = fmt.Fprintln(splitWOComponent, v)
			if err != nil {
				splitWOComponent.Close()
				return nil
			}
		}
	}
	return nil

}

func (r *panaCIMStorage) WtitePanaCIMDataComponentsToFile(in []InfoInstallLastJobId_View) (err error) {
	// logger := logging.GetLogger()
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
		writer.Write([]string{partNO, sumPlaceCount})
		writer.Comma = ','
		writer.Flush()
	}

	splitPanaCIM, err := os.OpenFile(panacimFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		r.logger.Errorf(err.Error()) //r.logger.Errorf(err.Error())
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

// ???????????? ?????????????????????? id ?? ??????-???? ?? ???????? ???? ???????????? ???? panacim
func (r *panaCIMStorage) WtitePanaCIMDataComponentsToFileUnpackId(in []InfoInstallLastJobId_View) (err error) {
	// logger := logging.GetLogger()
	unpack_id_path := os.Getenv("unpack_id")

	unpack_id_pathRemove := unpack_id_path
	if utils.FileExists(unpack_id_pathRemove) {
		os.Remove(unpack_id_pathRemove)
	}

	var id string = `id`   // reel_barcode
	var qty string = `qty` // PLACE_COUNT
	unpack_idFile := unpack_id_path
	if _, err := os.Stat(unpack_idFile); os.IsNotExist(err) {
		unpack_idFile, err := os.Create(unpack_idFile)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer unpack_idFile.Close()

		writer := csv.NewWriter(unpack_idFile)
		writer.Write([]string{id + `,` + qty})
		writer.Comma = ','
		writer.Flush()
	}

	splitUnpakId, err := os.OpenFile(unpack_idFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		r.logger.Errorf(err.Error())
		return nil
	}
	defer splitUnpakId.Close()

	for _, i := range in {
		var result = []string{i.ReelBarcode + "," + i.PlaceCount}

		for _, v := range result {
			_, err = fmt.Fprintln(splitUnpakId, v)
			if err != nil {
				splitUnpakId.Close()
				return nil
			}
		}
	}
	return nil
}

const querySelectGetScrapID = `
SELECT 
reel_barcode
, ([PICKUP_COUNT] - [PLACE_COUNT]) AS Delta
FROM [PanaCIM].dbo.InfoInstallLastJobId_View
WHERE ([PICKUP_COUNT] - [PLACE_COUNT]) > 0
order by reel_barcode`

// ???????????????? ???? ???? PanaCIM ???????????? ????????????????
func (r *panaCIMStorage) GetPanaCIMDBScrapID(jobid string) ([]InfoInstallLastJobId_View, error) {
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

	qr, err := r.DB.QueryContext(ctx, querySelectGetScrapID)
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
			&qrts.ReelBarcode,
			&qrts.Delta,
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

// ???????????? ???????????????????? GetPanaCIMDBScrapID ?? ????????
func (r *panaCIMStorage) WritePanaCIMDBScrapIDToFile(in []InfoInstallLastJobId_View) (err error) {
	scrapIDPath := os.Getenv("unpack_id_scrap")
	if utils.FileExists(scrapIDPath) {
		os.Remove(scrapIDPath)
	}

	var reel_id string = `id`
	var qty string = `qty`
	if _, err := os.Stat(scrapIDPath); os.IsNotExist(err) {
		scrap_id_f, err := os.Create(scrapIDPath)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer scrap_id_f.Close()

		writer := csv.NewWriter(scrap_id_f)
		writer.Write([]string{reel_id, qty})
		writer.Comma = ','
		writer.Flush()
	}

	splitScrapID, err := os.OpenFile(scrapIDPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		r.logger.Errorf(err.Error())
		return nil
	}
	defer splitScrapID.Close()

	for _, i := range in {
		var result = []string{i.ReelBarcode + "," + i.Delta}
		for _, v := range result {
			_, err = fmt.Fprintln(splitScrapID, v)
			if err != nil {
				splitScrapID.Close()
				return nil
			}
		}
	}
	return nil
}

const querySelectGetScrapSumANDLot = `
SELECT
[PART_NO]
, SUM([PICKUP_COUNT] - [PLACE_COUNT]) AS SUM_PLACE_COUNT
, LOT_NO
FROM [PanaCIM].dbo.InfoInstallLastJobId_View
WHERE ([PICKUP_COUNT] - [PLACE_COUNT]) > 0
group by PART_NO, LOT_NO;`

// ???????????????? ???? ???? PanaCIM ???????????? ????????-??????????, ??????????, ?????? ???? ??????????
func (r *panaCIMStorage) GetPanaCIMDBScrap(jobid string) ([]InfoInstallLastJobId_View, error) {
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

	qr, err := r.DB.QueryContext(ctx, querySelectGetScrapSumANDLot)
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
			&qrts.Lot,
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

func (r *panaCIMStorage) WritePanaCIMDBScrapToFile(in []InfoInstallLastJobId_View) (err error) {
	scrapPath := os.Getenv("scrap")
	if utils.FileExists(scrapPath) {
		os.Remove(scrapPath)
	}

	var partNO string = `PART_NO`
	var sum string = `SUM`
	var lot string = `Lot`

	if _, err := os.Stat(scrapPath); os.IsNotExist(err) {
		scrap_c, err := os.Create(scrapPath)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer scrap_c.Close()

		writer := csv.NewWriter(scrap_c)
		writer.Write([]string{partNO, sum, lot})
		writer.Comma = ','
		writer.Flush()
	}

	splitScrap, err := os.OpenFile(scrapPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		r.logger.Errorf(err.Error())
		return nil
	}
	defer splitScrap.Close()

	for _, i := range in {
		var result = []string{i.PartNo + "," + i.SumPlaceCount + "," + i.Lot}
		for _, v := range result {
			_, err = fmt.Fprintln(splitScrap, v)
			if err != nil {
				splitScrap.Close()
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

func (r *panaCIMStorage) GetPanaCIMixName(productid string) ([]ProductSetup, error) {
	// logger := logging.GetLogger()
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

// ?????????? ?????????? ?? ???????? ???? mix_name
func (r *panaCIMStorage) GetPanaCIMParts(mixname string) ([]SubstituteParts, error) {
	// logger := logging.GetLogger()
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

func (r *panaCIMStorage) WritePanaCIMPartsToFile(in []SubstituteParts) (err error) {
	// logger := logging.GetLogger()
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
		r.logger.Errorf(err.Error()) //r.logger.Errorf(err.Error())
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

// ???????????????? ?????????? ?? ?????????? ???????????? WO ?? unix-??????????????
func (r *panaCIMStorage) GetUnixTimeWO(jobid string) ([]Job_History, error) {
	// logger := logging.GetLogger()
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

type PatternStore struct {
	Code string
}

func (r *panaCIMStorage) GetSumPCBFromU03V2(startUnixTimeWO, finishUnixTimeWO, npm string) (sumstrPCBOrder string) {
	// logger := logging.GetLogger()
	npmToUp := strings.ToUpper(npm)
	// ?????????????????????? unix ?????????????? ???????????? ?? ???????????????????? ???????????? WO

	tStartWO, err := strconv.ParseInt(startUnixTimeWO, 10, 64)
	if err != nil {
		r.logger.Panic(err)
		panic(err)
	}
	tmStartWO := time.Unix(tStartWO, 0)

	//tmmStartWO := tmStartWO.Add(time.Minute * 3)
	p_tmStartWO, _ := time.Parse(layoutDate, tmStartWO.Format(layoutDate))
	fmt.Println("p_tmStartWO:", p_tmStartWO)

	tFinishWO, err := strconv.ParseInt(finishUnixTimeWO, 10, 64)
	if err != nil {
		panic(err)
	}
	tmFinishWO := time.Unix(tFinishWO, 0)
	p_tmFinishWO, _ := time.Parse(layoutDate, tmFinishWO.Format(layoutDate))
	fmt.Println("p_tmFinishWO: ", p_tmFinishWO)
	/*
		folderFromPanaNPM_1 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resourcePanaCIM/NPM-1/"
		folderFromPanaNPM_2 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resourcePanaCIM/NPM-2/"
		folderFromPanaNPM_3 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resourcePanaCIM/NPM-3/"
		folderFromPanaNPM_4 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resourcePanaCIM/NPM-4/"

		folderToCopyNPM_1 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resource/NPM-1/"
		folderToCopyNPM_2 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resource/NPM-2/"
		folderToCopyNPM_3 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resource/NPM-3/"
		folderToCopyNPM_4 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resource/NPM-4/"
	*/
	folderFromPanaNPM_1 := "/mnt/npm-1/"
	folderFromPanaNPM_2 := "/mnt/npm-2/"
	folderFromPanaNPM_3 := "/mnt/npm-3/"
	folderFromPanaNPM_4 := "/mnt/npm-4/"

	folderToCopyNPM_1 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/npm/NPM-1/processed/"
	folderToCopyNPM_2 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/npm/NPM-2/processed/"
	folderToCopyNPM_3 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/npm/NPM-3/processed/"
	folderToCopyNPM_4 := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/npm/NPM-4/processed/"

	//wg := sync.WaitGroup{}
	//wg.Add(2)
	//go func() {
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
	//	wg.Done()
	//}()

	// ???????????????? ???????????? ?????????? ?? ?????????????????????????? ????????????????????
	//resourcePath := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resource/"
	resourcePath := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/npm/"
	f_resource, err := ioutil.ReadDir(resourcePath)
	if err != nil {
		log.Panicln(err)
	}
	// ???????????????????????????????? ??????????????, ???????? ?????? ????????
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
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
		wg.Done()
	}()

	sumPCBOrder := 0
	checkDuble := map[string]bool{}
	inputCoreFolder, err := ioutil.ReadDir(resourcePath)
	if err != nil {
		log.Fatal(err)
	}

	storagPattern := []PatternStore{}

	go func() {
		for _, npmf := range inputCoreFolder {
			if npmf.IsDir() {
				// ???????????? ??????-???? ????????
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
													r.logger.Errorf("Fail to read file: %v", err)
													os.Exit(1)
												}
												dataFile := cfg.Section("Index").Key("Date").String()
												pdataFile, _ := time.Parse(layoutDate, dataFile)
												// p_tmStartWO.Add(-3*time.Minute) - ???????????????? ?????????????????? ?????????? ???? 3 ????????????
												if (pdataFile.After(p_tmStartWO.Add(-3*time.Minute)) && pdataFile.Before(p_tmFinishWO.Add(3*time.Minute))) &&
													(strings.EqualFold(strings.ToUpper(cfg.Section("Information").Key("LotName").String()), strings.ToUpper(npmToUp))) {
													if checkDuble[cfg.Section("Information").Key("Code").String()] {
														fmt.Println("Code Double: ", cfg.Section("Information").Key("Code").String())
													} else {
														checkDuble[cfg.Section("Information").Key("Code").String()] = true
														storagPattern = append(storagPattern, PatternStore{Code: cfg.Section("Information").Key("Code").String()})
													}

													//fmt.Printf("Pattern Store: %v", storagPattern)

													// ???????????????? ???? ?????????? ?? ???????????? ???? ?????????? Code checkDuble[cfg.Section("Information").Key("Code").String()] == true
													/*if checkDuble[cfg.Section("Information").Key("Code").String()] {
															fmt.Println("Code Double: ", cfg.Section("Information").Key("Code").String())
														} else {
															checkDuble[cfg.Section("Information").Key("Code").String()] = true
															fmt.Println("Code Single: \n", cfg.Section("Information").Key("Code").String())
															// ???????????????? ???????????? ???? ???????????????????? ????????????
															readSection := cfg.Section("MountQualityTrace").Body()
															// ???????????????????????? ?? ??????????
															writeByte := []byte(readSection)
															// ???????????????????? ???????????? ?? ????????
															if err := ioutil.WriteFile("internal/out", writeByte, 0644); err != nil {
																fmt.Printf("%v", err)
															}
															// ???????????????????? ???????????? ????????, ?????????? ?????? ?????????? ?????????????? ???????????? ?? ????????????
															filepcb := "internal/pcb"

															fileRemovePCB := "internal/pcb"

															if _, err := os.Stat(fileRemovePCB); os.IsNotExist(err) {
																pcbFile, err := os.Create(fileRemovePCB)
																if err != nil {
																	r.logger.Errorf(err.Error())
																}
																defer pcbFile.Close()
															}
															e_pcb := os.Remove(fileRemovePCB)
															if e_pcb != nil {
																log.Fatal(e_pcb)
															}
															filepcbRW, err := os.OpenFile(filepcb, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
															if err != nil {
																r.logger.Errorf(err.Error())
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
															// ?????????? ???????? ??????????????????
															getnumberpcb, err := readLines(filepcb)
															if err != nil {
																r.logger.Errorf(err.Error())
															}
															// ???????????? ?????????????????????????? ????????????
															resnumberpcb := removeDuplicatesinfile(getnumberpcb)
															// ???????????? ???????????????????? ???????????? ???????? ?? ?????????????????? ?? ?????????????? sumPCBOrder
															sumTest := 0
															for i := 0; i < len(resnumberpcb); i++ {
																//for _, i := range resnumberpcb {
																//fmt.Println("t ", i)
																sumPCBOrder++
																sumTest++

															}
															fmt.Printf("\nCode: %s, Sum: %v\n", cfg.Section("Information").Key("Code").String(), sumTest)
														}

													}*/

													for _, i := range storagPattern {
														if i.Code == cfg.Section("Information").Key("Code").String() {
															// ???????????????? ???????????? ???? ???????????????????? ????????????
															readSection := cfg.Section("MountQualityTrace").Body()
															// ???????????????????????? ?? ??????????
															writeByte := []byte(readSection)
															// ???????????????????? ???????????? ?? ????????
															if err := ioutil.WriteFile("internal/out", writeByte, 0644); err != nil {
																fmt.Printf("%v", err)
															}
															// ???????????????????? ???????????? ????????, ?????????? ?????? ?????????? ?????????????? ???????????? ?? ????????????
															filepcb := "internal/pcb"

															fileRemovePCB := "internal/pcb"

															if _, err := os.Stat(fileRemovePCB); os.IsNotExist(err) {
																pcbFile, err := os.Create(fileRemovePCB)
																if err != nil {
																	r.logger.Errorf(err.Error())
																}
																defer pcbFile.Close()
															}
															e_pcb := os.Remove(fileRemovePCB)
															if e_pcb != nil {
																log.Fatal(e_pcb)
															}
															filepcbRW, err := os.OpenFile(filepcb, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
															if err != nil {
																r.logger.Errorf(err.Error())
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
															// ?????????? ???????? ??????????????????
															getnumberpcb, err := readLines(filepcb)
															if err != nil {
																r.logger.Errorf(err.Error())
															}
															// ???????????? ?????????????????????????? ????????????
															resnumberpcb := removeDuplicatesinfile(getnumberpcb)
															// ???????????? ???????????????????? ???????????? ???????? ?? ?????????????????? ?? ?????????????? sumPCBOrder
															sumTest := 0
															for i := 0; i < len(resnumberpcb); i++ {
																//for _, i := range resnumberpcb {
																//fmt.Println("t ", i)
																sumPCBOrder++
																sumTest++

															}
															fmt.Printf("\nCode: %s, Sum: %v\n", cfg.Section("Information").Key("Code").String(), sumTest)
														}
														//fmt.Printf("test sum")
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
		wg.Done()
	}()

	wg.Wait()

	sumstrPCBOrder = strconv.Itoa(sumPCBOrder)

	return sumstrPCBOrder
}

func (r *panaCIMStorage) GetSumPCBFromU03(startUnixTimeWO, finishUnixTimeWO, npm string) (sumstrPCBOrder string) {
	// logger := logging.GetLogger()

	npmToUp := strings.ToUpper(npm)
	tStartWO, err := strconv.ParseInt(startUnixTimeWO, 10, 64)
	if err != nil {
		r.logger.Fatalf(err.Error()) //panic(err)
	}
	tmStartWO := time.Unix(tStartWO, 0)
	//fmt.Println("tmStartWO:", tmStartWO)
	p_tmStartWO, _ := time.Parse(layoutDate, tmStartWO.Format(layoutDate))
	fmt.Println("p_tmStartWO", p_tmStartWO)
	// ?????????????????? +3 ?? 15 ?????? ???? GMT
	//chH_tmStartWO := tmStartWO.Add(time.Hour * 3)
	//chM_tmStartWO := chH_tmStartWO.Add(time.Minute * 15)
	//fmt.Printf("chM_tmStartWO: %v\n", chM_tmStartWO)

	tFinishWO, err := strconv.ParseInt(finishUnixTimeWO, 10, 64)
	if err != nil {
		r.logger.Fatalf(err.Error()) //panic(err)
	}
	tmFinishWO := time.Unix(tFinishWO, 0)
	//fmt.Println("tmFinishWO: ", tmFinishWO)
	p_tmFinishWO, _ := time.Parse(layoutDate, tmFinishWO.Format(layoutDate))
	fmt.Println("p_tmFinishWO: ", p_tmFinishWO)
	//r.logger.Println("p_tmFinishWO: ", p_tmFinishWO)
	// ?????????????????? +3 ?? 15 ?????? ???? GMT
	//chH_tmFinishWO := tmFinishWO.Add(time.Hour * 3)
	//chM_tmFinishWO := chH_tmFinishWO.Add(time.Minute * 15)
	//fmt.Printf("chM_tmFinishWO: %v\n", chM_tmFinishWO)
	// ?????????????????????? ???????????????????? ???????????? ??????
	folderFrom := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/u03/processed/"
	folderToCopy := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/u03/processed_copy/"
	err = cp.Copy(folderFrom, folderToCopy)
	if err != nil {
		r.logger.Errorf(err.Error())
	}

	// ???????????????? ???????????? ?????????? ?? ?????????????????????????? ????????????????????
	files, err := ioutil.ReadDir(folderToCopy)
	if err != nil {
		r.logger.Fatalf(err.Error()) //log.Fatal(err)
	}
	// ?????????????? ?????????? ???????????????????? ?? ?????????????????????????? ????????????, ?????????????? ?????? ????????
	for _, f := range files {
		if f.IsDir() {
			folder, err := ioutil.ReadDir(folderToCopy + f.Name())
			if err != nil {
				r.logger.Fatalf(err.Error())
			}
			for _, g := range folder {
				if !g.IsDir() {

					if strings.Contains(g.Name(), ".gz") {

						// fmt.Printf("?????????? - %v, ???????? - %v \n", f.Name(), g.Name())
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
						//fmt.Printf("writeToFile: %v\n", writeToFile)
						//fmt.Printf("path %v, %v, %v\n", folderToCopy, f.Name(), writeToFile)
						/*_, err = os.Stat(folderToCopy + "/" + f.Name() + "/" + writeToFile)
						if err == nil {
							fmt.Printf("File %s already exists.", folderToCopy+"/"+f.Name()+"/"+writeToFile)
						}*/
						// Uncompress to a writer. We'll use a file writer

						outfileWriter, err := os.Create(folderToCopy + f.Name() + "/" + writeToFile)
						//fmt.Println("outfileWriter:", &outfileWriter)

						if err != nil {
							// r.r.logger.Fatalf(err.Error())
							//fmt.Printf("error: %v\n", err)
							r.logger.Printf("error: %v\n", err)
							r.logger.Debugf(err.Error())
							//r.r.logger.Errorf(err.Error())
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
	checkDubleCode := map[string]bool{}
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
							// fmt.Printf("Fail to read file: %v", err)
							r.logger.Printf("Fail to read file: %v", err)
							os.Exit(1)
						}
						dataFile := cfg.Section("Index").Key("Date").Value()
						pdataFile, _ := time.Parse(layoutDate, dataFile)
						lotnameToUpper := strings.ToUpper(cfg.Section("Information").Key("LotName").String())
						// if cfg.Section("Information").Key("LotName").String() == "NPM_915-00211_A_S"

						if (strings.EqualFold(lotnameToUpper, npmToUp)) && (pdataFile.After(p_tmStartWO.Add(-3*time.Minute)) && pdataFile.Before(p_tmFinishWO.Add(3*time.Minute))) {
							// ???????????????? ???? ?????????? ?? ???????????? ???? ?????????? Code
							if checkDubleCode[cfg.Section("Information").Key("Code").String()] {
								fmt.Println("Duble Code: ", cfg.Section("Information").Key("Code").String())
							} else {
								checkDubleCode[cfg.Section("Information").Key("Code").String()] = true
								//fmt.Printf("NPM_ %v, ?????????? - %v, ???????? - %v \n", npmToUp, ff.Name(), gg.Name())
								//fmt.Println("Date", cfg.Section("Index").Key("Date").Value())

								//fmt.Println("pdataFile: ", pdataFile.Format(layoutDate))
								//if pdataFile.After(chM_tmStartWO.AddDate(0, 0, -1)) && pdataFile.Before(chM_tmFinishWO.AddDate(0, 0, +1)) {
								// ???????????????? ???????????? ???? ???????????????????? ????????????
								readSection := cfg.Section("MountQualityTrace").Body()
								// ???????????????????????? ?? ??????????
								writeByte := []byte(readSection)
								// ???????????????????? ???????????? ?? ????????
								if err := ioutil.WriteFile("/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/u03/out", writeByte, 0644); err != nil {
									// fmt.Printf("%v", err)
									r.logger.Errorf(err.Error())
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
										r.r.logger.Errorf(err.Error())
									}
									defer outFile.Close()
								}*/
								numberPCBs := filereader.Readfileseekerspace("/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/u03/out")

								for _, i := range numberPCBs {
									data := U03{
										B: i[0],
									}
									// ?????????????? ???????????? ?? ????????????
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
									r.r.logger.Errorf(e_out.Error())
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

								//}
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

type ComponentAllData struct {
	SAP    string
	ReelID string
	Qty    string
	Lot    string
}

type WOComponentListPartNumberAndLot struct {
	SAP string
	Lot string
}

type WOComponentListPartNumberAndLotAndSum struct {
	SAP string
	Lot string
	Sum string
}

type PanacimComponentStore struct {
	SAP string
	Qty string // SUM_PLACE_COUNT
}

type PanacimComponentPartNumberStore struct {
	SAP string
}

const querySelectReelData = `
SELECT [PART_NO]
,[REEL_BARCODE]
,[LOT_NO]
FROM [PanaCIM].[dbo].[reel_data]
where REEL_BARCODE = `

func (r *panaCIMStorage) GetSumComponentFromU03(startUnixTimeWO, finishUnixTimeWO, npm string) error {
	// logger := logging.GetLogger()
	npmToUp := strings.ToUpper(npm)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// ?????????????????????? unix ?????????????? ???????????? ?? ???????????????????? ???????????? WO

	tStartWO, err := strconv.ParseInt(startUnixTimeWO, 10, 64)
	if err != nil {
		r.logger.Fatal(err)
	}
	tmStartWO := time.Unix(tStartWO, 0)
	p_tmStartWO, _ := time.Parse(layoutDate, tmStartWO.Format(layoutDate))
	fmt.Println("p_tmStartWO:", p_tmStartWO)

	tFinishWO, err := strconv.ParseInt(finishUnixTimeWO, 10, 64)
	if err != nil {
		r.logger.Fatal(err)
	}
	tmFinishWO := time.Unix(tFinishWO, 0)
	p_tmFinishWO, _ := time.Parse(layoutDate, tmFinishWO.Format(layoutDate))
	fmt.Println("p_tmFinishWO: ", p_tmFinishWO)
	// ???????????????? ???????????? ?????????? ?? ?????????????????????????? ????????????????????
	resourcePath := "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/source/resource/"
	// ??????????????????, ???? ???????????????????? ???? ?????????? ???????? reel_id
	fileReelId := "internal/reel_id"
	// ??????????????, ???????? ?????????? ?????? ????????????
	if utils.FileExists(fileReelId) {
		os.Remove(fileReelId)
	}
	/*rm_reedid := os.Remove(fileReelId)
	if rm_reedid != nil {
		r.logger.Fatalf("%v\n", rm_reedid)
	}*/
	if _, err := os.Stat(fileReelId); os.IsNotExist(err) {
		reelidFile, err := os.Create(fileReelId)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer reelidFile.Close()
	}

	fileReedIdScrap := "internal/reel_id_scrap"
	// ??????????????, ???????? ?????????? ?????? ????????????
	if utils.FileExists(fileReedIdScrap) {
		os.Remove(fileReedIdScrap)
	}
	if _, err := os.Stat(fileReedIdScrap); os.IsNotExist(err) {
		reedIdScrapFile, err := os.Create(fileReedIdScrap)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer reedIdScrapFile.Close()
	}

	/*rm_reedIdScrapFile := os.Remove(fileReedIdScrap)
	if rm_reedIdScrapFile != nil {
		r.logger.Fatalf("%v\n", rm_reedIdScrapFile)
	}*/

	//checkDubleComponent := map[string]bool{}
	inputCoreFolder, err := ioutil.ReadDir(resourcePath)
	if err != nil {
		r.logger.Fatal(err)
	}
	for _, npmf := range inputCoreFolder {
		fmt.Printf("Folder npmf %v\n", npmf.Name())
		if npmf.IsDir() {
			processedf, err := ioutil.ReadDir(resourcePath + npmf.Name())
			if err != nil {
				r.logger.Fatal(err)
			}
			for _, processed := range processedf {
				fmt.Printf("Folder processed %v Great!!!\n", processed.Name())
				if processed.IsDir() {

					dataf, err := ioutil.ReadDir(resourcePath + npmf.Name() + "/" + processed.Name())
					if err != nil {
						r.logger.Fatal(err)
					}
					for _, data := range dataf {
						fmt.Printf("Folder data %v Great!!!\n", data.Name())
						if data.IsDir() {
							fileu03f, err := ioutil.ReadDir(resourcePath + npmf.Name() + "/" + processed.Name() + "/" + data.Name())
							if err != nil {
								r.logger.Fatal(err)
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
										if (pdataFile.After(p_tmStartWO.Add(-3*time.Minute)) && pdataFile.Before(p_tmFinishWO.Add(3*time.Minute))) &&
											(strings.EqualFold(strings.ToUpper(cfg.Section("Information").Key("LotName").String()), strings.ToUpper(npmToUp))) {
											// ???????????????? ???? ?????????? ?? ???????????? ???? ?????????? Code checkDuble[cfg.Section("Information").Key("Code").String()] == true
											//if checkDubleComponent[cfg.Section("Information").Key("Code").String()] {
											//	fmt.Println("Code Double: ", cfg.Section("Information").Key("Code").String())
											//} else {
											//checkDubleComponent[cfg.Section("Information").Key("Code").String()] = true
											// ???????????????? ???????????? ???? ???????????????????? ????????????
											readSection := cfg.Section("MountQualityTrace").Body()
											// ???????????????????????? ?? ??????????
											writeByte := []byte(readSection)
											// ???????????????????? ???????????? ?? ????????
											if err := ioutil.WriteFile("internal/out", writeByte, 0644); err != nil {
												fmt.Printf("%v", err)
											}

											fileReelIdRW, err := os.OpenFile(fileReelId, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
											if err != nil {
												r.logger.Errorf(err.Error())
												return err
											}
											//defer fileReelIdRW.Close()

											fileReedIdScrapRW, err := os.OpenFile(fileReedIdScrap, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
											if err != nil {
												r.logger.Errorf(err.Error())
												return err
											}
											//defer fileReedIdScrapRW.Close()

											numberPCBs := filereader.Readfileseekerspace("internal/out")
											for _, i := range numberPCBs {
												data := U03{
													B:      i[0],
													ReelID: i[13],
													F:      i[14],
												}
												if data.F == "0" && data.ReelID != "" {
													var result = []string{data.ReelID}
													for _, v := range result {
														_, err := fmt.Fprintln(fileReelIdRW, v)
														if err != nil {
															fileReelIdRW.Close()
															return err
														}
													}
												}
												// ???????????? reel_id ???? ??????????????
												if data.F == "2" && data.ReelID != "" {
													var result = []string{data.ReelID}
													fmt.Printf("testcheck: %v, %v\n", data.ReelID, fileu03.Name())
													for _, v := range result {
														_, err := fmt.Fprintln(fileReedIdScrapRW, v)
														if err != nil {
															fileReedIdScrapRW.Close()
															return err
														}
													}
													/*if data.ReelID == "1000487433" {
														fmt.Printf("1000487433 YYYYYYYYYYYYYYYYY - %v, %v\n", data.ReelID, fileu03.Name())
													}*/
												}

												if data.F != "2" && data.ReelID != "" && data.F != "0" {
													fmt.Printf("testcheck no name: %v, %v, %v\n", data.F, data.ReelID, fileu03.Name())
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

	if utils.FileExists(file_reelid_unic) {
		os.Remove(file_reelid_unic)
	}
	//var reelid_unicFile *os.File
	if _, err := os.Stat(file_reelid_unic); os.IsNotExist(err) {
		reelid_unicFile, err := os.Create(file_reelid_unic)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer reelid_unicFile.Close()
	}

	/*rm_reelid_unicFile := os.Remove(file_reelid_unic)
	if rm_reelid_unicFile != nil {
		r.logger.Fatalf("%v\n", rm_reelid_unicFile)
	}*/

	// ?????????? ???????? ??????????????????
	get_reel_id, err := readLines(fileReelId)
	if err != nil {
		r.logger.Errorf(err.Error())
	}
	// ???????????? ?????????????????????????? ????????????
	rmdreelid_unic := removeDuplicatesinfile(get_reel_id)

	fileReelIdUnic, err := os.OpenFile(file_reelid_unic, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error creating file %v", err)
	}

	datawriter := bufio.NewWriter(fileReelIdUnic)

	for _, data := range rmdreelid_unic {
		_, err := datawriter.WriteString(data + "\n")
		if err != nil {
			r.logger.Errorf(err.Error())
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
	/*if len(reelIdStore) == 0 {
		r.logger.Errorf("reelIdStore is empty, %v", err.Error())
	}*/
	//fileReelIDData := "internal/testReelId"
	//reelIdData, err := os.Create(fileReelIDData)
	//if err != nil {
	//	r.logger.Errorf(err.Error())
	//	}
	//defer reelIdData.Close()
	//writer := csv.NewWriter(reelIdData)
	//writer.Write([]string{chapterReelID, chapterQty})
	//writer.Comma = ','
	//writer.Flush()
	unpack_id_path := os.Getenv("unpack_id")

	unpack_id_pathRemove := unpack_id_path
	if utils.FileExists(unpack_id_pathRemove) {
		os.Remove(unpack_id_pathRemove)
	}

	var chapterReelID string = `id`
	var chapterQty string = `qty`
	unpack_idFile := unpack_id_path
	if _, err := os.Stat(unpack_idFile); os.IsNotExist(err) {
		unpack_idFile, err := os.Create(unpack_idFile)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer unpack_idFile.Close()

		writer := csv.NewWriter(unpack_idFile)
		writer.Write([]string{chapterReelID, chapterQty})
		writer.Comma = ','
		writer.Flush()
	}

	// ?????????????? ???????? internal/pysaprfc/data/unpack_id.csv ???? ???????????? reelIdStore
	addDataComponentsUnpackId, err := os.OpenFile(unpack_idFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		r.logger.Errorf(err.Error())
		return err
	}
	defer addDataComponentsUnpackId.Close()

	for _, i := range reelIdStore {
		var result = []string{i.ReelID + "," + i.Qty}
		for _, v := range result {
			_, err = fmt.Fprintln(addDataComponentsUnpackId, v)
			if err != nil {
				addDataComponentsUnpackId.Close()
				return err
			}
		}
	}

	valuesReelId := []string{}
	for _, r := range reelIdStore {
		valuesReelId = append(valuesReelId, "'"+r.ReelID+"'")
	}
	// Join our string slice.
	resultReelId := strings.Join(valuesReelId, " or REEL_BARCODE = ")
	qrReelId, err := r.DB.QueryContext(ctx, querySelectReelData+fmt.Sprintln(resultReelId))
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			//return nil, err

		}
	}
	defer qrReelId.Close()
	var qrsRI []Reel_Data
	for qrReelId.Next() {
		var qrts Reel_Data
		if err := qrReelId.Scan(
			&qrts.PART_NO,
			&qrts.REEL_BARCODE,
			&qrts.LOT_NO,
		); err != nil {
			r.logger.Errorf(err.Error())
		}
		qrsRI = append(qrsRI, qrts)
	}
	if err = qrReelId.Err(); err != nil {
		//return qrs, err
		r.logger.Errorf(err.Error())
	}
	componentDataStore := []ComponentAllData{}
	for _, i := range reelIdStore {
		for _, j := range qrsRI {
			if i.ReelID == j.REEL_BARCODE {
				fmt.Printf("GOOD SAP: %v, ID: %v, Qty: %v, Lot: %v\n", j.PART_NO, i.ReelID, i.Qty, j.LOT_NO)
				componentDataStore = append(componentDataStore, ComponentAllData{SAP: j.PART_NO, ReelID: i.ReelID, Qty: i.Qty, Lot: j.LOT_NO})
			}
		}
	}
	fmt.Printf("componentDataStore: %v\n", componentDataStore)
	// ?????????????? ?????????? ???? ????????????????????, ?????????????????????????????? ?????????????? ????????-?????????? ?? ????????????
	checkDubleLot := map[string]bool{}
	woComponentPartNumberAndLotStore := []WOComponentListPartNumberAndLot{}
	for _, i := range componentDataStore {
		sum := 0
		if checkDubleLot[i.SAP+i.Lot] {
			//fmt.Printf("duble sap %v + lot %v\n", i.SAP, i.Lot)
		} else {
			checkDubleLot[i.SAP+i.Lot] = true
			fmt.Printf("no duble sap %v + lot %v + sum %s\n", i.SAP, i.Lot, strconv.Itoa(sum))
			woComponentPartNumberAndLotStore = append(woComponentPartNumberAndLotStore, WOComponentListPartNumberAndLot{SAP: i.SAP, Lot: i.Lot})
		}
	}
	// ?????????????????????? ???????????? ???? ???????????????? ????????-??????????, ???????????? ?? ?????????????????????????? ???????????????? ???? ??????-????
	woComponentPartNumberAndLotAndSumStore := []WOComponentListPartNumberAndLotAndSum{}
	for _, i := range woComponentPartNumberAndLotStore {
		var sum int
		for _, j := range componentDataStore {
			if i.SAP == j.SAP && i.Lot == j.Lot {
				sumInt, err := strconv.Atoi(j.Qty)
				if err != nil {
					r.logger.Errorf(err.Error())
					return nil
				}
				sum += sumInt
			}
		}
		fmt.Printf("HHHH SAP %v, Lot %v, SUM %d\n", i.SAP, i.Lot, sum)
		woComponentPartNumberAndLotAndSumStore = append(woComponentPartNumberAndLotAndSumStore,
			WOComponentListPartNumberAndLotAndSum{
				SAP: i.SAP, Lot: i.Lot, Sum: strconv.Itoa(sum)})
	}
	//fmt.Printf("HH: %v", woComponentPartNumberAndLotAndSumStore)
	// ?????????????????? ?????????????????? ?? ????????
	wo_component_path := os.Getenv("wo_component")

	wo_component_pathRemove := wo_component_path
	if utils.FileExists(wo_component_pathRemove) {
		os.Remove(wo_component_pathRemove)
	}
	var part_number string = `PART_NO`
	var sum string = `SUM`
	var lot string = `Lot`
	wo_componentFile := wo_component_path
	if _, err := os.Stat(wo_componentFile); os.IsNotExist(err) {
		wo_componentFile, err := os.Create(wo_componentFile)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer wo_componentFile.Close()

		writer := csv.NewWriter(wo_componentFile)
		writer.Write([]string{part_number, sum, lot})
		writer.Comma = ','
		writer.Flush()
	}

	addWOComponentFile, err := os.OpenFile(wo_componentFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		r.logger.Errorf(err.Error())
		return nil
	}
	defer addWOComponentFile.Close()

	for _, i := range woComponentPartNumberAndLotAndSumStore {
		var result = []string{i.SAP + "," + i.Sum + "," + i.Lot}

		for _, v := range result {
			_, err = fmt.Fprintln(addWOComponentFile, v)
			if err != nil {
				addWOComponentFile.Close()
				return nil
			}
		}
	}

	// ?????????????????? ???????????? ???? ???????????? ?????????????? ???? ???????????????? ????????-?????????? ?? ??????-????
	// ???????????????? ???????????? ???????????????????? ????????-????????????
	checkDublePartNumber := map[string]bool{}
	panacimPartNumberStore := []PanacimComponentPartNumberStore{}
	for _, i := range componentDataStore {
		if checkDublePartNumber[i.SAP] {

		} else {
			checkDublePartNumber[i.SAP] = true
			panacimPartNumberStore = append(panacimPartNumberStore,
				PanacimComponentPartNumberStore{SAP: i.SAP})
		}
	}
	// ?????????????????? ???????????? ???????????????????? ????????-?????????????? ?? ???????? ???? ??????????????????????
	panacimComponentstore := []PanacimComponentStore{}
	for _, i := range panacimPartNumberStore {
		var sum int
		for _, j := range woComponentPartNumberAndLotAndSumStore {
			if i.SAP == j.SAP {
				sumInt, err := strconv.Atoi(j.Sum)
				if err != nil {
					r.logger.Errorf(err.Error())
					return nil
				}
				sum += sumInt
			}
		}
		fmt.Printf("PANACIM SAP SUM: SAP %v, SUM %d\n", i.SAP, sum)
		panacimComponentstore = append(panacimComponentstore,
			PanacimComponentStore{
				SAP: i.SAP,
				Qty: strconv.Itoa(sum),
			})
	}
	// ???????????? ???????????????????? ?? ????????
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
		writer.Write([]string{partNO, sumPlaceCount})
		writer.Comma = ','
		writer.Flush()
	}

	addPanacimComponentToFile, err := os.OpenFile(panacimFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		r.logger.Errorf(err.Error()) //r.logger.Errorf(err.Error())
		return nil
	}
	defer addPanacimComponentToFile.Close()

	for _, i := range panacimComponentstore {
		var result = []string{i.SAP + "," + i.Qty}

		for _, v := range result {
			_, err = fmt.Fprintln(addPanacimComponentToFile, v)
			if err != nil {
				addPanacimComponentToFile.Close()
				return nil
			}
		}
	}
	// ?????????????????? ????????????  ???? ??????????????
	// ??????????????????, ???????????????????? ???? ???????? ??????????
	//fileReelIDScrap_unic := "internal/reelid_scrap_unic"
	fileReelIDScrap_unic := os.Getenv("reelid_scrap_unic")

	if utils.FileExists(fileReelIDScrap_unic) {
		os.Remove(fileReelIDScrap_unic)
	}
	if _, err := os.Stat(fileReelIDScrap_unic); os.IsExist(err) {
		reelid_scrap_unicFile, err := os.Create(fileReelIDScrap_unic)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer reelid_scrap_unicFile.Close()
	}
	/*
		rm_reelid_scrap_unicFile := os.Remove(fileReelIDScrap_unic)
		if rm_reelid_scrap_unicFile != nil {
			r.logger.Fatalf("%v\n", rm_reelid_scrap_unicFile)
		}*/
	// ?????????? ???????? ??????????????????
	get_reel_id_scrap, err := readLines(fileReedIdScrap)
	if err != nil {
		r.logger.Errorf(err.Error())
	}
	// ???????????? ?????????????????????????? ????????????
	rmreel_id_scrap := removeDuplicatesinfile(get_reel_id_scrap)

	fileReelIdUnicScrap, err := os.OpenFile(fileReelIDScrap_unic, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error creating file %v", err)
	}

	datawriterScrap := bufio.NewWriter(fileReelIdUnicScrap)

	for _, data := range rmreel_id_scrap {
		_, err := datawriterScrap.WriteString(data + "\n")
		if err != nil {
			r.logger.Errorf(err.Error())
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

	// ???????????? ???? ?????????????????? ?? ??????-?????? ?????? ????????????????????
	unpack_id__scrap_path := os.Getenv("unpack_id_scrap")
	unpackIDFileRemove := unpack_id__scrap_path

	if utils.FileExists(unpackIDFileRemove) {
		os.Remove(unpackIDFileRemove)
	}

	var idScrap string = `id`
	var qtyScrap string = `qty`
	unpackIDScrapFile := unpack_id__scrap_path
	if _, err := os.Stat(unpackIDScrapFile); os.IsNotExist(err) {
		unpackIDScrapFile, err := os.Create(unpackIDScrapFile)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer unpackIDScrapFile.Close()

		writer := csv.NewWriter(unpackIDScrapFile)
		writer.Write([]string{idScrap, qtyScrap})
		writer.Comma = ','
		writer.Flush()

	}

	addUnpackIDScrapFile, err := os.OpenFile(unpackIDScrapFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		r.logger.Errorf(err.Error()) //r.logger.Errorf(err.Error())
		return nil
	}
	defer addUnpackIDScrapFile.Close()

	for _, i := range reelIDScrapStore {
		var result = []string{i.ReelID + "," + i.Qty}

		for _, v := range result {
			_, err := fmt.Fprintln(addUnpackIDScrapFile, v)
			if err != nil {
				addUnpackIDScrapFile.Close()
				return nil
			}
		}
	}

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
			r.logger.Errorf(err.Error())
		}
		qrs = append(qrs, qrts)
	}
	if err = qr.Err(); err != nil {
		//return qrs, err
		r.logger.Errorf(err.Error())
	}
	fmt.Println(qrs)

	reelDataIDScrapStoreAll := []ComponentAllData{}
	for _, i := range reelIDScrapStore {
		//fmt.Printf("i.ReelID: %v\n", i.ReelID)
		for _, j := range qrs {
			//	fmt.Printf("j.REEL_BARCODE: %v\n", j.REEL_BARCODE)
			if i.ReelID == j.REEL_BARCODE {
				fmt.Printf("SAP SCRAP: %v, ID: %v, Qty: %v, Lot: %v\n", j.PART_NO, i.ReelID, i.Qty, j.LOT_NO)
				reelDataIDScrapStoreAll = append(reelDataIDScrapStoreAll,
					ComponentAllData{
						SAP:    j.PART_NO,
						ReelID: i.ReelID,
						Qty:    i.Qty,
						Lot:    j.LOT_NO,
					})
			}
		}
	}
	fmt.Println("reelDataIDScrapStoreAll", reelDataIDScrapStoreAll)
	checkScrapDoubleLot := map[string]bool{}
	scrapComponentPartNumberAndLotStore := []WOComponentListPartNumberAndLot{}
	for _, i := range reelDataIDScrapStoreAll {
		if checkScrapDoubleLot[i.SAP+i.Lot] {

		} else {
			checkScrapDoubleLot[i.SAP+i.Lot] = true
			fmt.Printf("scrap no duble sap %v + lot %v\n", i.SAP, i.Lot)
			scrapComponentPartNumberAndLotStore = append(scrapComponentPartNumberAndLotStore,
				WOComponentListPartNumberAndLot{
					SAP: i.SAP,
					Lot: i.Lot,
				})
		}
	}
	fmt.Println("scrapComponentPartNumberAndLotStore", scrapComponentPartNumberAndLotStore)
	scrapComponentPartNumberAndLotAndSumStore := []WOComponentListPartNumberAndLotAndSum{}
	for _, i := range scrapComponentPartNumberAndLotStore {
		var sum int
		for _, j := range reelDataIDScrapStoreAll {
			if i.SAP == j.SAP && i.Lot == j.Lot {
				sumInt, err := strconv.Atoi(j.Qty)
				if err != nil {
					r.logger.Errorf(err.Error())
					return nil
				}
				sum += sumInt
			}
		}

		scrapComponentPartNumberAndLotAndSumStore = append(scrapComponentPartNumberAndLotAndSumStore,
			WOComponentListPartNumberAndLotAndSum{
				SAP: i.SAP,
				Lot: i.Lot,
				Sum: strconv.Itoa(sum),
			})
	}
	fmt.Println("scrapComponentPartNumberAndLotAndSumStore", scrapComponentPartNumberAndLotAndSumStore)
	scrap_path := os.Getenv("scrap")

	scrap_pathRemove := scrap_path
	if utils.FileExists(scrap_pathRemove) {
		os.Remove(scrap_pathRemove)
	}
	var scrap_part_number string = `PART_NO`
	var scrap_sum string = `SUM`
	var scrap_lot string = `Lot`
	scrapFile := scrap_path
	if _, err := os.Stat(scrapFile); os.IsNotExist(err) {
		scrapFile, err := os.Create(scrapFile)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer scrapFile.Close()

		writer := csv.NewWriter(scrapFile)
		writer.Write([]string{scrap_part_number, scrap_sum, scrap_lot})
		writer.Comma = ','
		writer.Flush()
	}

	addScrapToFile, err := os.OpenFile(scrapFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		r.logger.Errorf(err.Error())
		return nil
	}
	defer addScrapToFile.Close()

	for _, i := range scrapComponentPartNumberAndLotAndSumStore {
		var result = []string{i.SAP + "," + i.Sum + "," + i.Lot}

		for _, v := range result {
			_, err = fmt.Fprintln(addScrapToFile, v)
			if err != nil {
				addScrapToFile.Close()
				return nil
			}
		}
	}
	return nil
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
		if encountered[elements[v]] {
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
