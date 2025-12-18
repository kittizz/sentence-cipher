package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	sentencecipher "github.com/kittizz/sentence-cipher"
)

const version = "1.0.0"

func main() {
	// Flags
	decodeFlag := flag.Bool("d", false, "Decode mode (default is encode)")
	naturalFlag := flag.Bool("n", false, "Use natural encoding (more varied sentences)")
	keyFlag := flag.String("k", "", "Encryption key (shuffles word lists)")
	inputFile := flag.String("i", "", "Input file (default: stdin)")
	outputFile := flag.String("o", "", "Output file (default: stdout)")
	versionFlag := flag.Bool("v", false, "Show version")
	helpFlag := flag.Bool("h", false, "Show help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `GrammarCipher - Hide messages in English sentences

Usage:
  grammarcipher [options] [text]

Options:
  -d          Decode mode (default is encode)
  -n          Use natural encoding (more varied sentences)
  -k KEY      Encryption key (shuffles word lists for added security)
  -i FILE     Read input from file
  -o FILE     Write output to file
  -v          Show version
  -h          Show help

Examples:
  # Encode text
  grammarcipher "Hello World"
  
  # Encode with key
  grammarcipher -k "my-secret-key" "Hello World"
  
  # Encode with natural mode and key
  grammarcipher -n -k "my-secret-key" "Secret message"
  
  # Decode text with key
  grammarcipher -d -k "my-secret-key" "Tom loves Mary books."
  
  # Encode from file
  grammarcipher -i secret.txt -o encoded.txt
  
  # Decode from stdin
  echo "Tom loves Mary books." | grammarcipher -d

`)
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Printf("GrammarCipher v%s\n", version)
		os.Exit(0)
	}

	// Get input
	var input string
	var err error

	if *inputFile != "" {
		// Read from file
		data, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		input = string(data)
	} else if flag.NArg() > 0 {
		// Read from arguments
		input = strings.Join(flag.Args(), " ")
	} else {
		// Read from stdin
		input, err = readStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
	}

	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Fprintln(os.Stderr, "Error: no input provided")
		flag.Usage()
		os.Exit(1)
	}

	// Create cipher (with or without key)
	var cipher *sentencecipher.Cipher
	if *keyFlag != "" {
		cipher, err = sentencecipher.NewCipher(*keyFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating cipher: %v\n", err)
			os.Exit(1)
		}
	} else {
		cipher = sentencecipher.NewDefaultCipher()
	}

	// Process
	var output string

	if *decodeFlag {
		// Decode
		var decoded string
		var err error
		if *naturalFlag {
			data, err := cipher.DecodeNatural(input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error decoding: %v\n", err)
				os.Exit(1)
			}
			decoded = string(data)
		} else {
			decoded, err = cipher.DecodeString(input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error decoding: %v\n", err)
				os.Exit(1)
			}
		}
		output = decoded
	} else {
		// Encode
		if *naturalFlag {
			output = cipher.EncodeNatural([]byte(input))
		} else {
			output = cipher.EncodeString(input)
		}
	}

	// Write output
	if *outputFile != "" {
		err := os.WriteFile(*outputFile, []byte(output+"\n"), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println(output)
	}
}

func readStdin() (string, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// Interactive mode - read single line
		reader := bufio.NewReader(os.Stdin)
		fmt.Fprint(os.Stderr, "Enter text: ")
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		return strings.TrimSpace(line), nil
	}

	// Pipe mode - read all
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
