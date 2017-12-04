package model

//Token structure defines all properties that can be present in a token
type Token struct {
	Key             string     `json:"key"`
	User            string     `json:"user"`
	UserIdentity    Identity   `json:"userIdentity"`
	GroupIdentities []Identity `json:"groupIdentities"`
	Authprovider    string     `json:"authProvider"`
}
