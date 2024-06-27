package semver

import (
	"strings"

	git "github.com/purpleclay/gitz"
)

// Increment defines the different types of increment that can be performed
// against a semantic version
type Increment string

type ParseOptions struct {
	TrimHeader bool
}

const (
	// NoIncrement represents no increment change to a semantic version
	NoIncrement Increment = "None"
	// PatchIncrement represents a patch increment (1.0.x) to a semantic version
	PatchIncrement Increment = "Patch"
	// MinorIncrement represents a minor increment (1.x.0) to a semantic version
	MinorIncrement Increment = "Minor"
	// MajorIncrement represents a major increment (x.0.0) to a semantic version
	MajorIncrement Increment = "Major"
)

const (
	colonSpace     = ": "
	featUpper      = "FEAT"
	fixUpper       = "FIX"
	breaking       = "BREAKING CHANGE: "
	breakingHyphen = "BREAKING-CHANGE: "
	breakingBang   = '!'
)

// ParseLog will identify the maximum semantic increment by parsing the commit
// log against the conventional commit standards defined, @see:
// https://www.conventionalcommits.org/en/v1.0.0/
func ParseLog(log []git.LogEntry) Increment {
	return ParseLogWithOptions(log, ParseOptions{TrimHeader: false})
}

func ParseLogWithOptions(log []git.LogEntry, options ParseOptions) Increment {
	mode := NoIncrement
	for _, entry := range log {
		// Check for the existence of a conventional commit type
		colonSpaceIdx := strings.Index(entry.Message, colonSpace)
		if colonSpaceIdx == -1 {
			continue
		}

		startIdx := 0
		// Commit messages may have leading lines before the conventional commit type
		if options.TrimHeader {
			startIdx = FindStartIdx(entry.Message)
		}

		leadingType := strings.ToUpper(entry.Message[startIdx:colonSpaceIdx])
		if leadingType[len(leadingType)-1] == breakingBang || multilineBreaking(entry.Message) {
			return MajorIncrement
		}

		// Only feat and fix types now make a difference. Both have the same first letter
		if leadingType[0] != featUpper[0] {
			continue
		}

		if mode == MinorIncrement {
			continue
		}

		if contains(leadingType, featUpper) {
			mode = MinorIncrement
		} else if contains(leadingType, fixUpper) {
			mode = PatchIncrement
		}
	}

	return mode
}

func contains(str, prefix string) bool {
	if str == prefix {
		return true
	}

	if strings.HasPrefix(str, prefix) {
		if len(str) > len(prefix) &&
			(str[len(prefix)] == '(' && str[len(str)-1] == ')') {
			return true
		}
	}

	return false
}

func multilineBreaking(msg string) bool {
	n := strings.Count(msg, "\n")
	if n == 0 {
		return false
	}

	idx := strings.LastIndex(msg, "\n")

	if idx == len(msg) {
		// There is a newline at the end of the string, so jump back one
		if idx = strings.LastIndex(msg[:len(msg)-1], "\n"); idx == -1 {
			return false
		}
	}

	footer := msg[idx+1:]
	return strings.HasPrefix(footer, "BREAKING CHANGE: ") ||
		strings.HasPrefix(footer, "BREAKING-CHANGE: ")
}

func FindStartIdx(msg string) int {
	colonIdx := strings.Index(msg, colonSpace)
	if colonIdx == -1 {
		return 0
	}

	trimmedMsg := msg[:colonIdx]
	leadingLineBreakIdx := strings.LastIndex(trimmedMsg, "\n")
	if leadingLineBreakIdx == -1 {
		return 0
	}

	return leadingLineBreakIdx + 1
}
