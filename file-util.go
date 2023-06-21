package tokenizer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

// This file provides functions to work with local dataset cache, ...
// Copied from `transformer` package

const (
	WeightName    = "pytorch_model.gt"
	ConfigName    = "config.json"
	TokenizerName = "tokenizer.json"

	// NOTE. URL form := `$HFpath/ModelName/resolve/main/WeightName`
	HFpath = "https://huggingface.co"
)

var (
	DUMMY_INPUT [][]int64 = [][]int64{
		{7, 6, 0, 0, 1},
		{1, 2, 3, 0, 0},
		{0, 0, 0, 4, 5},
	}
)

// CachedPath resolves and caches data based on input string, then returns fullpath to the cached data.
//
// Parameters:
// - `modelNameOrPath`: model name e.g., "bert-base-uncase" or path to directory contains model/config files.
// - `fileName`: model or config file name. E.g., "pytorch_model.py", "config.json"
//
// CachedPath does several things consequently:
// 1. Resolves input string to a  fullpath cached filename candidate.
// 2. Check it at `CachedPath`, if exists, then return the candidate. If not
// 3. Retrieves and Caches data to `CachedPath` and returns path to cached data
//
// NOTE. default `CachedDir` is at "{$HOME}/.cache/transformer"
// Custom `CachedDir` can be changed by setting with environment `GO_TRANSFORMER`
func CachedPath(modelNameOrPath, fileName string) (resolvedPath string, err error) {

	// Resolves to "candidate" filename at `CacheDir`
	cachedFileCandidate := fmt.Sprintf("%s/%s/%s", CachedDir, modelNameOrPath, fileName)

	// 1. Cached candidate file exists
	if _, err := os.Stat(cachedFileCandidate); err == nil {
		return cachedFileCandidate, nil
	}

	// 2. If valid fullpath to local file, caches it and return cached filename
	filepath := fmt.Sprintf("%s/%s", modelNameOrPath, fileName)
	if _, err := os.Stat(filepath); err == nil {
		err := copyFile(filepath, cachedFileCandidate)
		if err != nil {
			err := fmt.Errorf("CachedPath() failed at copying file: %w", err)
			return "", err
		}
		return cachedFileCandidate, nil
	}

	// 3. Cached candidate file NOT exist. Try to download it and save to `CachedDir`
	url := fmt.Sprintf("%s/%s/resolve/main/%s", HFpath, modelNameOrPath, fileName)
	// url := fmt.Sprintf("%s/%s/raw/main/%s", HFpath, modelNameOrPath, fileName)
	if isValidURL(url) {
		if _, err := http.Get(url); err == nil {
			err := downloadFile(url, cachedFileCandidate)
			if err != nil {
				err = fmt.Errorf("CachedPath() failed at trying to download file: %w", err)
				return "", err
			}

			return cachedFileCandidate, nil
		} else {
			err = fmt.Errorf("CachedPath() failed: Unable to parse '%v' as a URL or as a local path.\n", url)
			return "", err
		}
	}

	// Not resolves
	err = fmt.Errorf("CachedPath() failed: Unable to parse '%v' as a URL or as a local path.\n", url)
	return "", err
}

func isValidURL(url string) bool {

	// TODO: implement
	return true
}

// downloadFile downloads file from URL and stores it in local filepath.
// It writes to the destination file as it downloads it, without loading
// the entire file into memory. An `io.TeeReader` is passed into Copy()
// to report progress on the download.
func downloadFile(url string, filepath string) error {
	// Create path if not existing
	dir := path.Dir(filepath)
	filename := path.Base(filepath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	// Create the file with .tmp extension, so that we won't overwrite a
	// file until it's downloaded fully
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("bad status: %s(%v)", resp.Status, resp.StatusCode)

		if resp.StatusCode == 404 {
			// if filename == "rust_model.ot" {
			// msg := fmt.Sprintf("model weight file not found. That means a compatible pretrained model weight file for Go is not available.\n")
			// msg = msg + fmt.Sprintf("You might need to manually convert a 'pytorch_model.bin' for Go. ")
			// msg = msg + fmt.Sprintf("See tutorial at: 'example/convert'")
			// err = fmt.Errorf(msg)
			// } else {
			// err = fmt.Errorf("download file not found: %q for downloading", url)
			// }
			err = fmt.Errorf("download file not found: %q for downloading", url)
		} else {
			err = fmt.Errorf("download file failed: %q", url)
		}
		return err
	}

	// the total file size to download
	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	downloadSize := uint64(size)

	// Create our bytes counter and pass it to be used alongside our writer
	counter := &writeCounter{FileSize: downloadSize}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	fmt.Printf("\r%s... %s/%s completed", filename, byteCountIEC(counter.Total), byteCountIEC(counter.FileSize))
	// The progress use the same line so print a new line once it's finished downloading
	fmt.Println()

	// Rename the tmp file back to the original file
	err = os.Rename(filepath+".tmp", filepath)
	if err != nil {
		return err
	}

	return nil
}

// writeCounter counts the number of bytes written to it. By implementing the Write method,
// it is of the io.Writer interface and we can pass this into io.TeeReader()
// Every write to this writer, will print the progress of the file write.
type writeCounter struct {
	Total    uint64
	FileSize uint64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.printProgress()
	return n, nil
}

// PrintProgress prints the progress of a file write
func (wc writeCounter) printProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 50))

	// Return again and print current status of download
	fmt.Printf("\rDownloading... %s/%s", byteCountIEC(wc.Total), byteCountIEC(wc.FileSize))
}

// byteCountIEC converts bytes to human-readable string in binary (IEC) format.
func byteCountIEC(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

// CleanCache removes all files cached in transformer cache directory `CachedDir`.
//
// NOTE. custom `CachedDir` can be changed by setting environment `GO_TRANSFORMER`
func CleanCache() error {
	err := os.RemoveAll(CachedDir)
	if err != nil {
		err = fmt.Errorf("CleanCache() failed: %w", err)
		return err
	}

	return nil
}
