package mongo

import (
	"labix.org/v2/mgo"
	"github.com/murz/ego/db"
	"fmt"
)

type MongoDriver struct {
	session *mgo.Session
	cfg *db.Config
}

func (d *MongoDriver) Initialize(cfg *db.Config) {
	session, err := mgo.Dial(
		fmt.Sprintf("%v:%v@%v:%v/%v",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DBName))
	if (err != nil) {
		panic(err)
	}
	d.session = session
	d.cfg = cfg
}

func (d *MongoDriver) Dispose() {
	if d.session != nil {
		d.session.Close()
	}
}