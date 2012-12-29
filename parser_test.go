package parser

import (
	"bytes"
	"container/list"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Diff cut'n'paste from http://golang.org/src/cmd/gofmt/gofmt.go
func diff(b1, b2 []byte) (data []byte, err error) {
	f1, err := ioutil.TempFile("", "parser")
	if err != nil {
		return
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "parser")
	if err != nil {
		return
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(b1)
	f2.Write(b2)

	data, err = exec.Command("diff", "-u", f1.Name(), f2.Name()).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	return

}

func TestParser(t *testing.T) {
	var p PegParser
	if data, err := ioutil.ReadFile("./peg.peg"); err != nil {
		t.Fatalf("%s", err)
	} else {
		p.data = strings.NewReader(string(data))
		if !p.Grammar() {
			t.Fatalf("Didn't parse correctly")
		} else {
			if data, err := ioutil.ReadFile("./parser.go"); err != nil {
				t.Fatalf("%s", err)
			} else {
				data2 := []byte(GenerateParser(p.RootNode(), &GoGenerator{name: "PegParser"}))
				if cmp := bytes.Compare(data, data2); cmp != 0 {
					d, _ := diff(data, data2)
					t.Fatalf("Generated parser isn't equal to self: %d\n%s\n", cmp, string(d))
				}
			}
		}
	}
}

func BenchmarkParser(b *testing.B) {
	var p PegParser
	if data, err := ioutil.ReadFile("./peg.peg"); err != nil {
		b.Fatalf("%s", err)
	} else {
		p.data = strings.NewReader(string(data))
		for i := 0; i < b.N; i++ { //use b.N for looping
			p.data.Seek(0, 0)
			p.stack.Clear()
			p.currentNode.Children = list.List{}
			p.Grammar()
		}
	}
}

// func TestParser2(t *testing.T) {
// 	var p PegParser2
// 	if data, err := ioutil.ReadFile("./peg.peg"); err != nil {
// 		t.Fatalf("%s", err)
// 	} else {
// 		p.data = strings.NewReader(string(data))
// 		if !p.Grammar() {
// 			t.Fatalf("Didn't parse correctly")
// 		} else {
// 			//t.Log(p.currentNode)
// 			//p.Dump()
// 		}
// 	}
// }

// func BenchmarkParser2(b *testing.B) {
// 	var p PegParser
// 	if data, err := ioutil.ReadFile("./peg.peg"); err != nil {
// 		b.Fatalf("%s", err)
// 	} else {
// 		p.data = strings.NewReader(string(data))
// 		for i := 0; i < b.N; i++ { //use b.N for looping
// 			p.data.Seek(0, 0)
// 			//			p.stack.Clear()
// 			//			p.currentNode.Children = list.List{}
// 			p.Grammar()
// 		}
// 	}
// }
