package panacim

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/logging"
	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/utils"
)

type panaCIMStorage struct {
	DB     *sql.DB
	logger *logging.Logger
	//mu     *sync.Mutex
}

func NewPanaCIMRepository(db *sql.DB, logger *logging.Logger) PanaCIMRepository {
	return &panaCIMStorage{
		DB:     db,
		logger: logger,
	}
}

func (r *panaCIMStorage) Print() {
	r.logger.Info("TEST PRINT")
}

const querySelectInfoInstallJobId_ViewComponent = `
SELECT 
    [PART_NO]
    ,[LOT_NO] 
        , SUM([PLACE_COUNT]) AS SUM_PLACE_COUNT
FROM [PanaCIM].[dbo].[InfoInstallLastJobId_View]
 group by LOT_NO, PART_NO;`

func (r *panaCIMStorage) GetPanacimDataComponentsByJobIdSAP(jobid string) ([]InfoInstallLastJobId_View, error) {
	//// logger := logging.GetLogger()
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

	qr, err := r.DB.QueryContext(ctx, querySelectInfoInstallJobId_ViewComponent)
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
			&qrts.Lot,
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

func (r *panaCIMStorage) WritePanacimDataComponentsByJobIdSAPToFile(in []InfoInstallLastJobId_View) (err error) {
	// logger := logging.GetLogger()
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

	splitWOComponent, err := os.OpenFile(wo_componentFile, os.O_APPEND|os.O_WRONLY, 0644)
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

func (r *panaCIMStorage) WriteDataInfoOrderSAP(wo_name, sum string) error {
	// logger := logging.GetLogger()
	info_orderPath := os.Getenv("info_order")

	dateNow := time.Now()
	layoutDateNow := "20060102"
	dateNowFormat := dateNow.Format(layoutDateNow)

	info_orderRemove := info_orderPath
	if utils.FileExists(info_orderRemove) {
		os.Remove(info_orderRemove)
	}

	info_orderFile := info_orderPath
	var wo string = `WO`
	var date string = `Date`
	var qty string = `Qty`
	if _, err := os.Stat(info_orderFile); os.IsNotExist(err) {
		info_order, err := os.Create(info_orderFile)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer info_order.Close()

		writer := csv.NewWriter(info_order)
		writer.Write([]string{wo, date, qty})
		writer.Write([]string{wo_name, dateNowFormat, sum})
		writer.Comma = ','
		writer.Flush()
	}
	/*
		splitInfoOrder, err := os.OpenFile(info_orderFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			r.logger.Errorf(err.Error()) //r.logger.Errorf(err.Error())
			return nil
		}
		defer splitInfoOrder.Close()

		var result = wo_name + dateNowFormat + sum
		//	for _, v := range result {
		_, err = fmt.Fprintln(splitInfoOrder, result)
		if err != nil {
			splitInfoOrder.Close()
			return nil
		}
		//	}
	*/
	return nil
}
