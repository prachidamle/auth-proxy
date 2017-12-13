package api

import (
	"fmt"
	"strings"

	authnv1 "github.com/rancher/types/apis/management.cattle.io/v3"
	log "github.com/sirupsen/logrus"

	"github.com/rancher/auth-proxy/util"
	"github.com/rancher/types/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultSecret          = "secret"
	defaultTokenTTL        = "57600000"
	defaultRefreshTokenTTL = "7200000"
)

type apiServer struct {
	client *config.ClusterContext
}

//NewAPIServer sets the parameters necessary
func newAPIServer(clusterManagerCfg string, clusterCfg string, clusterName string) (*apiServer, error) {
	newclient, err := setupClient(clusterManagerCfg, clusterCfg, clusterName)
	if err != nil {
		return nil, fmt.Errorf("Failed to create a k8s client for  API server: %v", err)
	}

	apiServer := &apiServer{
		client: newclient,
	}

	return apiServer, nil
}

func getKubeConfig(clusterCfg string) (*rest.Config, error) {
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
	return kubeConfig, nil
}

func setupClient(clusterManagerCfg string, clusterCfg string, clusterName string) (*config.ClusterContext, error) {
	clusterManagementKubeConfig, err := clientcmd.BuildConfigFromFlags("", clusterManagerCfg)
	if err != nil {
		return nil, err
	}

	clusterKubeConfig, err := clientcmd.BuildConfigFromFlags("", clusterCfg)
	if err != nil {
		return nil, err
	}

	workload, err := config.NewClusterContext(*clusterManagementKubeConfig, *clusterKubeConfig, clusterName)
	if err != nil {
		return nil, err
	}

	return workload, nil

}

func NewTokenClient(clusterCfg string) (*authnv1.Client, error) {
	// build kubernetes config
	kubeConfig, err := getKubeConfig(clusterCfg)
	if err != nil {
		return nil, err
	}
	nclient, err := authnv1.NewForConfig(*kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to build cluster config: %v", err)
	}

	return nclient.(*authnv1.Client), nil
}

//CreateLoginToken will authenticate with provider and create a jwt token
func (s *apiServer) createLoginToken(jsonInput authnv1.LoginInput) (authnv1.Token, int, error) {

	log.Info("Create Token Invoked %v", jsonInput)
	authenticated := true

	/* Authenticate User
		if provider != nil {
		token, status, err := provider.GenerateToken(json)
		if err != nil {
			return model.TokenOutput{}, status, err
		}
	}*/

	if authenticated {

		if s.client != nil {

			key, err := util.GenerateKey()
			if err != nil {
				log.Info("Failed to generate token key: %v", err)
				return authnv1.Token{}, 0, fmt.Errorf("Failed to generate token key")
			}

			//check that there is no token with this key.
			payload := make(map[string]interface{})
			tokenValue, err := util.CreateTokenWithPayload(payload, defaultSecret)
			if err != nil {
				log.Info("Failed to generate token value: %v", err)
				return authnv1.Token{}, 0, fmt.Errorf("Failed to generate token value")
			}

			ttl := jsonInput.TTLMillis
			refreshTTL := jsonInput.IdentityRefreshTTLMillis
			if ttl == "" {
				ttl = defaultTokenTTL               //16 hrs
				refreshTTL = defaultRefreshTokenTTL //2 hrs
			}

			k8sToken := &authnv1.Token{
				TokenID:                  key,
				TokenValue:               tokenValue, //signed jwt containing user details
				IsDerived:                false,
				TTLMillis:                ttl,
				IdentityRefreshTTLMillis: refreshTTL,
				User:         "dummy",
				ExternalID:   "github_12346",
				AuthProvider: "github",
			}
			rToken, err := s.createK8sTokenCR(k8sToken)
			return rToken, 0, err
		}
		log.Info("Client nil %v", s.client)
		return authnv1.Token{}, 500, fmt.Errorf("No k8s Client configured")
	}

	return authnv1.Token{}, 0, fmt.Errorf("No auth provider configured")
}

//CreateDerivedToken will create a jwt token for the authenticated user
func (s *apiServer) createDerivedToken(jsonInput authnv1.Token, tokenID string) (authnv1.Token, int, error) {

	log.Info("Create Derived Token Invoked")

	token, err := s.getK8sTokenCR(tokenID)

	if err != nil {
		return authnv1.Token{}, 401, err
	}

	if s.client != nil {
		key, err := util.GenerateKey()
		if err != nil {
			log.Info("Failed to generate token key: %v", err)
			return authnv1.Token{}, 0, fmt.Errorf("Failed to generate token key")
		}

		ttl := jsonInput.TTLMillis
		refreshTTL := jsonInput.IdentityRefreshTTLMillis
		if ttl == "" {
			ttl = defaultTokenTTL               //16 hrs
			refreshTTL = defaultRefreshTokenTTL //2 hrs
		}

		k8sToken := &authnv1.Token{
			TokenID:                  key,
			TokenValue:               token.TokenValue, //signed jwt containing user details
			IsDerived:                true,
			TTLMillis:                ttl,
			IdentityRefreshTTLMillis: refreshTTL,
			User:         token.User,
			ExternalID:   token.ExternalID,
			AuthProvider: token.AuthProvider,
		}
		rToken, err := s.createK8sTokenCR(k8sToken)
		return rToken, 0, err

	}
	log.Info("Client nil %v", s.client)
	return authnv1.Token{}, 500, fmt.Errorf("No k8s Client configured")

}

