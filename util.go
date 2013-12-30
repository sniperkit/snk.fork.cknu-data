package data

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"os/exec"
	"unicode"
)

var Debug bool

// Shorthand printing functions.
func pErr(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func pOut(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}

func dErr(format string, a ...interface{}) {
	if Debug {
		pErr(format, a...)
	}
}

func dOut(format string, a ...interface{}) {
	if Debug {
		pOut(format, a...)
	}
}

// Checks whether string is a hash (sha1)
func isHash(hash string) bool {
	if len(hash) != 40 {
		return false
	}

	for _, r := range hash {
		if !unicode.Is(unicode.ASCII_Hex_Digit, r) {
			return false
		}
	}

	return true
}

func shortHash(hash string) string {
	return hash[:7]
}

func readerHash(r io.Reader) (string, error) {
	bf := bufio.NewReader(r)
	h := sha1.New()
	_, err := bf.WriteTo(h)
	if err != nil {
		return "", err
	}

	hex := fmt.Sprintf("%x", h.Sum(nil))
	return hex, nil
}

func copyFile(src string, dst string) error {
	cmd := exec.Command("cp", src, dst)
	return cmd.Run()
}

func set(slice []string) []string {
	dedup := []string{}
	elems := map[string]bool{}
	for _, elem := range slice {
		_, seen := elems[elem]
		if !seen {
			dedup = append(dedup, elem)
			elems[elem] = true
		}
	}
	return dedup
}

func validHashes(hashes []string) (valid []string, err error) {
	hashes = set(hashes)

	// append only valid hashes
	for _, hash := range hashes {
		if isHash(hash) {
			valid = append(valid, hash)
		} else {
			err = fmt.Errorf("invalid <hash>: %v", hash)
		}
	}

	return
}
