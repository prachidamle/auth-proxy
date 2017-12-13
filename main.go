package main

import (
	"net/http"
	"os"

	"github.com/rancher/auth-proxy/api"
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
			Name:  "httpHost",
			Usage: "host:port to listen on",
		},
	}

	app.Run(os.Args)
}

func StartService(c *cli.Context) {

	handler, err := api.NewTokenAndIdentityAPIHandler(c.String("cluster-config"), c.String("cluster-config"), "")
	if err != nil {
		log.Fatalf("Failed to get tokenAndIdentity handler: %v", err)
	}

	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	textFormatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(textFormatter)

	log.Info("Starting Rancher Auth proxy")

	httpHost := c.GlobalString("httpHost")
	server := &http.Server{
		Handler: handler,
		Addr:    httpHost,
	}
	log.Infof("Starting http server listening on %v.", httpHost)
	err = server.ListenAndServe()
	log.Infof("https server exited. Error: %v", err)

}
