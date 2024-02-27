package brute

import (
	"bufio"
	"os"
)

const chunkSize = 10 * 1024 * 1024

func (b *Brute) Start() error {

	file, err := os.Open(b.WordListPath)
	if err != nil {
		b.LogError("Error opening file", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	buf := make([]byte, chunkSize)
	scanner.Buffer(buf, chunkSize)

	for scanner.Scan() {

		line := scanner.Text()
		if line == "" {
			continue
		}

		err := b.Publish(line)
		if err != nil {
			b.LogError("Error while publishing message", err)
			continue
		}

	}

	if err := scanner.Err(); err != nil {
		b.LogError("Error while scanning file", err)
		return err
	}

	return nil
}
