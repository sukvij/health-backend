package query

import (
	"fmt"

	"gorm.io/gorm"
)

// creates new record in database
func CreateNewRecord(db *gorm.DB, val interface{}) (interface{}, error) {
	err := db.Create(val).Error
	if err != nil {
		return nil, err
	}
	return val, nil
}

func FirstRecordWithPrimaryKey(db *gorm.DB, model interface{}) (interface{}, error) {
	if model == nil {
		return nil, fmt.Errorf("model cannot be nil")
	}

	// Create a new instance of the same type as model to store the result
	// Since model is a pointer, we assume it's a struct pointer (e.g., *User)
	result := model

	err := db.Where(model).First(result).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find record: %v", err)
	}

	return result, nil
}

// Find finds all records matching given conditions conds
func FindAllRecordsWithoutCondition(db *gorm.DB, val interface{}) (interface{}, error) {
	err := db.Find(&val).Error
	if err != nil {
		return nil, err
	}
	return val, nil
}

// it will take primary key from model and update non null values from object val
// if model dont have primary key then it updates all the entries of model table  --> db.Model(&User{}).Updates()
// UpdateRecord updates an existing record or creates a new one if not found.
func UpdateRecord(db *gorm.DB, model interface{}, val interface{}) (interface{}, error) {
	// Find existing record
	existing, err := FirstRecordWithPrimaryKey(db, model)
	fmt.Printf("UpdateRecord: FirstRecordWithPrimaryKey returned err: %v\n", err)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new record if not found
			return CreateNewRecord(db, val)
		}
		return nil, fmt.Errorf("cannot find record to update: %v", err)
	}

	// Update the existing record with values from val
	err = db.Model(existing).Updates(val).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update record: %v", err)
	}

	return existing, nil
}

//S ave will save all fields when performing the Updating SQL and If the value contains no primary key, it performs Create
// it pushes new updates
// func Save(db *gorm.DB, val interface{}) (interface{}, error) {
// 	var result interface{}
// 	err := db.Save(val).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }
