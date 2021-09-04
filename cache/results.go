package cache

import (
	"strconv"
	"time"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/stringHelper"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//AddResult adds a result to the database
func AddResult(Result models.Result) error {
	//logger.Debug(loggingArea, "Starting transaction for result of item", Result.ItemID)
	tx, err := db.Begin()
	if err != nil {
		logger.Error(loggingArea, "Couldn't start transaction for result:", err)
		return err
	}

	_, err = tx.Exec(`INSERT INTO results (ItemID, CapturedAt, ValueString, ValueInt, Error) VALUES (?,?,?,?,?)`, Result.ItemID.Hex(), Result.CapturedAt, Result.ValueString, Result.ValueNumeric, Result.Error)
	if err != nil {
		logger.Error(loggingArea, "Couldn't exec query for result:", err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(loggingArea, "Couldn't commit transaction for result:", err)
		tx.Rollback()
	}
	//logger.Debug(loggingArea, "Result captured at", Result.CapturedAt, "persisted")
	return nil
}

//RetrieveCache returns all cached results and deltes them from the cache
func RetrieveCache() ([]models.Result, error) {
	logger.Debug(loggingArea, "Retrieving cached results")

	results := make([]models.Result, 0)

	tx, err := db.Begin()
	if err != nil {
		logger.Error(loggingArea, "Couldn't start transaction for cached results:", err)
		return results, err
	}

	rows, err := tx.Query(`SELECT ItemID,CapturedAt,ValueString,ValueInt,Error FROM results`)
	if err != nil {
		logger.Error(loggingArea, "Couldn't retrieve cached results:", err)
		return results, err
	}

	for rows.Next() {
		var cachedResult models.Result
		var capturedAtString string
		var itemIDString string
		var valueNumericString string

		err := rows.Scan(&itemIDString, &capturedAtString, &cachedResult.ValueString, &valueNumericString, &cachedResult.Error)
		if err != nil {
			//Dont abort here so the corrupted result doesnt block the result flow
			logger.Error(loggingArea, "Couldn't decode cached results:", err)
			continue
		}

		cachedResult.CapturedAt, _ = time.Parse("2006-01-02 15:04:05-07:00", capturedAtString)

		//Check if ValueInt was set (Empty if not)
		if !stringHelper.IsEmpty(valueNumericString) {
			cachedResult.ValueNumeric, _ = strconv.ParseFloat(valueNumericString, 64)
		}

		cachedResult.ItemID, _ = primitive.ObjectIDFromHex(itemIDString)

		results = append(results, cachedResult)
	}

	rows.Close()
	logger.Debug(loggingArea, "Retrieved", len(results), "cached results")
	logger.Debug(loggingArea, "Removing cached results")

	_, err = tx.Exec(`DELETE FROM results`)
	if err != nil {
		logger.Error(loggingArea, "Couldn't remove cached results:", err)
		return results, err
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(loggingArea, "Couldn't commit transaction for cached results:", err)
		return results, err
	}

	return results, nil
}
