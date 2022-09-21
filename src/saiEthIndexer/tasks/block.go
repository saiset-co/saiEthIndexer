package tasks

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/onrik/ethrpc"
	"github.com/saiset-co/saiEthIndexer/config"
	"github.com/saiset-co/saiEthIndexer/utils/saiStorageUtil"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type BlockManager struct {
	config  *config.Configuration
	storage saiStorageUtil.Database
	logger  *zap.Logger
}

type Block struct {
	ID int `json:"id"`
}

func NewBlockManager(c config.Configuration, logger *zap.Logger) *BlockManager {

	manager := &BlockManager{
		config:  &c,
		storage: saiStorageUtil.Storage(c.Specific.Storage.URL, c.Storage.Email, c.Storage.Password),
		logger:  logger,
	}

	return manager
}

func (bm *BlockManager) GetLastBlock(id int) (*Block, error) {
	var blocks []Block
	err, resultJsonString := bm.storage.Get("last_block", bson.M{}, bson.M{}, bm.config.Storage.Token)

	if err != nil {
		bm.logger.Error("tasks - BlockManager - get last block - get last_block from storage", zap.Error(err))
		return &Block{
			ID: id,
		}, nil
	}

	err = json.Unmarshal(resultJsonString, &blocks)

	if err != nil {
		bm.logger.Error("tasks - BlockManager - get last block - unmarshal result of get last_block", zap.Error(err))
		return &Block{
			ID: id,
		}, nil
	}

	var startBlock int
	if len(blocks) > 0 {
		startBlock = blocks[0].ID + 1
	} else if bm.config.StartBlock > 0 {
		startBlock = bm.config.StartBlock
	} else {
		startBlock = id
	}

	return &Block{
		ID: startBlock,
	}, nil
}

func (bm *BlockManager) SetLastBlock(blk *Block) error {
	var blocks []Block

	err, resultJsonString := bm.storage.Get("last_block", bson.M{}, bson.M{}, bm.config.Storage.Token)
	if err != nil {
		bm.logger.Error("tasks - BlockManager - set last block - get last_block from storage", zap.Error(err))
		return err
	}
	err = json.Unmarshal(resultJsonString, &blocks)
	if err != nil {
		bm.logger.Error("tasks - BlockManager - set last block - unmarshal result of get last_block", zap.Error(err))
		return err
	}

	if len(blocks) > 0 {
		err, _ = bm.storage.Update("last_block", bson.M{"id": bson.M{"$exists": true}}, blk, bm.config.Storage.Token)
		if err != nil {
			bm.logger.Error("tasks - BlockManager - set last block - update storage last_block", zap.Error(err))
			return err
		}
	} else {
		err, _ = bm.storage.Put("last_block", blk, bm.config.Storage.Token)
		if err != nil {
			bm.logger.Error("tasks - BlockManager - set last block - unmarshal result of get last_block", zap.Error(err))
			return err
		}
	}
	bm.logger.Sugar().Debugf("block %d was saved to storage", blk.ID)
	return nil
}

func (bm *BlockManager) HandleTransactions(trs []ethrpc.Transaction) {
	for j := 0; j < len(trs); j++ {
		for i := 0; i < len(bm.config.EthContracts.Contracts); i++ {
			if strings.ToLower(trs[j].From) != strings.ToLower(bm.config.EthContracts.Contracts[i].Address) && strings.ToLower(trs[j].To) != strings.ToLower(bm.config.EthContracts.Contracts[i].Address) {
				continue
			}

			raw, err := json.Marshal(trs[j])
			if err != nil {
				bm.logger.Error("block manager - handle transaction - marshal transaction", zap.String("transaction hash", trs[j].Hash), zap.Error(err))
				continue
			}

			data := bson.M{
				"Hash":   trs[j].Hash,
				"From":   trs[j].From,
				"To":     trs[j].To,
				"Amount": trs[j].Value,
			}

			decodedSig, err := hex.DecodeString(trs[j].Input[2:10])

			if err != nil {
				bm.logger.Error("block manager - handle transaction - decode transaction function idintifier", zap.String("transaction hash", trs[j].Hash), zap.Error(err))
				continue
			}

			abi, err := abi.JSON(strings.NewReader(bm.config.EthContracts.Contracts[i].ABI))
			if err != nil {
				bm.logger.Error("block manager - handle transaction - parse abi from config", zap.String("address", bm.config.EthContracts.Contracts[i].Address), zap.Error(err))
				continue
			}

			method, err := abi.MethodById(decodedSig)
			if err != nil {
				bm.logger.Error("block manager - handle transaction - MethodById", zap.String("transaction hash", trs[j].Hash), zap.Error(err))
				continue
			}

			decodedData, err := hex.DecodeString(trs[j].Input[2:])

			if err != nil {
				bm.logger.Error("block manager - handle transaction - decode input", zap.String("transaction hash", trs[j].Hash), zap.Error(err))
				continue
			}

			decodedInput := map[string]interface{}{}
			err = method.Inputs.UnpackIntoMap(decodedInput, decodedData[4:])

			if err != nil {
				bm.logger.Error("block manager - handle transaction - UnpackIntoMap", zap.String("transaction hash", trs[j].Hash), zap.Error(err))
				continue
			}

			data["Operation"] = method.Name
			data["Input"] = decodedInput

			for _, operation := range bm.config.Operations {
				if operation == method.Name {
					err = bm.SendWebSocketMsg(raw, method.Name)
					if err != nil {
						bm.logger.Error("block manager - handle transaction - SendWebSocketMsg", zap.String("transaction hash", trs[j].Hash), zap.Error(err))
						continue
					}
				}
			}
			err, _ = bm.storage.Put("transactions", data, bm.config.Storage.Token)

			if err != nil {
				bm.logger.Error("block manager - handle transaction - storage.Put", zap.String("transaction hash", trs[j].Hash), zap.Error(err))
				continue
			}

			bm.logger.Sugar().Infof("%d transaction from %s to %s has been updated.\n", trs[j].TransactionIndex, trs[j].From, trs[j].To)
		}
	}
}
