package main

import (
	"context"
	"log"

	. "github.com/takanoriyanagitani/go-fs-unix-stat/util"
	js "github.com/takanoriyanagitani/go-fs-unix-stat/writer/json/std"
)

var filenames2stats2json2writer js.FilenamesToStatsToJsonToWriter = js.
	FilenamesToStatsToJsonToStdoutDefault

var stdin2names2writer IO[Void] = filenames2stats2json2writer.
	StdinToNamesToWriter()

var sub IO[Void] = func(ctx context.Context) (Void, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return stdin2names2writer(ctx)
}

func main() {
	_, e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
