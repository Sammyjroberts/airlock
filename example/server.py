import argparse
import importlib.util
import os
import socket
import sys

import msgpack


def load_implementation(file_path: str, class_name: str):
    """Dynamically load implementation class from a Python file."""
    # Get absolute path
    abs_path = os.path.abspath(file_path)

    # Load module from file path
    spec = importlib.util.spec_from_file_location("implementation", abs_path)
    if not spec or not spec.loader:
        raise ImportError(f"Could not load {file_path}")

    module = importlib.util.module_from_spec(spec)
    sys.modules["implementation"] = module
    spec.loader.exec_module(module)

    # Get the implementation class
    if not hasattr(module, class_name):
        raise AttributeError(f"Class {class_name} not found in {file_path}")

    impl_class = getattr(module, class_name)
    return impl_class()


def serve_rpc(socket_path: str, handler: object):
    """Run RPC server with given handler."""
    # Cleanup socket if it exists
    if os.path.exists(socket_path):
        os.remove(socket_path)

    # Create server socket
    server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    server.bind(socket_path)
    server.listen(1)

    print(f"Server started on {socket_path}")
    print(
        f"Available methods: {[method for method in dir(handler) if not method.startswith('_')]}"
    )

    try:
        while True:
            print("\nWaiting for connection...")
            conn, addr = server.accept()

            try:
                # Receive data
                data = conn.recv(4096)  # Increased buffer size
                if data:
                    # Unpack the message
                    msg = msgpack.unpackb(data)
                    print(f"Received: {msg}")

                    try:
                        # Get the method from the handler
                        if not hasattr(handler, msg["func"]):
                            raise AttributeError(f"Method {msg['func']} not found")

                        method = getattr(handler, msg["func"])
                        result = method(*msg["args"])

                        response = {"result": result, "error": ""}
                    except Exception as e:
                        response = {"result": None, "error": str(e)}

                    print(f"Sending: {response}")
                    conn.send(msgpack.packb(response))
            finally:
                conn.close()
    except KeyboardInterrupt:
        print("\nShutting down server...")
    finally:
        server.close()
        if os.path.exists(socket_path):
            os.remove(socket_path)


def main():
    parser = argparse.ArgumentParser(description="Generic RPC Server")
    parser.add_argument(
        "implementation", help="Python file containing the implementation"
    )
    parser.add_argument("class_name", help="Name of the implementation class")
    parser.add_argument("--socket", default="rpc.sock", help="Unix socket path")

    args = parser.parse_args()

    try:
        handler = load_implementation(args.implementation, args.class_name)
        serve_rpc(args.socket, handler)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
