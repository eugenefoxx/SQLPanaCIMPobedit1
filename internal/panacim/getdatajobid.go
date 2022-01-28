package panacim

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/utils"
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
