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
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type node struct {
	lineNumber int
	text       string
	sep        string
	children   map[string]*node
	numMatched int // The number of matches of text
}

//
// If this node's children have any leaf, return true.
func (n *node) hasChildLeaves() bool {

	for _, c := range n.children {
		//fmt.Fprintf(os.Stderr, "%v#\n", c)
		if len(c.children) == 0 {
			return true
		}
	}
	return false
}

//
func sliceHasLeaves(nodes []*node) bool {

	for _, n := range nodes {
		//fmt.Fprintf(os.Stderr, "%v#\n", c)
		if len(n.children) == 0 {
			return true
		}
	}
	return false
}

func add(lineNumber int, tree *node, tok []string, sep []string, depth int) {
	//fmt.Printf("add %d %d node (%s) <= %s %s\n", lineNumber, depth, tree.text, tok, sep)
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
		//fmt.Printf("children %d node %s child %d %s\n", depth, tree.text, i, c.text)
		if tok[0] == c.text {
			c.numMatched++
			add(lineNumber, c, tok[1:], restSeps, depth+1)
			return
		}
	}
	// So not a match to the children. It's a new child.
	x := node{lineNumber, tok[0], firstSep, map[string]*node{}, 1}
	tree.children[tok[0]] = &x
	add(lineNumber, &x, tok[1:], restSeps, depth+1)
	//fmt.Printf("newchild %d %s\n", depth, tree)
}

func fprintchildslice(out io.Writer, childs []*node, depth int, orderBy string, parent *node) {

	childs = nodeGetChildrenSliceSorted(parent)

	for _, kc := range childs {
		fprintTree(out, kc, depth+1, orderBy) // print the children in order
	}
}

func nodeGetChildrenSlice(x *node) []*node {
	// Convert map to list for sorting
	childs := make([]*node, 0, len(x.children)) // list of childs
	for n := range x.children {
		childs = append(childs, x.children[n])
	}
	return childs
}
func nodeGetChildrenSliceSorted(x *node) []*node {

	childs := nodeGetChildrenSlice(x)
	switch orderBy {
	case "input":
		sort.SliceStable(childs, func(i, j int) bool {
			return childs[i].lineNumber < childs[j].lineNumber
		})

	case "alpha":
		sort.SliceStable(childs, func(i, j int) bool {
			return childs[i].text < childs[j].text
		})

	default:
		log.Fatalf("Error: unknown order option '%v'", orderBy)
	}
	return childs
}

func fprintTree(out io.Writer, t *node, depth int, orderBy string) {

	if depth+1 > printDepth {
		return
	}
	indent(out, depth)
	x := t                                // temp pointer
	for len(x.children) == 1 && !noFold { // Print singletons on the same line
		if !printSeparators && x == t { // First one
			fmt.Fprint(out, x.text)
		} else {
			fmt.Fprint(out, x.sep+x.text)
		}
		for k := range x.children { // Get first and only child, loops once.
			x = x.children[k]
		}
	}
	count := ""
	if printCounts {
		count = ": " + strconv.Itoa(x.numMatched)
	}
	if !printSeparators && x == t { // First one
		fmt.Fprintln(out, x.text+count)
	} else {
		fmt.Fprintln(out, x.sep+x.text+count)
	}

	childs := nodeGetChildrenSlice(x)
	fprintchildslice(out, childs, depth, orderBy, x)
}

