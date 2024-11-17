local crypto = require("crypto")
local aes = require("algo.aes") -- current standard algorithms will be provided
local keyExchange = require("keyExchange")
local dh = require("keyalgo.dh")
local rsa = require("keyalgo.rsa")

crypto.register("AES", aes.encrypt, aes.decrypt)

keyExchange.register("DHKC", dh.clientDiffieHellman, dh.serverDiffieHellman)
keyExchange.register("RSA", rsa.clientRSA, rsa.serverRSA)
