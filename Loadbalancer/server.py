import sys
from flask import Flask
app = Flask(__name__)

serverName = sys.argv[1]


@app.route('/')
def hello():
    return " 'Hello World !' from "+serverName


@app.route("/:name")
def hello_name(name):
    return "Hello "+name+" from "+serverName


if __name__ == "__main__":
    app.run(port=sys.argv[2])
