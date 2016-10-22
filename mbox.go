package mailparser

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

type emailData []byte

const _MAX_LINE_LEN = 1024

var crlf = []byte{'\r', '\n'}

// If debug is true, errors parsing messages will be printed to stderr. If
// false, they will be ignored. Either way those messages will not appear in
// the msgs slice.
func Read(r io.Reader, debug bool) (msgs []emailData, err error) {
	var mbuf *bytes.Buffer
	lastblank := true
	br := bufio.NewReaderSize(r, _MAX_LINE_LEN)
	l, _, err := br.ReadLine()
	for err == nil {
		fs := bytes.SplitN(l, []byte{' '}, 3)
		if len(fs) == 3 && string(fs[0]) == "From" && lastblank {
			// flush the previous message, if necessary
			if mbuf != nil {
				msgs = append(msgs, mbuf.Bytes())
			}
			mbuf = new(bytes.Buffer)
		} else {
			_, err = mbuf.Write(l)
			if err != nil {
				return
			}
			_, err = mbuf.Write(crlf)
			if err != nil {
				return
			}
		}
		if len(l) > 0 {
			lastblank = false
		} else {
			lastblank = true
		}
		l, _, err = br.ReadLine()
	}
	if err == io.EOF {
		msgs = append(msgs, mbuf.Bytes())
		err = nil
	}
	return
}

// If debug is true, errors parsing messages will be printed to stderr. If
// false, they will be ignored. Either way those messages will not appear in
// the msgs slice.
func ReadFile(filename string, debug bool) ([]emailData, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	msgs, err := Read(f, debug)
	f.Close()
	return msgs, err
}
