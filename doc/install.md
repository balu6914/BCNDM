# Bigsense Install

## MongoDB and BigchainDB
Requirements:
- MongoDB >= v3.4
- BigchainDB

BigchainDB installation instructions can be found [here](https://docs.bigchaindb.com/projects/server/en/latest/quickstart.html)

Start MongoDB with "bigchain-rs" replica set option (default configuration of BigchainDB):
```
sudo mongod --replSet=bigchain-rs
```

Start BigchainDB:
```
bigchaindb start
```

> Make sure that you previously configured BigchainDB with `bigchaindb -y configure mongodb`

## Bigsense
Install dependencies:
```
sudo pip --proxy $http_proxy install -r requirements.txt
```

Run:
```
gunicorn bigsense:app
```
