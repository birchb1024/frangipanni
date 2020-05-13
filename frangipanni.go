package main

//
// Usage:  cat file.txt | go run frangipanni.go
//
import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"
)

type node struct {
	lineNumber int
	prefix     string
	sep        string
	children   map[string]*node
}

func add(lineNumber int, tree *node, tok []string, sep []string, max int, depth int) {
	//fmt.Printf("add %d node (%s) <= %s %s\n", depth, tree.prefix, tok, sep)
	if len(tok) < 1 {
		return
	}
	for _, c := range tree.children {
		//fmt.Printf("children %d node %s child %d %s\n", depth, tree.prefix, i, c.prefix)
		if tok[0] == c.prefix {
			add(lineNumber, c, tok[1:], sep[1:], max, depth+1)
			return
		}
	}
	// So not a match to the children. It's a new child.
	x := node{lineNumber, tok[0], sep[0], map[string]*node{}}
	tree.children[tok[0]] = &x
	add(lineNumber, &x, tok[1:], sep[1:], max, depth+1)
	//fmt.Printf("newchild %d %s\n", depth, tree)
}

func ptree(t *node, depth int, orderBy string) {

	for i := 0; i < depth; i++ { // Indentation
		fmt.Printf("  ")
	}

	for len(t.children) == 1 { // Print singletons on the same line
		fmt.Printf("%s%s", t.sep, t.prefix)
		for k := range t.children { // Loops once because len() always == 1 ;-)
			t = t.children[k]
		}
	}
	fmt.Printf("%s%s\n", t.sep, t.prefix)

	// Convert map to list for sorting
	nodes := make([]*node, 0, len(t.children)) // list of nodes
	for n := range t.children {
		nodes = append(nodes, t.children[n])
	}
	switch orderBy {
	case "input":
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].lineNumber < nodes[j].lineNumber
		})

	case "alphabetic":
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].prefix < nodes[j].prefix
		})

	default:
		log.Fatalf("Error: unknown order option '%v'", orderBy)
	}

	for _, kc := range nodes {
		ptree(kc, depth+1, orderBy) // print the children in order
	}
}

func makeBooleanFlag(flagVar *bool, switchName string, desc string) {
	flag.BoolVar(flagVar, switchName, false, desc)
	flag.BoolVar(flagVar, string(switchName[0]), false, desc)
}

func main() {
	max := 2

	var help bool
	var orderBy string

	odf := struct {
		switchString string
		defaul       string
		description  string
	}{"order", "input", "Sort order: input|alphabetic"}

	flag.StringVar(&orderBy, string(odf.switchString[0]), odf.defaul, odf.description)
	flag.StringVar(&orderBy, odf.switchString, odf.defaul, odf.description)

	makeBooleanFlag(&help, "help", "Print helpful text.")

	flag.Parse()
	helpText(os.Stderr, help)

	file := os.Stdin
	defer file.Close()

	isSep := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	isNotSep := func(c rune) bool {
		return !isSep(c)
	}

	root := node{-1, "", "", map[string]*node{}}
	scanner := bufio.NewScanner(file)
	nr := 0
	for scanner.Scan() {
		line := scanner.Text()
		nr++
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
		// seps = append(seps, "$")
		add(nr, &root, t, seps, max, 0)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ptree(&root, 0, orderBy)
}

func helpText(out io.Writer, doOrNotDo bool) {
	if !doOrNotDo {
		return
	}
	usage := `
USAGE:

 $ frangipanni [-h|-help] [-o|-order input|alphabetic]

	-o -order :    Sort the nodes either in input order or via character ordering
	-h -help  :    Prints this text.
`
	fmt.Fprintln(out, usage)
	os.Exit(0)
}
