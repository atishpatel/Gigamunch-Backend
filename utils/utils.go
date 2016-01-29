package utils

import (
	"bytes"
	"strings"
)

var (
	banWords = []string{"anal", "anus", "ballsack", "blowjob", "dick", "dildo", "nigger", "penis", "vagina"}
)

// ContainsBanWord checks if a sentence contains a ban word
func ContainsBanWord(sentence string) bool {
	for index := range banWords {
		if strings.Contains(sentence, banWords[index]) {
			return false
		}
	}
	return true
}

// ConcatenateStrings fast way of concatenating strings
func ConcatenateStrings(args ...string) string {
	var buffer bytes.Buffer
	for _, value := range args {
		buffer.WriteString(value)
	}
	return buffer.String()
}
