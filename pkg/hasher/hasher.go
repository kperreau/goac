package hasher

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
)

func Files(files []string, hashPool *sync.Pool) (string, error) {
	h := hashPool.Get().(hash.Hash)
	defer hashPool.Put(h)
	h.Reset()

	hf := hashPool.Get().(hash.Hash)
	defer hashPool.Put(hf)

	files = append([]string(nil), files...)
	sort.Strings(files)
	for _, file := range files {
		if strings.Contains(file, "\n") {
			return "", errors.New("filenames with newlines are not supported")
		}
		r, err := os.Open(file)
		if err != nil {
			return "", err
		}

		hf.Reset()
		_, err = io.Copy(hf, r)
		r.Close()
		if err != nil {
			return "", err
		}
		fmt.Fprintf(h, "%x  %s\n", hf.Sum(nil), file)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func WithPool(hashPool *sync.Pool, s string) (string, error) {
	h := hashPool.Get().(hash.Hash)
	defer hashPool.Put(h)
	h.Reset()
	_, err := h.Write([]byte(s))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func NewPool() *sync.Pool {
	return &sync.Pool{
		New: func() any { return sha1.New() },
	}
}
