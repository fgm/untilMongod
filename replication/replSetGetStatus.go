// Initial attempt to use replication information, but a mere Dial() turned out
// to be sufficient. This will like be used to check for rs.status() as an
// extension of the command like untilMongod --isPrimary --replSet foo
package replication

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type optime struct {
	t  int64
	ts bson.MongoTimestamp
}

type member map[string]interface{}

type ReplSetStatus struct {
	date              time.Time
	heartbeatInterval int64
	members           []member
	myState           int
	ok                float64
	optimes           map[string]optime
	set               string
	term              int64
}

func ReplSetGetStatus(session *mgo.Session) ReplSetStatus {
	adminDb := session.DB("admin")

	result := bson.M{}
	err := adminDb.Run("replSetGetStatus", &result)
	if err != nil {
		panic(err)
	}

	status := ReplSetStatus{
		date:              result["date"].(time.Time),
		heartbeatInterval: result["heartbeatIntervalMillis"].(int64),
		members:           make([]member, 0),
		myState:           result["myState"].(int),
		ok:                result["ok"].(float64),
		optimes:           make(map[string]optime),
		set:               result["set"].(string),
		term:              result["term"].(int64),
	}

	for k, v := range result["optimes"].(bson.M) {
		mv := v.(bson.M)
		o := optime{
			t:  mv["t"].(int64),
			ts: mv["ts"].(bson.MongoTimestamp),
		}
		status.optimes[k] = o
	}

	for _, variantMember := range result["members"].([]interface{}) {
		m := member{}
		for k, v := range variantMember.(bson.M) {
			m[k] = v
		}
		status.members = append(status.members, m)
	}

	return status
}
