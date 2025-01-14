package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {

	//creating a json log handler to log to stdout
	jsonLog := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(jsonLog))

	//creating a context with a timeout of 2 minutes to stop the program
	//after 2 minutes if all producers aren't done
	ctx, cancelFn := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancelFn()

	buf := []chan int{}
	defer clearBuffers(buf)

	//creating 10 buffered channels
	for i := 0; i < 10; i++ {
		buf = append(buf, make(chan int, rand.Intn(10)))
	}

	var wg sync.WaitGroup

	// creating 10 producers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go produce(ctx, i, buf, &wg)
	}

	//creating 10 consumers
	for i := 0; i < 10; i++ {
		go consume(ctx, i, buf)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	stop := make(chan int)
	go func() {
		wg.Wait()
		stop <- 1
	}()

	for {
		select {
		case <-quit:
			slog.Info("received interrupt signal. stopping program")
		case <-ctx.Done():
			slog.Info("application timeout. stopping program")
			return
		case <-stop:
			slog.Info("all producers are done. stopping program")
			return
		}
	}
}

func produce(ctx context.Context, producerId int, buffs []chan int, g *sync.WaitGroup) {
	l := slog.With(slog.Int("producer", producerId))
	for j := 0; j < 100; j++ {
		select {
		case <-ctx.Done():
			return
		default:
			buffIndex := rand.Intn(len(buffs))
			l.Info(fmt.Sprintf("producer:%d peaking buffer:%d", producerId, buffIndex))
			b := buffs[buffIndex]

			select {
			case b <- rand.Intn(1000):
				l.Info(fmt.Sprintf("producer: %d pushed to buffer: %d", producerId, buffIndex))
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				break
			default:
				l.Warn(fmt.Sprintf("producer: %d buffer: %d is full. new attempt in 1 second", producerId, buffIndex))

				time.Sleep(time.Second)
			}
		}
	}
	g.Done()
}

func consume(ctx context.Context, consumerId int, buffs []chan int) {
	l := slog.With(slog.Int("consumer", consumerId))
	for {
		select {
		case <-ctx.Done():
			return
		default:
			buffIndex := rand.Intn(len(buffs))
			b := buffs[buffIndex]

			l.Info(fmt.Sprintf("consumer:%d peaked buffer:%d", consumerId, buffIndex), slog.Int("buffer", buffIndex), slog.Int("cap", cap(b)), slog.Int("len", len(b)))

			if len(b) == 0 {
				l.Warn(fmt.Sprintf("consumer:%d buffer:%d is empty. try new buffer", consumerId, buffIndex), slog.Int("cap", cap(b)), slog.Int("len", len(b)))
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				continue
			}
			l.Info(fmt.Sprintf("consumer:%d popped from buffer:%d value:%d", consumerId, buffIndex, <-b), slog.Int("buffer", buffIndex), slog.Int("cap", cap(b)), slog.Int("len", len(b)))
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}
}

func clearBuffers(ch []chan int) {
	for _, c := range ch {
		close(c)
	}
}
