package helpers

import "ecommerce/database"

var DBHelper *database.Mongo

func Push(){
	DBHelper.Insert()
}