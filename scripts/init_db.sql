-- 基于区块链的乒乓球拍防伪溯源系统 - 数据库初始化脚本
-- 使用方法: mysql -u root -p < scripts/init_db.sql

-- 创建数据库
CREATE DATABASE IF NOT EXISTS hzy_trace CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE hzy_trace;

-- ========================================
-- 用户表 (users)
-- ========================================
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID，自增主键',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希（bcrypt）',
    role ENUM('FACTORY', 'LOGISTICS', 'RETAILER') NOT NULL COMMENT '角色类型',
    public_key TEXT COMMENT '公钥（PEM格式）',
    company_name VARCHAR(100) COMMENT '企业/组织名称',
    phone VARCHAR(20) COMMENT '手机号',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE INDEX idx_username (username),
    INDEX idx_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ========================================
-- 产品表 (products)
-- ========================================
CREATE TABLE IF NOT EXISTS products (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '产品ID，自增主键',
    product_uid VARCHAR(19) NOT NULL COMMENT '基于雪花算法的唯一业务ID',
    brand VARCHAR(50) COMMENT '品牌名称',
    model VARCHAR(100) COMMENT '型号',
    material VARCHAR(100) COMMENT '底板材质',
    rubber_type VARCHAR(100) COMMENT '胶皮类型',
    batch_no VARCHAR(50) COMMENT '生产批次号',
    production_date DATE COMMENT '生产日期',
    quality_report_hash VARCHAR(64) COMMENT '质检报告哈希',
    factory_id BIGINT UNSIGNED NOT NULL COMMENT '生产厂商ID，外键→users.id',
    status ENUM('PRODUCED', 'IN_TRANSIT', 'IN_STOCK', 'SOLD') DEFAULT 'PRODUCED' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE INDEX idx_product_uid (product_uid),
    INDEX idx_factory_id (factory_id),
    INDEX idx_status (status),
    FOREIGN KEY (factory_id) REFERENCES users(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品表';

-- ========================================
-- 物流记录表 (logistics_records)
-- ========================================
CREATE TABLE IF NOT EXISTS logistics_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '记录ID，自增主键',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '产品ID，外键→products.id',
    action ENUM('INBOUND', 'OUTBOUND') NOT NULL COMMENT '动作类型',
    operator_id BIGINT UNSIGNED NOT NULL COMMENT '操作人ID，外键→users.id',
    warehouse_name VARCHAR(100) COMMENT '仓库名称',
    location VARCHAR(200) COMMENT '地理位置',
    carrier VARCHAR(100) COMMENT '承运方',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_product_id (product_id),
    INDEX idx_operator_id (operator_id),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (operator_id) REFERENCES users(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='物流记录表';

-- ========================================
-- 交易存证表 (tx_records)
-- ========================================
CREATE TABLE IF NOT EXISTS tx_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '记录ID，自增主键',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '产品ID，外键→products.id',
    tx_hash VARCHAR(128) COMMENT '链上交易哈希',
    tx_type ENUM('CREATE', 'TRANSFER', 'CONFIRM') NOT NULL COMMENT '交易类型',
    data_hash VARCHAR(64) COMMENT '链上存储的数据摘要',
    chain_status ENUM('PENDING', 'CONFIRMED', 'FAILED') DEFAULT 'PENDING' COMMENT '链上状态',
    block_height BIGINT DEFAULT 0 COMMENT '区块高度',
    operator_id BIGINT UNSIGNED NOT NULL COMMENT '操作人ID，外键→users.id',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_product_id (product_id),
    INDEX idx_tx_hash (tx_hash),
    INDEX idx_operator_id (operator_id),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (operator_id) REFERENCES users(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='交易存证表';

-- 完成
SELECT '数据库初始化完成！' AS message;
