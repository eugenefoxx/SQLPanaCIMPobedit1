package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/panacim"
	"github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/sortworkorders"
	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/filereader"
	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/logging"
	"github.com/joho/godotenv"
)

const (
	//value uint16 = 3000
	value int = 1040
)

var (
	//logger = logging.GetLogger()
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

	jobIdStorage := sortworkorders.OperationStorage{
		DB: db,
	}
	rr, err := jobIdStorage.TestQr()
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

	//	logger.Println("logger initialized")
	// берем для проверки три последних полученных job_id и если их нет,
	// вносим в файл processedwo.csv для последующей обработки
	jobIdStorage.Getclosedworkorders()

	res, err := jobIdStorage.GetLastJobIdValue1()
	if err != nil {
		logger.Errorf(err.Error())
	}
	if res != "" {
		logger.Infof(("res - %v"), res)

		//s := string.patterSlice[0]
		//t1 := strings.Replace(s, "{", "", -1)
		// запросы для формирования рецепта
		// получение product_id
		productIdSlice, err := panacimStorage.GetProductId(res)
		if err != nil {
			logger.Errorf(err.Error())
		}
		fmt.Printf("product id - %v\n", productIdSlice[0].Product_Id)
		productid := productIdSlice[0].Product_Id
		// получение product_name для рецепта
		productNameSlice, err := panacimStorage.GetProductName(productid)
		if err != nil {
			logger.Errorf(err.Error())
		}
		fmt.Printf("product name NPM - %v\n", productNameSlice[0].ProductName)
		npm := productNameSlice[0].ProductName
		lineSlice, err := panacimStorage.GetRouteId(productid)
		if err != nil {
			logger.Errorf(err.Error())
		}
		fmt.Printf("route %v\n", lineSlice[0].Route_Id)
		routeid := lineSlice[0].Route_Id
		if routeid == "1009" {
			app := "/home/a20272/Code/github.com/eugenefoxx/readDGSP1forKATE/readDGSP1forKATE"
			args := []string{"-L1", npm}
			cmd := exec.Command(app, args...)
			_, err = cmd.Output()

			if err != nil {
				//	println(err.Error())
				logger.Errorf(err.Error())
				return
			}
		}

		// получаем потребленные компоненты их кол-ва по job_id
		componentsSlice, err := panacimStorage.GetPanacimDataComponentsByJobId(res)
		if err != nil {
			logger.Errorf(err.Error())
		}
		fmt.Printf("componentsSlice %v %v\n", componentsSlice[0].PartNo, componentsSlice[0].SumPlaceCount)
		fmt.Printf("componentsSlice %v %v\n", componentsSlice[1].PartNo, componentsSlice[1].SumPlaceCount)
		if err := panacimStorage.WtitePanaCIMDataComponentsToFile(componentsSlice); err != nil {
			logger.Errorf(err.Error())
		}

		mixnameSlice, err := panacimStorage.GetPanaCIMixName(productid)
		if err != nil {
			logger.Errorf(err.Error())
		}
		fmt.Printf("mixname: %v\n", mixnameSlice[0].MixName)

		mixname := mixnameSlice[0].MixName
		partsSlice, err := panacimStorage.GetPanaCIMParts(mixname)
		if err != nil {
			logger.Errorf(err.Error())
		}
		fmt.Printf("%v\n", partsSlice[0].PrimaryPn)
		// запись полученных замен из БД в файл
		if err := panacimStorage.WritePanaCIMPartsToFile(partsSlice); err != nil {
			logger.Errorf(err.Error())
		}

		// объем выпуска ? пока вопрос корректности такого подсчета
		fmt.Println("get sum pattern")
		fmt.Println(panacimStorage.GetSumPattert(res))
		patterSlice, err := panacimStorage.GetSumPattert(res)
		if err != nil {
			logger.Errorf(err.Error())
		}
		fmt.Printf("кол-во м/з: %v\n", patterSlice[0].SumPattern)
		qtyPattern := patterSlice[0].SumPattern

		pcbSlice, err := panacimStorage.GetPatternForPanel()
		if err != nil {
			logger.Errorf(err.Error())
		}
		fmt.Printf("2222 кол-во плат в м/з: %v\n", pcbSlice[0])
		fmt.Printf("кол-во плат в м/з: %v\n", pcbSlice[0].PatternPerPanel)
		qtyPCB := pcbSlice[0].PatternPerPanel

		qtyPatternInt, err := strconv.Atoi(qtyPattern)
		if err != nil {
			logger.Errorf(err.Error())
		}

		qtyPCBInt, err := strconv.Atoi(qtyPCB)
		if err != nil {
			logger.Errorf(err.Error())
		}
		//qtyPCB
		valueLot := qtyPatternInt * qtyPCBInt
		fmt.Printf("valueLot: %v\n", valueLot)

		// конец блока расчета объема выпуска партии

		// ВСТАВКА
		recipe := os.Getenv("recipe")                              // internal/source/recipte.csv
		reportCsv := os.Getenv("report")                           // /internal/report/report.csv
		substituteCsv := os.Getenv("substitute")                   // /internal/source/parts.csv
		substituteCsvFormatted := os.Getenv("substituteFormatted") // /internal/source/partsFormatted.csv
		panacimCsv := os.Getenv("panacim")                         // /internal/source/panacim.csv
		reportSUMCsv := os.Getenv("reportSUM")                     // /internal/report/reportSumComponent.csv

		//npm := readfileseeker("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/NPM_910-00473_A_recipte.csv")
		npmRecipe := filereader.Readfileseeker(recipe)
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

		for _, iter := range npmRecipe {

			qtytotal, err := strconv.Atoi(iter[1])
			if err != nil {
				logger.Errorf(err.Error())
				return
			}

			//var result = []string{iter[0] + "," + iter[1] + "," + strconv.Itoa(int(uint16(qtytotal)*value))}
			var result = []string{iter[0] + "," + iter[1] + "," + strconv.Itoa(int(qtytotal)*valueLot)}
			//fmt.Println(result)
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}

		}
		////
		// читаем файл с заменами в массив строк, поскольку записи оригинал - замена могут неоднократно повторяться
		partsGet, err := readLines(substituteCsv)
		if err != nil {
			logger.Errorf(err.Error())
		}
		// передаем на проверку дублей
		arrFotmattedparts := removeDuplicatesinfile(partsGet)

		file, err := os.OpenFile(substituteCsvFormatted, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Errorf(("failed creating file: %s"), err)
		}
		// записываем полученный очищенный результат в файл
		datawriterFormatted := bufio.NewWriter(file)
		for _, data := range arrFotmattedparts {
			_, _ = datawriterFormatted.WriteString(data + "\n")
		}
		datawriterFormatted.Flush()
		file.Close()
		////
		//fmt.Println(nmpparts[0], nmpparts[1])
		reportDGS := filereader.Readfile(reportCsv)
		//reportParts := readfileseeker("/home/eugenearch/Code/github.com/eugenefoxx/SQLPanacimP1/csvfolder/parts.csv")
		//reportParts := filereader.Readfileseeker(substituteCsv) - меняю на ниже
		reportParts := filereader.Readfileseeker(substituteCsvFormatted)
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

		summaryReportComponents(reportSummary, res, strconv.Itoa(valueLot))

		//var i int
		/*
				Стопосто, [01.12.2021 13:03]
			Загнать в мапу и проверить длинну мапа с массивом

			Viacheslav Poturaev, [01.12.2021 13:03]
			либо отсортировать массивы и пробежать соседей

			map[int]string использовать. В качестве ключа - индекс в строке
		*/

		// очистка директории
		//directorypath := os.Getenv("operationdata")
		//directory := directorypath
		//removefiles.RemoveFiles(directory)

		// КОНЕЦ

	}
	res2, err := jobIdStorage.GetLastJobIdValue2()
	if err != nil {
		logger.Errorf(err.Error())
	}
	if res2 != "" {
		logger.Infof(("res2 - %v"), res2)
	}
	res3, err := jobIdStorage.GetLastJobIdValue3()
	if err != nil {
		logger.Errorf(err.Error())
	}
	if res3 != "" {
		logger.Infof(("res3 - %v"), res3)
	}

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

// Вычисляем расхожение с рецептом
func summaryReportComponents(reportSumRead [][]string, jobid, sum_googs string) {
	logger := logging.GetLogger()
	reportSummaryCsv := os.Getenv("reportSummary")

	reportSummary, err := os.Create(reportSummaryCsv)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	defer reportSummary.Close()

	writer := csv.NewWriter(reportSummary)
	writer.Write([]string{"jobid:" + jobid + "," + "sum_googs:" + sum_googs})
	writer.Comma = ','
	writer.Flush()

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

	// dirReportsummary
	reportpath := os.Getenv("dirReportsummary")
	reportfile, err := os.Create(reportpath + jobid + ".csv")
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	defer reportfile.Close()

	pathReportToFile := reportpath + jobid + ".csv"

	split2, err := os.OpenFile(pathReportToFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	defer split2.Close()

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
				_, err = fmt.Fprintln(split2, v)
				if err != nil {
					split2.Close()
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