func (s *apiServer) createK8sTokenCR(k8sToken *authnv1.Token) (authnv1.Token, error) {
	if s.client != nil {

		labels := make(map[string]string)
		labels["io.cattle.token.field.externalID"] = k8sToken.ExternalID

		k8sToken.APIVersion = "management.cattle.io/v3"
		k8sToken.Kind = "Token"
		k8sToken.ObjectMeta = metav1.ObjectMeta{
			Name:   strings.ToLower(k8sToken.TokenID),
			Labels: labels,
		}
		createdToken, err := s.client.Management.Management.Tokens("").Create(k8sToken)

		if err != nil {
			log.Info("Failed to create token resource: %v", err)
			return authnv1.Token{}, err
		}
		log.Info("Created Token %v", createdToken)
		return *createdToken, nil
	}

	return authnv1.Token{}, fmt.Errorf("No k8s Client configured")
}

func (s *apiServer) getK8sTokenCR(tokenID string) (*authnv1.Token, error) {
	if s.client != nil {
		storedToken, err := s.client.Management.Management.Tokens("").Get(strings.ToLower(tokenID), metav1.GetOptions{})

		if err != nil {
			log.Info("Failed to get token resource: %v", err)
			return nil, fmt.Errorf("Failed to retrieve auth token")
		}

		log.Info("storedToken token resource: %v", storedToken)

		return storedToken, nil
	}
	return nil, fmt.Errorf("No k8s Client configured")
}

//GetTokens will list all tokens of the authenticated user - login and derived
func (s *apiServer) getTokens(tokenID string) ([]authnv1.Token, int, error) {
	log.Info("GET Token Invoked")
	var tokens []authnv1.Token

	if s.client != nil {
		storedToken, err := s.client.Management.Management.Tokens("").Get(strings.ToLower(tokenID), metav1.GetOptions{})

		if err != nil {
			log.Info("Failed to get token resource: %v", err)
			return tokens, 401, fmt.Errorf("Failed to retrieve auth token")
		}
		log.Info("storedToken token resource: %v", storedToken)
		externalID := storedToken.ExternalID
		set := labels.Set(map[string]string{"io.cattle.token.field.externalID": externalID})
		tokenList, err := s.client.Management.Management.Tokens("").List(metav1.ListOptions{LabelSelector: set.AsSelector().String()})
		if err != nil {
			return tokens, 0, fmt.Errorf("Error getting tokens for user: %v selector: %v  err: %v", externalID, set.AsSelector().String(), err)
		}

		for _, t := range tokenList.Items {
			log.Info("List token resource: %v", t)
			tokens = append(tokens, t)
		}
		return tokens, 0, nil

	}
	log.Info("Client nil %v", s.client)
	return tokens, 500, fmt.Errorf("No k8s Client configured")
}

func (s *apiServer) deleteToken(tokenKey string) (int, error) {
	log.Info("DELETE Token Invoked")

	if s.client != nil {
		err := s.client.Management.Management.Tokens("").Delete(strings.ToLower(tokenKey), &metav1.DeleteOptions{})

		if err != nil {
			return 500, fmt.Errorf("Failed to delete token")
		}
		log.Info("Deleted Token")
		return 0, nil

	}
	log.Info("Client nil %v", s.client)
	return 500, fmt.Errorf("No k8s Client configured")
}

func (s *apiServer) getIdentities(tokenKey string) ([]authnv1.Identity, int, error) {
	var identities []authnv1.Identity

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

func getUserIdentity() authnv1.Identity {

	identity := authnv1.Identity{
		LoginName:      "dummy",
		DisplayName:    "Dummy User",
		ProfilePicture: "",
		ProfileURL:     "",
		Kind:           "user",
		Me:             true,
		MemberOf:       false,
	}
	identity.ObjectMeta = metav1.ObjectMeta{
		Name: "ldap://cn=dummy,dc=tad,dc=rancher,dc=io",
	}

	return identity
}

func getGroupIdentities() []authnv1.Identity {

	var identities []authnv1.Identity

	identity1 := authnv1.Identity{
		DisplayName:    "Admin group",
		LoginName:      "Administrators",
		ProfilePicture: "",
		ProfileURL:     "",
		Kind:           "group",
		Me:             false,
		MemberOf:       true,
	}
	identity1.ObjectMeta = metav1.ObjectMeta{
		Name: "ldap://cn=group1,dc=tad,dc=rancher,dc=io",
	}

	identity2 := authnv1.Identity{
		DisplayName:    "Dev group",
		LoginName:      "Developers",
		ProfilePicture: "",
		ProfileURL:     "",
		Kind:           "group",
		Me:             false,
		MemberOf:       true,
	}
	identity2.ObjectMeta = metav1.ObjectMeta{
		Name: "ldap://cn=group2,dc=tad,dc=rancher,dc=io",
	}

	identities = append(identities, identity1)
	identities = append(identities, identity2)

	return identities
}
