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

func isSetBrightness(lastVal, currVal, count *int) bool {
	last := *lastVal
	curr := *currVal
	ret := false
	if last != int(curr) {
		delta := last - curr
		if delta < 0 {
			delta = -delta
		}
		if delta > 5 {
			ret = true
		}
	} else {
		*count = 0
		return false
	}

	if ret == false {
		ret = *count > 9
		*count++
	}

	if ret == true {
		*lastVal = curr
		*count = 0
	}

	return ret
}
