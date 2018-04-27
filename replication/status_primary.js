const primary = {
    "set": "ranking",
    "date": ISODate("2018-04-27T14:50:56.167Z"),
    "myState": 1,
    "term": NumberLong(3),
    "heartbeatIntervalMillis": NumberLong(2000),
    "optimes": {
        "lastCommittedOpTime": {
            "ts": Timestamp(1524840647, 1),
            "t": NumberLong(3)
        },
        "appliedOpTime": {
            "ts": Timestamp(1524840647, 1),
            "t": NumberLong(3)
        },
        "durableOpTime": {
            "ts": Timestamp(1524840647, 1),
            "t": NumberLong(3)
        }
    },
    "members": [
        {
            "_id": 0,
            "name": "mongod_single:27017",
            "health": 1,
            "state": 1,
            "stateStr": "PRIMARY",
            "uptime": 5192,
            "optime": {
                "ts": Timestamp(1524840647, 1),
                "t": NumberLong(3)
            },
            "optimeDate": ISODate("2018-04-27T14:50:47Z"),
            "electionTime": Timestamp(1524835465, 1),
            "electionDate": ISODate("2018-04-27T13:24:25Z"),
            "configVersion": 1,
            "self": true
        }
    ],
    "ok": 1
};

