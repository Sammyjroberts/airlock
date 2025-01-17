import os
import socket

import msgpack

SOCKET_PATH = "math.sock"

# Cleanup socket if it exists
if os.path.exists(SOCKET_PATH):
    os.remove(SOCKET_PATH)

# Create server socket
server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
server.bind(SOCKET_PATH)
server.listen(1)

print("Waiting for connection...")
conn, addr = server.accept()

try:
    # Receive data
    data = conn.recv(1024)
    if data:
        # Unpack the message
        msg = msgpack.unpackb(data)
        print("Received:", msg)
        response = {"result": msg["args"][0] + msg["args"][1]}
        print("Sending:", response)
        conn.send(msgpack.packb(response))

finally:
    conn.close()
    os.remove(SOCKET_PATH)
