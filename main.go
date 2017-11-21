package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	authnv1 "github.com/rancher/types/apis/authentication.cattle.io/v1"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var VERSION = "v0.0.0-dev"

var client *authnv1.Client

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
	}

	app.Run(os.Args)
}

func StartService(c *cli.Context) {

	client, err := NewTokenClient(c.String("cluster-config"))
	if err != nil {
		log.Fatal("Failed to create token client: %v", err)
	}
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("POST").Path("/v1/token").Handler(http.HandlerFunc(CreateToken))

	log.Info("Listening on 9998", client)
	log.Fatal(http.ListenAndServe(":9998", router))
}

func NewTokenClient(clusterCfg string) (*authnv1.Client, error) {
	// build kubernetes config
	var kubeConfig *rest.Config
	var err error
	if clusterCfg != "" {
		log.Info("Using out of cluster config to connect to kubernetes cluster")
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", clusterCfg)
	} else {
		log.Info("Using in cluster config to connect to kubernetes cluster")
		kubeConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to build cluster config: %v", err)
	}

	/*if kubeConfig.NegotiatedSerializer == nil {
		configConfig := dynamic.ContentConfig()
		kubeConfig.NegotiatedSerializer = configConfig.NegotiatedSerializer
	}
	restClient, err := rest.UnversionedRESTClientFor(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to build cluster client: %v", err)
	}

	client = &authnv1.Client{
		restClient: restClient,
		tokenControllers: map[string]authnv1.TokenController{},
	}*/

	nclient, err := authnv1.NewForConfig(*kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to build cluster config: %v", err)
	}

	return nclient.(*authnv1.Client), nil
}

func CreateToken(w http.ResponseWriter, r *http.Request) {
	log.Info("Create Token Invoked")

	token := &authnv1.Token{
		TokenKey:   "dummy-token-key",
		TokenValue: "dummy-token-val",
		User:       "dummy",
		IsCLI:      false,
	}
	token.APIVersion = "authentication.cattle.io/v1"
	token.Kind = "Token"

	createdToken, err := client.Tokens("").Create(token)

	if err != nil {
		log.Info("Failed to create token resource: %v", err)
	} else {
		log.Info("Created Token %v", createdToken)
	}

}
