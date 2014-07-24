package main

import (
	"os"
	"regexp"
)

var regexSplitEnv = regexp.MustCompile(`^([^=]*)[=](.*)$`)

func EnvToMap(env []string) map[string]string {
	m := make(map[string]string)

	for _, v := range env {
		match := regexSplitEnv.FindStringSubmatch(v)
		if match != nil {
			//fmt.Printf("match = %#v\n", match)
			if len(match) != 3 {
				panic("regexSplitEnv must return two groups")
			}
			m[match[1]] = match[2]
		}
	}

	return m
}

func MapToEnv(m map[string]string) {
	for k, v := range m {
		InjectToEnv(k, v)
	}
}

func InjectToEnv(key, val string) {
	var err error
	err = os.Setenv(key, val)
	if err != nil {
		panic(err)
	}
}
