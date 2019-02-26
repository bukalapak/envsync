package envsync_test

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
	defer exec.Command("rm", "-f", result).Run()

	err := syncer.Sync("testdata/env.success", result)
	assert.Nil(t, err)

	sMap := fileToMap("testdata/env.success")
	sortedMap := fileToMap("testdata/env.sorted")
	tMap := fileToMap(result)

	for k, v := range sMap {
		r, ok := tMap[k]
		assert.True(t, ok)
		assert.Equal(t, v, r)
	}

	//test sorted map
	for k, v := range tMap {
		r, ok := sortedMap[k]
		assert.True(t, ok)
		assert.Equal(t, v, r)
	}
}

func TestSyncer_Sync_Fail_Write(t *testing.T) {
	syncer := &envsync.Syncer{}
	envsync.IOWriteString = func(w io.Writer, s string) (n int, err error) {
		return len(s), errors.New("fail write")
	}
	defer func() { envsync.IOWriteString = io.WriteString }()
	result := "testdata/env.result.error.corrupt"
	exec.Command("touch", result).Run()
	defer exec.Command("rm", "-f", result).Run()

	err := syncer.Sync("testdata/env.success", result)
	assert.Error(t, err)

}

func TestSyncer_Sync_FailBackup(t *testing.T) {
	syncer := &envsync.Syncer{}
	counter := 0
	envsync.ExecCommand = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", command}
		cs = append(cs, args...)
		name := os.Args[0]

		if counter == 0 { //forced first command failed
			name = "xxxx"
		}
		cmd := exec.Command(name, cs...)
		fmt.Println("call", cs, "->", cmd.Run())
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		counter++
		return cmd
	}
	defer func() { envsync.ExecCommand = exec.Command }()
	result := "testdata/env.result.error.corrupt"
	exec.Command("touch", result).Run()
	defer exec.Command("rm", "-f", result).Run()

	err := syncer.Sync("testdata/env.success", result)
	assert.Error(t, err)
}

func TestSyncer_Sync_Success_Rewrite(t *testing.T) {
	syncer := &envsync.Syncer{}

	result := "testdata/env.result.sorted"
	exec.Command("cp", "testdata/env.success", result).Run()
	defer exec.Command("rm", "-f", result).Run()

	err := syncer.Sync("testdata/env.sorted_append", result)
	assert.Nil(t, err)

	sMap := fileToMap("testdata/env.success")
	sortedMap := fileToMap("testdata/env.sorted_append")
	tMap := fileToMap(result)

	for k, v := range sMap {
		r, ok := tMap[k]
		assert.True(t, ok)
		assert.Equal(t, v, r)
	}

	//test sorted map
	for k, v := range tMap {
		r, ok := sortedMap[k]
		assert.True(t, ok)
		assert.Equal(t, v, r)
	}
}

func TestSyncer_Sync_CorruptSourceFormat(t *testing.T) {
	syncer := &envsync.Syncer{}

	result := "testdata/env.result.success.corrupt"
	exec.Command("touch", result).Run()
	defer exec.Command("rm", "-f", result).Run()

	err := syncer.Sync("testdata/env.error", result)
	assert.NotNil(t, err)
}

func TestSyncer_Sync_CorruptTargetFormat(t *testing.T) {
	syncer := &envsync.Syncer{}

	result := "testdata/env.result.error.corrupt"
	exec.Command("touch", result).Run()
	defer exec.Command("rm", "-f", result).Run()

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
