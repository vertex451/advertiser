package localstorage

import (
	"advertiser/shared/pkg/storage"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
)

type service struct {
	destDir string
}

func New(dataDir string) storage.Provider {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		zap.L().Fatal("failed to create storage dir", zap.Error(err))
	}

	return &service{
		destDir: dataDir,
	}
}
func (s *service) Store(name string, body io.Reader) (string, error) {
	// Create the file to save the photo
	savePath := filepath.Join(s.destDir, name)
	saveFile, err := os.Create(savePath)
	if err != nil {
		zap.L().Error("Error creating file", zap.Error(err))

		return "", errors.Errorf("Error creating file: %v", err)
	}
	defer saveFile.Close()

	// Save the photo data to the file
	_, err = io.Copy(saveFile, body)
	if err != nil {
		zap.L().Error("Error saving file", zap.Error(err))

		return "", errors.Errorf("Error saving file: %v", err)
	}

	return savePath, nil
}

func (s *service) Load() {

}
