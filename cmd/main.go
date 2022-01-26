package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/panacim"
	"github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/sortworkorders"
	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/filereader"
	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/logging"
	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/removefiles"
	"github.com/joho/godotenv"
)

const (
	//value uint16 = 3000
	value int = 1040
)

var (
	logger = logging.GetLogger()
	//logger logging.Logger
	db *sql.DB
)

func init() {
	logging.Init()
	logger := logging.GetLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err.Error)
	}
}

// NPM_910-00473_A_
func main() {
	//	logging.Init()
	logger := logging.GetLogger()
	var err error
	//connString := "sqlserver://pana-ro:gfhjkm123@10.1.14.21/Panacim?database=PanaCIM&encrypt=disable"
	connString := "sqlserver://cim:cim@10.1.14.21/Panacim?database=PanaCIM&encrypt=disable"
	db, err = sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Error creating connerction pool: " + err.Error())
	}
	defer db.Close()
	log.Printf("Connected!\n")
	/*db, err := mssql.NewMSSQL()
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	//defer db.Close()
	log.Printf("Connected!\n")
	err = db.Ping()
	if err != nil {
		panic("ping error: " + err.Error())
	}*/
	/*	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Errorf(err.Error())
		}
	}(db) */

	//opmopites.MSSQLComposite(db)
	SelectVersion()

	test1 := sortworkorders.OperationStorage{
		DB: db,
	}
	rr, err := test1.TestQr()
	if err != nil {
		logger.Errorf(err.Error())
	}
	fmt.Println(rr)

	panacimStorage := panacim.PanaCIMStorage{
		DB: db,
	}
	// получаем список из трех закрытых WO в моменте
	doLastListWO, err := panacimStorage.GetLastListWO()
	if err != nil {
		logger.Errorf(err.Error())
	}
	fmt.Println("do ", doLastListWO)

	// записываем результат doLastListWO в файл closedwo.csv
	if err := panacimStorage.WriteListWOToFile(doLastListWO); err != nil {
		logger.Errorf(err.Error())
	}

	/*test, err := sortworkorders.Sort.TestQr() //TestQr()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(test)*/
	//	logger.Println("logger initialized")
	/*
		sortworkorders.Getclosedworkorders()

		res, err := sortworkorders.GetLastJobIdValue1()
		if err != nil {
			logger.Errorf(err.Error())
		}
		if res != "" {
			logger.Infof(("res - %v"), res)
		}
		res2, err := sortworkorders.GetLastJobIdValue2()
		if err != nil {
			logger.Errorf(err.Error())
		}
		if res2 != "" {
			logger.Infof(("res2 - %v"), res2)
		}
		res3, err := sortworkorders.GetLastJobIdValue3()
		if err != nil {
			logger.Errorf(err.Error())
		}
		if res3 != "" {
			logger.Infof(("res3 - %v"), res3)
		}
	*/

	recipe := os.Getenv("recipe")
	reportCsv := os.Getenv("report")
	substituteCsv := os.Getenv("substitute")
	panacimCsv := os.Getenv("panacim")
	reportSUMCsv := os.Getenv("reportSUM")

	//npm := readfileseeker("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/NPM_910-00473_A_recipte.csv")
	npm := filereader.Readfileseeker(recipe)
	report, err := os.Create(reportCsv)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	defer report.Close()

	split, err := os.OpenFile(reportCsv, os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	defer split.Close()

	for _, iter := range npm {

		qtytotal, err := strconv.Atoi(iter[1])
		if err != nil {
			logger.Errorf(err.Error())
			return
		}

		//var result = []string{iter[0] + "," + iter[1] + "," + strconv.Itoa(int(uint16(qtytotal)*value))}
		var result = []string{iter[0] + "," + iter[1] + "," + strconv.Itoa(int(qtytotal)*value)}
		//fmt.Println(result)
		for _, v := range result {
			_, err = fmt.Fprintln(split, v)
			if err != nil {
				split.Close()
				return
			}
		}

	}
	//fmt.Println(nmpparts[0], nmpparts[1])
	reportDGS := filereader.Readfile(reportCsv)
	//reportParts := readfileseeker("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/parts.csv")
	reportParts := filereader.Readfileseeker(substituteCsv)
	//panacimdata := readfileseeker("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/panacim.csv")
	panacimdata := filereader.Readfileseeker(panacimCsv)

	for p := 0; p < len(reportDGS); p++ {
		parseParts(reportParts, reportDGS, panacimdata, reportDGS[p][0])
	}
	// формируем файлы
	for p := 0; p < len(reportDGS); p++ {
		insertPanacimDataQty(panacimdata, reportDGS[p][0])
	}
	//  формируем файлы с подсчетом Итого установленных компонентов оригинал + замена
	for p := 0; p < len(reportDGS); p++ {
		insertPanacimDataQtyTotal(reportDGS[p][0])
	}
	//reportSum, err := os.Create("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/reportSumComponent.csv")
	reportSum, err := os.Create(reportSUMCsv)
	if err != nil {
		//log.Println(err)
		logger.Errorf(err.Error())
	}
	defer reportSum.Close()
	//reportSumRead := filereader.Readfile("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/reportSumComponent.csv")
	reportSumRead := filereader.Readfile(reportSUMCsv)
	for r := 0; r < len(reportDGS); r++ {
		sumComponent(reportDGS, reportSumRead, reportDGS[r][0])
	}

	reportSummary := filereader.Readfile(reportSUMCsv)

	summaryReportComponents(reportSummary)

	//var i int
	/*
			Стопосто, [01.12.2021 13:03]
		Загнать в мапу и проверить длинну мапа с массивом

		Viacheslav Poturaev, [01.12.2021 13:03]
		либо отсортировать массивы и пробежать соседей

		map[int]string использовать. В качестве ключа - индекс в строке
	*/
	directorypath := os.Getenv("operationdata")
	directory := directorypath
	removefiles.RemoveFiles(directory)

	/*
		app := "/home/eugenearch/Code/github.com/eugenefoxx/test/readIni/readIni"
		args := []string{"-L1", "NPM_brain-1_p"}
		cmd := exec.Command(app, args...)
		_, err = cmd.Output()

		if err != nil {
			println(err.Error())
			return
		} */

}

