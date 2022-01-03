# Mandatory databases
curl -u ${COUCHDB_USER}:${COUCHDB_PASSWORD} -X PUT http://127.0.0.1:5984/_users
curl -u ${COUCHDB_USER}:${COUCHDB_PASSWORD} -X PUT http://127.0.0.1:5984/_replicator
curl -u ${COUCHDB_USER}:${COUCHDB_PASSWORD} -X PUT http://127.0.0.1:5984/_global_changes

# App database
curl -u ${COUCHDB_USER}:${COUCHDB_PASSWORD} -X PUT http://127.0.0.1:5984/cookbook

