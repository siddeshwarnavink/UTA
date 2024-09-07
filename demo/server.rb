require 'socket'

server_ip = '127.0.0.1'
server_port = 10000

server = TCPServer.new(server_ip, server_port)

puts "Server started on #{server_ip}:#{server_port}"

loop do
  client = server.accept
  puts "Client connected"

  while true
    data = client.recv(1024)
    break if data.empty?

    puts "Received: #{data}"

    response = "Hello from server!"
    client.send(response, 0)
    puts "Sent: #{response}"

    sleep(0.5)
  end

  client.close
  puts "Client disconnected"
end