package auxiliaryfunctions

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

func GetChecksumHex(r io.Reader) (string, error) {
	h := sha256.New()

	if _, err := io.Copy(h, r); err != nil {
		return "", fmt.Errorf("не удалось прочитать входной поток: %w", err)
	}

	sum := h.Sum(nil)
	return hex.EncodeToString(sum), nil
}
