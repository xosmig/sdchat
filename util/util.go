package util

import "bufio"

// DiscardLineFromReader discards all bytes until the end of line.
func DiscardLineFromReader(reader *bufio.Reader) error {
	for {
		_, isPrefix, err := reader.ReadLine()
		if err != nil {
			return err
		}
		if !isPrefix {
			return nil
		}
	}
}
