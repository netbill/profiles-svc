package tokenmanager

import (
	"fmt"

	"github.com/netbill/restkit/tokens"
)

func (m *Manager) ParseAccessClaims(token string) (tokens.AccountClaims, error) {
	data, err := tokens.ParseAccountJWT(token, m.accessSK)
	if err != nil {
		return tokens.AccountClaims{}, fmt.Errorf("failed to parse access token, cause: %w", err)
	}

	return data, nil
}
