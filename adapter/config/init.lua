local crypto = require("crypto")
local aes = require("algo.aes") -- current standard algorithms will be provided
local keyExchange = require("keyExchange")
local dh = require("keyalgo.dh")
local rsa = require("keyalgo.rsa")

crypto.register("AES", aes.encrypt, aes.decrypt)

keyExchange.register("DHKC", dh.clientDiffieHellman, dh.serverDiffieHellman)
keyExchange.register("RSA", rsa.clientRSA, rsa.serverRSA)

local ui = require("ui")
local mcq = require("ui.mcq")
local form = require("ui.form")

ui.add("MODE","Select the mode",{"Client", "Server"}," ",mcq.new)
ui.add("UNENCRYPTED_ADDRESS","Unencrypted Connection's Address",{" "},"127.0.0.1:8888",form.new)
ui.add("ENCRYPTED_ADDRESS","Encrypted Connection's Address",{" "},"127.0.0.1:9999",form.new)
ui.add("KEY_EXCHANGE","Select the key exchange algorithm",{"DHKC", "RSA"}," ",mcq.new)
ui.add("ENCRYPTION","Select the encryption algorithm",{"AES"}," ",mcq.new)
