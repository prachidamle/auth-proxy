package model

//Identity structure defines all properties that can be present in an identity
type Identity struct {
	ExternalId     string `json:"externalId,omitempty"`
	DisplayName    string `json:"displayName,omitempty"`
	Name           string `json:"name,omitempty"`
	ProfilePicture string `json:"profilePicture,omitempty"`
	ProfileUrl     string `json:"profileUrl,omitempty"`
	Kind           string `json:"kind,omitempty"`
	Me             bool   `json:"me,omitempty"`
	MemberOf       bool   `json:"me,omitempty"`
}
