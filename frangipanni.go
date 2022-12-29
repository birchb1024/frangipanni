package main

//
// Usage:  cat file.txt | go run frangipanni.go
//
import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	luajson "github.com/layeh/gopher-json"
	"github.com/yuin/gopher-lua"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var Version = "none"

type node struct {
	lineNumber int
	text       string
	sep        string
	children   map[string]*node
	numMatched int // The number of matches of text
	depth      int // The depth of this node in the tree
}

//
// If this node's children have any leaf, return true.
func (n *node) hasChildLeaves() bool {

	for _, c := range n.children {
		if len(c.children) == 0 {
			return true
		}
	}
	return false
}

//
func sliceHasLeaves(nodes []*node) bool {

	for _, n := range nodes {
		if len(n.children) == 0 {
			return true
		}
	}
	return false
}

func depthAdjust(n *node) {
	for _, c := range n.children {
		c.depth = n.depth + 1
		depthAdjust(c)
	}
}

func fold(n *node) *node {
	if len(n.children) == 0 {
		return n
	}
	for txt, c := range n.children {
		nc := fold(c)
		delete(n.children, txt)
		n.children[nc.text] = nc
	}
	if len(n.children) != 1 || n.depth == 0 { // Don't fold into the root node.
		return n
	}
	for _, c := range n.children {
		// contatenate this into the parent node
		n.text = n.text + c.sep + c.text
		n.children = c.children
	}
	depthAdjust(n)
	return n
}

func add(lineNumber int, n *node, tok []string, sep []string) {
	if len(tok) < 1 {
		return
	}
	firstSep := ""
	restSeps := []string{}
	if len(sep) > 0 {
		firstSep = sep[0]
		restSeps = sep[1:]
	}
	for _, c := range n.children {
		if tok[0] == c.text {
			c.numMatched++
			add(lineNumber, c, tok[1:], restSeps)
			return
		}
	}
	// So not a match to the children. It's a new child.
	x := node{lineNumber, tok[0], firstSep, map[string]*node{}, 1, n.depth + 1}
	n.children[tok[0]] = &x
	add(lineNumber, &x, tok[1:], restSeps)
}

