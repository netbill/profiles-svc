package cli

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/netbill/logium"
	"github.com/sirupsen/logrus"
)

func Run(args []string) bool {
	cfg, err := LoadConfig()
	if err != nil {
		logium.Fatalf("failed to load config: %v", err)
	}

	logium.SetLevel(logrus.DebugLevel)

	log := logium.New()

	lvl, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		lvl = logrus.InfoLevel
		log.WithField("bad_level", cfg.Log.Level).Warn("unknown log level, fallback to info")
	}
	log.SetLevel(lvl)

	switch {
	case cfg.Log.Format == "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	var (
		service    = kingpin.New("profiles-svc", "A service for managing user profiles")
		serviceCmd = service.Command("run service", "starting all service processes")

		migrateCmd  = service.Command("migrate", "migrate command flags --up or --down")
		migrateUp   = migrateCmd.Flag("up", "migrate db up").Bool()
		migrateDown = migrateCmd.Flag("down", "migrate db down").Bool()

		eventsCmd = service.Command("events", "events commands")

		// INBOX
		// examples:
		// profiles-svc events inbox cleanup processing --process-id=worker-1
		// profiles-svc events inbox cleanup processing --process-id=worker-1 --process-id=worker-2
		// profiles-svc events inbox cleanup failed
		eventsInbox        = eventsCmd.Command("inbox", "inbox events")
		eventsInboxCleanup = eventsInbox.Command("cleanup", "cleanup inbox events")
		eventsInboxFailed  = eventsInboxCleanup.Command(
			"failed", "cleanup inbox events in failed state",
		)
		eventsInboxProcessing = eventsInboxCleanup.Command(
			"processing", "cleanup inbox events stuck in processing state",
		)
		eventsInboxProcessingProcessIDs = eventsInboxProcessing.Flag(
			"process-id", "cleanup only events reserved by this process id (repeatable)",
		).Strings()

		// OUTBOX
		// examples:
		// profiles-svc events inbox cleanup processing --process-id=worker-1
		// profiles-svc events inbox cleanup processing --process-id=worker-1 --process-id=worker-2
		// profiles-svc events inbox cleanup failed
		eventsOutbox        = eventsCmd.Command("outbox", "outbox events")
		eventsOutboxCleanup = eventsOutbox.Command("cleanup", "cleanup outbox events")
		eventsOutboxFailed  = eventsOutboxCleanup.Command(
			"failed", "cleanup outbox events in failed state",
		)
		eventsOutboxProcessing = eventsOutboxCleanup.Command(
			"processing", "cleanup outbox events stuck in processing state",
		)
		eventsOutboxProcessingProcessIDs = eventsOutboxProcessing.Flag(
			"process-id", "cleanup only events reserved by this process id (repeatable)",
		).Strings()
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	command, err := service.Parse(args[1:])
	if err != nil {
		log.WithError(err).Error("failed to parse arguments")
		return false
	}

	switch command {
	case serviceCmd.FullCommand():
		Start(ctx, cfg, log, &wg)
	case migrateCmd.FullCommand():
		err = Migrate(ctx, cfg.Database.SQL.URL, *migrateUp, *migrateDown)
	case eventsOutboxFailed.FullCommand():
		err = CleanupOutboxFailed(ctx, cfg, log)
	case eventsOutboxProcessing.FullCommand():
		err = CleanupOutboxProcessing(ctx, cfg, log, *eventsOutboxProcessingProcessIDs...)
	case eventsInboxFailed.FullCommand():
		err = CleanupInboxFailed(ctx, cfg, log)
	case eventsInboxProcessing.FullCommand():
		err = CleanupInboxProcessing(ctx, cfg, log, *eventsInboxProcessingProcessIDs...)
	default:
		log.Errorf("unknown command %s", command)
		return false
	}
	if err != nil {
		log.WithError(err).Error("failed to exec cmd")
		return false
	}

	wgch := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgch)
	}()

	select {
	case <-ctx.Done():
		log.Warnf("Interrupt signal received: %v", ctx.Err())
		<-wgch
	case <-wgch:
		log.Warnf("All services stopped")
	}

	return true
}