func escapeJSON(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func indent(out io.Writer, depth int) {
	for i := 0; i < depth; i++ {
		for ts := 0; ts < indentWidth; ts++ {
			out.Write([]byte(" "))
		}
	}
}
func fprintNodeChildrenListJSON(out io.Writer, childs []*node, depth int) {

	if depth+1 > printDepth {
		//fmt.Fprint(out, "null")
		return
	}
	if len(childs) == 0 {
		return
	}
	if len(childs) == 1 {
		fprintNodeJSON(out, childs[0], depth)
		return
	}
	fmt.Fprint(out, "[")
	for i, c := range childs {
		fprintNodeJSON(out, c, depth)
		if i < len(childs)-1 {
			fmt.Fprint(out, ",\n")
		}
	}
	fmt.Fprint(out, "]\n")

}

func fprintNodeChildrenMapJSON(out io.Writer, childs []*node, depth int, parent *node) {

	if depth+1 > printDepth {
		fmt.Fprint(out, "null")
		return
	}
	if len(childs) == 0 {
		return
	}
	fmt.Fprint(out, "{")
	for i, c := range childs {
		ctext := escapeJSON(c.text)
		if printSeparators {
			ctext = escapeJSON(c.sep + c.text)
		}

		fmt.Fprint(out, ctext+" : ")
		fprintNodeChildrenJSON(out, c.children, depth+1, c)
		if i < len(childs)-1 {
			fmt.Fprint(out, ",\n")
		}
	}
	fmt.Fprint(out, "}\n")
}

func fprintNodeChildrenJSON(out io.Writer, nodemap map[string]*node, depth int, parent *node) {

	if depth+1 > printDepth {
		fmt.Fprint(out, "null")
		return
	}
	if len(nodemap) == 0 {
		return
	}

	childs := nodeGetChildrenSlice(parent)

	if sliceHasLeaves(childs) {
		fprintNodeChildrenListJSON(out, childs, depth)
		return
	}
	fprintNodeChildrenMapJSON(out, childs, depth, parent)
}

func fprintNodeJSON(out io.Writer, n *node, depth int) {
	if depth+1 > printDepth {
		fmt.Fprint(out, "null")
		return
	}
	ntext := escapeJSON(n.text)
	if printSeparators {
		ntext = escapeJSON(n.sep + n.text)
	}
	if len(n.children) == 0 { // No children, it's a leaf
		if !printSeparators {
			if number, err := strconv.Atoi(n.text); err == nil {
				fmt.Fprint(out, number)
				return
			}
		}
		fmt.Fprint(out, ntext)
		return
	}
	fmt.Fprint(out, "{"+ntext+" : ")
	fprintNodeChildrenJSON(out, n.children, depth+1, n)
	fmt.Fprint(out, "}\n")
}

func fakeCounts(n *node) {
	tag := "_count_"
	if n.lineNumber == -42 {
		return
	}
	for _, c := range n.children {
		fakeCounts(c)
	}
	value := &node{
		lineNumber: -42,
		text:       strconv.Itoa(n.numMatched),
		children:   map[string]*node{},
		numMatched: 0}
	key := &node{
		lineNumber: -42,
		text:       tag,
		children:   map[string]*node{},
		numMatched: 0}
	key.children[strconv.Itoa(n.numMatched)] = value
	n.children[tag] = key

	//fmt.Printf("%+v\n", n)
}

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
//var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

// Nasty Globals for options ;-)
var printSeparators bool
var noFold bool
var fieldSeparators string // List of characters to split line on, e.g. "/:"
var orderBy string
var format string
var maxLevel int
var splitOnCharacters bool
var printCounts bool
var printDepth int
var indentWidth int

func main() {

	var stdoutBuffered *bufio.Writer
	stdoutBuffered = bufio.NewWriter(os.Stdout)
	defer stdoutBuffered.Flush()

	flag.BoolVar(&printSeparators, "separators", false, "Print leading separators.")
	flag.StringVar(&orderBy, "order", "input", "Sort order input|alpha. Sort the childs either in input order or via character ordering")
	flag.StringVar(&format, "format", "indent", "Format of output: indent|json")
	flag.StringVar(&fieldSeparators, "breaks", "", "Characters to slice lines with.")
	flag.BoolVar(&noFold, "no-fold", false, "Don't fold into one line.")
	flag.IntVar(&maxLevel, "level", math.MaxInt32, "Analyse down to this level (positive integer).")
	flag.BoolVar(&splitOnCharacters, "chars", false, "Slice line after each character.")
	flag.BoolVar(&printCounts, "counts", false, "Print number of matches at the end of the line.")
	flag.IntVar(&printDepth, "depth", math.MaxInt32, "Maximum tree depth to print.")
	flag.IntVar(&indentWidth, "indent", 4, "Number of spaces to indent per level.")

	flag.Parse()
	if maxLevel < 0 {
		log.Fatalf("Error: %d is negative.\n", maxLevel)
	}
	if indentWidth < 0 {
		log.Fatalf("Error: %d is negative.\n", indentWidth)
	}
	if fieldSeparators != "" && splitOnCharacters {
		log.Fatalln("Breaks option incompatible with chars option.")
	}
	printSeparators = printSeparators || splitOnCharacters

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
		if fieldSeparators == "" {
			return !unicode.IsLetter(c) && !unicode.IsNumber(c)
		}
		return strings.ContainsRune(fieldSeparators, c)
	}

	isNotSep := func(c rune) bool {
		return !isSep(c)
	}

	root := node{-1, "stdin", "", map[string]*node{}, 1}
	scanner := bufio.NewScanner(file)
	nr := 0
	t := make([]string, 1024)
	seps := make([]string, 1024)
	for scanner.Scan() {
		line := scanner.Text()
		nr++
		if len(line) == 0 {
			continue // skip empty lines
		}
		if splitOnCharacters {
			t = strings.Split(line, "")
			for i := 0; i < len(t); i++ {
				seps[i] = ""
			}

		} else {
			t = strings.FieldsFunc(line, isSep)
			seps = strings.FieldsFunc(line, isNotSep)
			if isNotSep(rune(line[0])) {
				// line didn't start with a seperator, so insert a fake one
				seps = append(seps, "") // add space at the end
				copy(seps[1:], seps)    // shift right
				seps[0] = ""            // inject fake at the front
			}
		}
		if len(t) <= maxLevel {
			add(nr, &root, t, seps, 0)
		} else {
			// Don't use the tokens beyond maxLevel - concatenate the remainder into one
			childs := make([]string, maxLevel+1)
			separators := make([]string, maxLevel+1)
			for i := 0; i < maxLevel && i < len(t); i++ {
				childs[i] = t[i]
				separators[i] = seps[i]
			}
			for i := maxLevel; i < len(t) && i < len(seps); i++ {
				childs[maxLevel] = childs[maxLevel] + seps[i] + t[i]
			}
			add(nr, &root, childs, separators, 0)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	switch format {
	case "indent":
		childs := nodeGetChildrenSlice(&root)
		fprintchildslice(stdoutBuffered, childs, -1, orderBy, &root)

	case "json":
		if printCounts {
			fakeCounts(&root)
		}
		fprintNodeChildrenJSON(stdoutBuffered, root.children, 0, &root)

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
