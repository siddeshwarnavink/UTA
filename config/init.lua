local crypto = require "crypto"

function encrypt(data, key)
   return data
end

function decrypt(data, key)
   return data
end

crypto.register("AES", encrypt, decrypt)
