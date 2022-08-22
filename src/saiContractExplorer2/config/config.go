package config

import (
	valid "github.com/asaskevich/govalidator"
	configinternal "github.com/webmakom-com/saiBoilerplate/internal/config-internal"
)

type Configuration struct {
	Common   configinternal.Common `json:"common"` // built-in framework config
	Specific `json:"specific"`
}

// Specific - specific for current microservice settings
type Specific struct {
	GethServer string `json:"geth_server"`
	Storage    `json:"storage"`
	Contracts  []Contract `json:"contracts"`
	StartBlock int        `json:"start_block"`
	Operations []string   `json:"operations"`
	Sleep      int        `json:"sleep"`
}

// settings for saiStorage
type Storage struct {
	Token    string `json:"token"`
	URL      string `json:"url"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Contract struct {
	Address    string `json:"address" valid:",required"`
	ABI        string `json:"abi" valid:",required"`
	StartBlock int    `json:"start_block" valid:",required"`
}

func (r *Contract) Validate() error {
	_, err := valid.ValidateStruct(r)

	return err
}
