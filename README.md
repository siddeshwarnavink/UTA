

# UTA - Universal TCP Adapter

A highly customized TCP adapter to port any system to up-to-date security standards


## How to run it

    go run . --server -enc 127.0.0.1:9999 -dec 127.0.0.1:10000 --prot dhkc --algo AES
    go run . --client -dec 127.0.0.1:8888 -enc 127.0.0.1:9999 --prot dhkc --algo AES
    go run demo/server.go
    ruby demo/client.rb

