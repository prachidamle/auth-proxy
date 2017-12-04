package server

import (
	"fmt"
	"strings"

	authnv1 "github.com/rancher/types/apis/authentication.cattle.io/v1"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/rancher/auth-proxy/model"
	"github.com/rancher/auth-proxy/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	Client *authnv1.Client
)

//SetEnv sets the parameters necessary
func SetEnv(c *cli.Context) {

	newclient, err := NewTokenClient(c.String("cluster-config"))
	if err != nil {
		log.Fatal("Failed to create token client: %v", err)
	}
	Client = newclient
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

	nclient, err := authnv1.NewForConfig(*kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to build cluster config: %v", err)
	}

	return nclient.(*authnv1.Client), nil
}

//CreateToken will authenticate with provider and create a jwt token
func CreateToken(json map[string]string) (model.Token, int, error) {

	log.Info("Create Token Invoked")

	authenticated := true

	/* Authenticate User
		if provider != nil {
		token, status, err := provider.GenerateToken(json)
		if err != nil {
			return model.Token{}, status, err
		}
	}*/

	if authenticated {

		if Client != nil {

			key, err := util.GenerateKey()
			if err != nil {
				log.Info("Failed to generate token key: %v", err)
				return model.Token{}, 0, fmt.Errorf("Failed to generate token key")
			}

			k8sToken := &authnv1.Token{
				TokenID:      strings.ToLower(key),
				TokenValue:   "dummy-token-val", //signed jwt containing user details
				User:         "dummy",
				IsDerived:    false,
				AuthProvider: "ActiveDirectory",
			}
			k8sToken.APIVersion = "authentication.cattle.io/v1"
			k8sToken.Kind = "Token"
			k8sToken.ObjectMeta = metav1.ObjectMeta{
				Name: strings.ToLower(key),
			}

			createdToken, err := Client.Tokens("").Create(k8sToken)

			if err != nil {
				log.Info("Failed to create token resource: %v", err)
			} else {
				log.Info("Created Token %v", createdToken)

				rToken := model.Token{
					Key:             createdToken.TokenID,
					User:            createdToken.User,
					UserIdentity:    getUserIdentity(),
					GroupIdentities: getGroupIdentities(),
					Authprovider:    createdToken.AuthProvider,
				}

				return rToken, 0, nil
			}
		} else {
			log.Info("Client nil %v", Client)
			return model.Token{}, 500, fmt.Errorf("No k8s Client configured")
		}
	}

	return model.Token{}, 0, fmt.Errorf("No auth provider configured")
}

func GetToken(tokenKey string) (model.Token, int, error) {
	log.Info("GET Token Invoked")

	if Client != nil {
		storedToken, err := Client.Tokens("").Get(tokenKey, metav1.GetOptions{})

		if err != nil {
			log.Info("Failed to get token resource: %v", err)
			return model.Token{}, 500, fmt.Errorf("Failed to get token")
		} else {
			rToken := model.Token{
				Key:             storedToken.TokenID,
				User:            storedToken.User,
				UserIdentity:    getUserIdentity(),
				GroupIdentities: getGroupIdentities(),
				Authprovider:    storedToken.AuthProvider,
			}
			log.Info("Got Token %v", storedToken)
			return rToken, 0, nil
		}
	} else {
		log.Info("Client nil %v", Client)
		return model.Token{}, 500, fmt.Errorf("No k8s Client configured")
	}

}

func DeleteToken(tokenKey string) (int, error) {
	log.Info("DELETE Token Invoked")

	if Client != nil {
		err := Client.Tokens("").Delete(tokenKey, &metav1.DeleteOptions{})

		if err != nil {
			return 500, fmt.Errorf("Failed to delete token")
		} else {
			log.Info("Deleted Token")
			return 0, nil
		}
	} else {
		log.Info("Client nil %v", Client)
		return 500, fmt.Errorf("No k8s Client configured")
	}

}

func GetIdentities(tokenKey string) ([]model.Identity, int, error) {
	var identities []model.Identity

	/*token, status, err := GetToken(tokenKey)

	if err != nil {
		return identities, 401, err
	} else {
		identities = append(identities, token.UserIdentity)
		identities = append(identities, token.GroupIdentities...)

		return identities, status, nil
	}*/

	identities = append(identities, getUserIdentity())
	identities = append(identities, getGroupIdentities()...)

	return identities, 0, nil

}

func getUserIdentity() model.Identity {

	identity := model.Identity{
		ExternalId:     "ldap://cn=dummy,dc=tad,dc=rancher,dc=io",
		Name:           "dummy",
		DisplayName:    "Dummy User",
		ProfilePicture: "",
		ProfileUrl:     "",
		Kind:           "user",
		Me:             true,
		MemberOf:       true,
	}

	return identity
}

func getGroupIdentities() []model.Identity {

	var identities []model.Identity

	identity1 := model.Identity{
		ExternalId:     "ldap://cn=group1,dc=tad,dc=rancher,dc=io",
		DisplayName:    "Admin group",
		Name:           "Administrators",
		ProfilePicture: "",
		ProfileUrl:     "",
		Kind:           "group",
		Me:             false,
		MemberOf:       false,
	}

	identity2 := model.Identity{
		ExternalId:     "ldap://cn=group2,dc=tad,dc=rancher,dc=io",
		DisplayName:    "Dev group",
		Name:           "Developers",
		ProfilePicture: "",
		ProfileUrl:     "",
		Kind:           "group",
		Me:             false,
		MemberOf:       false,
	}

	identities = append(identities, identity1)
	identities = append(identities, identity2)

	return identities
}
