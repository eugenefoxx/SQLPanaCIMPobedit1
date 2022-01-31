package panacim

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/utils"
)

/*
Создать представление из [PanaCIM].[dbo].[work_orders] и [PanaCIM].[dbo].[job_history]
для добавления к JOB_ID параметра CLOSING_TYPE = '0'

*/

/*const queryLastListWO = `SELECT TOP 3 [WORK_ORDER_ID],[WORK_ORDER_NAME],[LOT_SIZE],[JOB_ID],
[MASTER_WORK_ORDER_ID],[COMMENTS] FROM [PanaCIM].[dbo].[work_orders] order by [JOB_ID] desc;`*/
const queryLastListWO = `
SELECT 
DISTINCT([PanaCIM].[dbo].[work_orders].WORK_ORDER_ID) AS WORK_ORDER_ID
,[PanaCIM].[dbo].[work_orders].WORK_ORDER_NAME
,[PanaCIM].[dbo].[work_orders].LOT_SIZE
,[PanaCIM].[dbo].[work_orders].JOB_ID
,[PanaCIM].[dbo].[job_history].CLOSING_TYPE
,[PanaCIM].[dbo].[job_history].SETUP_ID
FROM [PanaCIM].[dbo].[work_orders]
INNER JOIN  [PanaCIM].[dbo].[job_history] ON [PanaCIM].[dbo].[work_orders].JOB_ID=[PanaCIM].[dbo].[job_history].JOB_ID
WHERE [PanaCIM].[dbo].[job_history].CLOSING_TYPE='0'
order by [PanaCIM].[dbo].[work_orders].JOB_ID desc 
`

type LastWOData struct {
	WORKORDERID          string         `db:"WORK_ORDER_ID"`
	WORKORDERNAME        string         `db:"WORK_ORDER_NAME"`
	LOTSIZE              string         `db:"LOT_SIZE"`
	JOBID                string         `db:"JOB_ID"`
	MASTER_WORK_ORDER_ID string         `db:"MASTER_WORK_ORDER_ID"`
	COMMENTS             sql.NullString `db:"COMMENTS"`
}

func (r PanaCIMStorage) GetLastListWO() ([]LastWOData, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	qr, err := r.DB.QueryContext(ctx, queryLastListWO)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Error(err.Error())
			return nil, err
		}
	}
	defer qr.Close()

	var qrs []LastWOData
	for qr.Next() {
		var qrts LastWOData
		if err := qr.Scan(
			&qrts.WORKORDERID,
			&qrts.WORKORDERNAME,
			&qrts.LOTSIZE,
			&qrts.JOBID,
			&qrts.MASTER_WORK_ORDER_ID,
			&qrts.COMMENTS); err != nil {
			return qrs, err
		}
		qrs = append(qrs, qrts)
	}
	if err = qr.Err(); err != nil {
		return qrs, err
	}
	return qrs, nil

}

func (r PanaCIMStorage) WriteListWOToFile(in []LastWOData) (err error) {

	dirWOpath := os.Getenv("dirWO")
	closedWORemovepath := os.Getenv("closedWORemove")

	dirWO := dirWOpath
	if _, err := os.Stat(dirWO); os.IsNotExist(err) {
		os.Mkdir(dirWO, 0755)
	}
	closedWORemove := closedWORemovepath

	if utils.FileExists(closedWORemove) {
		os.Remove(closedWORemove)
	}

	closedWO := closedWORemovepath
	if _, err := os.Stat(closedWO); os.IsNotExist(err) {
		clwo, err := os.Create(closedWO)
		if err != nil {
			r.logger.Errorf(err.Error())
		}
		defer clwo.Close()

		writer := csv.NewWriter(clwo)
		writer.Write([]string{"0"})
		writer.Comma = ','
		writer.Flush()
	}

	splitWO, err := os.OpenFile(closedWO, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		r.logger.Errorf(err.Error()) //logger.Errorf(err.Error())
		return nil
	}
	defer splitWO.Close()

	for y, i := range in {
		fmt.Println("test JobId", i.JOBID)
		var result = []string{i.JOBID}
		// обрезаем select до первых трех строк по порядку
		if y < 3 {
			for _, v := range result {
				_, err = fmt.Fprintln(splitWO, v)
				if err != nil {
					splitWO.Close()
					return nil
				}
			}
		}

	}
	return nil
}
