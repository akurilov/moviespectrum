package util

import (
	"io"
	"os"
)

func WriteToFile(in *io.ReadCloser, outFileName string) (int64, error) {
	var size int64 = 0
	var err error = nil
	outFile, err := os.OpenFile(outFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err == nil {
		defer outFile.Close()
		size, err = io.Copy(outFile, *in)
	}
	return size, err
}
