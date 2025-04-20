package generator

import (
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Options struct {
	InputFirstName    []string
	InputLastName     []string
	InputBirthday     []string
	InputRelatedWords []string
	InputMinLength    string
	InputMaxLength    string
	EnableLeet        bool
	EnableCapitalize  bool
}

var leetMap = map[rune]string{
	'a': "4",
	'e': "3",
	'i': "1",
	'o': "0",
	's': "5",
}

func Run(opts Options) error {
	inputs, minLength, maxLength := collectAllInputs(opts)
	words := []string{}
	words = append(words, inputs...)

	for n := 2; n <= 3; n++ {
		words = append(words, combineWordsN(inputs, n)...)
	}

	if opts.EnableLeet {
		words = append(words, leetVariants(words)...)
	}
	if opts.EnableCapitalize {
		words = append(words, caseVariants(words)...)
	}

	words = removeDuplicates(words)
	words = filterWordsByLength(words, minLength, maxLength)

	err := saveToFile(words)
	if err != nil {
		return err
	}
	return nil
}

func capitalize(word string) string {
	if len(word) == 0 {
		return ""
	}
	return strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
}

func collectAllInputs(opts Options) ([]string, int, int) {
	words := []string{}

	appendWithCap := func(list []string) {
		for _, w := range list {
			words = append(words, w)
			capW := capitalize(w)
			if capW != w {
				words = append(words, capW)
			}
		}
	}

	appendWithCap(opts.InputFirstName)
	appendWithCap(opts.InputLastName)
	appendWithCap(opts.InputRelatedWords)

	if len(opts.InputBirthday) > 0 {
		fullBirthday := strings.Join(opts.InputBirthday, "")
		words = append(words, fullBirthday)
		words = append(words, opts.InputBirthday...)
	}

	minLength := 6
	maxLength := 12
	if opts.InputMinLength != "" {
		minLength, _ = strconv.Atoi(opts.InputMinLength)
	}
	if opts.InputMaxLength != "" {
		maxLength, _ = strconv.Atoi(opts.InputMaxLength)
	}
	return words, minLength, maxLength
}

func combineWordsN(words []string, n int) []string {
	var result []string
	var combine func(word []string, used []bool)
	combine = func(word []string, used []bool) {
		if len(word) == n {
			combined := ""
			for _, w := range word {
				combined += w
			}
			result = append(result, combined)
			return
		}
		for i, w := range words {
			if !used[i] {
				used[i] = true
				combine(append(word, w), used)
				used[i] = false
			}
		}
	}
	used := make([]bool, len(words))
	combine([]string{}, used)
	return removeDuplicates(result)
}

func leetVariants(words []string) []string {
	var result []string
	for _, word := range words {
		leetWord := ""
		for _, char := range word {
			lowerChar := unicode.ToLower(char)
			if leetChar, ok := leetMap[lowerChar]; ok {
				leetWord += leetChar
			} else {
				leetWord += string(char)
			}
		}
		if leetWord != word {
			result = append(result, leetWord)
		}
	}
	return removeDuplicates(result)
}

func caseVariants(words []string) []string {
	var result []string
	for _, word := range words {
		result = append(result, word)
		upperWord := ""
		for _, char := range word {
			if unicode.IsLower(char) {
				upperWord += string(unicode.ToUpper(char))
			} else {
				upperWord += string(unicode.ToLower(char))
			}
		}
		if upperWord != word {
			result = append(result, upperWord)
		}
	}
	return removeDuplicates(result)
}

func removeDuplicates(words []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, word := range words {
		if _, ok := seen[word]; !ok {
			seen[word] = struct{}{}
			result = append(result, word)
		}
	}
	return result
}

func filterWordsByLength(words []string, minLength, maxLength int) []string {
	filtered := []string{}
	for _, word := range words {
		if len(word) >= minLength && len(word) <= maxLength {
			filtered = append(filtered, word)
		}
	}
	return filtered
}

func saveToFile(words []string) error {
	file, err := os.Create("wordlist.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	for _, word := range words {
		_, err := file.WriteString(word + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
