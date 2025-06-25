package auxiliaryfunctions

import (
	"crypto/sha256"
	"io"
	"log"
)

func GetChecksum(f io.Reader) []byte {
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return h.Sum(nil)
}
