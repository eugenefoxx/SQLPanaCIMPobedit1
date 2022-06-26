package filereader

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/logging"
)

var (
//	logger *logging.Logger
)

func Readfileseeker(name string) (cr [][]string, err error) {
	logger := logging.GetLogger()
	f, err := os.Open(name)
	if err != nil {
		//	fmt.Println(err)
		logger.Println(err.Error())
		return nil, err
	}
	defer f.Close()

	cr, err = readseeker(f)
	if err != nil {
		//log.Fatalf("error read %s", err)
		logger.Fatalf("error read %s", err)
		//logger.Errorf("error read %s", err.Error())
		return nil, err
	}
	return cr, nil
}

func Readfileseekerspace(name string) [][]string {
	logger := logging.GetLogger()
	f, err := os.Open(name)
	if err != nil {
		//	fmt.Println(err)
		logger.Println(err.Error())
	}
	defer f.Close()

	cr, err := readseekerspace(f)
	if err != nil {
		//log.Fatalf("error read %s", err)
		logger.Fatalf("error read %s", err)
	}

	return cr
}

func readseeker(rs io.ReadSeeker) ([][]string, error) {
	logger := logging.GetLogger()
	row1, err := bufio.NewReader(rs).ReadSlice('\n')
	if err != nil {
		return nil, err
	}

	_, err = rs.Seek(int64(len(row1)), io.SeekStart)
	if err != nil {
		return nil, err
	}

	lines := csv.NewReader(rs) //.ReadAll()
	//lines.Comma = '|'
	lines.Comma = ','
	lines.LazyQuotes = true

	CSVdata, err := lines.ReadAll()
	if err != nil {
		//	fmt.Println(err)
		//	os.Exit(1)
		//log.Fatal(err)
		logger.Fatalf("error read %s", err)
		//logger.Errorf("error read %s", err.Error())
	}

	return CSVdata, nil
}

func readseekerspace(rs io.ReadSeeker) ([][]string, error) {
	logger := logging.GetLogger()
	row1, err := bufio.NewReader(rs).ReadSlice('\n')
	if err != nil {
		return nil, err
	}

	_, err = rs.Seek(int64(len(row1)), io.SeekStart)
	if err != nil {
		return nil, err
	}

	lines := csv.NewReader(rs) //.ReadAll()
	//lines.Comma = '|'
	lines.Comma = ' '
	lines.LazyQuotes = true

	CSVdata, err := lines.ReadAll()
	if err != nil {
		//	fmt.Println(err)
		//	os.Exit(1)
		//log.Fatal(err)
		logger.Fatalf("error read %s", err)
	}

	return CSVdata, nil
}

func Readfile(name string) [][]string {
	logger := logging.GetLogger()
	f, err := os.Open(name)
	if err != nil {
		//fmt.Println(err)
		logger.Printf("error read %s", err.Error())
	}
	defer f.Close()

	cr := csv.NewReader(f)

	cr.LazyQuotes = true
	cr.Comma = ','
	//	cr.FieldsPerRecord = 10

	CSVdata, err := cr.ReadAll()
	if err != nil {
		//log.Fatalf("readfile %s", err)
		logger.Fatalf("error read %s", err)
	}
	return CSVdata
}
