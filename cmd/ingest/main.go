package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	// DatasetSourceURL points to the AOL Search Logs dataset archive
	DatasetSourceURL string = "https://archive.org/download/academictorrents_cd339bddeae7126bb3b15f3a72c903cb0c401bd1/AOL_search_data_leak_2006.zip"

	// FinalDataset is the output CSV file containing unique queries with their frequencies
	FinalDataset string = "./data/dataset.csv"
)

// datasetFilePattern matches data files within the compressed archive
var datasetFilePattern = regexp.MustCompile(`^AOL-user-ct-collection/user-ct-test-collection-\d\d\.txt\.gz$`)

func main() {
	// Download and decompress the dataset from source
	// Note: ZIP format requires random access so it must be loaded into memory
	response, err := http.Get(DatasetSourceURL)
	if err != nil {
		log.Fatalf("downloading compressed dataset: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %v", err)
	}
	compressed, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("reading entire compressed dataset in-memory: %v", err)
	}

	// Extract all dataset files from the ZIP archive
	readerAt := bytes.NewReader(compressed)
	zipReader, err := zip.NewReader(readerAt, int64(len(compressed)))
	if err != nil {
		log.Fatalf("creating zip reader over compressed dataset: %v", err)
	}

	// Index all queries with their frequency counts
	// This map accumulates frequency data as we process the dataset
	queryFreq := make(map[string]int)

	// Process all files inside the compressed dataset.
	for _, file := range zipReader.File {
		// Ignore files which do not contain the data that we are interested in.
		if !datasetFilePattern.MatchString(file.Name) {
			continue
		}

		if err := processFile(queryFreq, file); err != nil {
			log.Fatalf("processing decompressed file containing query data: %s: %v", file.Name, err)
		}
	}

	// Create the final dataset csv file to write the queries and their frequencies into.
	dataset, err := os.Create(FinalDataset)
	if err != nil {
		log.Fatalf("create final dataset file: %v", err)
	}
	defer dataset.Close()

	// Create a writer over the final dataset file.
	csvWriter := csv.NewWriter(dataset)
	defer csvWriter.Flush()

	// Error check if flush fails
	if err := csvWriter.Error(); err != nil {
		log.Fatalf("flush csv writer: %v", err)
	}

	// Define and write the columns at the top of the csv file
	header := []string{"Query", "Frequency"}
	if err := csvWriter.Write(header); err != nil {
		log.Fatalf("write columns to the final dataset: %v", err)
	}

	// Process all the unique queries and write all of them
	// including the frequency count into the final dataset.
	for query, frequency := range queryFreq {
		if err := csvWriter.Write([]string{query, strconv.Itoa(frequency)}); err != nil {
			log.Fatalf("write query and frequency to final dataset: %v", err)
		}
	}
}

// processFile processes a decompressed ZIP file by iterating over all queries
// and updating the queryFreq map for each query.
func processFile(queryFreq map[string]int, file *zip.File) error {
	// Get the reader over the file inside the compressed dataset.
	zipDecompressed, err := file.Open()
	if err != nil {
		return fmt.Errorf("create file reader for reading file from compressed dataset: %v", err)
	}
	defer zipDecompressed.Close()

	// Each file in this compressed dataset is a gzip-compressed txt file.
	gzipReader, err := gzip.NewReader(zipDecompressed)
	if err != nil {
		return fmt.Errorf("creating gzip reader over file from compressed dataset: %v", err)
	}
	defer gzipReader.Close()

	// Each gzip-compressed txt file contains rows of tab separated values.
	// We are only interested in the queries and hence we read the records line by line
	// and only process the queries.
	lineReader := bufio.NewScanner(gzipReader)
	for lineReader.Scan() {
		line := lineReader.Text()
		fields := strings.Split(line, "\t")
		if len(fields) >= 2 {
			queryFreq[fields[1]]++
		}
	}

	if err := lineReader.Err(); err != nil {
		return err
	}
	return nil
}
