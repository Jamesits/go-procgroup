//go:build windows

package go_procgroup

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestGroup(t *testing.T) {
	group, err := NewGroup()
	assert.NoError(t, err)
	assert.NotNil(t, group)

	cmd, err := group.NewCmd()
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	cmd.Path = filepath.Join(os.Getenv("WINDIR"), "system32", "notepad.exe")
	err = cmd.Start()
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = cmd.Wait()
		assert.NoError(t, err)
		wg.Done()
	}()

	time.Sleep(1 * time.Second)
	err = group.Terminate(0)
	assert.NoError(t, err)

	wg.Wait()
}
