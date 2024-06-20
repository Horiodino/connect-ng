package collectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArchitectureCollectorsRun(t *testing.T) {
	assert := assert.New(t)
	testObj := Architecture{}
	expectedResult := Result{"arch": ARCHITECTURE_ARM64}

	result, err := testObj.run(ARCHITECTURE_ARM64)

	assert.NoError(err)
	assert.Equal(expectedResult, result)
}

func TestFallBackToUnameM(t *testing.T) {
	assert := assert.New(t)

	Uname = func(flag string) (string, error) {
		if flag == "-i" {
			return "unknown", nil
		}
		return "aarch64", nil
	}

	arch, err := DetectArchitecture()
	assert.NoError(err)
	assert.Equal(arch, "aarch64")
}