package onkyo

import "testing"
import "fmt"

func Test_Detect(t *testing.T) {
	for _, s := range OnkyoAutoDetect(1000) {
		fmt.Printf(">> %v\n", s)
	}
}
