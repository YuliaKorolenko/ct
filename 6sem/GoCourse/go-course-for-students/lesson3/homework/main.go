package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Options struct {
	From      string
	To        string
	Offset    int
	Limit     int
	Conv      convArray
	BlockSize int
}
type convArray []Converter

type Converter interface {
	Convert(string) string
}

type OffsetReader struct {
	reader io.Reader
}

func (r *OffsetReader) read(offset int, block int) error {
	forOffset := make([]byte, block)
	for i := 0; i < offset; {
		n, err := r.reader.Read(forOffset)
		i += n
		if err == io.EOF && i < offset {
			return fmt.Errorf("offset is more than file length: %w", err)
		}
		if err != nil && err != io.EOF {
			return fmt.Errorf("error in readBuf: %w", err)
		}
	}
	return nil
}

type LowerCaseConverter struct {
}

func (*LowerCaseConverter) Convert(string1 string) string {
	return strings.ToLower(string1)
}

type UpperCaseConverter struct {
}

func (*UpperCaseConverter) Convert(string1 string) string {
	return strings.ToUpper(string1)
}

type TrimConverter struct {
	isBegin bool
	spaces  []rune
}

func (t *TrimConverter) Convert(string1 string) string {
	if t.isBegin {
		string1 = strings.TrimLeftFunc(string1, unicode.IsSpace)
	}
	currentSpaces := make([]rune, 0)
	runes := []rune(string1)
	for i := len(runes) - 1; i >= 0; i-- {
		if unicode.IsSpace(runes[i]) {
			currentSpaces = append(currentSpaces, runes[i])
			runes = runes[:i]
		} else {
			break
		}
	}
	if string(runes) != "" {
		t.isBegin = false
		runes = append(t.spaces, runes...)
		t.spaces = currentSpaces
	} else if !t.isBegin {
		t.spaces = append(t.spaces, currentSpaces...)
	}
	return string(runes)

}

func ParseFlags() (*Options, error) {
	var opts Options
	var convString string

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", 0, "offset. by default - 0")
	flag.IntVar(&opts.Limit, "limit", math.MaxInt, "limit. by default - -1")
	flag.StringVar(&convString, "conv", "", "file to write. by default - stdout")
	flag.IntVar(&opts.BlockSize, "block-size", 1000, "block size. by default - 48")
	flag.Parse()

	opts.Conv = make([]Converter, 0)

	if opts.Offset < 0 {
		return &opts, errors.New("offset is less than 0")
	}

	if convString == "" {
		return &opts, nil
	}

	flags := strings.Split(convString, ",")
	isContr := false
	for _, element := range flags {
		if element == "upper_case" && !isContr {
			isContr = true
			opts.Conv = append(opts.Conv, &UpperCaseConverter{})
		} else if element == "lower_case" && !isContr {
			isContr = true
			opts.Conv = append(opts.Conv, &LowerCaseConverter{})
		} else if element == "trim_spaces" {
			opts.Conv = append(opts.Conv, &TrimConverter{true, make([]rune, 0)})
		} else if isContr && (element == "upper_case" || element == "lower_case") {
			return &opts, errors.New("contradictory flags")
		} else {
			return &opts, errors.New("wrong flag")
		}
	}
	return &opts, nil
}

func chooseReader(opts *Options) (*os.File, error) {
	if opts.From == "" {
		return os.Stdin, nil
	}
	f, err := os.Open(opts.From)
	if err != nil {
		return nil, errors.New("can not open file")
	}
	return f, nil

}

func chooseWriter(opts *Options) (*os.File, error) {
	if opts.To == "" {
		return os.Stdout, nil
	}
	_, err := os.Stat(opts.To)
	if !os.IsNotExist(err) {
		return nil, errors.New("file already exist")
	}
	f, err := os.Create(opts.To)
	if err != nil {
		return nil, fmt.Errorf("can not open file :%w", err)

	}
	return f, nil
}

func convMethods(opts *Options, convResult string) string {
	for _, element := range opts.Conv {
		convResult = element.Convert(convResult)
	}
	return convResult
}

func writeBuf(changeString string, writeFile *os.File) error {
	_, err := writeFile.Write([]byte(changeString))

	if err != nil {
		_, err = fmt.Fprintln(os.Stderr, "can not write:", err)
		return err
	}

	return nil
}

func decodeArray(bytesSafe []byte, block []byte) (string, []byte) {
	allBytes := append(bytesSafe, block...)
	prefix := make([]byte, 0)
	buffer := make([]byte, 0)
	runes := make([]rune, 0)
	for len(allBytes) > 0 {
		myRune, n := utf8.DecodeRune(allBytes)
		if myRune == utf8.RuneError {
			buffer = append(buffer, allBytes[:n]...)
		} else {
			if len(runes) == 0 {
				prefix = buffer
				buffer = make([]byte, 0)
			}
			runes = append(runes, myRune)
		}
		allBytes = allBytes[n:]
	}
	answer := string(prefix) + string(runes)
	return answer, buffer
}

func readBuf(opts *Options) error {
	readFile, err := chooseReader(opts)
	if err != nil {
		return err
	}
	defer readFile.Close()

	writeFile, err := chooseWriter(opts)
	if err != nil {
		return fmt.Errorf("error in writeBuf: %w", err)
	}
	defer writeFile.Close()

	offsetReader := OffsetReader{io.LimitReader(readFile, int64(opts.Offset))}
	err = offsetReader.read(opts.Offset, opts.BlockSize)
	if err != nil {
		return fmt.Errorf("error in writeBuf: %w", err)
	}

	buffer := make([]byte, opts.BlockSize)
	bytesSafe := make([]byte, 0)
	limitedReader := io.LimitReader(readFile, int64(opts.Limit))
	for {
		n, err := limitedReader.Read(buffer)
		runeStr, bytesSafeTemp := decodeArray(bytesSafe, buffer[:n])
		bytesSafe = bytesSafeTemp
		changeString := convMethods(opts, runeStr)

		errBuf := writeBuf(changeString, writeFile)

		if err == io.EOF {
			if len(bytesSafe) > 0 {
				_ = writeBuf(string(bytesSafe), writeFile)
			}
			return nil
		}
		if errBuf != nil {
			return err
		}

	}
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}
	err = readBuf(opts)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "read error", err)
		os.Exit(1)
	}

}
