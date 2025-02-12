package utils_test

import (
	"fmt"
	"testing"

	"github.com/enenisme/definger/utils"
)

func TestProbesContent2ProbesStruct(t *testing.T) {
	content := utils.ProbesForGetTitle
	probes := utils.ProbesContent2ProbesStruct(content)
	fmt.Println(probes)
}
