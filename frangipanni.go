package main

import (
    "bufio"
    "fmt"
    "log"
	"os"
	"strings"
)
type Node struct {
	prefix string
	children []Node
}

func add(tree Node, tok []string) {
	fmt.Printf("node " + tree.prefix + " : %s\n", tok)
	if( tree.prefix == "/") {
		return 
	}
	if len(tok) == 1 && tok[0] == tree.prefix {
		// Duplicate
		fmt.Println("duplicate")
		return
	} else {
		for _, c := range (tree.children) {
			fmt.Println("node " + tree.prefix + "child " + c.prefix)
			if tok[0] == c.prefix {
				add(c, tok[1:])
				return
			}
		}
	}
	x := Node{tok[0], nil}
	if len(tok) > 1 {
		add(x, tok[1:])
	}
	tree.children = append(tree.children, x)
}

func main() {
    file, err := os.Open("./file.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

	root := Node{"/", nil}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
		l := scanner.Text()
		t := strings.Split(l, "/")
        fmt.Println(t)
		add(root, t)
    }
	fmt.Println(root)
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}