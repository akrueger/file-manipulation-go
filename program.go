package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type label struct {
	value string
	dcn   string
}

var labelMap = map[string]label{
	"Proposal 1-4": {
		"",
		"",
	},
	"Proposal 5-8": {
		"",
		"",
	},
	"Proposal 9-12": {
		"",
		"",
	},
	"Proposal 13-16": {
		"",
		"",
	},
	"Proposal 17-18": {
		"",
		"",
	},
	"Proposal 19-20": {
		"",
		"",
	},
	"Proposal 21-24": {
		"",
		"",
	},
	"Proposal 25-28": {
		"",
		"",
	},
	"Proposal 29-32": {
		"",
		"",
	},
	"Proposal 33": {
		"",
		"",
	},
	"Proposal 34-37": {
		"",
		"",
	},
	"Proposal 38-41": {
		"",
		"",
	},
	"Proposal 42-45": {
		"",
		"",
	},
	"Proposal 46-49": {
		"",
		"",
	},
	"Proposal 50-53": {
		"",
		"",
	},
	"Proposal 54-57": {
		"",
		"",
	},
	"Legal Proxy": {
		"",
		"",
	},
	"For All Except": {
		"",
		"",
	},
	"Signature 1": {
		"",
		"",
	},
	"Signature 2": {
		"",
		"",
	},
	"ERROR": {
		"",
		"",
	},
}

type proposals struct {
	one string
	two string
}

func main() {
	file1 := os.Args[1]
	file2 := os.Args[2]

	fileOneLines := getLines(file1)
	fileTwoLines := getLines(file2)

	if len(fileOneLines) != len(fileTwoLines) {
		errors.New("files do not contain an equal number of lines")
	}

	var outputLines []string

	for index := range fileOneLines {
		if index%2 != 0 {
			for key := range labelMap {
				fileOneLabel := fileOneLines[index][key]
				fileTwoLabel := fileTwoLines[index][key]
				if len(fileOneLabel.value) > 0 && len(fileTwoLabel.value) > 0 {

					if strings.Contains(key, "Proposal") {
						p := proposals{fileOneLabel.value, fileTwoLabel.value}
						pPrime := formatProposals(p)

						fileOneLabel.value = pPrime.one
						fileTwoLabel.value = pPrime.two
					}

					if fileOneLabel.value != fileTwoLabel.value {
						fileOneValueOutput := fmt.Sprintf("%v, %v, %v, %v, %v", filepath.Base(file1), index-2, fileOneLabel.dcn, key, fileOneLabel.value)
						fileTwoValueOutput := fmt.Sprintf("%v, %v, %v, %v, %v", filepath.Base(file2), index-2, fileTwoLabel.dcn, key, fileTwoLabel.value)

						outputLines = append(outputLines, fileOneValueOutput)
						outputLines = append(outputLines, fileTwoValueOutput)
					}
				}
			}
		}
	}

	writeFile(outputLines, file1, file2)
}

func getLines(filename string) []map[string]label {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var textLines []string

	for scanner.Scan() {
		textLines = append(textLines, scanner.Text())
	}

	file.Close()

	lineNumber := 1

	var lineSlice []map[string]label

	for _, line := range textLines {

		fileLabelMap := make(map[string]label)

		for key, value := range labelMap {
			fileLabelMap[key] = value
		}

		for key := range fileLabelMap {
			if key == "ERROR" {
				fileLabelMap[key] = label{findLabelValues(line, fmt.Sprintf("[%s]", key), "|"), findLabelValues(line, "[DCN]", "[")}
			} else if key == "For All Except" {
				s := findLabelValues(line, fmt.Sprintf("[%s]", key), "[")
				if len(s) > 0 {
					s = string(s[0])
					fileLabelMap[key] = label{s, findLabelValues(line, "[DCN]", "[")}
				}
			} else {
				fileLabelMap[key] = label{findLabelValues(line, fmt.Sprintf("[%s]", key), "["), findLabelValues(line, "[DCN]", "[")}
			}
		}
		lineSlice = append(lineSlice, fileLabelMap)

		lineNumber += 1
	}

	return lineSlice
}

func findLabelValues(line string, left string, right string) string {
	rx := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(left) + `(.*?)` + regexp.QuoteMeta(right))
	matches := rx.FindAllStringSubmatch(line, -1)
	for _, v := range matches {
		return v[1]
	}
	return ""
}

func writeFile(lines []string, file1 string, file2 string) {
	filename1 := strings.Split(filepath.Base(file1), ".")[0]
	filename2 := strings.Split(filepath.Base(file2), ".")[0]

	err := os.MkdirAll("output", os.ModePerm)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	outputDirectory := "output"
	outputFile := "diff-" + filename1 + "-" + filename2 + ".csv"

	file, err := os.Create(filepath.Join(outputDirectory, filepath.Base(outputFile)))
	if err != nil {
		fmt.Println(err)
		file.Close()
		return
	}

	for _, value := range lines {
		fmt.Fprintln(file, value)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s written successfully\n", outputFile)
}

func formatProposals(p proposals) proposals {
	s1 := strings.ReplaceAll(p.one, ",", "")
	s2 := strings.ReplaceAll(p.two, ",", "")
	s1a := strings.Split(s1, ";")
	s2a := strings.Split(s2, ";")

	for index := range s1a {
		s1a[index] = s1a[index][1:]
		s2a[index] = s2a[index][1:]

		if s1a[index] == s2a[index] {
			s1a[index] = "***"
			s2a[index] = "***"
		}
	}

	s1b := strings.Join(s1a, "|")
	s2b := strings.Join(s2a, "|")

	return proposals{s1b, s2b}
}
