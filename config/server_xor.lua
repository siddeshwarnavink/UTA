local crypto = require("crypto")

-- XOR Encryption/Decryption Function
local function xor_cipher(key, val)
    local result = {}
    local key_len = #key
    for i = 1, #val do
        -- XOR each byte of the value with the corresponding byte of the key (repeating the key if necessary)
        local val_byte = string.byte(val, i)
        local key_byte = string.byte(key, (i - 1) % key_len + 1)
        result[i] = string.char(bit.bxor(val_byte, key_byte))
    end
    return table.concat(result)
end

-- Register XOR as a simple encryption method
crypto.register("XOR", xor_cipher, xor_cipher)

-- Call the crypto functions with the XOR cipher
conf.crypto("XOR")
conf.keyExchange("DHKC")
conf.serverMode(true)
conf.decryptPort("127.0.0.1:10000")
conf.encryptPort("127.0.0.1:9999")
