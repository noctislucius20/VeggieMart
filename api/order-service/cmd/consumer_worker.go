package cmd

import (
	"order-service/worker"

	"github.com/labstack/gommon/color"
	"github.com/spf13/cobra"
)

var workerConsumerCmd = &cobra.Command{
	Use:   "consumer-worker",
	Short: "Menjalankan worker untuk consume RabbitMQ dan index ke Elasticsearch",
	Run: func(cmd *cobra.Command, args []string) {
		color.Println(color.Green("Consumer worker is running..."))
		worker.StartConsumerWorker()
	},
}
