local crypto = require "crypto"
local aes = require "algo.aes" -- current standard algorithms will be provided

-- if a new and better algorithm comes, we can code that algorithm
-- here in lua code

crypto.register("AES", aes.encrypt, aes.decrypt)
