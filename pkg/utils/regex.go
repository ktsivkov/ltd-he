package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func RegexpMatchString(exp *regexp.Regexp, content string) (string, error) {
	matches := exp.FindStringSubmatch(content)
	if matches == nil {
		return "", fmt.Errorf("could not match pattern: %s", exp)
	}
	return matches[1], nil
}

func RegexpMatchInt(exp *regexp.Regexp, content string) (int, error) {
	matches := exp.FindStringSubmatch(content)
	if matches == nil {
		return 0, fmt.Errorf("could not match pattern: %s", exp)
	}
	return strconv.Atoi(matches[1])
}

func RegexpMatchAllInt(exp *regexp.Regexp, content string) ([]int, error) {
	matches := exp.FindStringSubmatch(content)
	if matches == nil {
		return nil, fmt.Errorf("could not match pattern: %s", exp)
	}

	res := make([]int, 0, len(matches)-1)
	for _, match := range matches[1:] {
		parsed, err := strconv.Atoi(match)
		if err != nil {
			return nil, err
		}
		res = append(res, parsed)
	}
	return res, nil
}
