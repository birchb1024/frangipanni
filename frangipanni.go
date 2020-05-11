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
	"unicode"
)

type node struct {
	prefix   string
	children []*node
}

func ptree(t *node, max int, d int) {
	if d > max {
		return
	}
	for i := 0; i < d; i++ {
		fmt.Printf("  ")
	}
	for len(t.children) == 1 {
		fmt.Printf("%s ", t.prefix)
		t = t.children[0]
	}
	fmt.Printf("%s:\n", t.prefix)
	for _, c := range t.children {
		ptree(c, max, d+1)
	}
}
func add(tree *node, tok []string, max int, depth int) {
	//fmt.Printf("add %d node %s %s\n", depth, tree.prefix, tok)
	if len(tok) < 1 {
		return
	}
	if depth == max {
		x := node{strings.Join(tok, " "), []*node{}}
		tree.children = append(tree.children, &x)
		return	
	}
	for _, c := range tree.children {
		//fmt.Printf("children %d node %s child %d %s\n", depth, tree.prefix, i, c.prefix)
		if tok[0] == c.prefix {
			add(c, tok[1:], max, depth+1)
			return
		}
	}
	// So not a match to the children. It's a new child.
	x := node{tok[0], []*node{}}
	tree.children = append(tree.children, &x)
	add(&x, tok[1:], max, depth+1)
	//fmt.Printf("newchild %d %s\n", depth, tree)
}

func main() {
	max := 100
	file := os.Stdin
	defer file.Close()
//	fs := " "
	f := func(c rune) bool {
		return /* c != '-' && */ !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	root := node{".", nil}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		t := strings.FieldsFunc(l, f)
		// t := strings.Split(l, fs)
		//fmt.Printf("read %s\n", t)
		add(&root, t, max, 0)
		//ptree(&root, 0)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ptree(&root, max+1, -1)
}
