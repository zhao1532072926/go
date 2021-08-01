package week03

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func handleHello(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "Hello Go!")
}

func serverStart(server *http.Server) error {
	http.HandleFunc("/hello", handleHello)
	return server.ListenAndServe()
}

func serverStop(ctx context.Context, server *http.Server) error {
	<-ctx.Done()
	server.Shutdown(ctx)
	return nil
}

func listenSignal(ctx context.Context) error {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		return nil
	case sig := <-signalChan:
		return errors.Errorf("get os signal: %v", sig)
	}
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	server := &http.Server{
		Addr: "0.0.0.0:8080",
	}

	g.Go(func() error {
		return serverStart(server)
	})

	g.Go(func() error {
		return serverStop(ctx, server)
	})

	g.Go(func() error {
		return listenSignal(ctx)
	})

	err := g.Wait()
	if err != nil {
		log.Printf("errgroup stoped, err: %+v \n", err)
	}
	log.Println("stop")
	return

}
