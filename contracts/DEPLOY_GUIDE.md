# hzy_trace 智能合约部署指南

> 本指南用于在百度超级链开放网络（XuperChain）上手动部署 `hzy_trace` C++ WASM 智能合约。
> 由于当前 Windows 环境无法本地编译 WASM，推荐通过[XuperChain 开放网络控制台](https://xchain.baidu.com/)完成在线编译与部署。

---

## 一、部署前准备

1. **登录控制台**
   - 打开 <https://xchain.baidu.com/> 并使用已有百度超级链账户登录。

2. **确认合约账户与余额**
   - 本任务使用的合约账户为：`XC4103761871843472@xuper`
   - 部署 WASM 合约需要消耗 gas（Xuper），请确保该合约账户余额充足；如余额不足，请先充值。

3. **确认私钥与密码**
   - 私钥文件路径：`d:\Study\1\private.key`
   - 私钥密码通过环境变量 `XUPER_KEY_PASSWORD` 注入，**请勿将密码写入代码或文档**。
   - 当前私钥对应的 admin 地址示例：`Vh7gwgrwwJvrTnYQorVAN2uAJZSmnjWEb`，请以控制台实际显示的私钥地址为准。

4. **确认合约源码**
   - 合约源码位于：`contracts/hzy_trace.cc`
   - 部署前请确认文件内容完整，包含 `initialize`、`createGoods`、`updateGoods`、`queryRecords` 四个方法。

---

## 二、控制台部署步骤

### 步骤 1：进入合约管理

登录控制台后，依次点击 **开放网络 -> 合约管理**。

### 步骤 2：创建智能合约

1. 点击 **创建智能合约**（如已有 `paddledledle_trace` 合约条目，可选择编辑后重新部署）。
2. 填写基本信息：
   - **合约名称**：`paddledledle_trace`
   - **合约账户**：`XC4103761871843472@xuper`
   - **合约语言**：`C++`

### 步骤 3：粘贴合约代码

1. 打开本地文件 `contracts/hzy_trace.cc`。
2. 复制全部内容。
3. 粘贴到控制台代码编辑器中。

### 步骤 4：编译合约

1. 点击 **编译**。
2. 等待控制台提示编译成功。如编译失败，请检查：
   - 是否包含 `xchain/xchain.h` 头文件；
   - 方法签名是否符合 XuperChain C++ WASM 规范；
   - 是否正确定义了 `DEFINE_METHOD`。

### 步骤 5：安装（部署）合约

1. 编译成功后，点击 **安装**。
2. 在预执行/初始化页面，调用 `initialize` 方法，参数如下：
   - `admin`：你的私钥对应地址（例如 `Vh7gwgrwwJvrTnYQorVAN2uAJZSmnjWEb`，以控制台显示为准）
3. 确认安装，等待合约状态变为 **安装成功**。

---

## 三、部署后验证

### 3.1 控制台验证

1. 在控制台 **合约管理** 中找到 `hzy_trace`，点击 **调用**。
2. 调用 `createGoods` 方法，参数：
   - `id`：`test_product_12345`
   - `desc`：`{"product_uid":"test_product_12345","data_hash":"abc123","operator_role":"FACTORY"}`
3. 调用 `queryRecords` 方法，参数：
   - `id`：`test_product_12345`
4. 如果能查询到类似如下记录，说明合约部署成功：

   ```text
   goodsId=test_product_12345,updateRecord=0,reason=CREATE
   ```

5. 再次调用 `queryRecords`，参数：
   - `id`：`nonexistent_product_99999`
6. 应返回 `the id not exist` 错误，说明查询逻辑正常。

### 3.2 SDK 验证脚本

项目同时提供了 Go 语言验证脚本，部署后可在命令行执行：

```powershell
# 切换到 contracts 目录
cd d:\Study\1\paddle-traceability\contracts

# 设置环境变量（PowerShell 示例）
$env:XUPER_KEY_PASSWORD = "你的私钥密码"

# 查询默认测试商品
go run verify_deploy.go

# 查询指定商品 ID
go run verify_deploy.go -id test_product_12345
```

---

## 四、SDK 部署脚本（可选）

如果你已经拿到了编译后的 `hzy_trace.wasm` 文件，可以使用 `contracts/deploy.go` 通过 SDK 部署：

```powershell
cd d:\Study\1\paddle-traceability\contracts

# 设置环境变量
$env:XUPER_KEY_PASSWORD = "你的私钥密码"

# 默认部署同目录下的 hzy_trace.wasm，并调用 initialize
go run deploy.go

# 或指定 wasm 文件路径与 admin 地址
go run deploy.go -wasm ./hzy_trace.wasm -admin Vh7gwgrwwJvrTnYQorVAN2uAJZSmnjWEb
```

> 注意：当前仓库未包含 `hzy_trace.wasm` 文件，因此 `deploy.go` 默认不会自动运行，仅作为后续参考。

---

## 五、后端真实上链验证

合约部署并验证成功后，建议通过后端 API 进行一次完整的上链测试：

1. 启动后端服务（确保已设置 `XUPER_KEY_PASSWORD`、`MYSQL_PASSWORD` 等环境变量）。
2. 登录厂商账号，调用：

   ```http
   POST /api/v1/products
   Content-Type: application/json
   Authorization: Bearer <你的 JWT Token>

   {
     "name": "测试球拍",
     "model": "TEST-001",
     "serial_number": "SN20260618001",
     "factory_info": "测试厂商"
   }
   ```

3. 记录返回的 `product_uid`。
4. 调用验证接口：

   ```http
   GET /api/v1/verify/{product_uid}
   ```

5. 若返回的 `blockchain_records` 中包含链上记录，则证明后端与链上合约交互正常。

---

## 六、安全提醒

- **严禁**将私钥密码、数据库密码等敏感信息写入代码、文档或 Git 仓库。
- 私钥文件 `private.key` 请妥善保管，建议仅存储在本地安全目录。
- 生产环境部署前，请确认合约账户 `XC4103761871843472@xuper` 的私钥权限配置正确。
