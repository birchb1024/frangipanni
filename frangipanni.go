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
	//	"runtime"
	//	"runtime/pprof"
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
	firstSep := ""
	restSeps := []string{}
	if len(sep) > 0 {
		firstSep = sep[0]
		restSeps = sep[1:]
	}
	for _, c := range tree.children {
		//fmt.Printf("children %d node %s child %d %s\n", depth, tree.prefix, i, c.prefix)
		if tok[0] == c.prefix {
			add(lineNumber, c, tok[1:], restSeps, max, depth+1)
			return
		}
	}
	// So not a match to the children. It's a new child.

	x := node{lineNumber, tok[0], firstSep, map[string]*node{}}
	tree.children[tok[0]] = &x
	add(lineNumber, &x, tok[1:], restSeps, max, depth+1)
	//fmt.Printf("newchild %d %s\n", depth, tree)
}

func fprintTree(out io.Writer, t *node, depth int, orderBy string) {

	for i := 0; i < depth; i++ { // Indentation
		fmt.Fprint(out, "  ")
	}
	x := t // temp pointer
	for len(x.children) == 1 { // Print singletons on the same line
		if !printSeparators && x == t { // First one
			fmt.Fprint(out, x.prefix)
		} else {
			fmt.Fprint(out, x.sep+x.prefix)
		}
		for k := range x.children { // Get first and only child, loops once.
			x = x.children[k]
		}
	}
	if !printSeparators && x == t { // First one
		fmt.Fprintln(out, x.prefix)
	} else {
		fmt.Fprintln(out, x.sep+x.prefix)
	}
	// Convert map to list for sorting
	nodes := make([]*node, 0, len(x.children)) // list of nodes
	for n := range x.children {
		nodes = append(nodes, x.children[n])
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
		fprintTree(out, kc, depth+1, orderBy) // print the children in order
	}
}

func fprintTreeJSON(out io.Writer, t *node, depth int, orderBy string) {

	if t.lineNumber < 0 { // root node
		fmt.Fprint(out, "{")
	}
	if len(t.children) == 0 {
		fmt.Fprint(out, "\""+t.sep+t.prefix+"\": null")
		return
	}
	fmt.Fprint(out, "\""+t.sep+t.prefix+"\" : ")

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

	fmt.Fprint(out, "{")
	for i, kc := range nodes {
		fprintTreeJSON(out, kc, depth+1, orderBy) // print the children in order
		if i < len(nodes)-1 {
			fmt.Fprint(out, ",")
		}
	}
	fmt.Fprint(out, "}")
	if t.lineNumber < 0 { // root node
		fmt.Fprint(out, "}")
	}

}

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
//var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

// Nasty Globals for options ;-)
var printSeparators bool
var fieldSeparators string
var orderBy string
var format string

func main() {
	max := 2

	var stdoutBuffered *bufio.Writer
	stdoutBuffered = bufio.NewWriter(os.Stdout)
	defer stdoutBuffered.Flush()

	flag.BoolVar(&printSeparators, "separators", false, "Print leading separators.")
	flag.StringVar(&orderBy, "order", "input", "Sort order input|alphabetic. Sort the nodes either in input order or via character ordering")
	flag.StringVar(&format, "format", "indent", "Format of output: indent|json")
	flag.Parse()

	/* 	if *cpuprofile != "" {
	   		f, err := os.Create(*cpuprofile)
	   		if err != nil {
	   			log.Fatal("could not create CPU profile: ", err)
	   		}
	   		defer f.Close() // error handling omitted for example
	   		if err := pprof.StartCPUProfile(f); err != nil {
	   			log.Fatal("could not start CPU profile: ", err)
	   		}
	   		defer pprof.StopCPUProfile()
	   	}
	*/

	file := os.Stdin
	defer file.Close()

	isSep := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	isNotSep := func(c rune) bool {
		return !isSep(c)
	}

	root := node{-1, "stdin", "", map[string]*node{}}
	scanner := bufio.NewScanner(file)
	nr := 0
	for scanner.Scan() {
		line := scanner.Text()
		nr++
		if len(line) == 0 {
			continue // skip empty lines
		}
		t := strings.FieldsFunc(line, isSep)
		seps := []string{}
		seps = strings.FieldsFunc(line, isNotSep)
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
	switch format {
	case "indent":
		fprintTree(stdoutBuffered, &root, 0, orderBy)

	case "json":
		fprintTreeJSON(stdoutBuffered, &root, 0, orderBy)

	default:
		log.Fatalf("Error: unknown format option '%v'", format)
	}

	/* 	if *memprofile != "" {
	   		f, err := os.Create(*memprofile)
	   		if err != nil {
	   			log.Fatal("could not create memory profile: ", err)
	   		}
	   		defer f.Close() // error handling omitted for example
	   		runtime.GC()    // get up-to-date statistics
	   		if err := pprof.WriteHeapProfile(f); err != nil {
	   			log.Fatal("could not write memory profile: ", err)
	   		}
	   	}
	*/
}
