local crypto = require("crypto")
local aes = require("algo.aes") -- current standard algorithms will be provided
local keyExchange = require("keyExchange")
local dh = require("keyalgo.dh")
local rsa = require("keyalgo.rsa")

crypto.register("AES", aes.encrypt, aes.decrypt)

keyExchange.register("DHKC", dh.clientDiffieHellman, dh.serverDiffieHellman)
keyExchange.register("RSA", rsa.clientRSA, rsa.serverRSA)

-- local mode = require("mode")

local ui = require("ui")
local mcq = require("ui.mcq")
-- local form = require("ui.question")

-- ui.add(
    mcq.new("Select the mode", {"Client", "Server"},"")
-- )
-- OR
-- ui.add(mcq.new("Select the mode", mode.list))

-- ui.add(question.new("Unencrypted Connection's Address"))
-- ui.add(question.new(""Encrypted Connection's Address"))
-- ui.add(
    mcq.new("Select the key exchange algorithm", keyExchange.list)
-- )
-- ui.add(
    mcq.new("Select the encryption algorithm", crypto.list)
-- )


