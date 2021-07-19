package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/workspace/evoting/ev-webservice/internal/utils"
)

func Uint64() uint64 {
	buf := make([]byte, 10)
	rand.Read(buf) // Always succeeds, no need to check error
	return binary.LittleEndian.Uint64(buf)
}
func main() {
	fmt.Println(utils.GenBetterGUID())
	fmt.Println(utils.GenSonyflake())
	// fmt.Println(10000000000 + rand.Intn(99999999999-10000000000))
}
