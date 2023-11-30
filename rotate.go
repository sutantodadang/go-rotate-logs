package gorotatelogs

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Directory string

	Filename string

	MaxSize int // Megabyte MB

	BackupName string

	UsingTime bool // if true it will add "-" + FormatTime

	FormatTime string // if no format it will default using rfc3399
}

type RotateLogsWriter struct {
	config Config

	mut sync.Mutex

	file *os.File
}

func New(config Config) *RotateLogsWriter {

	w := &RotateLogsWriter{config: config}

	err := w.Rotate()
	if err != nil {
		return nil
	}

	return w

}

// write func to satisfy io.writer interface
func (r *RotateLogsWriter) Write(output []byte) (int, error) {

	r.mut.Lock()

	defer r.mut.Unlock()

	return r.file.Write(output)

}

func (r *RotateLogsWriter) Rotate() (err error) {

	r.mut.Lock()

	defer r.mut.Unlock()

	if r.file != nil {

		err = r.file.Close()

		r.file = nil

		if err != nil {
			return
		}

	}

	err = os.MkdirAll(r.config.Directory, os.ModePerm)
	if err != nil {
		return
	}

	str := strings.Split(r.config.Filename, ".log")

	if r.config.UsingTime {

		if r.config.FormatTime == "" {
			r.config.FormatTime = time.RFC3339
		}

		r.config.Filename = str[0] + "-" + time.Now().Format(r.config.FormatTime) + ".log"

	}

	pathFile := filepath.Join(r.config.Directory, r.config.Filename)

	dir, err := os.ReadDir(r.config.Filename)
	if err != nil {
		return
	}

	info, err := os.Stat(pathFile)
	if err == nil {

		// if file size over maxsize rename to backup and create new file
		if (info.Size() / 1000000) > int64(r.config.MaxSize) {

			newStr := strings.Split(r.config.Filename, ".log")

			i := strconv.Itoa(len(dir))

			err = os.Rename(pathFile, filepath.Join(r.config.Directory, newStr[0]+"-"+r.config.BackupName+"-"+i+".log"))
			if err != nil {
				return
			}

		}

	} else {
		return
	}

	r.file, err = os.OpenFile(pathFile, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return
	}

	return

}
