package main

//
// Usage:  cat file.txt | go run frangipanni.go
//
import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"
)

type node struct {
	prefix   string
	sep      string
	children map[string]*node
}

func ptree(t *node, d int) {
	// Indentation
	for i := 0; i < d; i++ {
		fmt.Printf("  ")
	}
	// Print singletons on the same line
	for len(t.children) == 1 {
		fmt.Printf("%s%s", t.sep, t.prefix)
		for k := range t.children {
			// Loops once because len() always == 1 ;-)
			t = t.children[k]
		}
	}
	fmt.Printf("%s%s\n", t.sep, t.prefix)
	// print in sorted order
	keys := make([]string, 0, len(t.children)) // list of keys
	for k := range t.children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, kc := range keys {
		ptree(t.children[kc], d+1)
	}
}
func add(tree *node, tok []string, sep []string, max int, depth int) {
	//fmt.Printf("add %d node (%s) <= %s %s\n", depth, tree.prefix, tok, sep)
	if len(tok) < 1 {
		return
	}
	for _, c := range tree.children {
		//fmt.Printf("children %d node %s child %d %s\n", depth, tree.prefix, i, c.prefix)
		if tok[0] == c.prefix {
			//if c.sep == "" {
			//	c.sep = sep[0]
			//}
			add(c, tok[1:], sep[1:], max, depth+1)
			return
		}
	}
	// So not a match to the children. It's a new child.
	x := node{tok[0], sep[0], map[string]*node{}}
	tree.children[tok[0]] = &x
	add(&x, tok[1:], sep[1:], max, depth+1)
	//fmt.Printf("newchild %d %s\n", depth, tree)
}

func main() {
	max := 2
	file := os.Stdin
	defer file.Close()
	//	fs := " "
	isSep := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	isNotSep := func(c rune) bool {
		return !isSep(c)
	}

	root := node{"", "", map[string]*node{}}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue // skip empty lines
		}
		t := strings.FieldsFunc(line, isSep)
		seps := strings.FieldsFunc(line, isNotSep)
		if isNotSep(rune(line[0])) {
			// line didn't start with a seperator, so insert a fake one
			seps = append(seps, "") // add space at the end
			copy(seps[1:], seps)    // shift right
			seps[0] = ""            // inject fake at the front
		}
		seps = append(seps, "$")
		// t := strings.Split(l, fs)
		//fmt.Printf("read %s\n", t)
		add(&root, t, seps, max, 0)
		//ptree(&root, 0)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ptree(&root, -1)
}
