package utils_test

import (
	"os"
	"path/filepath"
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

func TestAddon_GetVersion(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		addon       utils.Addon
		config      utils.Config
		tocContent  string
		want        string
		wantErr     bool
		expectedErr string
	}{
		{
			name:  "version standard",
			addon: utils.Addon{Folder: "TestAddon"},
			config: utils.Config{
				WowAddonsDirectory: tmpDir,
			},
			tocContent: "## Interface: 100200\n## Version: 1.2.3\n## Title: Test Addon",
			want:       "1.2.3",
			wantErr:    false,
		},
		{
			name:  "version avec préfixe v",
			addon: utils.Addon{Folder: "TestAddon"},
			config: utils.Config{
				WowAddonsDirectory: tmpDir,
			},
			tocContent: "## Interface: 100200\n## Version: v2.0.0\n## Title: Test Addon",
			want:       "2.0.0",
			wantErr:    false,
		},
		{
			name:  "pas de version",
			addon: utils.Addon{Folder: "TestAddon"},
			config: utils.Config{
				WowAddonsDirectory: tmpDir,
			},
			tocContent:  "## Interface: 100200\n## Title: Test Addon",
			want:        "",
			wantErr:     true,
			expectedErr: "version non trouvée dans le fichier TOC",
		},
		{
			name:  "version vide",
			addon: utils.Addon{Folder: "TestAddon"},
			config: utils.Config{
				WowAddonsDirectory: tmpDir,
			},
			tocContent:  "## Interface: 100200\n## Version: \n## Title: Test Addon",
			want:        "",
			wantErr:     true,
			expectedErr: "version trouvée mais vide",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			// Créer le dossier de l'addon
			addonDir := filepath.Join(tmpDir, tt.addon.Folder)
			err := os.MkdirAll(addonDir, 0755)
			assert.NoError(err)

			// Créer le fichier .toc
			tocPath := filepath.Join(addonDir, tt.addon.Folder+"_Mainline.toc")
			err = os.WriteFile(tocPath, []byte(tt.tocContent), 0644)
			assert.NoError(err)

			// Exécuter le test
			got, err := tt.addon.GetVersion(tt.config)

			if tt.wantErr {
				assert.Error(err)
				if tt.expectedErr != "" {
					assert.Equal(tt.expectedErr, err.Error())
				}
			} else {
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}
