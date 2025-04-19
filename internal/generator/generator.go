package generator

import (
	"os"
	"strconv"
)

type Options struct {
	InputFirstName    []string
	InputLastName     []string
	InputBirthday     []string
	InputRelatedWords []string
	InputMinLength    string
	InputMaxLength    string
}

func Run(opts Options) error {
	inputs, minLength, maxLength := collectAllInputs(opts)
	words := []string{}
	words = append(words, inputs...)

	for n := 2; n <= len(inputs); n++ {
		words = append(words, combineWordsN(inputs, n)...)
	}

	words = removeDuplicates(words)
	words = filterWordsByLength(words, minLength, maxLength)

	err := saveToFile(words)
	if err != nil {
		return err
	}
	return nil
}

func collectAllInputs(opts Options) ([]string, int, int) {
	words := []string{}
	words = append(words, opts.InputFirstName...)
	words = append(words, opts.InputLastName...)
	words = append(words, opts.InputBirthday...)
	words = append(words, opts.InputRelatedWords...)
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
	return result
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
