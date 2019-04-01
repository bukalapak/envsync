package envsync

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

var (
	//ExecCommand cmd executor
	ExecCommand = exec.Command
	//IOWriteString alias command to write string
	IOWriteString = io.WriteString
)

// EnvSyncer describes some contracts to synchronize env.
type EnvSyncer interface {
	// Sync synchronizes source and target.
	// Source is the default env or the sample env.
	// Target is the actual env.
	// Both source and target are string and indicate the location of the files.
	//
	// Any values in source that aren't in target will be written to target.
	// Any values in source that are in target won't be written to target.
	Sync(source, target string) error
}

// Syncer implements EnvSyncer.
type Syncer struct {
}

// Sync implements EnvSyncer.
// Sync will read the file line by line.
// It will read the first '=' character.
// All characters prior to the first '=' character is considered as the key.
// All characters after the first '=' character until a newline character is considered as the value.
//
// e.g: FOO=bar.
// FOO is the key and bar is the value.
//
// Before the synchronization process, it will create .bak from target file as backup.
// During the synchronization process, there may be an error.
// If error occurred, target file will be restored from previously created .bak file.
// If the synchronization succeed, the key-values will be sorted and header will be added based on key's prefix
// e.g: for FOO_BAR, '#FOO' header will be added before the first key of FOO_
func (s *Syncer) Sync(source, target string) error {
	var err error
	backupFile := fmt.Sprintf("%s.bak", target)

	// open the source file
	sFile, err := os.Open(source)
	if err != nil {
		return errors.Wrap(err, "couldn't open source file")
	}
	defer sFile.Close()

	// open the target file
	tFile, err := os.OpenFile(target, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		return errors.Wrap(err, "couldn't open target file")
	}
	defer tFile.Close()

	sMap, err := godotenv.Parse(sFile)
	if err != nil {
		return err
	}

	tMap, err := godotenv.Parse(tFile)
	if err != nil {
		return err
	}
	if err = ExecCommand("cp", "-f", target, backupFile).Run(); err != nil {
		return err
	}
	newEnv, additionalEnv := s.appendNewEnv(sMap, tMap)
	if len(additionalEnv) > 0 {
		fmt.Printf("New env added:\n%s\n", s.toString(additionalEnv))
	}

	//clear current file
	tFile.Truncate(0)
	tFile.Seek(0, 0)
	b := s.toString(newEnv)
	_, err = IOWriteString(tFile, b)
	if err != nil {
		ExecCommand("cp", "-f", backupFile, target).Run()
	}
	ExecCommand("rm", "-f", backupFile).Run()
	return errors.Wrap(err, "couldn't write target file")
}

func (s *Syncer) appendNewEnv(sMap, tMap map[string]string) (map[string]string, map[string]string) {
	addedEnv := make(map[string]string)
	for k, v := range sMap {
		if _, found := tMap[k]; !found {
			tMap[k] = v
			addedEnv[k] = v
		}
	}
	return tMap, addedEnv
}

func (s *Syncer) prefix(key string) string {
	return strings.Split(key, "_")[0]
}

func (s *Syncer) toString(env map[string]string) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys) // sort env before write
	group, groupComment := "", ""

	var buff bytes.Buffer

	for i, k := range keys {
		if g := s.prefix(k); g != group {
			if i == 0 {
				groupComment = "# %s\n"
			} else {
				groupComment = "\n# %s\n"
			}
			buff.WriteString(fmt.Sprintf(groupComment, g))
			group = g
		}

		if strings.Contains(env[k], "\n") {
			env[k] = strconv.Quote(env[k])
		}
		buff.WriteString(fmt.Sprintf("%s=%s\n", k, env[k]))
	}

	return buff.String()
}
