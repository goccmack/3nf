package ast

import (
	"fmt"
	"os"
)

func fail(pos *Position, fmtStr string, args ...any) {
	fmt.Fprintf(os.Stderr, "%s: %s\n",
		pos,
		fmt.Sprintf(fmtStr, args...))
	os.Exit(1)
}
