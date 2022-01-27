package panacim

import (
	"context"
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

type JobProducts struct {
	SetupId string `db:"SETUP_ID"`
}

type ProductSetup struct {
	Product_Id string `db:"PRODUCT_ID"`
	Route_Id   string `db:"ROUTE_ID"`
}

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

type ProductData struct {
	ProductName     string `db:"PRODUCT_NAME"`
	PatternPerPanel string `db:"PATTERN_COMBINATIONS_PER_PANEL"`
}

type ProductDataLink []ProductData

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
