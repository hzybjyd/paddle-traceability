// verify_deploy.go
// 部署后验证脚本：通过 xuper-sdk-go/v2 调用 hzy_trace.queryRecords。
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuper-sdk-go/v2/xuper"
)

func main() {
	var (
		productID       = flag.String("id", "test_product_12345", "要查询的商品 ID")
		nodeAddr        = flag.String("node", "39.156.69.83:37100", "XuperChain 开放网络节点地址")
		contractName    = flag.String("name", "hzy_trace", "合约名称")
		contractAccount = flag.String("account", "XC4103761871843472@xuper", "合约账户")
		privateKeyPath  = flag.String("key", `d:\Study\1\private.key`, "私钥文件路径")
	)
	flag.Parse()

	password := os.Getenv("XUPER_KEY_PASSWORD")
	if password == "" {
		log.Fatal("请设置环境变量 XUPER_KEY_PASSWORD")
	}

	acc, err := loadAccount(*privateKeyPath, password, *contractAccount)
	if err != nil {
		log.Fatalf("加载账户失败: %v", err)
	}

	xclient, err := xuper.New(*nodeAddr)
	if err != nil {
		log.Fatalf("创建 XuperChain 客户端失败: %v", err)
	}

	args := map[string]string{"id": *productID}
	tx, err := xclient.QueryWasmContract(acc, *contractName, "queryRecords", args)
	if err != nil {
		log.Fatalf("查询合约失败: %v", err)
	}

	fmt.Printf("查询商品 ID: %s\n", *productID)
	if tx != nil && tx.ContractResponse != nil {
		fmt.Printf("查询结果:\n%s\n", string(tx.ContractResponse.Body))
	} else {
		fmt.Println("查询结果为空")
	}
}

// loadAccount 从私钥文件加载账户并设置合约账户。
func loadAccount(keyPath, password, contractAccount string) (*account.Account, error) {
	keyDir := filepath.Dir(keyPath)
	if !strings.HasSuffix(keyDir, string(filepath.Separator)) {
		keyDir += string(filepath.Separator)
	}

	acc, err := account.GetAccountFromFile(keyDir, password)
	if err != nil {
		return nil, fmt.Errorf("从文件加载账户失败: %w", err)
	}

	if err := acc.SetContractAccount(contractAccount); err != nil {
		return nil, fmt.Errorf("设置合约账户失败: %w", err)
	}

	return acc, nil
}
