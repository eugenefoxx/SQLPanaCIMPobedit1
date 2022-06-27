package pysaprfc

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/eugenefoxx/SQLPanaCIMPobedit1/pkg/logging"
)

func PyExec(path string) error {
	logger := logging.GetLogger()
	cmd := exec.Command("python3", path)
	// cmd := exec.Command("python", "script.py", "--input-file", "documents/doc.png")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	err = cmd.Start()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	go copyOutput(stdout)
	go copyOutput(stderr)

	cmd.Wait()

	return err
}

func PyExecArg(path, arg string) error {
	logger := logging.GetLogger()
	cmd := exec.Command("python3", path)
	// cmd := exec.Command("python", "script.py", "--input-file", "documents/doc.png")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	err = cmd.Start()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	go copyOutput(stdout)
	go copyOutput(stderr)

	cmd.Wait()

	return err
}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
