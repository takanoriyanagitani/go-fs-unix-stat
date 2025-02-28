package stat

import (
	"golang.org/x/sys/unix"
)

type FilenameToStat func(filename string, stat *unix.Stat_t) error

var FilenameToStatDefault FilenameToStat = unix.Stat
