package domain

import (
	"core/internal/common/interfaces"
	"core/internal/di"
)

func Execute(link string) (<-chan []byte, string) {
	lib := di.Inject[interfaces.DownloaderLib]()

	return lib.DownloadStream(link, "12")
}
