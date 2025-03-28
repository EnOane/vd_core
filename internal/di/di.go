package di

import (
	"core/internal/common/interfaces"
	"core/internal/infr/lib"
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

var Container *dig.Container

func MakeDIContainer() {
	Container = dig.New()

	makeProviders()
}

func makeProviders() {
	Container.Provide(func() interfaces.DownloaderLib {
		return lib.NewLib()
	})
}

func Inject[T any]() T {
	var dep T

	err := Container.Invoke(func(d T) { dep = d })
	if err != nil {
		log.Fatal().Err(err)
	}

	return dep
}
