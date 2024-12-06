local crypto = require("crypto")
local aes = require("algo.aes") 
local keyExchange = require("keyExchange")
local dh = require("keyalgo.dh")
local conf = require("config")

crypto.register("AES", aes.encrypt, aes.decrypt)
keyExchange.register("DHKC", dh.clientDiffieHellman, dh.serverDiffieHellman)

conf.serverMode(false) 
conf.encryptPort("server-1:9999") 
conf.decryptPort("client-1:8888")
conf.crypto("AES")
conf.keyExchange("DHKC")