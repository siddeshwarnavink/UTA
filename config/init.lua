local crypto = require "crypto"

function encrypt(key, data)
   return data
end

function decrypt(key, data)
   return data
end

crypto.register("AES", encrypt, decrypt)
