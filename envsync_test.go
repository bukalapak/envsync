package envsync_test

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/bukalapak/envsync"
	"github.com/stretchr/testify/assert"
)

func TestSyncer_Sync_ErrorOpenSourceFile(t *testing.T) {
	syncer := &envsync.Syncer{}

	err := syncer.Sync("testdata/env.empty", "env.result")
	assert.NotNil(t, err)
}

func TestSyncer_Sync_ErrorOpenTargetFile(t *testing.T) {
	syncer := &envsync.Syncer{}

	err := syncer.Sync("testdata/env.success", "env.empty")
	assert.NotNil(t, err)
}

func TestSyncer_Sync_Success(t *testing.T) {
	syncer := &envsync.Syncer{}

	result := "testdata/env.result.error.corrupt"
	exec.Command("touch", result).Run()
	defer exec.Command("rm", "-rf", result).Run()

	err := syncer.Sync("testdata/env.success", result)
	assert.Nil(t, err)

	sMap := fileToMap("testdata/env.success")
	tMap := fileToMap(result)

	for k, v := range sMap {
		r, ok := tMap[k]
		assert.True(t, ok)
		assert.Equal(t, v, r)
	}
}

func TestSyncer_Sync_CorruptSourceFormat(t *testing.T) {
	syncer := &envsync.Syncer{}

	result := "testdata/env.result.success.corrupt"
	exec.Command("touch", result).Run()
	defer exec.Command("rm", "-rf", result).Run()

	err := syncer.Sync("testdata/env.error", result)
	assert.NotNil(t, err)
}

func TestSyncer_Sync_CorruptTargetFormat(t *testing.T) {
	syncer := &envsync.Syncer{}

	result := "testdata/env.result.error.corrupt"
	exec.Command("touch", result).Run()
	defer exec.Command("rm", "-rf", result).Run()

	file, _ := os.OpenFile(result, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	defer file.Close()
	file.WriteString("THIS_SHOULD_RAISE_ERROR\n")

	err := syncer.Sync("testdata/env.success", result)
	assert.NotNil(t, err)
}

func fileToMap(loc string) map[string]string {
	file, _ := os.OpenFile(loc, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer file.Close()

	res := make(map[string]string)

	sc := bufio.NewScanner(file)
	sc.Split(bufio.ScanLines)

	for sc.Scan() {
		if sc.Text() != "" {
			sp := strings.SplitN(sc.Text(), "=", 2)
			res[sp[0]] = sp[1]
		}
	}

	return res
}