func SelectVersion() {
	//ctx, _ := context.WithTimeout(context.Background(), 5*time.Millisecond)
	//time.Sleep(1 * time.Second)
	//context.Background()

	//err := db.PingContext(ctx)
	err := db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}
	var result string
	//err = db.QueryRowContext(ctx, "SELECT @@version").Scan(&result)
	err = db.QueryRow("SELECT @@version").Scan(&result)
	if err != nil {
		log.Fatal("Scan failed: ", err.Error())
	}
	fmt.Printf("%s\n", result)
}

func summaryReportComponents(reportSumRead [][]string) {
	logger := logging.GetLogger()
	reportSummaryCsv := os.Getenv("reportSummary")

	reportSummary, err := os.Create(reportSummaryCsv)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	defer reportSummary.Close()

	split, err := os.OpenFile(reportSummaryCsv, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	defer split.Close()

	for i := 0; i < len(reportSumRead); i++ {
		total1, err := strconv.Atoi(reportSumRead[i][2])
		if err != nil {
			logger.Errorf(err.Error())
			return
		}
		total2, err := strconv.Atoi(reportSumRead[i][3])
		if err != nil {
			logger.Errorf(err.Error())
			return
		}
		if total2-total1 != 0 {

			//	fmt.Printf("read Отклонение от DGS delta reportSummaryComponent %s %d\n", reportSumRead[i][0], total2-total1)
			var result = []string{reportSumRead[i][0] + "," + strconv.Itoa(total2-total1)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
		}
	}

}

func parseParts(reportParts, reportDGS, panacimdata [][]string, parts string) {
	logger := logging.GetLogger()
	subtitutepath := os.Getenv("parts")
	//report, err := os.Create("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/" + parts + ".csv")
	report, err := os.Create(subtitutepath + parts + ".csv")
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	defer report.Close()

	//split, err := os.OpenFile("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/"+parts+".csv", os.O_APPEND|os.O_WRONLY, 0644)
	split, err := os.OpenFile(subtitutepath+parts+".csv", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	defer split.Close()
	// просматриваем наличие оригинала по dgs из установленных компонентов по panacim
	for p := 0; p < len(panacimdata); p++ {
		if panacimdata[p][0] == parts {
			//var resultp = []string{pars + "," + panacimdata[p][1]}
			var resultp = []string{parts}
			//fmt.Println(resultp)
			for _, v := range resultp {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
		}
	}
	//  просматриваем замены по оригиналу
	//	for p := 0; p < len(panacimdata); p++ {
	for i := 0; i < len(reportParts); i++ {
		for ii := 0; ii < len(reportDGS); ii++ {
			if reportParts[i][0] == parts {
				//if panacimdata[p][0] == reportParts[i][0] {
				//	var result = []string{reportParts[i][1] + "," + panacimdata[p][1]}
				var result = []string{reportParts[i][1]}
				//fmt.Println(result)
				for _, v := range result {
					_, err = fmt.Fprintln(split, v)
					if err != nil {
						split.Close()
						return
					}
				}
				break
			}
		}
	}
}

// подставляем кол-во установленного компонента по отчету panacim
func insertPanacimDataQty(panacimdata [][]string, parts string) {
	logger := logging.GetLogger()
	subtitutepath := os.Getenv("parts")

	//pp := filereader.Readfile("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/" + pars + ".csv")
	component := filereader.Readfile(subtitutepath + parts + ".csv")
	//report, err := os.Create("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/" + parts + "pana.csv")
	report, err := os.Create(subtitutepath + parts + "pana.csv")
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	defer report.Close()
	//split2, err := os.OpenFile("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/"+parts+"pana.csv", os.O_APPEND|os.O_WRONLY, 0644)
	split, err := os.OpenFile(subtitutepath+parts+"pana.csv", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	defer split.Close()
	for p := 0; p < len(panacimdata); p++ {
		for s := 0; s < len(component); s++ {
			if panacimdata[p][0] == component[s][0] {
				var result = []string{component[s][0] + "," + panacimdata[p][1]}
				//	fmt.Println(result)
				for _, v := range result {
					_, err = fmt.Fprintln(split, v)
					if err != nil {
						split.Close()
						return
					}
				}
			}
		}
	}
}

// суммируем все кол-ва установленного компонента по ключу оригинала в файле собранных данных с кол-вом установленных компонентов
// по оригиналу и ключу
func insertPanacimDataQtyTotal(pars string) {
	logger := logging.GetLogger()
	subtitutepath := os.Getenv("parts")
	//readFile := filereader.Readfile("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/" + pars + "pana.csv")
	readFile := filereader.Readfile(subtitutepath + pars + "pana.csv")

	//report, err := os.Create("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/" + pars + "panatotal.csv")
	report, err := os.Create(subtitutepath + pars + "panatotal.csv")
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	defer report.Close()
	//split, err := os.OpenFile("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/"+pars+"panatotal.csv", os.O_APPEND|os.O_WRONLY, 0644)
	split, err := os.OpenFile(subtitutepath+pars+"panatotal.csv", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	defer split.Close()
	sumCol := 0
	for i := 0; i < len(readFile); i++ {
		convertsumCol, err := strconv.Atoi(readFile[i][1])
		if err != nil {
			logger.Errorf(err.Error())
			return
		}
		sumCol += (convertsumCol)
	}

	var result = []string{"Total:" + "," + strconv.Itoa(sumCol)}
	for _, v := range result {
		_, err = fmt.Fprintln(split, v)
		if err != nil {
			split.Close()
			return
		}
	}

}

func sumComponent(reportDGS, reportSumRead [][]string, component string) {
	logger := logging.GetLogger()
	subtitutepath := os.Getenv("parts")
	reportSUMCsv := os.Getenv("reportSUM")
	/*	report, err := os.Create("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/reportSumComponent.csv")
		if err != nil {
			log.Println(err)
		}
		defer report.Close()*/

	//split, err := os.OpenFile("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/reportSumComponent.csv", os.O_APPEND|os.O_WRONLY, 0644)
	split, err := os.OpenFile(reportSUMCsv, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	defer split.Close()

	//parts := filereader.Readfile("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/" + component + "panatotal.csv")
	parts := filereader.Readfile(subtitutepath + component + "panatotal.csv")
	//fmt.Printf("readfile TEST %s %v\n", component, parts)

	for rp := 0; rp < len(reportDGS); rp++ {
		for p := 0; p < len(parts); p++ {

			sumc, err := strconv.Atoi(parts[p][1])
			if err != nil {
				logger.Errorf(err.Error())
				return
			}
			if reportDGS[rp][0] == component {
				//fmt.Printf("reportDGS Test %v\n", reportDGS[rp][0])

				//fmt.Printf("reportDGS Test Sum %v %v\n", reportDGS[rp][0], sumc)
				var result = []string{reportDGS[rp][0] + "," + reportDGS[rp][1] + "," + reportDGS[rp][2] + "," + strconv.Itoa(sumc)}
				//var result = []string{reportDGS[rp][0] + "," + reportDGS[rp][1] + "," + reportDGS[rp][2] + "," + strconv.Itoa(sum)}
				//	fmt.Println("result TEST", result)
				for _, v := range result {
					_, err = fmt.Fprintln(split, v)
					if err != nil {
						split.Close()
						return
					}
				}
			}
		}
	}
	if parts == nil {
		fmt.Println("nil", component)
		for rp := 0; rp < len(reportDGS); rp++ {
			if reportDGS[rp][0] == component {
				//	fmt.Printf("reportDGS Test %v\n", reportDGS[rp][0])

				//	fmt.Printf("reportDGS Test Sum %v %v\n", reportDGS[rp][0], "0")
				var result = []string{reportDGS[rp][0] + "," + reportDGS[rp][1] + "," + reportDGS[rp][2] + "," + "0"}
				//fmt.Println("result TEST", result)
				for _, v := range result {
					_, err = fmt.Fprintln(split, v)
					if err != nil {
						split.Close()
						return
					}
				}
			}

		}
	}
}
