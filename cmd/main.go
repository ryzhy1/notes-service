package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"testovoe/internal/app"
	"testovoe/internal/config"
)

const (
	envDev   = "dev"
	envProd  = "prod"
	envLocal = "local"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	fmt.Println(`                                                $$\                               $$\   
                                                $$ |                            $$$$ |  
 $$$$$$\        $$\   $$\       $$$$$$$$\       $$$$$$$\        $$\   $$\       \_$$ |  
$$  __$$\       $$ |  $$ |      \____$$  |      $$  __$$\       $$ |  $$ |        $$ |  
$$ |  \__|      $$ |  $$ |        $$$$ _/       $$ |  $$ |      $$ |  $$ |        $$ |  
$$ |            $$ |  $$ |       $$  _/         $$ |  $$ |      $$ |  $$ |        $$ |  
$$ |            \$$$$$$$ |      $$$$$$$$\       $$ |  $$ |      \$$$$$$$ |      $$$$$$\ 
\__|             \____$$ |      \________|      \__|  \__|       \____$$ |      \______|
                $$\   $$ |                                      $$\   $$ |              
                \$$$$$$  |                                      \$$$$$$  |              
                 \______/                                        \______/               `)

	log.Info("Starting http", "env", cfg.Env)

	application := app.New(log, cfg.Server.Port, cfg.Storage, cfg.TokenTTL)

	go application.HTTPServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("Application stopped", slog.String("signal", sign.String()))

	application.HTTPServer.Stop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
