package clerk

import "encoding/json"

type Invitation struct {
	APIResource
	Object         string          `json:"object"`
	ID             string          `json:"id"`
	EmailAddress   string          `json:"email_address"`
	PublicMetadata json.RawMessage `json:"public_metadata"`
	Revoked        bool            `json:"revoked,omitempty"`
	Status         string          `json:"status"`
	URL            string          `json:"url,omitempty"`
	ExpiresAt      *int64          `json:"expires_at"`
	CreatedAt      int64           `json:"created_at"`
	UpdatedAt      int64           `json:"updated_at"`
}

type Invitations struct {
	APIResource
	Invitations []*Invitation `json:"data"`
}

type InvitationList struct {
	APIResource
	Invitations []*Invitation `json:"data"`
	TotalCount  int64         `json:"total_count"`
}
