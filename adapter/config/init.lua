local crypto = require("crypto")
local aes = require("algo.aes") -- current standard algorithms will be provided
local keyExchange = require("keyExchange")
local dh = require("keyalgo.dh")
local ecdh = require("keyalgo.ecdh")

crypto.register("AES", aes.encrypt, aes.decrypt)

keyExchange.register("dhkc", dh.clientDiffieHellman, dh.serverDiffieHellman)
keyExchange.register("ecdhkc", ecdh.clientECDH, ecdh.serverECDH)
