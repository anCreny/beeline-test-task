package utils

import (
	"os"
	"strconv"
)

func GetCLArgWithDefault(flag string, d string) string {
	if !isFlag(flag) {
		return d
	}

	args := os.Args

	for index, value := range args {
		if flag == value && len(args)-1 > index && !isFlag(args[index+1]) {
			return args[index+1]
		}
	}

	return d
}

func GetCLArgWithDefaultAsInt(flag string, d int) int {
	if !isFlag(flag) {
		return d
	}

	args := os.Args

	for index, value := range args {
		if flag == value && len(args)-1 > index && !isFlag(args[index+1]) {
			res, err := strconv.Atoi(args[index+1])
			if err == nil {
				return res
			}
		}
	}

	return d
}

func isFlag(arg string) bool {
	if arg != "" && rune([]byte(arg)[0]) == '-' {
		return true
	}

	return false
}
