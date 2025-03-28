package broker

import (
	"context"
	"core/internal/domain"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
)

func MustConnect() {
	nc, err := nats.Connect("nats://192.168.0.147:4222")
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка подключения к NATS")
	}
	log.Info().Msg("Подключено к NATS на 192.168.0.147:4222")

	ctx := context.TODO()

	// Create a JetStream management interface
	js, _ := jetstream.New(nc)

	// Create a stream
	js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "DOWNLOAD",
		Subjects: []string{"DOWNLOADER.*"},
	})

	// Подписываемся на запросы загрузки
	_, err = nc.Subscribe("download.request", func(msg *nats.Msg) {
		log.Info().Msgf("Получен запрос: %s", string(msg.Data))
		go handle(ctx, js, msg)
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка подписки на download.request")
	}

	select {} // Блокируем основной поток, чтобы подписка работала
}

func handle(ctx context.Context, js jetstream.JetStream, msg *nats.Msg) {
	msg.Ack()

	in, _ := domain.Execute(string(msg.Data))

	for chunk := range in {
		log.Debug().Msgf("отправлено %d", len(chunk))

		asd, err := js.Publish(ctx, "DOWNLOADER.stream", chunk)
		if err != nil {
			log.Fatal().Err(err).Msg("Ошибка публикации чанка в JetStream")
		}

		fmt.Println(asd.Sequence)

	}

	// Отправляем сигнал завершения в отдельный канал
	_, err := js.Publish(ctx, "DOWNLOADER.complete", []byte("done"))
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка отправки завершения")
	}

	log.Debug().Msg("все отправлено")
}
