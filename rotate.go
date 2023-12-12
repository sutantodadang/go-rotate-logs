package gorotatelogs

import (
	"errors"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	// required
	Directory string

	// required. name file end .log extension
	Filename string

	// required in Megabyte MB. etc: 10 is equal to 10MB
	MaxSize int

	// if not provide it will default use "backup"
	BackupName string

	// if true it will add "-" + FormatTime
	UsingTime bool

	// required UsingTime True. if no format it will default using rfc3399
	FormatTime string

	// if true it will remove file with mod time before MaxAge
	CleanOldFiles bool

	// required CleanOldFiles true. if not provide it will default using 7. etc: equal to 7 days before now
	MaxAge int
}

type RotateLogsWriter struct {
	mut sync.Mutex

	Config Config

	file *os.File
}

// write func to satisfy io.writer interface
func (w *RotateLogsWriter) Write(p []byte) (n int, err error) {

	w.mut.Lock()

	defer w.mut.Unlock()

	if w.file != nil {

		err = w.file.Close()

		w.file = nil

		if err != nil {
			return
		}

	}

	err = w.Rotate(p)
	if err != nil {
		return
	}

	return w.file.Write(p)

}

func (w *RotateLogsWriter) Rotate(p []byte) (err error) {

	tempFilename := w.Config.Filename

	if w.Config.Directory == "" {
		err = errors.New("no dir. plase provided")
		return
	}

	if w.Config.MaxSize == 0 {
		err = errors.New("no maxsize. plase provided")
		return
	}

	err = os.MkdirAll(w.Config.Directory, os.ModePerm)
	if err != nil {
		return
	}

	if w.Config.CleanOldFiles {

		if w.Config.MaxAge == 0 {
			w.Config.MaxAge = 7
		}

		go func() {

			err = w.clean()

		}()

	}

	str := strings.Split(w.Config.Filename, ".log")

	if w.Config.UsingTime && !strings.Contains(w.Config.Filename, time.Now().Format(w.Config.FormatTime)) {

		if w.Config.FormatTime == "" {
			w.Config.FormatTime = time.RFC3339
		}

		w.Config.Filename = str[0] + "-" + time.Now().Format(w.Config.FormatTime) + ".log"

	}

	if w.Config.BackupName == "" {
		w.Config.BackupName = "backup"
	}

	pathFile := filepath.Join(w.Config.Directory, w.Config.Filename)

	info, err := os.Stat(pathFile)

	if err == nil {

		// if file size over maxsize rename to backup and create new file
		if ((info.Size() + int64(len(p))) / 1000000) > int64(w.Config.MaxSize) {

			err = w.backup(pathFile)
			if err != nil {
				return
			}

			return

		}

	} else if os.IsNotExist(err) {

		w.file, err = os.Create(pathFile)
		if err != nil {

			return
		}

		// reset filename
		w.Config.Filename = tempFilename

		return

	}

	w.file, err = os.OpenFile(pathFile, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}

	// reset filename
	w.Config.Filename = tempFilename

	return

}

func (w *RotateLogsWriter) backup(pathFile string) (err error) {

	var count int

	newStr := strings.Split(w.Config.Filename, ".log")

	dir, err := os.ReadDir(w.Config.Directory)
	if err != nil {
		return
	}

	count = len(dir)

	if w.Config.UsingTime {

		count = 0

		for _, v := range dir {

			if strings.Contains(v.Name(), time.Now().Format(w.Config.FormatTime)) {

				count++

			}

		}

	}

	i := strconv.Itoa(count)

	err = os.Rename(pathFile, filepath.Join(w.Config.Directory, newStr[0]+"-"+w.Config.BackupName+"-"+i+".log"))
	if err != nil {
		return
	}

	w.file, err = os.Create(pathFile)
	if err != nil {
		return
	}

	return
}

func (w *RotateLogsWriter) clean() (err error) {

	ageTime := time.Now().AddDate(0, 0, -int(math.Abs(float64(w.Config.MaxAge))))

	dir, err := os.ReadDir(w.Config.Directory)
	if err != nil {
		return
	}

	for _, v := range dir {

		info, errInfo := v.Info()
		if errInfo != nil {
			err = errInfo
			return
		}

		if ageTime.After(info.ModTime()) {

			err = os.Remove(v.Name())
			if err != nil {
				return
			}

		}

	}

	return
}
