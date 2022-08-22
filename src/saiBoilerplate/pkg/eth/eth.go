package eth

import (
	"github.com/onrik/ethrpc"
)

func GetClient(address string) (*ethrpc.EthRPC, error) {
	ethClient := ethrpc.New(address)

	_, err := ethClient.Web3ClientVersion()
	if err != nil {
		return nil, err
	}
	return ethClient, nil
}
