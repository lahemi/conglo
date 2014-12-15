package scraps

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Return 12 chars long "title" string, with alternating
// consonants and vowels, hastebin style.
func GenRandTitle() (title string) {
	var (
		co = "bcdfghjklmnpqrstvwxyz"
		vo = "aeiou"
		r  = rand.New(rand.NewSource(time.Now().UnixNano()))
	)
	for i := 0; i < 12; i++ {
		if (i % 2) == 0 {
			title += string(co[r.Intn(len(co))])
		} else {
			title += string(vo[r.Intn(len(vo))])
		}
	}
	return
}

func clock() string {
	return fmt.Sprint(time.Now())
}
func Print(str ...interface{}) {
	fmt.Fprintf(os.Stdout, "%s %v\n", clock(), str)
}
func PrintErr(str ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s %v\n", clock(), str)
}
