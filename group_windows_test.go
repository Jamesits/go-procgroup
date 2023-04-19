//go:build windows

package go_procgroup

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

type TestProgram struct {
	Path string
	cmd  *Cmd
}

var programs = []TestProgram{
	{
		Path: filepath.Join(os.Getenv("WINDIR"), "system32", "notepad.exe"),
	},
	{
		Path: filepath.Join("testdata", "win32", "access_violation.exe"),
	},
	{
		Path: filepath.Join("testdata", "win32", "messagebox.exe"),
	},
}

func TestGroup(t *testing.T) {
	group, err := NewGroup()
	assert.NoError(t, err)
	assert.NotNil(t, group)

	var wg sync.WaitGroup

	// create test programs
	for _, p := range programs {
		cmd, err := group.NewCmd()
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		cmd.Path = p.Path
		err = cmd.Start()
		assert.NoError(t, err)

		wg.Add(1)
		p := p
		go func() {
			err = cmd.Wait()
			fmt.Printf("Program exit: %s\n", p.Path)
			wg.Done()
		}()
	}

	time.Sleep(3 * time.Second)
	err = group.Terminate(0)
	assert.NoError(t, err)

	wg.Wait()
}
