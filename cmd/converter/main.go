package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/shooketh/sakura/module/converter"
)

func main() {
	var ids int64SliceValue
	flag.Var(&ids, "ids", "ids to convert")
	flag.Parse()
	for _, id := range ids {
		ts, did, wid, seq := converter.Dump(id)
		fmt.Printf("id: %d\ntime: %s\ndatacenterID: %d\nworkerID: %d\nsequence: %d\n\n", id, ts, did, wid, seq)
	}
}

type int64SliceValue []int64

func (v *int64SliceValue) Set(s string) error {
	strs := strings.Split(s, ",")
	int64s := make([]int64, len(strs))
	var err error
	for i, v := range strs {
		int64s[i], err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return errors.New("parse error")
		}
	}
	*v = append(*v, int64s...)
	return nil
}
func (v *int64SliceValue) String() string {
	return ""
}
