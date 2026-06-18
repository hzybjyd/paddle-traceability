# 基于区块链的乒乓球拍防伪溯源系统

基于百度超级链开放网络（XuperChain）的乒乓球拍全生命周期防伪溯源系统。通过将产品关键信息上链存证，实现消费者输入产品 ID 即可验证真伪、追溯流通链路。

> 本项目为核心功能为 **上链存证 + 真伪验证**，其余功能已按简化设计处理。

---

## 目录

- [基于区块链的乒乓球拍防伪溯源系统](#基于区块链的乒乓球拍防伪溯源系统)
  - [目录](#目录)
  - [功能特性](#功能特性)
  - [系统架构](#系统架构)
  - [技术栈](#技术栈)
  - [目录结构](#目录结构)
  - [环境准备](#环境准备)
    - [必需环境](#必需环境)
    - [准备私钥](#准备私钥)
  - [安装与部署](#安装与部署)
    - [前置要求](#前置要求)
    - [方式一：本地开发环境（推荐）](#方式一本地开发环境推荐)
      - [1. 安装 MySQL 8.0](#1-安装-mysql-80)
      - [2. 初始化数据库](#2-初始化数据库)
      - [3. 配置环境变量](#3-配置环境变量)
      - [4. 部署智能合约](#4-部署智能合约)
      - [5. 启动后端](#5-启动后端)
      - [6. 启动前端](#6-启动前端)
  - [环境变量说明](#环境变量说明)
  - [智能合约部署](#智能合约部署)
    - [部署参数](#部署参数)
    - [部署步骤](#部署步骤)
    - [部署后验证](#部署后验证)
  - [部署后验证](#部署后验证-1)
  - [调试与排查](#调试与排查)
    - [后端无法启动](#后端无法启动)
    - [合约调用失败 / gas 不足](#合约调用失败--gas-不足)
    - [验证返回 `verified: false`](#验证返回-verified-false)
    - [前端无法连接后端](#前端无法连接后端)
  - [安全与隐私](#安全与隐私)
    - [绝不能提交到 Git 的内容](#绝不能提交到-git-的内容)
    - [路径安全](#路径安全)
    - [生产环境建议](#生产环境建议)
  - [常见问题](#常见问题)

---

## 功能特性

- **产品生产上链**：厂商创建产品时，将产品关键信息哈希写入百度超级链。
- **物流流转上链**：物流/经销商添加出入库记录时，同步上链存证。
- **销售确认上链**：产品销售状态变更时更新链上记录。
- **防伪验证**：消费者无需登录，输入产品 ID 即可查询链上链下一致性。
- **链上链下哈希校验**：验证时比对链上 `data_hash` 与链下交易记录保存的 `data_hash`。
- **角色隔离**：厂商只能查看自己生产的产品，物流/经销商可查看全部产品用于操作。

---

## 系统架构

```
前端 Vue3 (Element Plus)
    │  HTTP/JSON
    ▼
后端 Go 1.21 + Gin
    │
    ├──► MySQL 8.0（链下完整数据）
    │
    └──► XuperChain Go SDK v2
            │
            ▼
        百度超级链开放网络（节点 39.156.69.83:37100）
            │
            └──► 智能合约 hzy_trace（C++ / WASM）
                    ├─ initialize(admin)
                    ├─ createGoods(id, desc)
                    ├─ updateGoods(id, reason)
                    └─ queryRecords(id)
```

**混合存储**：

- 链上：仅存 `product_uid`、`data_hash`、`timestamp`、`operator` 等关键摘要。
- 链下 MySQL：完整产品详情、物流记录、操作历史。

---

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 / Element Plus / Vite / Pinia |
| 后端 | Go 1.21+ / Gin / GORM / JWT |
| 数据库 | MySQL 8.0 |
| 区块链 | 百度超级链 XuperChain 开放网络 |
| 智能合约 | C++ / WASM |

---

## 目录结构

```
paddle-traceability/
├── backend/              # Go 后端 API 服务
│   ├── blockchain/       # XuperChain SDK 封装
│   ├── config/           # 配置加载（仅读取环境变量）
│   ├── handlers/         # HTTP 接口处理器
│   ├── middleware/       # JWT 认证、CORS 等中间件
│   ├── models/           # 数据模型
│   ├── services/         # 业务逻辑
│   └── conf/             # SDK 配置文件（本地使用，不提交）
├── contracts/            # 智能合约源码与部署辅助脚本
│   ├── hzy_trace.cc      # 合约源码
│   ├── DEPLOY_GUIDE.md   # 合约部署指南
│   └── conf/             # SDK 配置文件（本地使用，不提交）
├── frontend/             # Vue 3 前端
│   ├── src/views/        # 页面组件
│   └── dist/             # 构建产物（.gitignore 忽略）
├── scripts/
│   └── init_db.sql       # MySQL 初始化脚本
└── README.md             # 本文件
```

---

## 环境准备

### 必需环境

- Go 1.21+（后端开发/运行）
- Node.js 18+ 与 npm（前端开发/构建）
- MySQL 8.0
- 百度超级链开放网络账户、私钥文件、合约账户

### 准备私钥

1. 在百度超级链开放网络控制台下载私钥文件（通常为 `private.key`）。
2. 将私钥文件保存在项目**外部**的安全目录，例如：
   - Windows: `D:\secure\xuper\private.key`
   - Linux/macOS: `/secure/xuper/private.key`
3. 记录私钥密码，后续通过环境变量 `XUPER_KEY_PASSWORD` 注入。

> ⚠️ **安全警告**：`private.key` 和私钥密码严禁提交到 Git。项目 `.gitignore` 已配置忽略 `*.key`、`.env*` 及 SDK 配置文件。

---

## 安装与部署

本系统使用**本地安装的 MySQL 8.0** 作为链下数据库。

### 前置要求

| 软件 | 版本要求 | 用途 |
|------|----------|------|
| MySQL | 8.0+ | 链下数据存储 |
| Go | 1.21+ | 后端服务 |
| Node.js | 18+ | 前端构建与开发 |

---

### 方式一：本地开发环境（推荐）

#### 1. 安装 MySQL 8.0

1. 访问 [MySQL Installer 下载页](https://dev.mysql.com/downloads/installer/)，下载离线安装包 `mysql-installer-community-8.0.XX.msi`（约 500MB）。
2. 以管理员身份运行安装程序，选择 **Server only** 或 **Full**。
3. 安装类型选择 **Development Machine**。
4. 身份验证方式选择 **Use Strong Password Encryption**。
5. 设置 root 密码为：**`hzy5414410`**。
6. 勾选 **Configure MySQL Server as a Windows Service** 和 **Add to PATH**。

默认安装路径：

```text
C:\Program Files\MySQL\MySQL Server 8.0\
```

安装完成后，打开新的 PowerShell 窗口验证：

```powershell
# 查看 MySQL 版本
mysql --version

# 使用 root 登录（输入密码 hzy5414410）
mysql -u root -p
```

> **数据库 root 密码属于敏感信息，仅用于本地开发环境。生产环境请创建独立应用账户并限制权限。**

#### 2. 初始化数据库

在项目根目录 `paddle-traceability/` 下执行：

```powershell
# 创建数据库 hzy_trace 及业务表
mysql -u root -p < scripts/init_db.sql
```

输入 root 密码 `hzy5414410`。脚本会创建数据库 `hzy_trace` 以及 `users`、`products`、`logistics_records`、`tx_records` 四张业务表。

#### 3. 配置环境变量

在项目根目录创建 `.env` 文件（已加入 `.gitignore`，不会提交）：

```bash
# MySQL 配置（本地 root 账户）
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=hzy5414410
DB_NAME=hzy_trace

# JWT 密钥（生产环境请使用 ≥32 位的随机字符串）
JWT_SECRET=your-random-jwt-secret-at-least-32-characters

# 百度超级链私钥配置
XUPER_PRIVATE_KEY_PATH=D:\secure\xuper\private.key
XUPER_KEY_PASSWORD=your-private-key-password
XUPER_CONTRACT_ACCOUNT=XC4103761871843472@xuper

# 区块链节点与合约配置
BLOCKCHAIN_NODE_ADDR=39.156.69.83:37100
BLOCKCHAIN_CONTRACT_NAME=hzy_trace
```

> `.env` 文件包含敏感信息，请妥善保管，不要上传到 GitHub。

#### 4. 部署智能合约

参见 [智能合约部署](#智能合约部署)。若已部署，确认合约名称为 `hzy_trace`。

#### 5. 启动后端

**Windows PowerShell:**

```powershell
cd backend

# 从 .env 加载环境变量
Get-Content ..\.env | ForEach-Object {
    if ($_ -match '^\s*([^#][^=]*)\s*=\s*(.*)\s*$') {
        [System.Environment]::SetEnvironmentVariable($matches[1], $matches[2], 'Process')
    }
}

# 启动后端
go mod tidy
go run main.go
```

**Linux / macOS Bash:**

```bash
cd backend

# 从 .env 加载环境变量
export $(grep -v '^#' ../.env | xargs)

# 启动后端
go mod tidy
go run main.go
```

后端默认监听 `:8080`。启动成功后，日志会显示类似 `Server started on :8080`。

#### 6. 启动前端

打开新的终端窗口，执行：

```bash
cd frontend
npm install
npm run dev
```

开发服务器默认监听 `:5173`，打开浏览器访问 `http://localhost:5173` 即可。

---

## 环境变量说明

| 变量名 | 是否必填 | 说明 | 示例 |
|--------|----------|------|------|
| `DB_HOST` | 否 | MySQL 主机地址 | `127.0.0.1` |
| `DB_PORT` | 否 | MySQL 端口 | `3306` |
| `DB_USER` | 否 | MySQL 用户名 | `root` |
| `DB_PASSWORD` | 是 | MySQL 密码 | 无默认值 |
| `DB_NAME` | 否 | 数据库名 | `hzy_trace` |
| `JWT_SECRET` | 是 | JWT 签名密钥 | 无默认值，建议 ≥32 位随机字符串 |
| `XUPER_PRIVATE_KEY_PATH` | 是 | 私钥文件绝对路径 | 无默认值 |
| `XUPER_KEY_PASSWORD` | 是 | 私钥解密密码 | 无默认值 |
| `XUPER_CONTRACT_ACCOUNT` | 否 | 合约账户 | `XC4103761871843472@xuper` |
| `BLOCKCHAIN_NODE_ADDR` | 否 | 开放网络节点地址 | `39.156.69.83:37100` |
| `BLOCKCHAIN_CONTRACT_NAME` | 否 | 合约名 | `hzy_trace` |

> 标注为**必填**的变量若缺失，后端启动时会 `log.Fatalf` 直接退出，避免使用默认值导致安全隐患。

---

## 智能合约部署

由于 Windows 环境缺少本地编译 C++ WASM 工具链，推荐在百度超级链开放网络控制台完成在线编译与部署。

### 部署参数

- **合约源码**：`contracts/hzy_trace.cc`
- **合约名称**：`hzy_trace`
- **合约账户**：`XC4103761871843472@xuper`
- **合约语言**：C++
- **初始化参数 `admin`**：私钥对应的管理员地址

### 部署步骤

详见 [`contracts/DEPLOY_GUIDE.md`](contracts/DEPLOY_GUIDE.md)。

### 部署后验证

部署完成后，可通过后端创建产品并验证，确认链上记录正常写入。也可使用辅助脚本：

```bash
cd contracts
XUPER_KEY_PASSWORD=your-key-password go run verify_deploy.go -id <product_uid>
```

---

## 部署后验证

1. 注册并登录厂商账号：

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "factory_demo",
  "password": "123456",
  "role": "FACTORY",
  "company_name": "Demo Factory"
}
```

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "factory_demo",
  "password": "123456"
}
```

2. 创建产品（使用返回的 JWT）：

```http
POST /api/v1/products
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>

{
  "brand": "DemoBrand",
  "model": "DEMO-001",
  "material": "Wood",
  "batch_no": "DEMO-2026-001",
  "production_date": "2026-06-18"
}
```

3. 记录返回的 `product_uid`，进行防伪验证：

```http
GET /api/v1/verify/{product_uid}
```

成功时返回：

```json
{
  "code": 200,
  "verified": true,
  "data": {
    "chain_verified": true,
    "data_hash_matched": true,
    "status": "PRODUCED",
    "trace_summary": [...]
  }
}
```

---

## 调试与排查

### 后端无法启动

1. 检查必填环境变量是否全部设置：

```powershell
# PowerShell
Get-ChildItem env: | Where-Object { $_.Name -match "DB_PASSWORD|JWT_SECRET|XUPER" }
```

2. 检查 MySQL 是否可连接：

```bash
docker exec -it paddle-mysql mysql -u paddle -p -e "SELECT 1"
```

3. 检查私钥路径是否正确，后端进程是否有读取权限。

### 合约调用失败 / gas 不足

百度超级链开放网络每次合约调用需支付 gas。若出现 `gas not enough` 或交易未上链，请登录开放网络控制台为合约账户充值。

### 验证返回 `verified: false`

- 确认产品 `product_uid` 输入正确。
- 确认合约已部署且名称、账户、节点地址配置正确。
- 查看后端日志中 `QueryRecords` 的原始返回，确认链上是否存在记录。

### 前端无法连接后端

- 检查 `frontend/src/api/request.js` 中的 `baseURL` 是否指向正确的后端地址。
- 检查后端 CORS 配置是否允许前端域名。

---

## 安全与隐私

### 绝不能提交到 Git 的内容

- 私钥文件（`*.key`）
- 私钥密码、数据库密码、JWT 密钥
- `.env`、`.env.local` 等环境变量文件
- `backend/conf/sdk.yaml`、`contracts/conf/sdk.yaml` 等含真实节点地址的配置文件
- 前端 `dist/`、后端 `*.exe`、日志文件

项目 `.gitignore` 已对上述内容做忽略配置，提交前请使用 `git status` 再次确认无敏感文件。

### 路径安全

- 私钥文件请保存在项目目录之外，使用绝对路径通过环境变量引用。
- Windows 路径在 PowerShell 中使用双反斜杠或单斜杠均可，但推荐使用原始字符串避免转义问题。

### 生产环境建议

- 使用 HTTPS 部署前端与后端。
- 使用独立的数据库账户，限制权限（仅对 `hzy_trace` 库有读写权）。
- 定期更换 `JWT_SECRET`。
- 私钥文件建议存储在密钥管理服务（KMS）或硬件加密机中。

---

## 常见问题

**Q: 为什么 Docker Compose 中 `XUPER_PRIVATE_KEY_PATH` 要使用环境变量而不是固定路径？**

A: 固定路径会暴露个人目录结构，且在不同机器上不一致。通过环境变量注入可保证隐私与可移植性。

**Q: 前端 `npm run build` 后为什么没有 `dist/` 目录提交？**

A: `dist/` 是构建产物，已加入 `.gitignore`。生产部署时可在服务器本地构建，或配合 CI/CD 流水线生成。

**Q: 合约名称 `hzy_trace` 与项目名 `paddle-traceability` 为什么不一致？**

A: `paddle-traceability` 是仓库/模块名，`hzy_trace` 是实际部署在百度超级链开放网络上的合约名。代码中通过 `BLOCKCHAIN_CONTRACT_NAME` 配置，默认即为 `hzy_trace`。

**Q: 是否支持本地单节点 XuperChain？**

A: 当前配置面向百度超级链开放网络。若需本地节点，请修改 `BLOCKCHAIN_NODE_ADDR`、`XUPER_CONTRACT_ACCOUNT` 及 SDK 配置，并重新部署合约。

---


