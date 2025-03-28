package lib

import (
	"core/internal/common/interfaces"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

type exp struct {
	id       string
	fileId   string
	provider string
	link     string
	hash     string
}

var expectedTable = [...]exp{
	{
		id:       "X-xPsJfIWK0",
		fileId:   uuid.New().String(),
		provider: "youtube",
		link:     "https://youtube.com/shorts/X-xPsJfIWK0",
		hash:     "79f060e8a9eb17ed207656f3af90d1441a994415b207a717b99ac38d1e4e0648",
	},
	{
		id:       "-46638176_456239535",
		fileId:   uuid.New().String(),
		provider: "vk",
		link:     "https://vk.com/clip-46638176_456239535",
		hash:     "a9b27902a6b38d8fd37b6f9f3824a92eabbed5df82ab260bc8db9beff4154866",
	},
	{
		id:       "ce0d3b5fddbb6829282d7a406f9df882",
		fileId:   uuid.New().String(),
		provider: "rutube",
		link:     "https://rutube.ru/shorts/ce0d3b5fddbb6829282d7a406f9df882",
		hash:     "1ea4515c67d8e2f5f96f573580efb2c3a1b5c14bb519242c44b13042413a7e59",
	},
}

type TestLib struct {
	*Lib
}

func (r *TestLib) DownloadStream(v, f string) (<-chan []byte, string) {
	ch := make(chan []byte)

	go func() {
		defer close(ch)
		ch <- []byte(v)
	}()

	return ch, "test.mp4"
}

func (r *TestLib) DownloadAndSave(v, f, p string) (string, error) {
	file, _ := os.CreateTemp(p, f)
	defer file.Close()

	_, err := file.WriteString(v)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func (r *TestLib) GetVideoMetadata(v string) (*interfaces.VideoMetadata, error) {
	m := &interfaces.VideoMetadata{Title: v}

	if strings.Contains(v, "youtube") {
		m.Id = expectedTable[0].id
	}

	if strings.Contains(v, "vk") {
		m.Id = expectedTable[1].id
	}

	if strings.Contains(v, "rutube") {
		m.Id = expectedTable[2].id
	}

	return m, nil
}

var l = &TestLib{}

func TestDownloadStreamParallel(t *testing.T) {
	for _, e := range expectedTable {
		t.Run("TestDownloadStreamParallel "+e.provider, func(t *testing.T) {
			t.Parallel()

			in, filename := l.DownloadStream(e.link, e.fileId)

			actual := make([]byte, 0, len(e.link))
			for bytes := range in {
				actual = append(actual, bytes...)
			}

			assert.Equal(t, int64(len(e.link)), int64(len(actual)))
			assert.NotEmpty(t, filename)
		})
	}
}

func TestDownloadAndSaveParallel(t *testing.T) {
	tmpDir := t.TempDir()

	for _, e := range expectedTable {
		t.Run("DownloadAndSaveParallel "+e.provider, func(t *testing.T) {
			t.Parallel()

			actual, err := l.DownloadAndSave(e.link, e.fileId, tmpDir)
			hash, _ := l.GetHashVideo(actual)

			assert.Nil(t, err)
			assert.Equal(t, e.hash, hash)
		})
	}
}

func TestGetVideoMetadataParallel(t *testing.T) {
	for _, e := range expectedTable {
		t.Run("TestGetHashVideoParallel "+e.provider, func(t *testing.T) {
			t.Parallel()

			actual, err := l.GetVideoMetadata(e.link)

			assert.Nil(t, err)
			assert.Equal(t, e.id, actual.Id)
		})
	}
}

func TestGetHashVideoParallel(t *testing.T) {
	tmpDir := t.TempDir()

	for _, e := range expectedTable {
		t.Run("TestGetHashVideoParallel "+e.provider, func(t *testing.T) {
			t.Parallel()

			filePath, err := l.DownloadAndSave(e.link, e.fileId, tmpDir)
			actual, err := l.GetHashVideo(filePath)

			assert.Nil(t, err)
			assert.Equal(t, e.hash, actual)
		})
	}
}

func TestGetVideoFileSizeParallel(t *testing.T) {
	tmpDir := t.TempDir()

	for _, e := range expectedTable {
		t.Run("TestGetVideoFileSizeParallel "+e.provider, func(t *testing.T) {
			t.Parallel()

			filePath, err := l.DownloadAndSave(e.link, e.fileId, tmpDir)
			actual, err := l.GetVideoFileSize(filePath)

			assert.Nil(t, err)
			assert.Equal(t, int64(len(e.link)), actual)
		})
	}
}
