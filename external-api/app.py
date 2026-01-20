from flask import Flask, request, jsonify
from xmlrpc.server import SimpleXMLRPCServer, SimpleXMLRPCRequestHandler
import threading

app = Flask(__name__)

STATE_TAX = {
    "Kentucky": 0.07,
    "California": 0.085,
    "Florida": 0.06,
    "New York": 0.088,
}

@app.route("/tax", methods=["GET"])
def get_tax():
    state = request.args.get("state")
    if not state:
        return jsonify({"error": "state is required"}), 400

    rate = STATE_TAX.get(state, 0.05)  # default
    return jsonify({
        "state": state,
        "taxRate": rate
    })

# ---- XML-RPC Server (Protocolo 3) ----
class TaxRPCHandler(SimpleXMLRPCRequestHandler):
    rpc_paths = ('/RPC2',)

def get_tax_rate_rpc(state: str) -> dict:
    """XML-RPC method to get tax rate for a state"""
    rate = STATE_TAX.get(state, 0.05)
    return {
        "state": state,
        "taxRate": rate,
        "source": "xmlrpc"
    }

def start_xml_rpc_server():
    """Start XML-RPC server in a separate thread"""
    server = SimpleXMLRPCServer(('0.0.0.0', 8091), requestHandler=TaxRPCHandler)
    server.register_function(get_tax_rate_rpc, 'getTaxRate')
    print("[XML-RPC] Server started on :8091")
    server.serve_forever()

if __name__ == "__main__":
    rpc_thread = threading.Thread(target=start_xml_rpc_server, daemon=True)
    rpc_thread.start()
    
    app.run(host="0.0.0.0", port=8090)
