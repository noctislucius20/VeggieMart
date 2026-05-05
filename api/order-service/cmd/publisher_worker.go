package cmd

import (
	"order-service/worker"

	"github.com/labstack/gommon/color"
	"github.com/spf13/cobra"
)

var workerPublisherCmd = &cobra.Command{
	Use:   "publisher-worker",
	Short: "Menjalankan worker untuk publish RabbitMQ",
	Run: func(cmd *cobra.Command, args []string) {
		color.Println(color.Green("Publisher worker is running..."))
		worker.StartPublisherWorker()
	},
}
