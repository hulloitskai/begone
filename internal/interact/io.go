package interact

import (
	"fmt"
)

// Printf is like fmt.Printf, except it writes to p.Out.
func (p *Prompter) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(p.Out, format, a...)
}

// Println is like fmt.Println, except it writes to p.Out.
func (p *Prompter) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(p.Out, a...)
}

// Errf is like fmt.Printf, except it writes to p.Err.
func (p *Prompter) Errf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(p.Out, format, a...)
}

// Errln is like fmt.Println, except it writes to p.Err.
func (p *Prompter) Errln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(p.Err, a...)
}

func (p *Prompter) scanf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fscanf(p.In, format, a...)
}
