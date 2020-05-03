package cat

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type Cat struct {
	ShowAll bool `clop:"-A; --show-all" usage:"equivalent to -vET"`

	NumberNonblank bool `clop:"-b; --number-nonblank"
	                     usage:"number nonempty output lines, overrides"`

	E bool `clop:"-e"; usage:"equivalent to -vE"`

	ShowEnds bool `clop:"-E;--show-ends"
	               usage:"display $ at end of each line"`

	Number bool `clop:"-n;--number"
	             usage:"number all output lines"`

	T bool `clop:"-t" usage:"equivalent to -vT" `

	SqueezeBlank bool `clop:"-s;--squeeze-blank"
	                   usage:"suppress repeated empty output lines"`

	ShowTabs bool `clop:"-T;--show-tabs"
	              usage:"display TAB characters as ^I"`

	ShowNonprinting bool `clop:"-v;--show-nonprinting"
	                      usage:"use ^ and M- notation, except for LFD and TAB" `

	Files []string `clop:"args=files"`

	oldNew []string
}

func writeNonblank(l []byte) []byte {
	var out bytes.Buffer

	for _, c := range l {
		switch {
		case c == 9: // '\t'
			out.WriteByte(c)
		case c >= 0 && c <= 8 || c > 10 && c <= 31:
			out.Write([]byte{'^', c + 64})
		case c >= 32 && c <= 126 || c == 10: // 10 is '\n'
			out.WriteByte(c)
		case c == 127:
			out.Write([]byte{'^', c - 64})
		case c >= 128 && c <= 159:
			out.Write([]byte{'M', '-', '^', c - 64})
		case c >= 160 && c <= 254:
			out.Write([]byte{'M', '-', c - 128})
		default:
			out.Write([]byte{'M', '-', '^', 63})
		}
	}

	return out.Bytes()
}

func (c *Cat) changOpt() {

	if c.ShowAll {
		c.ShowNonprinting = true
		c.ShowEnds = true
		c.ShowTabs = true
	}

	if c.E {
		c.ShowNonprinting = true
		c.ShowEnds = true
	}

	if c.T {
		c.ShowNonprinting = true
		c.ShowTabs = true
	}

	if c.ShowEnds {
		c.setEnds()
	}

	if c.ShowTabs {
		c.setTab()
	}

}

func (c *Cat) setTab() {
	c.oldNew = append(c.oldNew, "\t", "^I")
}

func (c *Cat) setEnds() {
	c.oldNew = append(c.oldNew, "\n", "$\n")
}

func (c *Cat) Cat(rs io.ReadSeeker, w io.Writer) {
	br := bufio.NewReader(rs)
	replacer := strings.NewReplacer(c.oldNew...)
	isSpace := 0

	for count := 1; ; count++ {

		l, e := br.ReadBytes('\n')
		if e != nil && len(l) == 0 {
			break
		}

		if c.SqueezeBlank {
			if len(bytes.TrimSpace(l)) == 0 {
				isSpace++
			} else {
				isSpace = 0
			}

			if isSpace > 1 {
				count--
				continue
			}
		}

		if len(c.oldNew) > 0 {
			l = []byte(replacer.Replace(string(l)))
		}

		if c.ShowNonprinting {
			l = writeNonblank(l)
		}

		if c.NumberNonblank || c.Number {

			if c.NumberNonblank && len(l) == 1 {
				count--
			}

			if !(c.NumberNonblank && len(l) == 1) {
				l = append([]byte(fmt.Sprintf("%6d\t", count)), l...)
			}
		}

		w.Write(l)
	}
}

func (c *Cat) Main() {

	c.changOpt()

	if len(c.Files) > 0 {
		for _, fileName := range c.Files {
			f, err := openFile(fileName)
			if err != nil {
				die("cat: %s\n", err)
			}

			c.Cat(f, os.Stdout)
			f.Close()
		}
		return
	}

	c.Cat(os.Stdin, os.Stdout)
}

func die(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

type file struct {
	*os.File
}

func (f *file) Close() {
	if f.File == os.Stdin {
		return
	}

	f.File.Close()
}

func openFile(fileName string) (*file, error) {
	if fileName == "-" {
		return &file{File: os.Stdin}, nil
	}

	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return &file{File: fd}, nil
}
