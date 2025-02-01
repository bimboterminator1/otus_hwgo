package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	//nolint:depguard
	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("cannot access file %s: %w", fromPath, err)
	}
	defer fromFile.Close()

	fileStat, err := fromFile.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	if !fileStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fileSize := fileStat.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("cannot open file %s: %w", toPath, err)
	}
	defer toFile.Close()

	if offset > 0 {
		_, err = fromFile.Seek(offset, 0)
		if err != nil {
			os.Remove(toPath)
			return err
		}
	}

	if limit == 0 || offset+limit > fileSize {
		limit = fileSize - offset
	}

	progressBar := pb.New(int(limit)).SetUnits(pb.U_BYTES)
	progressBar.Start()
	toWriter := io.MultiWriter(toFile, progressBar)
	_, err = io.CopyN(toWriter, fromFile, limit)
	if err != nil {
		os.Remove(toPath)
		return err
	}
	progressBar.Finish()
	return nil
}
