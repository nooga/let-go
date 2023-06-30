package rt

import (
	"bufio"
	"os"
)

func DumpNamespaces(fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	w.WriteString("letgo---")
	w.Flush()
}
