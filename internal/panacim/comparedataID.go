package panacim

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/filereader"
)

type PanaDataID struct {
	id          string
	qty         string
	part_number string
	lot         string
	current_qty string
}

type SAPDataID struct {
	id    string
	sap   string
	lot   string
	qty   string
	stock string
}

const querySelectReelDataCompare = `
SELECT [PART_NO]
,[REEL_BARCODE]
,[LOT_NO]
,[CURRENT_QUANTITY]
FROM [PanaCIM].[dbo].[reel_data]
where REEL_BARCODE = `

func (r *panaCIMStorage) CompareDataID(spp interface{}) (response bool, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fileReelId := "internal/reel_id"
	fileReedIdScrap := "internal/reel_id_scrap"
	arrReelIDOK := filereader.Readfile(fileReelId)
	arrReelIDNG := filereader.Readfile(fileReedIdScrap)

	reelIDALLList := []PanaDataID{}

	for _, i := range arrReelIDOK {
		reelIDALLList = append(reelIDALLList, PanaDataID{
			id: i[0],
		})
	}

	for _, i := range arrReelIDNG {
		reelIDALLList = append(reelIDALLList, PanaDataID{
			id: i[0],
		})
	}

	//fmt.Printf("reel_id_comp: %v\n", reelIDALLList)
	unicReelID := map[string]bool{}

	arrUnicReelID := []PanaDataID{}
	for _, i := range reelIDALLList {
		if unicReelID[i.id] {

		} else {
			unicReelID[i.id] = true
			arrUnicReelID = append(arrUnicReelID, PanaDataID{
				id: i.id,
			})
		}
	}

	//fmt.Printf("arrUnicReelID: %v\n", arrUnicReelID)

	arrUnicReelIDAndQty := []PanaDataID{}
	for _, i := range arrUnicReelID {
		sum := 0
		for _, j := range reelIDALLList {
			if i.id == j.id {
				sum += 1
			}
		}
		arrUnicReelIDAndQty = append(arrUnicReelIDAndQty, PanaDataID{
			id:  i.id,
			qty: strconv.Itoa(sum),
		})

	}

	//fmt.Printf("arrUnicReelID QTY: %v\n", arrUnicReelIDAndQty)

	valuesIDText := []string{}
	for _, i := range arrUnicReelIDAndQty {
		valuesIDText = append(valuesIDText, "'"+i.id+"'")
	}

	resultReelId := strings.Join(valuesIDText, " or REEL_BARCODE = ")

	qr, err := r.DB.QueryContext(ctx, querySelectReelDataCompare+fmt.Sprintln(resultReelId))
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return
		}
	}
	defer qr.Close()
	var qrsID []Reel_Data
	for qr.Next() {
		var qrts Reel_Data
		if err := qr.Scan(
			&qrts.PART_NO,
			&qrts.REEL_BARCODE,
			&qrts.LOT_NO,
			&qrts.CURRENT_QUANTITY,
		); err != nil {
			r.logger.Errorf(err.Error())
		}
		qrsID = append(qrsID, qrts)
	}
	if err = qr.Err(); err != nil {
		r.logger.Errorf(err.Error())
		return
	}

	reelIDPanaStore := []PanaDataID{}
	for _, i := range arrUnicReelIDAndQty {
		for _, j := range qrsID {
			if i.id == j.REEL_BARCODE {
				reelIDPanaStore = append(reelIDPanaStore, PanaDataID{
					id:          i.id,
					qty:         i.qty,
					part_number: j.PART_NO,
					lot:         j.LOT_NO,
					current_qty: j.CURRENT_QUANTITY,
				})
			}
		}
	}
	fmt.Printf("reelIDPanaStore: %v\n", reelIDPanaStore)

	sapid := filereader.Readfile("/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/sap_id_info.csv")
	// fmt.Printf("%v\n", sapid)
	for _, i := range sapid {
		for _, j := range reelIDPanaStore {

			value, err := strconv.ParseFloat(i[3], 32)
			if err != nil {
				// do something sensible
				r.logger.Errorf(err.Error())
			}
			float := float32(value)
			var sap_current_qty int = int(float)
			id, err := strconv.Atoi(i[0])
			if err != nil {
				r.logger.Errorf(err.Error())
			}
			pn, err := strconv.Atoi(i[1])
			if err != nil {
				r.logger.Errorf(err.Error())
			}
			pana_id, err := strconv.Atoi(j.id)
			if err != nil {
				r.logger.Errorf(err.Error())
			}
			pana_pn, err := strconv.Atoi(j.part_number)
			if err != nil {
				r.logger.Errorf(err.Error())
			}
			if id == pana_id {
				if i[4] == "7814" {
					if i[5] == spp {
						fmt.Printf("ID have: %v %v %v\n", pana_id, pn, sap_current_qty)
						if i[2] == j.lot {
							if pana_pn == pn {
								pana_qty, err := strconv.Atoi(j.qty)
								if err != nil {
									r.logger.Errorf(err.Error())
								}
								pana_curren_qty, err := strconv.Atoi(j.current_qty)
								if err != nil {
									r.logger.Errorf(err.Error())
								}
								result := sap_current_qty - pana_qty
								// ????????????????, ?????? ?????????????????????? ?????????????? ?????????? ?????? ???? ???????????? ????????
								if result >= 0 {
									// ????????????????, ?????? ?????????????? ?? ???? ?????????????? ???????????? ?????? ?????????? ????????????????????
									if pana_curren_qty <= result {
										response := true
										return response, nil
									} else {
										r.logger.Errorf("??????-???? ?????????? ???????????????????????? ???? ???? ???????????????????????? ???????????????????? ????????????????: ???? %v, ??????-???? ?????????????? ?? panacim %v, ???????????????????? ?????????????? %v", pana_id, j.current_qty, result)
										response := false
										return response, nil
									}

								} else {
									r.logger.Errorf("??????-???? ?????????? ???????????????????????? ???? ???? ???????????????????????? ???????????????????? ????????????????, ???????????? ????????: ???? %v, ??????-???? ?????????????? ?? panacim %v, ???????????????????? ?????????????? %v", pana_id, j.current_qty, result)
									response := false
									return response, nil
								}
							} else {
								r.logger.Errorf("????????-?????????? ???? ?????????????????? ?? ???? %v, ????????-?????????? sap %v, ????????-?????????? panacim %v", pana_id, pn, pana_pn)
								response := false
								return response, nil
							}
						} else {
							r.logger.Errorf("???????????? ???? ?????????????????? ?? ???? %v, ???????????? sap %v, ???????????? panacim %v", pana_id, i[2], j.lot)
							response := false
							return response, nil
						}
					} else {
						r.logger.Errorf("??????, ?????? ?????????????? ???? ???????????????????????? %v", i[5])
						response := false
						return response, nil
					}

				} else {
					r.logger.Errorf("?????????? ???? ???????????????????????? %v", i[4])
					response := false
					return response, nil
				}
			} else {
				r.logger.Errorf("???? PanaCIM ?????? ?????????????? ?? SAP: %v", pana_id)
				response := false
				return response, nil
			}
			//fmt.Printf("id:%v, pn:%v, lot:%v, qty:%v, w:%v\n", id, pn, i[2], sap_current_qty, i[4])
		}
	}
	return response, nil
}

