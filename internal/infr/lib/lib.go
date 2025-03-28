package lib

import (
	"bytes"
	"core/internal/common/interfaces"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"os/exec"
)

// TODO: проверка кодов ошибок yt-dlp
// TODO: преобразование в формат mp4

const format = "mp4"

type Lib struct{}

func NewLib() interfaces.DownloaderLib {
	return &Lib{}
}

// DownloadAndSave скачивание видео в файл
func (r *Lib) DownloadAndSave(videoUrl, filename, destPath string) (string, error) {
	template := filename + ".%(ext)s"

	cmd := exec.Command("yt-dlp", "-f", format, "-o", template, "--path", destPath, videoUrl)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("download video err: %v, %w", string(output), err)
	}

	return destPath + "/" + filename + "." + format, nil
}

// DownloadStream скачивание видео в поток
func (r *Lib) DownloadStream(videoUrl, filename string) (<-chan []byte, string) {
	out := make(chan []byte)

	go func() {
		defer close(out)

		cmd := exec.Command("yt-dlp", "-f", format, "-o", "-", videoUrl)
		pipe, err := cmd.StdoutPipe()
		if err != nil {
			log.Warn().Err(err)
		}

		if err := cmd.Start(); err != nil {
			log.Warn().Err(err)
		}

		buffer := make([]byte, 1024*1024)
		for {
			n, err := pipe.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Warn().Err(err)
				return
			}

			chunk := make([]byte, n)
			copy(chunk, buffer[:n])
			out <- chunk
		}

		if err := cmd.Wait(); err != nil {
			log.Warn().Err(err)
			return
		}
	}()

	return out, filename + "." + format
}

// GetVideoMetadata возвращает метаданные видео
func (r *Lib) GetVideoMetadata(videoUrl string) (*interfaces.VideoMetadata, error) {
	var out bytes.Buffer

	cmd := exec.Command("yt-dlp", "--dump-json", videoUrl)
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Error().Err(err).Msgf("error get metadata for: %v", videoUrl)
		return nil, fmt.Errorf("error get metadata for: %w; videoUrl: %v", err, videoUrl)
	}

	var videoInfo interfaces.VideoMetadata

	err = json.Unmarshal(out.Bytes(), &videoInfo)
	if err != nil {
		log.Error().Err(err).Msgf("error parse metadata to json - %v", videoUrl)
		return nil, fmt.Errorf("error parse metadata to json: %w; videoUrl: %v", err, videoUrl)
	}

	log.Info().Msgf("received metadata for: %v; %v", videoInfo, videoUrl)

	return &videoInfo, err
}

// GetHashVideo возвращает hash видео
func (r *Lib) GetHashVideo(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Error().Err(err).Msgf("error open file - %v", filePath)
		return "", fmt.Errorf("error open file: %w; path: %v", err, filePath)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		log.Error().Err(err).Msgf("error copy hash from file - %v", filePath)
		return "", fmt.Errorf("error copy hash from file: %w", err)
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	log.Info().Msgf("calc hash: %v for: %v", hash, filePath)

	return hash, nil
}

// GetVideoFileSize возвращает размер файла
func (r *Lib) GetVideoFileSize(filePath string) (int64, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		log.Error().Err(err).Msgf("error get file size - %v", filePath)
		return 0, fmt.Errorf("error get file size: %w", err)
	}

	size := stat.Size()

	log.Info().Msgf("get file size %v for: %v", size, filePath)

	return size, err
}
