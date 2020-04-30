package main
//
// Usage:  cat file.txt | go run frangipanni.go
//
import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type node struct {
	prefix   string
	children []*node
}

func ptree(t *node, d int) {
	for i := 0; i < d; i++ {
		fmt.Printf("  ")
	}
	fmt.Printf("%s:\n", t.prefix)
	for _, c := range t.children {
		ptree(c, d+1)
	}
}
func add(tree *node, tok []string, depth int) {
	//fmt.Printf("add %d node %s %s\n", depth, tree.prefix, tok)
	if len(tok) < 1 {
		return
	}
	for _, c := range tree.children {
		//fmt.Printf("children %d node %s child %d %s\n", depth, tree.prefix, i, c.prefix)
		if tok[0] == c.prefix {
			add(c, tok[1:], depth+1)
			return
		}
	}
	// So not a match to the children. It's a new child.
	x := node{tok[0], []*node{}}
	tree.children = append(tree.children, &x)
	add(&x, tok[1:], depth+1)
	//fmt.Printf("newchild %d %s\n", depth, tree)
}

func main() {
	file := os.Stdin
	defer file.Close()
	fs := "/"

	root := node{".", nil}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		t := strings.Split(l, fs)
		//fmt.Printf("read %s\n", t)
		add(&root, t, 0)
		//ptree(&root, 0)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ptree(&root, -1)
}
