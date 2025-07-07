--[[
    返回值
    -1: 系统问题
    -2: 发送成功
    0: 请求发送太频繁（1分钟内多次请求）
]]
local key = KEYS[1]
-- 验证码可验证的次数 key
local cntKey = key .. ":cnt"
-- 准备存储的验证码
local val = ARGV[1]
-- 验证码过期时间
local ttl = tonumber(redis.call("ttl", key))
local EXPIRE_SECONDS = 300 -- 验证码 5 分钟过期
local VERIFY_TIMES = 3  -- 一个验证码，最多可以输入3次
local SEND_LIMIT_SECONDS = 60  -- 发送间隔控制时间: 1分钟
if ttl == -1 then
    -- key 存在，但是没有过期时间。说明有人调整过这个key，或者key冲突。属于系统错误
    return -1
elseif ttl == -2 or ttl < EXPIRE_SECONDS - SEND_LIMIT_SECONDS then
    -- key 不存在，可以发送验证码
    redis.call("set", key, val)
    redis.call("expire", key, EXPIRE_SECONDS)
    redis.call("set", cntKey, VERIFY_TIMES)
    redis.call("expire", cntKey, EXPIRE_SECONDS)
    -- 发送成功
    return 0
else
    -- 已经发送了一个验证码，但是还不到一分钟。同一个手机号码，一分钟内只能发送一次
    return -2
end
