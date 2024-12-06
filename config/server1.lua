local crypto = require("crypto")
local aes = require("algo.aes") 
local keyExchange = require("keyExchange")
local dh = require("keyalgo.dh")
local conf = require("config")

crypto.register("AES", aes.encrypt, aes.decrypt)
keyExchange.register("DHKC", dh.clientDiffieHellman, dh.serverDiffieHellman)

conf.serverMode(true) 
conf.decryptPort("server-1:10000")
conf.encryptPort("server-1:9999") 
conf.crypto("AES")
conf.keyExchange("DHKC")