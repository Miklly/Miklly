package model

import (
	"common"
	"encoding/json"
	"reflect"
)

type base struct {
	Key string
}

func (this base) Save() error {
	db, err := bolt.Open(DataBaseFile, 0600, nil)
	defer db.Close()

	if err != nil {
		return err
	}
	var b []byte
	b, err = json.Marshal(this)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		n := reflect.TypeOf(this).Name()
		b, err := tx.CreateBucketIfNotExists([]byte(n))
		if err != nil {
			return fmt.Errorf("create %s: %s", n, err)
		}
		return b.Put([]byte(this.Key), b)
	})
	return err
}
