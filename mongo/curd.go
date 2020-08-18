package mongo

import "context"

func Insert(collection string,record interface{}) error {
	collectionSelected := Client.Database(DataBase).Collection(collection)
	_, err := collectionSelected.InsertOne(context.TODO(),record)
	return err
}
