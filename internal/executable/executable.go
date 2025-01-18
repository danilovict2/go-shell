package executable

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
)


func Execute(command command.Command) ([]byte, error) {
	if len(FindExecutableFilePaths(command.Name)) == 0 {
		return []byte{}, fmt.Errorf("command not found")
	}

	comm := exec.Command(command.Name, command.Args...)
	return comm.Output()
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