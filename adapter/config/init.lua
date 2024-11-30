local crypto = require("crypto")
local aes = require("algo.aes") -- current standard algorithms will be provided
local keyExchange = require("keyExchange")
local dh = require("keyalgo.dh")
local rsa = require("keyalgo.rsa")

crypto.register("AES", aes.encrypt, aes.decrypt)

keyExchange.register("DHKC", dh.clientDiffieHellman, dh.serverDiffieHellman)
keyExchange.register("RSA", rsa.clientRSA, rsa.serverRSA)



-- local ui = require("ui")
-- local mcq = require("ui.mcq")
-- local form = require("ui.form")

-- ui.add(mcq.new("Select the mode", {"Client", "Server"}))
-- ui.add(form.new("Enter the Connection Addresses", {"Unencrypted Connection's Address", "Encrypted Connection's Address"}))
-- ui.add(mcq.new("Select the key exchange algorithm", keyExchange.list))
-- ui.add(mcq.new("Select the encryption algorithm", crypto.list))
