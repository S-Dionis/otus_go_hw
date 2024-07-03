package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrFileNotExists         = errors.New("file is not exists")
	ErrWrongOffsetValue      = errors.New("parameter offset is not correct")
	ErrWrongLimitValue       = errors.New("parameter limit is not correct")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

var closeFile = func(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Println(err)
	}
}

func isSameFile(fromPath, toPath string) (bool, error) {
	toPathStat, err := os.Stat(toPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	fromPathStat, err := os.Stat(fromPath)
	if err != nil {
		return false, err
	}

	return os.SameFile(fromPathStat, toPathStat), nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrWrongOffsetValue
	}

	if limit < 0 {
		return ErrWrongLimitValue
	}

	sameFile, err := isSameFile(fromPath, toPath)
	if err != nil {
		return err
	}

	if sameFile {
		return ErrUnsupportedFile
	}

	file, err := os.Open(fromPath)
	defer closeFile(file)

	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotExists
		}
		return ErrUnsupportedFile
	}

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	if fileStat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		limit = fileStat.Size()
	}

	to, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer closeFile(to)

	err = copyFromTo(file, to, offset, limit)
	if err != nil {
		return err
	}

	return nil
}

func copyFromTo(from *os.File, to *os.File, offset, limit int64) error {
	_, err := from.Seek(offset, 0)
	if err != nil {
		return err
	}

	var bufSize int64 = 1024

	if limit < bufSize {
		bufSize = limit
	}

	buff := make([]byte, bufSize)
	count := limit / bufSize
	rest := limit % bufSize

	total := int(count)
	if rest > 0 {
		total++
	}

	bar := pb.StartNew(total)

	for i := 0; i < int(count); i++ {
		n, err := from.Read(buff)
		if err != nil {
			return err
		}

		if n == 0 || errors.Is(err, io.EOF) {
			break
		}

		_, err = to.Write(buff[:n])
		if err != nil {
			return err
		}
		bar.Increment()
	}

	if rest > 0 {
		buff = make([]byte, rest)
		n, err := from.Read(buff)
		if err != nil {
			return err
		}
		_, err = to.Write(buff[:n])
		if err != nil {
			return err
		}
		bar.Increment()
	}

	bar.Finish()

	return nil
}
