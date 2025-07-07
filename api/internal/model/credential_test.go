package model

import (
	"encoding/json"
	"testing"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/stretchr/testify/assert"
)

func TestCredential_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "valid data",
			data:    validCredentialData(),
			wantErr: false,
		},
		{
			name:    "nil data",
			data:    nil,
			wantErr: false,
		},
		{
			name:    "invalid data",
			data:    []byte("invalid data"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Credential{
				Data: tt.data,
			}
			err := c.Unmarshal()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.data != nil {
					assert.NotNil(t, c.Cred)
					assert.Equal(t, c.CredID, c.Cred.ID)
				}
			}
		})
	}
}

func validCredentialData() []byte {
	cred := webauthn.Credential{
		ID: []byte("test-cred-id"),
	}
	data, _ := json.Marshal(cred)
	return data
}
