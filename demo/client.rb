require 'socket'

server_ip = '127.0.0.1'
server_port = 8888
# server_port = 10000

client = TCPSocket.new(server_ip, server_port)

puts "Connected to server at #{server_ip}:#{server_port}"

message = "Hello from client!"

10.times do |i|
  client.send(message, 0)
  puts "Sent: #{message}"

  response = client.recv(1024)
  puts "Received: #{response}"

  sleep 1
end

client.close
puts "Connection closed"