func fprintchildslice(out io.Writer, childs []*node) {

	for _, kc := range childs {
		fprintTree(out, kc) // print the children in order
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
func reOrder(o bool) bool {
	if sortDescending {
		return !o
	}
	return o
}

func nodeGetChildrenSliceSorted(x *node) []*node {

	children := nodeGetChildrenSlice(x)
	switch sortBy {
	case "input":
		sort.SliceStable(children, func(i, j int) bool {
			return reOrder(children[i].lineNumber < children[j].lineNumber)
		})

	case "alpha":
		sort.SliceStable(children, func(i, j int) bool {
			if children[i].text == children[j].text {
				return reOrder(children[i].lineNumber < children[j].lineNumber)
			}
			return reOrder(children[i].text < children[j].text)
		})

	case "counts":
		sort.SliceStable(children, func(i, j int) bool {
			if children[i].numMatched == children[j].numMatched {
				return reOrder(children[i].lineNumber < children[j].lineNumber)
			}
			return reOrder(children[i].numMatched < children[j].numMatched)
		})

	default:
		log.Fatalf("Error: unknown sort value '%v'", sortBy)
	}
	return children
}

func fprintTree(out io.Writer, x *node) {
	if x.depth > printDepth {
		return
	}

	if x.depth != 0 { // Special case for the empty root node - do not print it
		indent(out, x.depth)

		count := ""
		spacer := ": "
		if indentString != " " {
			spacer = indentString
		}
		if printCounts {
			count = spacer + strconv.Itoa(x.numMatched)
		}
		if !printSeparators {
			fmt.Fprintln(out, x.text+count)
		} else {
			fmt.Fprintln(out, x.sep+x.text+count)
		}
	}
	childs := nodeGetChildrenSliceSorted(x)
	fprintchildslice(out, childs)
}

func escapeJSON(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func indent(out io.Writer, depth int) {
	for i := 0; i < depth-1; i++ {
		for ts := 0; ts < indentWidth; ts++ {
			out.Write([]byte(indentString))
		}
	}
}
func fprintNodeChildrenListJSON(out io.Writer, childs []*node, depth int) {

	if depth+1 > printDepth {
		return
	}
	if len(childs) == 0 {
		return
	}
	if len(childs) == 1 {
		fprintNodeJSON(out, childs[0])
		return
	}
	if depth > 0 {
		fmt.Fprint(out, "\n")
	}
	indent(out, depth+1)
	fmt.Fprint(out, "[")
	for i, c := range childs {
		if i > 0 {
			fmt.Fprint(out, "\n")
			indent(out, depth+1)
		}
		fprintNodeJSON(out, c)
		if i < len(childs)-1 {
			fmt.Fprint(out, ",")
		}
	}
	fmt.Fprint(out, "]")

}

func fprintNodeChildrenMapJSON(out io.Writer, childs []*node, depth int) {

	if depth+1 > printDepth {
		fmt.Fprint(out, "null")
		return
	}
	if len(childs) == 0 {
		return
	}
	if depth > 0 {
		fmt.Fprint(out, "\n")
	}
	indent(out, depth+1)
	fmt.Fprint(out, "{")
	for i, c := range childs {
		ctext := escapeJSON(c.text)
		if printSeparators {
			ctext = escapeJSON(c.sep + c.text)
		}

		if i > 0 {
			fmt.Fprint(out, "\n")
			indent(out, depth+1)
		}
		fmt.Fprint(out, ctext+" : ")
		fprintNodeChildrenJSON(out, c)
		if i < len(childs)-1 {
			fmt.Fprint(out, ",")
		}
	}
	fmt.Fprint(out, "}")
}

func fprintNodeChildrenJSON(out io.Writer, n *node) {

	if n.depth >= printDepth {
		fmt.Fprint(out, "null")
		return
	}
	if len(n.children) == 0 {
		return
	}

	childs := nodeGetChildrenSliceSorted(n)

	if sliceHasLeaves(childs) {
		fprintNodeChildrenListJSON(out, childs, n.depth)
		return
	}
	fprintNodeChildrenMapJSON(out, childs, n.depth)
}

func fprintNodeJSON(out io.Writer, n *node) {
	if n.depth > printDepth {
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
	fprintNodeChildrenJSON(out, n)
	fmt.Fprint(out, "}")
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

}

func makeLuaTableFromNode(L *lua.LState, n *node) *lua.LTable {
	tb := L.CreateTable(6, 6)

	tb.RawSet(lua.LString("lineNumber"), lua.LNumber(n.lineNumber))
	tb.RawSet(lua.LString("text"), lua.LString(n.text))
	tb.RawSet(lua.LString("sep"), lua.LString(n.sep))
	tb.RawSet(lua.LString("numMatched"), lua.LNumber(n.numMatched))
	tb.RawSet(lua.LString("depth"), lua.LNumber(n.depth))
	ch := L.CreateTable(len(n.children), len(n.children))
	tb.RawSet(lua.LString("children"), ch)
	for k, v := range n.children {
		ch.RawSet(lua.LString(k), makeLuaTableFromNode(L, v))
	}
	return tb
}

func makeLuaTableFromFlags(L *lua.LState) *lua.LTable {
	tb := L.CreateTable(6, 6)
	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == "order" { // Deprecated
			return
		}
		tb.RawSet(lua.LString(f.Name), lua.LString(f.Value.String()))
	})
	args := L.CreateTable(len(flag.Args()), len(flag.Args()))
	for i, a := range flag.Args() {
		args.RawSetInt(i+1, lua.LString(a))
	}
	tb.RawSet(lua.LString("args"), args)
	return tb
}

func luaRun(root *node) {
	L := lua.NewState()
	luajson.Preload(L)
	L.SetGlobal("frangipanni_args", makeLuaTableFromFlags(L))
	L.SetGlobal("frangipanni", makeLuaTableFromNode(L, root))
	defer L.Close()
	if err := L.DoFile(luaFile); err != nil {
		log.Fatal(err)
	}
}

// Nasty Globals for options ;-)
var printSeparators bool
var noFold bool
var fieldSeparators string // List of characters to split line on, e.g. "/:"
var sortBy string
var format string
var maxLevel int
var splitOnCharacters bool
var printCounts bool
var printDepth int
var indentWidth int
var indentString string
var luaFile string
var skipLevel int
var sortDescending bool
var printVersion bool
var csvInput bool

func main() {

	var stdoutBuffered *bufio.Writer
	stdoutBuffered = bufio.NewWriter(os.Stdout)
	defer stdoutBuffered.Flush()

	flag.BoolVar(&printSeparators, "separators", false, "Print leading separators.")
	flag.StringVar(&sortBy, "sort", "input", "Sort by input|alpha|counts. Sort the branches either by input order, or via alphanumeric character ordering, or the branch frequency count.")
	flag.StringVar(&sortBy, "order", "input", "** Deprecated, use the -sort option **")
	flag.StringVar(&format, "format", "indent", "Format of output: indent|json")
	flag.StringVar(&fieldSeparators, "breaks", "", "Characters to slice lines with.")
	flag.BoolVar(&noFold, "no-fold", false, "Don't fold into one line.")
	flag.IntVar(&maxLevel, "level", math.MaxInt32, "Analyse down to this level (positive integer).")
	flag.BoolVar(&splitOnCharacters, "chars", false, "Slice line after each character.")
	flag.BoolVar(&printCounts, "counts", false, "Print number of matches at the end of the line.")
	flag.IntVar(&printDepth, "depth", math.MaxInt32, "Maximum tree depth to print.")
	flag.IntVar(&indentWidth, "indent", 4, "Number of spaces to indent per level.")
	flag.StringVar(&indentString, "spacer", " ", "Characters to indent lines with.")
	flag.StringVar(&luaFile, "lua", "", "Lua Script to run")
	flag.IntVar(&skipLevel, "skip", 0, "Number of leading fields to skip.")
	flag.BoolVar(&sortDescending, "down", false, "Sort branches in descending order. (default ascending)")
	flag.BoolVar(&printVersion, "version", false, "Print frangipanni's version number and exit.")
	flag.BoolVar(&csvInput, "csv", false, "Input is CSV format.")

	flag.Parse()
	if printVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

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

	root := node{-1, "", "", map[string]*node{}, 1, 0}
	var err error
	if csvInput {
		err = scanCSV(file, root)
	} else {
		err = scanText(file, isSep, isNotSep, root)
	}
	if err != nil {
		log.Fatal(err)
	}
	// calculate root nodse count
	root.numMatched = 0
	for _, c := range root.children {
		root.numMatched += c.numMatched
	}
	// fold the root node
	froot := &root
	if !noFold {
		froot = fold(&root)
	}

	if luaFile != "" {
		luaRun(&root)
		os.Exit(0)
	}

	switch format {
	case "indent":
		fprintTree(stdoutBuffered, froot)

	case "json":
		if printCounts {
			fakeCounts(froot)
		}
		fprintNodeChildrenJSON(stdoutBuffered, froot)
		fmt.Fprintln(stdoutBuffered)

	default:
		log.Fatalf("Error: unknown format option '%v'", format)
	}
}

func scanText(file *os.File, isSep func(c rune) bool, isNotSep func(c rune) bool, root node) error {
	s := bufio.NewScanner(file)
	nr := 0
	t := make([]string, 1024)
	seps := make([]string, 1024)
	for s.Scan() {
		line := s.Text()
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
				// line didn't start with a separator, so insert a fake one
				seps = append(seps, "") // add space at the end
				copy(seps[1:], seps)    // shift right
				seps[0] = ""            // inject fake at the front
			}
		}
		//
		// Skip leading fields if required
		for s := skipLevel; s > 0 && len(t) > 1 && len(seps) > 1; s-- {
			t = t[1:]
			seps = seps[1:]
		}

		if len(t) <= maxLevel {
			add(nr, &root, t, seps)
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
			add(nr, &root, childs, separators)
		}
	}
	return s.Err()
}
func scanCSV(file *os.File, root node) error {
	r := csv.NewReader(file)
	if fieldSeparators != "" {
		r.Comma = rune(fieldSeparators[0])
	}
	r.LazyQuotes = true
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = -1
	nr := 0
	t := make([]string, 1024)
	for {
		record, err := r.Read()
		seps := make([]string, len(record))
		for i:=0;i<len(record);i++ {
			seps[i] = ","
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) == 0 {
			continue // empty line
		}

		nr++
		t = record

		// Skip leading fields if required
		for s := skipLevel; s > 0 && len(t) > 1 ; s-- {
			t = t[1:]
		}

		if len(t) <= maxLevel {
			add(nr, &root, t, seps)
		} else {
			// Don't use the tokens beyond maxLevel - concatenate the remainder into one
			childs := make([]string, maxLevel+1)
			for i := 0; i < maxLevel && i < len(t); i++ {
				childs[i] = t[i]
			}
			for i := maxLevel; i < len(t) ; i++ {
				childs[maxLevel] = childs[maxLevel] + "," + t[i]
			}
			add(nr, &root, childs, seps)
		}
	}
	return nil
}
