# go-rotate-logs

[![Go Reference](https://pkg.go.dev/badge/github.com/sutantodadang/go-rotate-logs.svg)](https://pkg.go.dev/github.com/sutantodadang/go-rotate-logs)

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)

Go-Rotate-Logs is helper to rolling your logfile, cleanup old log file and backup your log file if more than size you want. its works on different logging that implement io.writer.

## Warning ⚠️

This project is used by my companies on production and there is not enough time to make unit testing so maybe something not work for you. use at your own risk.

## Installation

```bash
  go get github.com/sutantodadang/go-rotate-logs
```

## Config

- Directory

  Required. dir your log file

- Filename

  Required. filename must end with .log ekstension

- MaxSize

  Required in Megabyte MB. etc: 10 is equal to 10MB

- BackupName

  Optional. if not provide it will default use "backup". etc: foo-backup-1.log

- UsingTime

  Optional. if true it will add "-" + FormatTime on your file. etc: foo-02-01-2006.log

- FormatTime

  Optional. required UsingTime True. if no format it will default using rfc3399

- CleanOldFiles

  Optional. if true it will remove file with mod time before MaxAge

- MaxAge

  Optional. required CleanOldFiles true. if not provide it will default using 7. etc: equal to 7 days before now

## Usage/Examples

- std log

```go
import (
    "log"
    rotate "github.com/sutantodadang/go-rotate-logs"
)

logFile := &rotate.RotateLogsWriter{
		Config: rotate.Config{
			Directory:     "path/dir",
			Filename:      "foo.log",
			MaxSize:       10,
			UsingTime:     true,
			FormatTime:    "02-01-2006",
			CleanOldFiles: true,
			MaxAge:        30,
		},
	}

mw := io.MultiWriter(logFile,os.Stdout)

log.SetOutput(mw)

```

- zerolog

```go
import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    rotate "github.com/sutantodadang/go-rotate-logs"
)

logFile := &rotate.RotateLogsWriter{
		Config: rotate.Config{
			Directory:     "path/dir",
			Filename:      "foo.log",
			MaxSize:       10,
			UsingTime:     true,
			FormatTime:    "02-01-2006",
			CleanOldFiles: true,
			MaxAge:        30,
		},
	}

mw := zerolog.MultiWriter(logFile,os.Stdout)

log.Logger = zerolog.New(mw).With().Caller().Logger()


```

## Authors

- [@sutantodadang](https://www.github.com/sutantodadang)
