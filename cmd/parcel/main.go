package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/yardnsm/leeches/pkg/parcel"
	"golang.org/x/term"
)

const usage = `Usage:
    parcel [flags]

Options:
    -e, --encrypt         Encrypt the input to the output. Default if omitted.
    -d, --decrypt         Decrypt the input to the output.
`

var (
	decryptFlag bool
	encryptFlag bool
)

func init() {
	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	flag.BoolVar(&decryptFlag, "d", false, "decrypt the input")
	flag.BoolVar(&decryptFlag, "decrypt", false, "decrypt the input")
	flag.BoolVar(&encryptFlag, "e", false, "encrypt the input")
	flag.BoolVar(&encryptFlag, "encrypt", false, "encrypt the input")
	flag.Parse()
}

func main() {

	// TODO allow passing in files
	var in io.Reader = os.Stdin
	var out io.Writer = os.Stdout

	fmt.Fprint(os.Stderr, "Enter password: ")
	password, err := term.ReadPassword(int(syscall.Stderr))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading password:", err)
		os.Exit(1)
	}

	input, err := io.ReadAll(in)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading standard input:", err)
		os.Exit(1)
	}

	if encryptFlag {
		enc, _ := parcel.EncryptWithNonce(password, input)
		fmt.Fprintln(out, enc)
		return
	}

	if decryptFlag {
		dec, _ := parcel.DecryptWithNonce(password, string(input))
		fmt.Fprintln(out, string(dec))
	}
}
