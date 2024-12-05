package utils_test

import (
	"testing"

	"github.com/HLerman/elvui/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsElvuiPresent(t *testing.T) {
	cf := utils.Config{
		ElvuiFolders: []string{"ElvUI"},
	}
	t.Run("present", func(t *testing.T) {
		b, _ := utils.IsElvuiPresent(cf, []utils.Addon{
			{
				Folder: "toto",
			},
			{
				Folder: "ElvUI",
			},
		})
		assert.Equal(t, true, b)
	})

	t.Run("not_present", func(t *testing.T) {
		b, _ := utils.IsElvuiPresent(cf, []utils.Addon{
			{
				Folder: "toto",
			},
		})
		assert.Equal(t, false, b)
	})
}
