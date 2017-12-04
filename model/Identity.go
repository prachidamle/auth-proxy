package model

//Identity structure defines all properties that can be present in an identity
type Identity struct {
	ExternalId     string `json:"externalId,omitempty"`
	ExternalIdType string `json:"externalIdType,omitempty"` //ldap/github
	Login          string `json:"login,omitempty"`
	Name           string `json:"name,omitempty"`
	ProfilePicture string `json:"profilePicture,omitempty"`
	ProfileUrl     string `json:"profileUrl,omitempty"`
	IsUser         bool   `json:"isUser,omitempty"`
}
