package cli

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/netbill/logium"
	"github.com/netbill/profiles-svc/internal/boot"
	"github.com/netbill/profiles-svc/internal/messenger/cleaning"
	"github.com/netbill/profiles-svc/migrations"
)

func Run(args []string) bool {
	cfg, err := boot.LoadConfig()
	if err != nil {
		logium.Fatalf("failed to load config: %v", err)
	}

	log := boot.NewLogger(cfg.Log)

	var (
		service    = kingpin.New("profiles-svc", "A service for managing user profiles")
		runCmd     = service.Command("run", "run command flags: service")
		serviceCmd = runCmd.Command("service", "starting all service processes")

		migrateCmd     = service.Command("migrate", "migrate command")
		migrateUpCmd   = migrateCmd.Command("up", "migrate db up")
		migrateDownCmd = migrateCmd.Command("down", "migrate db down")

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
		boot.RunService(ctx, log, &wg, cfg)
	case migrateUpCmd.FullCommand():
		err = migrations.MigrateUp(ctx, log, cfg.Database.SQL.URL)
	case migrateDownCmd.FullCommand():
		err = migrations.MigrateDown(ctx, log, cfg.Database.SQL.URL)
	case eventsOutboxFailed.FullCommand():
		err = cleaning.CleanupOutboxFailed(ctx, log, cfg.Database.SQL.URL)
	case eventsOutboxProcessing.FullCommand():
		err = cleaning.CleanupOutboxProcessing(ctx, log, cfg.Database.SQL.URL, *eventsOutboxProcessingProcessIDs...)
	case eventsInboxFailed.FullCommand():
		err = cleaning.CleanupInboxFailed(ctx, log, cfg.Database.SQL.URL)
	case eventsInboxProcessing.FullCommand():
		err = cleaning.CleanupInboxProcessing(ctx, log, cfg.Database.SQL.URL, *eventsInboxProcessingProcessIDs...)
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
		log.Infof("Interrupt signal received: %v", ctx.Err())
		<-wgch
	case <-wgch:
		log.Info("All services stopped")
	}

	return true
}
