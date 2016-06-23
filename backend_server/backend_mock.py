import flask
from flask import Flask
import pickle
from grpc.beta import implementations
import handler_manager_pb2
import os.path

app = Flask(__name__)

@app.route("/")
def root_request():
	return "Hello world!"

TOTALLY_NON_RANDOM_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"

@app.route("/key")
def key_request():
	ret = {}
	ret['grant_type'] = 'refresh_token'
	ret['code'] = 'randomcode'
	ret['redirect_url'] = 'http://localhost:1234/some/random/app'
	ret['access_token'] = TOTALLY_NON_RANDOM_TOKEN
	ret['refresh_token'] = TOTALLY_NON_RANDOM_TOKEN
	ret['email'] = 'doe@example.com'
	ret['expires'] = 999999999999
	ret['expires_in'] = 999999999999
	ret['client_id'] = 'ttn_fake'
	return flask.jsonify(ret)

@app.route("/users/token", methods = ['GET', 'POST'])
def user_login():
	ret = {}
	ret['grant_type'] = 'refresh_token'
	ret['code'] = 'randomcode'
	ret['redirect_url'] = 'http://localhost:1234/some/random/app'
	ret['client_id'] = 'ttn_fake'
	ret['access_token'] = TOTALLY_NON_RANDOM_TOKEN
	ret['refresh_token'] = TOTALLY_NON_RANDOM_TOKEN
	ret['email'] = 'doe@example.com'
	ret['expires'] = 999999999999
	ret['expires_in'] = 999999999999
	ret['client_id'] = 'ttn_fake'
	return flask.jsonify(ret)

applications = list()
handlers = [ {'host': 'localhost', 'port': 1782 }]
APP_FILE = 'applications.pickle'

@app.route("/applications", methods = ['GET'])
def list_applications():
	ret = list()
	for app in applications:
		ret.append(app)
		#{'EUI': 'BEEFBABEBEEFBABE', 'Name': 'The mocked application', 'Owner': 'doe@example.com', 'AccessKeys': ['token'], 'Valid': True});
#	"EUI", "Name", "Owner", "Access Keys", "Valid"
	return flask.jsonify(ret);

def generate_app_key():
	return 'APPKEYAPPKEYAPPKEYAPPKEY'

def generate_access_key():
	return 'ACCESSKEYACCESSKEY'

def generate_eui():
	return 'EUIEUI'

def save_apps():
	pickle.dump(applications, APP_FILE)

def load_apps():
	if os.path.exists(APP_FILE):
		applications = pickle.load(APP_FILE)

def update_handlers(newapp):
	for handler in handlers:
		channel = implementations.insecure_channel(handler['host'], handler['port'])
		params = handler_manager_pb2.SetDefaultDeviceReq()
		params.Token = 'ignored'
		params.AppEUI = newapp['EUI']
		params.AppKey = newapp['AppKey']
		stub = handler_manager_pb2.beta_create_HandlerManager_stub(channel)
		stub.SetDefaultDevice(params)

@app.route("/applications", methods = ['POST'])
def create_application():
	app = flask.request.form
	app_key = generate_app_key()
	newapp = {'EUI': generate_eui(), 'Name': app['name'], 'Owner': 'lora@telenordigital.com', 'AccessKeys': [ generate_access_key() ], 'Valid': True, 'AppKey': app_key }
	update_handlers(newapp)
	applications.append(newapp)
	save_apps()
	return flask.jsonify(newapp), 201


if __name__ == "__main__":
	load_apps()
	app.run()
