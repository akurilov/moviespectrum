package util

import (
	"io"
	"os"
)

func WriteToFile(in *io.ReadCloser, outFileName string) (int64, error) {
	var size int64 = 0
	var err error = nil
	outFile, err := os.Open(outFileName)
	if err == nil {
		defer outFile.Close()
		size, err = io.Copy(outFile, *in)
	}
	return size, err
}
