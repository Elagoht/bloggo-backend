package models

import cryptography "bloggo/internal/utils/cryptography"

// -- Create New User -- //
type UserCreateParams struct {
	Name           string
	Email          string
	Avatar         string
	PassphraseHash string
	RoleId         int64
}

func (model *RequestUserCreate) HashUserPassphrase() (*UserCreateParams, error) {
	hashedPassphrase, err := cryptography.HashPassphrase(model.Passphrase)
	if err != nil {
		return nil, err
	}

	return &UserCreateParams{
		Name:           model.Name,
		Email:          model.Email,
		Avatar:         model.Avatar,
		PassphraseHash: hashedPassphrase,
		RoleId:         model.RoleId,
	}, nil
}
