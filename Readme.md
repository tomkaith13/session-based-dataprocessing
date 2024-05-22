# Intro
This repo is a POC that looks at setting up a session based table using MongoDB to do some processing as long as the 
data resides in the database. We use TTL Indexes to ensure the data expires.
## How to run
You can build the containers and the get the service running using:
```
make buildup
```
This runs the container in as a foreground process.
Please use `-d` option to the `buildup` rule in Makefile to ensure it runs as a daemon.
Once you are done, clean up can be done using
```
make clean
```

If you want to clear the DB, please delete the Volume setup.

## How to login to mongosh to check contents
1. We first attach a shell to the mongodb container of the service
2. We then initiate `mongosh` command
```
# mongosh
Current Mongosh Log ID: 664ce644383fd660af99ea71
Connecting to:          mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.2.5
Using MongoDB:          7.0.9
Using Mongosh:          2.2.5

For mongosh info see: https://docs.mongodb.com/mongodb-shell/

test> 
```
3. We then switch to `admin` and authenticate by doing:
```
test> use admin
switched to db admin
admin> db.auth("rootuser", "rootpass")
{ ok: 1 }
admin>
```
4. we can then finally switch to our DB and then execute commands
```bash
admin> use mydatabase
switched to db mydatabase
mydatabase> show collections
persons
mydatabase> db.persons.getIndexes()
[
  { v: 2, key: { _id: 1 }, name: '_id_' },
  {
    v: 2,
    key: { createdAt: 1 },
    name: 'createdAt_1',
    expireAfterSeconds: 600
  }
]
mydatabase> db.persons.find()
[
  {
    _id: '4b6af7c5-c124-42d5-9538-5d0935fc2f82',
    name: 'name0',
    city: 'Toronto',
    age: 21,
    createdAt: ISODate('2024-05-21T18:15:50.377Z')
  }
]
```
