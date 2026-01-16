from flask import Flask, request, jsonify

app = Flask(__name__)

# Mock simples (mais que suficiente para o TP)
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

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8090)
