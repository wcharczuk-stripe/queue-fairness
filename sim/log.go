package sim

import (
	"fmt"
	"strings"
)

type logTag struct {
	K string
	V any
}

func log(message string, tags ...logTag) {
	var tagStrings []string
	for _, t := range tags {
		tagStrings = append(tagStrings, fmt.Sprintf("%s=%v", t.K, t.V))
	}
	fmt.Println(message, strings.Join(tagStrings, " "))
}
