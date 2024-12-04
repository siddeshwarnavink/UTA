# UTA Running Commands (Locally)

## Demo

### Server Application
cd demo/server-1
go run . --local --port 10000

### Client Application
cd demo/client-1
go run . --local --port 8888

## Adapter

### Client Proxy
go run ./adapter --config adapter/config/client.lua

### Server Proxy
go run ./adapter --config adapter/config/server.lua

## Wizard

### Build
cd wizard
npm install # if you haven't
npm run dev

### Run
go run ./wizard