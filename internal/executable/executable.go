package executable

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
)

func Execute(command command.Command, stdout, stderr io.Writer) error {
	if GetExecutableFilePath(command.Name) == "" {
		return fmt.Errorf("command not found")
	}

	comm := exec.Command(command.Name, command.Args...)
	comm.Stdout = stdout
	comm.Stderr = stderr

	return comm.Run()
}

func GetExecutableFilePath(command string) string {
	executables := FindExecutables()
	idx := slices.IndexFunc(executables, func(executable string) bool {
		return command == executable || command == filepath.Base(executable)
	})

	if idx == -1 {
		return ""
	}

	return executables[idx]
}

func FindExecutables() []string {
	executables := make([]string, 0)
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

				if !info.IsDir() && info.Mode().Perm()&0100 != 0 {
					executables = append(executables, fPath)
				}

				return nil
			})
		}()
	}
	
	wg.Wait()
	return executables
}