const querySelectReelDataCompareFromView = `
SELECT
[PART_NO]
, reel_barcode
, LOT_NO
, PLACE_COUNT
, CURRENT_QUANTITY
FROM [PanaCIM].dbo.InfoInstallLastJobId_View
group by PART_NO, reel_barcode, LOT_NO, PLACE_COUNT, CURRENT_QUANTITY;`

func (r *panaCIMStorage) CompareDataIDFromDB(spp interface{}, jobid string) (response bool, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	qrDel, err := r.DB.Query(queryDelObjInfoInstallJobId_View)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return false, err
		}
	}
	defer qrDel.Close()

	qrFunc, err := r.DB.ExecContext(ctx, queryCreateInfoInstallJobId_View1+jobid+queryCreateInfoInstallJobId_View2)
	if err != nil {
		if err.Error() != "sql: function no create" {
			r.logger.Errorf(err.Error())
			return false, err
		}
	}
	defer qrFunc.RowsAffected()

	qr, err := r.DB.QueryContext(ctx, querySelectReelDataCompareFromView)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			r.logger.Errorf(err.Error())
			return
		}
	}
	defer qr.Close()
	var qrs []InfoInstallLastJobId_View
	for qr.Next() {
		var qrts InfoInstallLastJobId_View
		if err := qr.Scan(
			&qrts.PartNo,
			&qrts.ReelBarcode,
			&qrts.Lot,
			&qrts.PlaceCount,
			&qrts.CurrentQuantity,
		); err != nil {
			return false, err
		}
		qrs = append(qrs, qrts)
	}
	if err = qr.Err(); err != nil {
		return false, err
	}

	sapid := filereader.Readfile("/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/sap_id_info.csv")
	for _, i := range sapid {
		for _, j := range qrs {
			value, err := strconv.ParseFloat(i[3], 32)
			if err != nil {
				// do something sensible
				r.logger.Errorf(err.Error())
			}
			float := float32(value)
			var sap_current_qty int = int(float)
			id, err := strconv.Atoi(i[0])
			if err != nil {
				r.logger.Errorf(err.Error())
			}
			pn, err := strconv.Atoi(i[1])
			if err != nil {
				r.logger.Errorf(err.Error())
			}
			pana_id, err := strconv.Atoi(j.ReelBarcode)
			if err != nil {
				r.logger.Errorf(err.Error())
			}
			pana_pn, err := strconv.Atoi(j.PartNo)
			if err != nil {
				r.logger.Errorf(err.Error())
			}
			if id == pana_id {
				if i[4] == "7814" {
					if i[5] == spp {
						fmt.Printf("ID have: %v %v %v\n", pana_id, pn, sap_current_qty)
						if i[2] == j.Lot {
							if pana_pn == pn {
								pana_qty, err := strconv.Atoi(j.PlaceCount)
								if err != nil {
									r.logger.Errorf(err.Error())
								}
								pana_curren_qty, err := strconv.Atoi(j.CurrentQuantity)
								if err != nil {
									r.logger.Errorf(err.Error())
								}
								result := sap_current_qty - pana_qty
								// ????????????????, ?????? ?????????????????????? ?????????????? ?????????? ?????? ???? ???????????? ????????
								if result >= 0 {
									// ????????????????, ?????? ?????????????? ?? ???? ?????????????? ???????????? ?????? ?????????? ????????????????????
									if pana_curren_qty <= result {
										response := true
										return response, nil
									} else {
										r.logger.Errorf("??????-???? ?????????? ???????????????????????? ???? ???? ???????????????????????? ???????????????????? ????????????????: ???? %v, ??????-???? ?????????????? ?? panacim %v, ???????????????????? ?????????????? %v", pana_id, j.CurrentQuantity, result)
										response := false
										return response, nil
									}
								} else {
									r.logger.Errorf("??????-???? ?????????? ???????????????????????? ???? ???? ???????????????????????? ???????????????????? ????????????????, ???????????? ????????: ???? %v, ??????-???? ?????????????? ?? panacim %v, ???????????????????? ?????????????? %v", pana_id, j.CurrentQuantity, result)
									response := false
									return response, nil
								}
							} else {
								r.logger.Errorf("????????-?????????? ???? ?????????????????? ?? ???? %v, ????????-?????????? sap %v, ????????-?????????? panacim %v", pana_id, pn, pana_pn)
								response := false
								return response, nil
							}
						} else {
							r.logger.Errorf("???????????? ???? ?????????????????? ?? ???? %v, ???????????? sap %v, ???????????? panacim %v", pana_id, i[2], j.Lot)
							response := false
							return response, nil
						}
					} else {
						r.logger.Errorf("??????, ?????? ?????????????? ???? ???????????????????????? %v", i[5])
						response := false
						return response, nil
					}
				} else {
					r.logger.Errorf("?????????? ???? ???????????????????????? %v", i[4])
					response := false
					return response, nil
				}
			} else {
				r.logger.Errorf("???? PanaCIM ?????? ?????????????? ?? SAP: %v", pana_id)
				response := false
				return response, nil
			}
		}
	}
	return response, nil
}
