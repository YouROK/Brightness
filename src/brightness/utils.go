package brightness

import (
	"io/ioutil"
	"strconv"
	"strings"
)

func readFileInt(name string) int {
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	ret, err := strconv.Atoi(strings.TrimSpace(string(buf)))
	if err != nil {
		panic(err)
	}
	return ret
}

func writeFileInt(name string, val int) {
	strVal := strconv.Itoa(val) + "\n"
	err := ioutil.WriteFile(name, []byte(strVal), 0666)
	if err != nil {
		panic(err)
	}
}
