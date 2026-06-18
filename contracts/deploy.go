// deploy.go
// 通过 xuper-sdk-go/v2 部署 hzy_trace WASM 合约并调用 initialize。
// 注意：当前仓库未包含 hzy_trace.wasm，本脚本默认不会自动运行，仅作为后续参考。
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
		wasmPath        = flag.String("wasm", "./hzy_trace.wasm", "编译后的 WASM 合约文件路径")
		admin           = flag.String("admin", "Vh7gwgrwwJvrTnYQorVAN2uAJZSmnjWEb", "合约管理员地址，请替换为控制台实际显示的私钥地址")
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

	code, err := os.ReadFile(*wasmPath)
	if err != nil {
		log.Fatalf("读取 WASM 文件失败 (%s): %v", *wasmPath, err)
	}

	acc, err := loadAccount(*privateKeyPath, password, *contractAccount)
	if err != nil {
		log.Fatalf("加载账户失败: %v", err)
	}

	xclient, err := xuper.New(*nodeAddr)
	if err != nil {
		log.Fatalf("创建 XuperChain 客户端失败: %v", err)
	}

	args := map[string]string{"admin": *admin}
	tx, err := xclient.DeployWasmContract(acc, *contractName, code, args)
	if err != nil {
		log.Fatalf("部署合约失败: %v", err)
	}

	fmt.Printf("合约部署成功\n")
	fmt.Printf("合约名称: %s\n", *contractName)
	fmt.Printf("合约账户: %s\n", *contractAccount)
	if tx != nil && tx.Tx != nil && len(tx.Tx.Txid) > 0 {
		fmt.Printf("交易哈希: %x\n", tx.Tx.Txid)
	}
	fmt.Printf("initialize 参数 admin: %s\n", *admin)
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
