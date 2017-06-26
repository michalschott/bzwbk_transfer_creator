package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
)

type Transfer struct {
	transfer_type    int
	transfer_account string
	transfer_name    string
	transfer_address string
	transfer_value   float32
	transfer_title   string
}

var Config struct {
	your_account_number string
	input_file          string
	output_file         string
}

func read_template_file(Transfers []Transfer) ([]Transfer, error) {
	fmt.Printf("==> Starting to parse template %s\n\n", Config.input_file)
	// Open destination file
	f, err := os.Open(Config.input_file)
	if err != nil {
		return Transfers, errors.New("Can't open intput_file\n")
	}
	defer f.Close()

	// Read file line by line
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	Config.your_account_number = scanner.Text()

	for scanner.Scan() {
		// Split line
		line := strings.Split(scanner.Text(), "|")

		// Do necessary format convertions
		transfer_type, err := strconv.Atoi(line[0])
		if err != nil {
			return Transfers, errors.New("Can't convert string to int\n")
		}
		transfer_value, err := strconv.ParseFloat(line[4], 32)
		if err != nil {
			return Transfers, errors.New("Can't convert string to float\n")
		}

		// Fill array with transfers
		Transfers = append(Transfers, Transfer{transfer_type, line[1], line[2], line[3], float32(transfer_value), line[5]})
		fmt.Printf("Imported transfer no #%d\n", len(Transfers))
	}

	return Transfers, nil
}

func render_transfers_to_file(Transfers []Transfer) (int, error) {
	fmt.Printf("\n\n==> Starting to render output file %s\n\n", Config.output_file)
	// Initial checks for limits
	if len(Transfers) > 20 {
		return 1, errors.New("Too many transfers\n")
	}

	// Create destination file
	f, err := os.Create(Config.output_file)
	if err != nil {
		return 1, errors.New("Can't create output_file\n")
	}
	defer f.Close()

	// Create write buffer to file
	w := bufio.NewWriter(f)

	// Write data to buffer
	fmt.Fprintln(w, "4120414")
	for c := 0; c < len(Transfers); c++ {
		fmt.Fprintf(w, "%d|%s|%s|%s|%s|%f|1|%s||\n", Transfers[c].transfer_type, Config.your_account_number, Transfers[c].transfer_account, Transfers[c].transfer_name, Transfers[c].transfer_address, Transfers[c].transfer_value, Transfers[c].transfer_title)
		fmt.Fprintf(os.Stdout, "Rendered transfer no #%d\n", c+1)
	}

	// Flush buffers
	w.Flush()

	return 0, nil
}

func configure() error {
	user, err := user.Current()
	if err != nil {
		return errors.New("Couldn't get current user properties\n")
	}

	in := []string{user.HomeDir, "bank_transfers_input"}
	out := []string{user.HomeDir, "bank_transfers_output"}

	Config.input_file = strings.Join(in, "/")
	Config.output_file = strings.Join(out, "/")

	return nil
}

func main() {
	log.SetPrefix("")
	os.Exit(realMain())
}

func realMain() int {
	var err error
	var Transfers []Transfer

	err = configure()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't set configuration: %s", err)
		return 1
	}

	Transfers, err = read_template_file(Transfers[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read template: %s", err)
		return 1
	}

	_, err = render_transfers_to_file(Transfers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't write destination file: %s", err)
		return 1
	}

	return 0
}
