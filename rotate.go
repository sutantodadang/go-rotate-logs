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
func (w *RotateLogsWriter) Write(output []byte) (int, error) {

	w.mut.Lock()

	defer w.mut.Unlock()

	return w.file.Write(output)

}

func (w *RotateLogsWriter) Rotate() (err error) {

	w.mut.Lock()

	defer w.mut.Unlock()

	if w.file != nil {

		err = w.file.Close()

		w.file = nil

		if err != nil {
			return
		}

	}

	err = os.MkdirAll(w.config.Directory, os.ModePerm)
	if err != nil {
		return
	}

	str := strings.Split(w.config.Filename, ".log")

	if w.config.UsingTime {

		if w.config.FormatTime == "" {
			w.config.FormatTime = time.RFC3339
		}

		w.config.Filename = str[0] + "-" + time.Now().Format(w.config.FormatTime) + ".log"

	}

	pathFile := filepath.Join(w.config.Directory, w.config.Filename)

	dir, err := os.ReadDir(w.config.Filename)
	if err != nil {
		return
	}

	info, err := os.Stat(pathFile)
	if err == nil {

		// if file size over maxsize rename to backup and create new file
		if (info.Size() / 1000000) > int64(w.config.MaxSize) {

			newStr := strings.Split(w.config.Filename, ".log")

			i := strconv.Itoa(len(dir))

			err = os.Rename(pathFile, filepath.Join(w.config.Directory, newStr[0]+"-"+w.config.BackupName+"-"+i+".log"))
			if err != nil {
				return
			}

		}

	} else {
		return
	}

	w.file, err = os.OpenFile(pathFile, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return
	}

	return

}
