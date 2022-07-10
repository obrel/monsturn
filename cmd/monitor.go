package cmd

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/obrel/go-lib/pkg/log"
	"github.com/obrel/monsturn/config"
	"github.com/obrel/monsturn/internal/app/monitor"
	"github.com/obrel/monsturn/internal/app/worker"
	"github.com/obrel/monsturn/internal/pkg/rdb"
	"github.com/obrel/monsturn/internal/pkg/util"
	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:     "worker",
	Short:   "Start monitor",
	Long:    "Start monsturn monitor worker.",
	PreRunE: monitorPrerun,
	RunE:    monitorRun,
}

func monitorPrerun(cmd *cobra.Command, args []string) error {
	var result error

	config.Init()

	if err := rdb.Init(config.Get().Redis); err != nil {
		result = multierror.Append(result, err)
	}

	/*
		if err := mdb.Init(config.Get().Db); err != nil {
			result = multierror.Append(result, err)
		}
	*/

	return result
}

func monitorRun(cmd *cobra.Command, args []string) error {
	log.Info("Starting monitor...")

	task := monitor.NewMonitorTask(rdb.GetConn(), config.Get().Redis.Topics, onMessage)
	wrk := worker.NewWorker("worker")
	wrk.Add("monitor", task)

	err := wrk.Start(context.Background())
	if err != nil {
		log.For("monitor", "run").Error(err)
		return err
	}

	return nil
}

func onMessage(channel string, data []byte) error {
	stat, err := util.MessageParser(string(data[:]))
	if err != nil {
		return err
	}

	fmt.Println(stat)
	return nil
}
