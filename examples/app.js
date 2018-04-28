const pick = require('lodash/pick');
const client = require('mongodb').MongoClient;

const url = process.env.MONGO_URL;

const getReplicationStatus = (db, cb) => {
    let adminDb = db.db('admin');
    adminDb.command({ replSetGetStatus: 1}, cb);
    db.close();
};

const showReplicationStatus = (err, res, db) => {
    if (err) {
        if (err.constructor.name !== 'MongoError') {
            throw err;
        }

        // Just like the mongo shell.
        res = pick(err, ['ok', 'errmsg', 'code', 'codeName']);
    }

    console.log(res);
};

client.connect(url, (err, db) => {
    if (err) {
        throw err;
    }

    getReplicationStatus(db, showReplicationStatus)
});
