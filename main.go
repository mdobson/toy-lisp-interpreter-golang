package main

import (
	"bytes"
	"fmt"
	"strings"
)

type token string

type atom struct {
	token  token
	typing int
}

type sexpr struct {
	atoms []atom
}

type node struct {
	typing string
	token  token
	tree   []*node
}

const (
	EXPRESSION = "EXPRESSION"
	CHARACTER  = "CHARACTER"
	OPERATOR   = "OPERATOR"
)

const (
	READY = iota
	READING
	STRREAD
	ESCAPE
	COMMENT
)

func main() {
	//programOne := "(defun adder(x y) (+ x y))"
	tokens := readTokens("+ 1 (+ 1 (+ (+ 1 1) 1) 1)")

	program := &node{
		tree: []*node{},
	}
	var curr *node
	var prev *node
	curr = program
	for _, token := range tokens {

		switch token {
		case "(":
			newExpression := &node{
				typing: EXPRESSION,
				token:  "",
				tree:   []*node{},
			}

			addToTree(curr, newExpression)
			prev = curr
			curr = newExpression
			newNode := &node{
				typing: CHARACTER,
				token:  token,
				tree:   []*node{},
			}
			addToTree(curr, newNode)
		case "+":
			newNode := &node{
				typing: OPERATOR,
				token:  token,
				tree:   []*node{},
			}
			addToTree(curr, newNode)
		case ")":
			newNode := &node{
				typing: CHARACTER,
				token:  token,
				tree:   []*node{},
			}
			addToTree(curr, newNode)
			curr = prev
		default:
			newNode := &node{
				typing: CHARACTER,
				token:  token,
				tree:   []*node{},
			}
			addToTree(curr, newNode)
		}
	}

	crawl(program, 0)
}

func crawl(tree *node, level int) {
	var limiter string
	if level > 0 {
		for index := 0; index < level; index++ {
			limiter += ">"
		}
	}
	level++

	for _, child := range tree.tree {
		fmt.Printf("%s Token <%s>: %s\n", limiter, child.typing, child.token)

		if len(child.tree) > 0 {
			crawl(child, level)
		}
	}
}

func readTokens(script string) []token {
	var tmp bytes.Buffer
	var tokens []token

	//f, _ := os.Open("script.ls")
	//r := bufio.NewReader(f)

	r := strings.NewReader(script)
	rn, _, err := r.ReadRune()
	if err != nil {
		fmt.Printf("Error %s", err.Error())
	}

	const WS = "\r\t\n "
	const PARENS = "()"
	const NONTOKEN = WS + PARENS

	state := READY

	for err == nil {
		switch state {
		//Entry point to string
		case READY:
			if strings.ContainsRune(PARENS, rn) {
				tokens = append(tokens, token(rn))
				state = READY
			} else if strings.ContainsRune(WS, rn) {
				//Ignore ws
			} else if rn == '"' {
				tmp.WriteRune(rn)
				state = STRREAD
			} else {
				tmp.WriteRune(rn)
				state = READING
			}
		//Non single token being read in
		case READING:
			if strings.ContainsRune(NONTOKEN, rn) {
				tokens = append(tokens, token(tmp.String()))
				tmp.Reset()
				state = READY
				r.UnreadRune()
			} else {
				tmp.WriteRune(rn)
			}
		case STRREAD:
			tmp.WriteRune(rn)
			if rn == '"' {
				tokens = append(tokens, token(tmp.String()))
				tmp.Reset()
				state = READY
			}
		}

		rn, _, err = r.ReadRune()
	}

	return tokens
}

func addToTree(t *node, n *node) {
	//fmt.Printf("ADD TOKEN:%s TO TREE.\n", n.token)
	t.tree = append(t.tree, n)
}
