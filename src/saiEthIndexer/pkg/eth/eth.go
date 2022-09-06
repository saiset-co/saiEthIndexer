package eth

import (
	"github.com/onrik/ethrpc"
	"go.uber.org/zap"
)

func GetClient(address string, logger *zap.Logger) (*ethrpc.EthRPC, error) {
	ethClient := ethrpc.New(address)

	version, err := ethClient.Web3ClientVersion()
	if err != nil {
		return nil, err
	}

	logger.Info("Connected to geth server", zap.String("client version", version))

	return ethClient, nil
}
