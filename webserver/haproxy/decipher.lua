---
--- Created by mma.
---
--- args[1] = TLS data collector URL
--- args[2] = Max queue size
---

local q = core.queue()
local qSize = 0
local args = table.pack(...)

local function watcher()
    local httpclient = core.httpclient()
    while true do
        local data = q.pop_wait(q)
        httpclient:post { url = args[1], body = data }
    end
end

local function InitQueueSize()
    if args[2] == nil then
        qSize = 1024
    else
        qSize = tonumber(args[2])
    end
    core.Debug("Max queue size set to " .. tostring(qSize))
end

function serialize(data)
    local sData = core.concat()
    for k, v in pairs(data) do
        sData:add("\"")
        sData:add(k)
        sData:add("\":\"")
        sData:add(v)
        sData:add("\",")
    end
    local sData2 = core.concat()
    sData2:add("{")
    sData2:add(sData:dump():sub(1, -2))
    sData2:add("}")
    return sData2:dump()
end

local function SendTlsDataFE(txn)
    local t = {}
    t.v = txn.sf:ssl_fc_protocol()
    t.cr = txn.sf:ssl_fc_client_random()
    t.cr = txn.sc:hex(t.cr)
    t.ssk = txn.sf:ssl_fc_session_key()
    t.ssk = txn.sc:hex(t.ssk)
    t.cets = txn.sf:ssl_fc_client_early_traffic_secret()
    t.chts = txn.sf:ssl_fc_client_handshake_traffic_secret()
    t.shts = txn.sf:ssl_fc_server_handshake_traffic_secret()
    t.cts0 = txn.sf:ssl_fc_client_traffic_secret_0()
    t.sts0 = txn.sf:ssl_fc_server_traffic_secret_0()
    t.ees = txn.sf:ssl_fc_early_exporter_secret()
    t.es = txn.sf:ssl_fc_exporter_secret()
    local json = serialize(t)
    if q:size() < qSize then
        q.push(q, json)
    else
        core.Debug("TLS deciphering: the queue is full, dropping message")
    end
end

local function SendTlsDataBE(txn)
    local t = {}
    t.v = txn.sf:ssl_bc_protocol()
    t.cr = txn.sf:ssl_bc_client_random()
    t.cr = txn.sc:hex(t.cr)
    t.ssk = txn.sf:ssl_bc_session_key()
    t.ssk = txn.sc:hex(t.ssk)
    t.cets = txn.sf:ssl_bc_client_early_traffic_secret()
    t.chts = txn.sf:ssl_bc_client_handshake_traffic_secret()
    t.shts = txn.sf:ssl_bc_server_handshake_traffic_secret()
    t.cts0 = txn.sf:ssl_bc_client_traffic_secret_0()
    t.sts0 = txn.sf:ssl_bc_server_traffic_secret_0()
    t.ees = txn.sf:ssl_bc_early_exporter_secret()
    t.es = txn.sf:ssl_bc_exporter_secret()
    local json = serialize(t)
    if q:size() < qSize then
        q.push(q, json)
    else
        core.Debug("TLS deciphering: the queue is full, dropping message")
    end
end

core.register_init(InitQueueSize)
core.register_task(watcher);
core.register_action("SendFeTlsData", { "tcp-req", "http-req" }, SendTlsDataFE, 0);
core.register_action("SendBeTlsData", { "tcp-res", "http-res" }, SendTlsDataBE, 0);
