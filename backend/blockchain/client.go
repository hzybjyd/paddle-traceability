package blockchain

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"paddle-traceability/config"

	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuper-sdk-go/v2/xuper"
)

// ChainClient XuperChain integration client
type ChainClient struct {
	cfg     *config.BlockchainConfig
	xclient *xuper.XClient
	account *account.Account
}

type ChainResult struct {
	TxHash      string `json:"tx_hash"`
	BlockHeight int64  `json:"block_height"`
}

func NewChainClient(cfg *config.BlockchainConfig) (*ChainClient, error) {
	keyDir := filepath.Dir(cfg.PrivateKeyPath)
	if !strings.HasSuffix(keyDir, string(filepath.Separator)) {
		keyDir += string(filepath.Separator)
	}
	acc, err := account.GetAccountFromFile(keyDir, cfg.PrivateKeyPassword)
	if err != nil {
		return nil, fmt.Errorf("load xuper account from %s failed: %w", cfg.PrivateKeyPath, err)
	}

	if err := acc.SetContractAccount(cfg.ContractAccount); err != nil {
		return nil, fmt.Errorf("set contract account %s failed: %w", cfg.ContractAccount, err)
	}

	xclient, err := xuper.New(cfg.NodeAddr)
	if err != nil {
		return nil, fmt.Errorf("create xuper client for node %s failed: %w", cfg.NodeAddr, err)
	}

	log.Printf("blockchain client initialized, node: %s, contract: %s, contract account: %s", cfg.NodeAddr, cfg.ContractName, cfg.ContractAccount)
	return &ChainClient{
		cfg:     cfg,
		xclient: xclient,
		account: acc,
	}, nil
}

// CreateGoods create goods record on chain
func (c *ChainClient) CreateGoods(productUID, descJSON string) (*ChainResult, error) {
	tx, err := c.xclient.InvokeWasmContract(
		c.account,
		c.cfg.ContractName,
		"createGoods",
		map[string]string{"id": productUID, "desc": descJSON},
	)
	if err != nil {
		return nil, fmt.Errorf("invoke createGoods failed: %w", err)
	}

	return c.toChainResult(tx), nil
}

// UpdateGoods update goods record on chain
func (c *ChainClient) UpdateGoods(productUID, reasonJSON string) (*ChainResult, error) {
	tx, err := c.xclient.InvokeWasmContract(
		c.account,
		c.cfg.ContractName,
		"updateGoods",
		map[string]string{"id": productUID, "reason": reasonJSON},
	)
	if err != nil {
		return nil, fmt.Errorf("invoke updateGoods failed: %w", err)
	}

	return c.toChainResult(tx), nil
}

// QueryRecords query goods records from chain
func (c *ChainClient) QueryRecords(productUID string) (string, error) {
	tx, err := c.xclient.QueryWasmContract(
		c.account,
		c.cfg.ContractName,
		"queryRecords",
		map[string]string{"id": productUID},
	)
	if err != nil {
		return "", fmt.Errorf("query queryRecords failed: %w", err)
	}

	if tx == nil || tx.ContractResponse == nil {
		return "", nil
	}

	return string(tx.ContractResponse.Body), nil
}

// InvokeContract invoke smart contract (write), kept for service layer compatibility
func (c *ChainClient) InvokeContract(method string, args map[string]string) (*ChainResult, error) {
	log.Printf("[ChainClient] InvokeContract - method: %s, args: %v", method, args)

	switch method {
	case "CreateProduct":
		desc := map[string]string{
			"product_uid": args["product_uid"],
			"data_hash":   args["data_hash"],
		}
		descJSON, err := json.Marshal(desc)
		if err != nil {
			return nil, fmt.Errorf("marshal createGoods desc failed: %w", err)
		}
		return c.CreateGoods(args["product_uid"], string(descJSON))
	case "TransferProduct", "ConfirmSale":
		reason := map[string]string{
			"product_uid": args["product_uid"],
			"new_status":  args["new_status"],
			"data_hash":   args["data_hash"],
		}
		reasonJSON, err := json.Marshal(reason)
		if err != nil {
			return nil, fmt.Errorf("marshal updateGoods reason failed: %w", err)
		}
		return c.UpdateGoods(args["product_uid"], string(reasonJSON))
	default:
		return nil, fmt.Errorf("unsupported invoke method: %s", method)
	}
}

// QueryContract query smart contract (read), kept for service layer compatibility
func (c *ChainClient) QueryContract(method string, args map[string]string) (map[string]interface{}, error) {
	log.Printf("[ChainClient] QueryContract - method: %s, args: %v", method, args)

	if method != "queryRecords" {
		return nil, fmt.Errorf("unsupported query method: %s", method)
	}

	resp, err := c.QueryRecords(args["id"])
	if err != nil {
		return nil, err
	}

	if strings.Contains(resp, "the id not exist") {
		return map[string]interface{}{"found": false}, nil
	}

	return map[string]interface{}{
		"found":   true,
		"records": resp,
	}, nil
}

func (c *ChainClient) toChainResult(tx *xuper.Transaction) *ChainResult {
	result := &ChainResult{
		TxHash:      "",
		BlockHeight: 0,
	}

	if tx != nil && tx.Tx != nil && len(tx.Tx.Txid) > 0 {
		result.TxHash = hex.EncodeToString(tx.Tx.Txid)
	}

	return result
}
