#include "xchain/xchain.h"
#include <cstdio>
#include <cstdlib>
#include <memory>
#include <string>

// 商品溯源合约：hzy_trace
// 基于 XuperChain C++ WASM 模板实现
class SourceTrace {
public:
    // 初始化商品溯源合约，指定管理员地址
    virtual void initialize() = 0;
    // 创建一个新的商品，仅 admin 可调用
    virtual void createGoods() = 0;
    // 变更商品信息，仅 admin 可调用
    virtual void updateGoods() = 0;
    // 查询商品变更记录，公开查询
    virtual void queryRecords() = 0;
};

struct SourceTraceDemo : public SourceTrace, public xchain::Contract {
public:
    const std::string GOODS = "GOODS_";
    const std::string GOODSRECORD = "GOODSSRECORD_";
    const std::string GOODSRECORDTOP = "GOODSSRECORDTOP_";
    const std::string CREATE = "CREATE";

    void initialize() {
        xchain::Context* ctx = this->context();
        const std::string& admin = ctx->arg("admin");
        if (admin.empty()) {
            ctx->error("missing admin address");
            return;
        }

        ctx->put_object("admin", admin);
        ctx->ok("initialize success");
    }

    bool isAdmin(xchain::Context* ctx, const std::string& caller) {
        std::string admin;
        if (!ctx->get_object("admin", &admin)) {
            return false;
        }
        return (admin == caller);
    }

    void createGoods() {
        xchain::Context* ctx = this->context();
        const std::string& caller = ctx->initiator();
        if (caller.empty()) {
            ctx->error("missing initiator");
            return;
        }

        if (!isAdmin(ctx, caller)) {
            ctx->error("only the admin can create new goods");
            return;
        }

        const std::string& id = ctx->arg("id");
        if (id.empty()) {
            ctx->error("missing 'id' as goods identity");
            return;
        }

        const std::string& desc = ctx->arg("desc");
        if (desc.empty()) {
            ctx->error("missing 'desc' as goods desc");
            return;
        }

        std::string goodsKey = GOODS + id;
        std::string value;
        if (ctx->get_object(goodsKey, &value)) {
            ctx->error("the id is already exist, please check again");
            return;
        }
        ctx->put_object(goodsKey, desc);

        std::string goodsRecordsKey = GOODSRECORD + id + "_0";
        std::string goodsRecordsTopKey = GOODSRECORDTOP + id;
        ctx->put_object(goodsRecordsKey, CREATE);
        ctx->put_object(goodsRecordsTopKey, "0");
        ctx->ok(id);
    }

    void updateGoods() {
        xchain::Context* ctx = this->context();
        const std::string& caller = ctx->initiator();
        if (caller.empty()) {
            ctx->error("missing initiator");
            return;
        }

        if (!isAdmin(ctx, caller)) {
            ctx->error("only the admin can update goods");
            return;
        }

        const std::string& id = ctx->arg("id");
        if (id.empty()) {
            ctx->error("missing 'id' as goods identity");
            return;
        }

        const std::string& reason = ctx->arg("reason");
        if (reason.empty()) {
            ctx->error("missing 'reason' as update reason");
            return;
        }

        std::string goodsKey = GOODS + id;
        std::string goodsValue;
        if (!ctx->get_object(goodsKey, &goodsValue)) {
            ctx->error("the id not exist, please check again");
            return;
        }

        std::string goodsRecordsTopKey = GOODSRECORDTOP + id;
        std::string value;
        ctx->get_object(goodsRecordsTopKey, &value);
        int topRecord = 0;
        topRecord = atoi(value.c_str()) + 1;

        char topRecordVal[32];
        snprintf(topRecordVal, 32, "%d", topRecord);
        std::string goodsRecordsKey = GOODSRECORD + id + "_" + topRecordVal;

        ctx->put_object(goodsRecordsKey, reason);
        ctx->put_object(goodsRecordsTopKey, topRecordVal);
        ctx->ok(topRecordVal);
    }

    void queryRecords() {
        xchain::Context* ctx = this->context();
        const std::string& id = ctx->arg("id");
        if (id.empty()) {
            ctx->error("missing 'id' as goods identity");
            return;
        }

        std::string goodsKey = GOODS + id;
        std::string value;
        if (!ctx->get_object(goodsKey, &value)) {
            ctx->error("the id not exist, please check again");
            return;
        }

        std::string goodsRecordsKey = GOODSRECORD + id;
        std::unique_ptr<xchain::Iterator> iter =
            ctx->new_iterator(goodsRecordsKey, goodsRecordsKey + "~");

        std::string result = "\n";
        while (iter->next()) {
            std::pair<std::string, std::string> res;
            iter->get(&res);
            if (res.first.length() > goodsRecordsKey.length()) {
                std::string goodsRecord = res.first.substr(GOODSRECORD.length());
                std::string::size_type pos = goodsRecord.find("_");
                std::string goodsId = goodsRecord.substr(0, pos);
                std::string updateRecord = goodsRecord.substr(pos + 1, goodsRecord.length());
                std::string reason = res.second;

                result += "goodsId=" + goodsId + ",updateRecord=" + updateRecord +
                          ",reason=" + reason + '\n';
            }
        }
        ctx->ok(result);
    }
};

DEFINE_METHOD(SourceTraceDemo, initialize);
static void cxx_initialize(SourceTraceDemo& self) { self.initialize(); }

DEFINE_METHOD(SourceTraceDemo, createGoods);
static void cxx_createGoods(SourceTraceDemo& self) { self.createGoods(); }

DEFINE_METHOD(SourceTraceDemo, updateGoods);
static void cxx_updateGoods(SourceTraceDemo& self) { self.updateGoods(); }

DEFINE_METHOD(SourceTraceDemo, queryRecords);
static void cxx_queryRecords(SourceTraceDemo& self) { self.queryRecords(); }
