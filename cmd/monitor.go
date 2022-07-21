package cmd

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/obrel/go-lib/pkg/log"
	"github.com/obrel/go-lib/pkg/wrk"
	"github.com/obrel/monsturn/config"
	"github.com/obrel/monsturn/internal/app/task"
	"github.com/obrel/monsturn/internal/pkg/pdb"
	"github.com/obrel/monsturn/internal/pkg/rdb"
	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:     "monitor",
	Short:   "Start monitor",
	Long:    "Start monsturn monitor worker.",
	PreRunE: monitorPrerun,
	RunE:    monitorRun,
}

// All configs an dependencies must be initiate here
func monitorPrerun(cmd *cobra.Command, args []string) error {
	var result error

	config.Init()

	if err := rdb.Init(config.Get().Redis); err != nil {
		result = multierror.Append(result, err)
	}

	if err := pdb.Init(config.Get().Db); err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

func monitorRun(cmd *cobra.Command, args []string) error {
	log.Info("Starting monitor...")

	wrk := wrk.NewWorker("worker")

	// Initiate monitor task
	mon := task.NewMonitorTask(rdb.GetConn(), config.Get().Redis.Topics)
	wrk.Add("monitor", mon)

	err := wrk.Start(context.Background())
	if err != nil {
		log.For("monitor", "run").Error(err)
		return err
	}

	return nil
}
