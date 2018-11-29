package mongo

import (
	logger "mtcomm/log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoClient interface {
	Insert(db string, c string, structBsonPrt interface{}) error
	Update(db string, c string, _id string, m map[string]interface{}) error
	Save(db string, c string, _id string, m map[string]interface{}) error
	SearchById(db string, c string, _id string, structBsonPrt interface{}) error
}

type mongoClient struct {
	session *mgo.Session
	log     *logger.Logger
}

func NewMongoClient(info *MongoServerInfo) MongoClient {
	log := logger.GetDefaultLogger()
	session, err := mgo.Dial(info.MongoHost)
	if err != nil {
		panic(err)
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return &mongoClient{
		session: session,
		log:     log,
	}
}

type MongoServerInfo struct {
	MongoHost string
}

func (client *mongoClient) Insert(db string, collection string, structBsonPrt interface{}) error {
	c := client.session.DB(db).C(collection)
	err := c.Insert(structBsonPrt)
	if err != nil {
		client.log.Error(err)
	}
	return err
}

func (client *mongoClient) Update(db string, collection string, _id string, m map[string]interface{}) error {
	c := client.session.DB(db).C(collection)
	err := c.Update(bson.M{"_id": _id}, bson.M{"$set": bson.M(m)})
	if err != nil {
		client.log.Error(err)
	}
	return err
}

func (client *mongoClient) Save(db string, collection string, _id string, m map[string]interface{}) error {
	c := client.session.DB(db).C(collection)
	_, err := c.Upsert(bson.M{"_id": _id}, bson.M{"$set": bson.M(m)})
	if err != nil {
		client.log.Error(err)
	}
	return err
}

func (client *mongoClient) SearchById(db string, collection string, _id string, structBsonPrt interface{}) error {
	c := client.session.DB(db).C(collection)
	err := c.Find(bson.M{"_id": _id}).One(structBsonPrt)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		client.log.Error(err)
		return err
	}
	return nil
}
