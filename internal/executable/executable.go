package executable

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/codecrafters-io/shell-starter-go/internal/reader"
)


func Execute(command reader.Command) error {
	executableFilePaths := FindExecutableFilePaths(command.Name)
	if len(executableFilePaths) == 0 {
		return fmt.Errorf("command not found")
	}

	comm := exec.Command(executableFilePaths[0], command.Args...)
	comm.Stdout = os.Stdout
	comm.Stderr = os.Stderr

	return comm.Run()
}

func FindExecutableFilePaths(executable string) []string {
	executableFilePaths := make([]string, 0)
	paths := strings.Split(os.Getenv("PATH"), ":")
	wg := sync.WaitGroup{}

	for _, path := range paths {
		wg.Add(1)
		go func() {
			defer wg.Done()
			filepath.Walk(path, func(fPath string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() && info.Name() == executable {
					executableFilePaths = append(executableFilePaths, fPath)
				}

				return nil
			})
		}()
	}

	wg.Wait()

	return executableFilePaths
}