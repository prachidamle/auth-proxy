package main

import (
	"net/http"
	"os"

	"github.com/rancher/auth-proxy/server"
	"github.com/rancher/auth-proxy/service"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var VERSION = "v0.0.0-dev"

func main() {
	app := cli.NewApp()
	app.Name = "auth-proxy"
	app.Version = VERSION
	app.Usage = "You need help!"
	app.Action = StartService

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "cluster-config",
			Usage: "Kube config for accessing cluster",
		},
		cli.StringFlag{
			Name:  "cluster-name",
			Usage: "name of the cluster",
		},
		cli.StringFlag{
			Name:  "listen",
			Usage: "port to listen on",
		},
	}

	app.Run(os.Args)
}

func StartService(c *cli.Context) {

	server.SetEnv(c)

	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	textFormatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(textFormatter)

	log.Info("Starting Rancher Auth proxy")

	router := service.NewRouter()

	log.Info("Listening on ", c.GlobalString("listen"))
	log.Fatal(http.ListenAndServe(c.GlobalString("listen"), router))

}
