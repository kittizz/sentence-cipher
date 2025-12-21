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
	var inputData []byte
	var inputText string
	var err error
	isFileInput := false

	if *inputFile != "" {
		// Read from file as raw bytes
		inputData, err = os.ReadFile(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		isFileInput = true
		if *decodeFlag {
			// For decoding, input is text (encoded sentences)
			inputText = strings.TrimSpace(string(inputData))
		}
	} else if flag.NArg() > 0 {
		// Read from arguments
		inputText = strings.Join(flag.Args(), " ")
		inputData = []byte(inputText)
	} else {
		// Read from stdin
		inputText, err = readStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		inputText = strings.TrimSpace(inputText)
		inputData = []byte(inputText)
	}

	if len(inputData) == 0 && inputText == "" {
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
	var outputText string
	var outputData []byte
	isBinaryOutput := false

	if *decodeFlag {
		// Decode - output is raw bytes
		if *naturalFlag {
			outputData, err = cipher.DecodeNatural(inputText)
		} else {
			outputData, err = cipher.Decode(inputText)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding: %v\n", err)
			os.Exit(1)
		}
		isBinaryOutput = true
	} else {
		// Encode - input is raw bytes, output is text
		if *naturalFlag {
			outputText, err = cipher.EncodeNatural(inputData)
		} else {
			outputText, err = cipher.Encode(inputData)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding: %v\n", err)
			os.Exit(1)
		}
	}

	// Write output
	if *outputFile != "" {
		var writeData []byte
		if isBinaryOutput {
			writeData = outputData
		} else {
			writeData = []byte(outputText + "\n")
		}
		err := os.WriteFile(*outputFile, writeData, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}
	} else {
		if isBinaryOutput {
			// Write raw bytes to stdout (for binary decode without -o)
			os.Stdout.Write(outputData)
		} else {
			fmt.Println(outputText)
		}
	}
	_ = isFileInput
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
