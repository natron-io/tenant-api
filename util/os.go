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
		Status = "Error"
		return err
	}

	if CLIENT_SECRET = os.Getenv("CLIENT_SECRET"); CLIENT_SECRET == "" {
		err = errors.New("CLIENT_SECRET is not set")
		ErrorLogger.Println(err)
		Status = "Error"
		return err
	}

	if CALLBACK_URL = os.Getenv("CALLBACK_URL"); CALLBACK_URL == "" {
		WarningLogger.Println("CALLBACK_URL is not set")
		CALLBACK_URL = "http://localhost:3000"
		InfoLogger.Printf("CALLBACK_URL set using default: %s", CALLBACK_URL)
	} else {
		InfoLogger.Printf("CALLBACK_URL set using env: %s", CALLBACK_URL)
	}

	if SECRET_KEY = os.Getenv("SECRET_KEY"); SECRET_KEY == "" {
		WarningLogger.Println("SECRET_KEY is not set")
		// setting random key
		SECRET_KEY = RandomStringBytes(32)
		InfoLogger.Printf("SECRET_KEY is not set, using random key: %s", SECRET_KEY)
	}

	if FRONTENDAUTH_ENABLED, err = strconv.ParseBool(os.Getenv("FRONTENDAUTH_ENABLED")); err != nil {
		WarningLogger.Println("FRONTENDAUTH_ENABLED is not set")
		FRONTENDAUTH_ENABLED = false
		InfoLogger.Printf("FRONTENDAUTH_ENABLED set using default: %t", FRONTENDAUTH_ENABLED)
	} else {
		InfoLogger.Printf("FRONTENDAUTH_ENABLED set using env: %t", FRONTENDAUTH_ENABLED)
	}

	if TENANT_LABEL = os.Getenv("TENANT_LABEL"); TENANT_LABEL == "" {
		WarningLogger.Println("TENANT_LABEL is not set")
		TENANT_LABEL = "natron.io/tenant"
		InfoLogger.Printf("TENANT_LABEL set using default: %s", TENANT_LABEL)
	} else {
		InfoLogger.Printf("TENANT_LABEL set using env: %s", TENANT_LABEL)
	}

	if DISCOUNT_LABEL = os.Getenv("DISCOUNT_LABEL"); DISCOUNT_LABEL == "" {
		WarningLogger.Println("DISCOUNT_LABEL is not set")
		DISCOUNT_LABEL = "natron.io/discount"
		InfoLogger.Printf("DISCOUNT_LABEL set using default: %s", DISCOUNT_LABEL)
	} else {
		InfoLogger.Printf("DISCOUNT_LABEL set using env: %s", DISCOUNT_LABEL)
	}

	if CPU_COST, err = strconv.ParseFloat(os.Getenv("CPU_COST"), 64); CPU_COST == 0 || err != nil {
		WarningLogger.Println("CPU_COST is not set or invalid float value")
		CPU_COST = 1.00
		InfoLogger.Printf("CPU_COST set using default: %f", CPU_COST)
	} else {
		InfoLogger.Printf("CPU_COST set using env: %f", CPU_COST)
	}

	if MEMORY_COST, err = strconv.ParseFloat(os.Getenv("MEMORY_COST"), 64); MEMORY_COST == 0 || err != nil {
		WarningLogger.Println("MEMORY_COST is not set or invalid float value")
		MEMORY_COST = 1.00
		InfoLogger.Printf("MEMORY_COST set using default: %f", MEMORY_COST)
	} else {
		InfoLogger.Printf("MEMORY_COST set using env: %f", MEMORY_COST)
	}

	// get every env variable starting with STORAGE_COST_ and parse it to STORAGE_COST with the storage class name after STORAGE_COST_ as key
	tempStorageCost := make(map[string]map[string]float64)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "STORAGE_COST_") {
			// split env variable to key and value
			keyValue := strings.Split(env, "=")
			// split key to storage class name and cost
			key := strings.Split(keyValue[0], "_")
			// parse value to float
			value, err := strconv.ParseFloat(keyValue[1], 64)
			if err != nil {
				err = errors.New("STORAGE_COST_" + key[2] + " is not set or invalid float value")
				ErrorLogger.Println(err)
				Status = "Error"
				return err
			}
			// add to tempStorageCost
			tempStorageCost[key[2]] = map[string]float64{"cost": value}
			InfoLogger.Printf("storage class %s set to cost value: %f", key[2], value)
		}
	}
	STORAGE_COST = tempStorageCost

	if STORAGE_COST == nil {
		WarningLogger.Println("STORAGE_COST is not set")
		STORAGE_COST = map[string]map[string]float64{
			"default": {"cost": 1.00},
		}
		InfoLogger.Printf("cost for storage class default set using default: %f", STORAGE_COST["default"]["cost"])
	}

	if INGRESS_COST, err = strconv.ParseFloat(os.Getenv("INGRESS_COST"), 64); INGRESS_COST == 0 || err != nil {
		WarningLogger.Println("INGRESS_COST is not set or invalid float value")
		INGRESS_COST = 1.00
		InfoLogger.Printf("INGRESS_COST set using default: %f", INGRESS_COST)
	} else {
		InfoLogger.Printf("INGRESS_COST set to: %f", INGRESS_COST)
	}

	return nil
}
