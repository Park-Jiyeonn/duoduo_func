package main

import (
	"bufio"
	"os"
)

const (
	N   = 1e6 + 5
	mod = 1e9 + 7
)

var in *bufio.Reader

type node struct {
	x, y float64
}
type nodeList []node

func (l nodeList) Len() int      { return len(l) }
func (l nodeList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

type sortByX struct{ nodeList }

func (l sortByX) Less(i, j int) bool { return l.nodeList[i].x < l.nodeList[j].x }

func solve() {

}

func main() {
	var out *bufio.Writer
	in = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
	defer out.Flush()

	tt := 1
	//fmt.Fscan(in, &tt)
	for i := 0; i < tt; i++ {
		solve()
	}
}
