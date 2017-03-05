package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
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

const (
	NUMBER = 1
	STRING = 2
	SYMBOL = 3
)

func main() {
	const (
		READY = iota
		READING
		STRREAD
		ESCAPE
		COMMENT
	)

	var tmp bytes.Buffer
	var tokens []token

	f, _ := os.Open("script.ls")

	r := bufio.NewReader(f)
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

	var expression sexpr
	for _, tok := range tokens {
		if tok == "(" {
			expression = sexpr{
				atoms: []atom{},
			}
		} else if tok == ")" {
			result := eval(expression)
			fmt.Printf("Script Result: %s\n", result.token)
		} else if strings.Count(string(tok), "\"") == 2 {
			singleAtom := atom{
				token:  tok,
				typing: STRING,
			}
			expression.atoms = append(expression.atoms, singleAtom)
		} else if string(tok) == "+" {
			singleAtom := atom{
				token:  tok,
				typing: SYMBOL,
			}
			expression.atoms = append(expression.atoms, singleAtom)
		} else {
			singleAtom := atom{
				token:  tok,
				typing: NUMBER,
			}
			expression.atoms = append(expression.atoms, singleAtom)
		}
	}
}

func eval(expression sexpr) atom {
	isConcat := false
	isAddition := false
	for _, a := range expression.atoms {
		if a.typing == STRING {
			isConcat = true
		} else if a.typing == NUMBER {
			isAddition = true
		}
	}

	if !isAddition && isConcat {
		result := concat(expression)
		return atom{
			token:  token(result),
			typing: STRING,
		}
	} else if isAddition && !isConcat {
		result := add(expression)
		return atom{
			token:  token(strconv.Itoa(result)),
			typing: STRING,
		}
	}
	return atom{
		token:  "0",
		typing: STRING,
	}
}

func concat(expression sexpr) string {
	values := []string{}
	for _, atom := range expression.atoms {
		if atom.token != "+" {
			unpackedValue := strings.Replace(string(atom.token), "\"", "", 2)
			values = append(values, unpackedValue)
		}
	}
	return fmt.Sprintf("\"%s\"", strings.Join(values, ""))
}

func add(expression sexpr) int {
	r := 0
	for _, atom := range expression.atoms {
		if atom.token != "+" {
			n, err := strconv.Atoi(string(atom.token))
			if err != nil {
				panic(err)
			}
			r += n
		}
	}
	return r
}
