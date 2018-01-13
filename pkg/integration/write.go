package integration

import (
	"fmt"
	"os"
	"path"
)

func Report(testName string, time float64) {
	f, err := os.OpenFile(path.Join("../../../../", "foo.csv"), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s,%f\n", testName, time)); err != nil {
		panic(err)
	}
}
