// gcat
// An MIT-licensed cat clone written in Go.

package main

import (
	"errors"
	"fmt"
	"os"
)

var helpMessage = `Usage: cat [OPTION]... [FILE]...
Concatenate FILE(s) to standard output.

With no FILE, or when FILE is -, read standard input.

  -A, --show-all           equivalent to -vET
  -b, --number-nonblank    number nonempty output lines, overrides -n
  -e                       equivalent to -vE
  -E, --show-ends          display $ at end of each line
  -n, --number             number all output lines
  -s, --squeeze-blank      suppress repeated empty output lines
  -t                       equivalent to -vT
  -T, --show-tabs          display TAB characters as ^I
  -u                       (ignored)
  -v, --show-nonprinting   use ^ and M- notation, except for LFD and TAB
      --help     display this help and exit
      --version  output version information and exit

Examples:
  cat f - g  Output f's contents, then standard input, then g's contents.
  cat        Copy standard input to standard output.

For more information, please visit https://github.com/jessehorne
`

// Options is used to determine what options are currently active
// See validSingleArgs for details on each option. O* where * is the option. Example, OA means the -A option.
type Options struct {
	OA bool // TODO
	Ob bool
	Oe bool // TODO
	OE bool // TODO
	On bool // TODO
	Os bool // TODO
	Ot bool // TODO
	OT bool
	Ou bool // TODO
	Ov bool // TODO
}

var validSingleArgs = []string{
	"A", // Equivalent to -vET.
	"b", // Number all nonempty output lines, starting with 1.
	"e", // Equivalent to -vE.
	"E", // Display a '$' after the end of each line. The '\r\n' combination is shown as '^m$'.
	"n", // Number all output lines, starting with 1. This option is ignored if -b is in effect.
	"s", // Suppress repeated adjacent blank lines; output just one empty line instead of several.
	"t", // Equivalent to -vT.
	"T", // Display TAB characters as '^I'.
	"u", // Ignored; for POSIX compatibility.
	"v", // Display control characters except for LDF and TAB using '^' notation and precede characters that have
	// the high bit set with 'M-'.
}

var validDoubleArgs = []string{
	"--show-all",         // Equivalent to -vET.
	"--number-nonblank",  // Number all nonempty output lines, starting with 1.
	"--show-ends",        // Display a '$' after the end of each line. The '\r\n' combination is shown as '^m$'.
	"--number",           // Number all output lines, starting with 1. This option is ignored if -b is in effect.
	"--squeeze-blank",    // Suppress repeated adjacent blank lines; output just one empty line instead of several.
	"--show-tabs",        // Display TAB characters as '^I'.
	"--show-nonprinting", // Display control characters except for LDF and TAB using '^' notation and precede characters that have
	// the high bit set with 'M-'.
}

var argConversions = map[string]string{
	"--show-all":         "A",
	"--number-nonblank":  "b",
	"--show-ends":        "E",
	"--number":           "n",
	"--squeeze-blank":    "s",
	"--show-tabs":        "T",
	"--show-nonprinting": "v",
}

// getOptions returns the single letter option(s) if exists or an error
func getOptions(o string) ([]string, error) {
	// check if -- argument
	convertedArg, ok := argConversions[o]
	if ok {
		return []string{convertedArg}, nil
	}

	// get all options from the option string
	var opts []string
	for i := range o {
		c := o[i]

		if c == '-' {
			continue
		}

		var foundValidOpt bool
		for _, opt := range validSingleArgs {
			if opt == string(c) {
				opts = append(opts, opt)
				foundValidOpt = true
			}
		}

		if !foundValidOpt {
			return opts, errors.New(fmt.Sprintf("invalid option %s", string(c)))
		}
	}

	return opts, nil
}

// parseArgs takes an array of string arguments and returns an array of files and an array of options to apply when
// performing the cat function on each file. It returns a non-nil error if the file(s) don't exist or the option(s) are
// invalid.
func parseArgs(args []string) ([]string, []string, error) {
	var files []string
	var options []string

	for _, a := range args {
		// determine if arg is a valid option or valid file path
		if a[0] == '-' {
			// if so, add it to the list of options if its valid
			opt, err := getOptions(a)
			if err != nil {
				return files, options, err
			}
			options = append(options, opt...)
		} else {
			// check if path exists
			f, err := os.Stat(a)
			if err != nil {
				return files, options, err
			}

			// check if path is a file
			m := f.Mode()
			if !m.IsRegular() {
				return files, options, errors.New(fmt.Sprintf("file '%s' cannot be read", a))
			}
			files = append(files, a)
		}
	}

	return files, options, nil
}

// gcat takes files and options and outputs the files contents according to the options provided
func gcat(files []string, options Options) error {
	for _, f := range files {
		dat, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		var last byte
		lineNumber := 1

		// check if numbering non-empty line numbers, if so, put first one
		if options.Ob {
			fmt.Print("     1 ")
			lineNumber += 1
		}

		for iChar, char := range dat {
			var next byte
			var showNonPrinting bool

			// used for determining if a line is empty
			if iChar < len(dat)-1 {
				next = dat[iChar+1]
			}

			var atEnd bool
			if iChar == len(dat)-1 {
				atEnd = true
			}

			if options.Ov || options.Oe || options.Ot {
				showNonPrinting = true
			}

			if showNonPrinting {
				if char >= 32 {
					if char < 127 {
						fmt.Print(string(char))
					} else if char == 127 {
						fmt.Print("^?")
					} else {
						fmt.Print("M-")

						if char >= 128+32 {
							if char < 128+127 {
								fmt.Print(string(char - 128))
							} else {
								fmt.Print("^?")
							}
						} else {
							fmt.Print("^")
							fmt.Print(string(char - 128 + 64))
						}
					}
				} else if char == '\t' {
					if options.OT {
						fmt.Print("^I")
					} else {
						fmt.Print('\t')
					}
				} else if char == '\n' && !atEnd {
					// TODO: gnu cat has a break here and line count?
					fmt.Println()
				} else {
					fmt.Print('^')
					fmt.Print(string(char - 128 + 64))
				}
			} else {
				// don't show nonprinting
				if char == '\t' {
					if options.OT {
						fmt.Print("^I")
					} else {
						fmt.Print(string(char))
					}
				} else if char != '\n' {
					if char == '\r' && last == '\n' && options.OE {
						fmt.Println("^M")
					} else {
						fmt.Print(string(char))
					}
				} else if char == '\n' && !atEnd {
					fmt.Println()
				}
			}

			if char == '\n' && options.Ob && next != '\n' && !atEnd {
				fmt.Print("     ", lineNumber, " ")
				lineNumber += 1
			}

			last = char
		}

		fmt.Println() // add empty line to end
	}

	return nil
}

// toOptions takes a []string of options and turns it into Options
func toOptions(opts []string) Options {
	var options Options

	for _, o := range opts {
		switch o {
		case "A":
			options.OA = true
			options.Ov = true
			options.OE = true
			options.OT = true
		case "b":
			options.Ob = true
		case "e":
			options.Oe = true
			options.Ov = true
			options.OE = true
		case "E":
			options.OE = true
		case "n":
			options.On = true
		case "s":
			options.Os = true
		case "t":
			options.Ot = true
			options.Ov = true
			options.OT = true
		case "T":
			options.OT = true
		case "u":
			options.Ou = true
		case "v":
			options.Ov = true
		}
	}

	return options
}

func main() {
	files, options, err := parseArgs(os.Args[1:])

	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
	}

	opts := toOptions(options)

	if err := gcat(files, opts); err != nil {
		fmt.Println(err)
	}
}
