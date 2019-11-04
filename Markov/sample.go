package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

//Prefix don't mind this one
type Prefix []string

//toString
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

//Shift removes the first word and append another one to the bottom
func (p Prefix) Shift(word string) {
	//[a b c d] -> [b c d d]
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of [prefixLen] words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	chain     map[string][]string
	prefixLen int
}

//NewChain returns a new empty Chain with prefixes of prefixLen words
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

//Build reads text from the provided Reader and
//parses it into prefixes and suffixes that are stored in Chain
//pass into r io.Reader so it reads context from a file txt or an assigned string value and returns respective value.
//And the reason for passing it in instead of creating a new one is that only 1 reader is needed to be created.
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	fmt.Println("Prefix:")
	for {
		var s string
		/*Fscan only reads one word at a time*/
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		/**************************/
		key := p.String()
		c.chain[key] = append(c.chain[key], s) //remember that 1 prefix has many suffix
		fmt.Printf("c.chain[%s]: ", key)
		fmt.Println(c.chain[key])
		p.Shift(s)
	}
}

//Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

func main() {
	//settings const value and param to pass into
	//EX: echo "a frog a horse and a chipmunk" \ | ./sample -prefix=1 -words=3
	numWords := flag.Int("words", 100, "maximum number of words to print")
	prefixLen := flag.Int("prefix", 2, "prefix length in words")

	flag.Parse()
	//apparently this piece of shit means something :/ but which part of the code did it affect though, still don't know lolz
	rand.Seed(time.Now().UnixNano())

	c := NewChain(*prefixLen)
	c.Build(os.Stdin)
	text := c.Generate(*numWords)
	fmt.Println("Result: ", text)
}
