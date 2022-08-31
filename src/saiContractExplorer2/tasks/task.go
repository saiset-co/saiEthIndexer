package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/onrik/ethrpc"
	"github.com/webmakom-com/saiBoilerplate/config"
	"github.com/webmakom-com/saiBoilerplate/pkg/eth"
	"go.uber.org/zap"
)

const (
	configPath    = "./config/config.json"
	contractsPath = "./config/contracts.json"
)

type TaskManager struct {
	Config       *config.Configuration
	EthClient    *ethrpc.EthRPC
	Logger       *zap.Logger
	BlockManager *BlockManager
}

func NewManager(config *config.Configuration) (*TaskManager, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	// todo:handle error
	ethClient, _ := eth.GetClient(config.Specific.GethServer, logger)
	if err != nil {
		return nil, err
	}

	blockManager := NewBlockManager(*config, logger)

	return &TaskManager{
		Config:       config,
		EthClient:    ethClient,
		Logger:       logger,
		BlockManager: blockManager,
	}, nil
}

// Process blocks, which got from geth-server
func (t *TaskManager) ProcessBlocks() {

	for {
		blockID, err := t.EthClient.EthBlockNumber()
		if err != nil {
			t.Logger.Error("tasks - ProcessBlocks - get block number from eth client", zap.Error(err))
			continue
		}

		blk, err := t.BlockManager.GetLastBlock(blockID)
		if err != nil {
			t.Logger.Error("tasks - ProcessBlocks - get last block from block manager", zap.Error(err))
			continue
		}

		for i := blk.ID; i <= blockID; i++ {
			blkInfo, err := t.EthClient.EthGetBlockByNumber(i, true)
			if err != nil {
				t.Logger.Error("tasks - ProcessBlocks - get block by number from server", zap.Error(err))
				i--
				continue
			}

			if len(blkInfo.Transactions) == 0 {
				t.Logger.Info("tasks - ProcessBlocks - get block by number from server - transactions - no transactions found", zap.Int("current block id in for cycle", i), zap.Int("current block id from eth server", blockID))
				continue
			}

			t.Logger.Info("tasks - ProcessBlocks - transactions found", zap.Int("block id", i), zap.Int("transactions count", len(blkInfo.Transactions)))

			t.BlockManager.HandleTransactions(blkInfo.Transactions)
		}
		blk.ID = blockID
		t.BlockManager.SetLastBlock(blk)
		time.Sleep(time.Duration(t.Config.Specific.Sleep) * time.Second)

	}
}

func (t *TaskManager) AddContract(contracts []config.Contract) error {
	t.Config.EthContracts.Mutex.Lock()
	defer t.Config.EthContracts.Mutex.Unlock()
	t.Config.EthContracts.Contracts = append(t.Config.EthContracts.Contracts, contracts...)

	for _, contract := range t.Config.EthContracts.Contracts {
		if contract.StartBlock < t.Config.StartBlock {
			t.Config.StartBlock = contract.StartBlock
		}
	}

	data, err := json.MarshalIndent(&t.Config.EthContracts, "", "	")
	if err != nil {
		return err
	}

	f, err := os.OpenFile(contractsPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	err = t.ReloadContracts()
	if err != nil {
		return err
	}

	t.Logger.Sugar().Infof("active config : %+v\n", t.Config)

	return nil
}

// reload config after contracts was added via http add_contracts endpoint
func (t *TaskManager) ReloadContracts() error {
	contracts := config.EthContracts{}
	b, err := os.ReadFile(contractsPath)
	if err != nil {
		return fmt.Errorf("contracts json read error: %w", err)
	}
	err = json.Unmarshal(b, &contracts)
	if err != nil {
		return fmt.Errorf("contracts json unmarshal error: %w", err)
	}

	t.Config.EthContracts = contracts
	t.BlockManager.config.EthContracts = contracts
	return nil
}
