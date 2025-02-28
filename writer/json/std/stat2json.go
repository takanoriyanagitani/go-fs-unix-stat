package sjson

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"iter"
	"log"
	"os"

	us "github.com/takanoriyanagitani/go-fs-unix-stat"
	. "github.com/takanoriyanagitani/go-fs-unix-stat/util"
	"golang.org/x/sys/unix"
)

func StatToJsonStd(s *unix.Stat_t) ([]byte, error) { return json.Marshal(s) }

type StatsWriter func(iter.Seq2[*unix.Stat_t, error]) IO[Void]

func WriterToStatsWriter(wtr io.Writer) StatsWriter {
	return func(stats iter.Seq2[*unix.Stat_t, error]) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var bw *bufio.Writer = bufio.NewWriter(wtr)
			defer bw.Flush()

			var enc *json.Encoder = json.NewEncoder(bw)
			for stat, e := range stats {
				select {
				case <-ctx.Done():
					return Empty, ctx.Err()
				default:
				}

				if nil != e {
					return Empty, e
				}

				e := enc.Encode(stat)
				if nil != e {
					return Empty, e
				}
			}

			return Empty, nil
		}
	}
}

type FilenamesToStatsToJsonToWriter func(names iter.Seq[string]) IO[Void]

func (w StatsWriter) ToFilenamesToStatsToJsonToWriter(
	f2s us.FilenameToStat,
) FilenamesToStatsToJsonToWriter {
	return func(names iter.Seq[string]) IO[Void] {
		var i iter.Seq2[*unix.Stat_t, error] = func(
			yield func(*unix.Stat_t, error) bool,
		) {
			var s unix.Stat_t
			for name := range names {
				e := f2s(name, &s)
				if !yield(&s, e) {
					return
				}
			}
		}
		return w(i)
	}
}

var StatsWriterStdout StatsWriter = WriterToStatsWriter(os.Stdout)

var FilenamesToStatsToJsonToStdoutDefault = StatsWriterStdout.
	ToFilenamesToStatsToJsonToWriter(us.FilenameToStatDefault)

func (f FilenamesToStatsToJsonToWriter) StdinToNamesToWriter() IO[Void] {
	var names iter.Seq[string] = func(
		yield func(string) bool,
	) {
		var s *bufio.Scanner = bufio.NewScanner(os.Stdin)
		for s.Scan() {
			var filename string = s.Text()
			if !yield(filename) {
				return
			}
		}

		e := s.Err()
		if nil != e {
			log.Printf("error while scanning filenames: %v\n", e)
		}
	}
	return f(names)
}
