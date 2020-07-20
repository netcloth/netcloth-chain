package types

import (
	"testing"
)

func Test_getCurAppVersionFromGenesisFile(t *testing.T) {
	v := getCurAppVersionFromGenesisFile("/Users/sun")
	t.Log(v)
}
