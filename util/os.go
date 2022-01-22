package util

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

var (
	err error
)

func LoadEnv() error {
	if CLIENT_ID = os.Getenv("CLIENT_ID"); CLIENT_ID == "" {
		err = errors.New("CLIENT_ID is not set")
		ErrorLogger.Println(err)
		return err
	}

	if CLIENT_SECRET = os.Getenv("CLIENT_SECRET"); CLIENT_SECRET == "" {
		err = errors.New("CLIENT_SECRET is not set")
		ErrorLogger.Println(err)
		return err
	}

	if CALLBACK_URL = os.Getenv("CALLBACK_URL"); CALLBACK_URL == "" {
		WarningLogger.Println("CALLBACK_URL is not set")
		CALLBACK_URL = "http://localhost:3000"
		InfoLogger.Printf("CALLBACK_URL set using default: %s", CALLBACK_URL)
	}

	if SECRET_KEY = os.Getenv("SECRET_KEY"); SECRET_KEY == "" {
		WarningLogger.Println("SECRET_KEY is not set")
		// setting random key
		SECRET_KEY = RandomStringBytes(32)
		InfoLogger.Printf("SECRET_KEY is not set, using random key: %s", SECRET_KEY)
	}

	if LABELSELECTOR = os.Getenv("LABELSELECTOR"); LABELSELECTOR == "" {
		WarningLogger.Println("LABELSELECTOR is not set")
		LABELSELECTOR = "natron.io/tenant"
		InfoLogger.Printf("LABELSELECTOR set using default: %s", LABELSELECTOR)
	}

	if CPU_COST, err = strconv.ParseFloat(os.Getenv("CPU_COST"), 64); CPU_COST == 0 || err != nil {
		WarningLogger.Println("CPU_COST is not set or invalid float value")
		CPU_COST = 1.00
		InfoLogger.Printf("CPU_COST set using default: %f", CPU_COST)
	}

	if MEMORY_COST, err = strconv.ParseFloat(os.Getenv("MEMORY_COST"), 64); MEMORY_COST == 0 || err != nil {
		WarningLogger.Println("MEMORY_COST is not set or invalid float value")
		MEMORY_COST = 1.00
		InfoLogger.Printf("MEMORY_COST set using default: %f", MEMORY_COST)
	}

	// get every env variable starting with STORAGE_COST_ and parse it to STORAGE_COST with the storage class name after STORAGE_COST_ as key
	tempStorageCost := make(map[string]float64)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "STORAGE_COST_") {
			// split env variable to key and value
			keyValue := strings.Split(env, "=")
			// split key to storage class name and cost
			storageClassCost := strings.Split(keyValue[0], "_")
			// get cost
			cost, err := strconv.ParseFloat(keyValue[1], 64)
			if err != nil {
				WarningLogger.Printf("Invalid float value for %s", keyValue[0])
				return err
			}
			// add storage class name and cost
			tempStorageCost[storageClassCost[1]] = cost

			InfoLogger.Printf("Added storage class %s with cost %f", storageClassCost[2], cost)
		}
	}
	// set STORAGE_COST
	STORAGE_COST = tempStorageCost

	if STORAGE_COST == nil {
		InfoLogger.Println("No storage class cost set")

		// add default storage class cost
		STORAGE_COST = make(map[string]float64)
		STORAGE_COST["default"] = 1.00
		InfoLogger.Printf("Added default storage class with cost %f", STORAGE_COST["default"])
	}

	return nil
}

func GetStatus() string {
	if err != nil {
		return "error"
	}
	return "ok"
}
