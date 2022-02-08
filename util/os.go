package util

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

var (
	err  error
	CORS string
)

// LoadEnv loads OS environment variables
func LoadEnv() error {
	if CLIENT_ID = os.Getenv("CLIENT_ID"); CLIENT_ID == "" {
		err = errors.New("CLIENT_ID is not set")
		ErrorLogger.Println(err)
		Status = "Error: CLIENT_ID is not set"
		return err
	}

	if CLIENT_SECRET = os.Getenv("CLIENT_SECRET"); CLIENT_SECRET == "" {
		err = errors.New("CLIENT_SECRET is not set")
		ErrorLogger.Println(err)
		Status = "Error: CLIENT_SECRET is not set"
		return err
	}

	if CALLBACK_URL = os.Getenv("CALLBACK_URL"); CALLBACK_URL == "" {
		WarningLogger.Println("CALLBACK_URL is not set")
		CALLBACK_URL = "http://localhost:3000"
		InfoLogger.Printf("CALLBACK_URL set using default: %s", CALLBACK_URL)
	} else {
		InfoLogger.Printf("CALLBACK_URL set using env: %s", CALLBACK_URL)
	}

	if CORS = os.Getenv("CORS"); CORS == "" {
		WarningLogger.Println("CORS is not set")
		CORS = "*"
		InfoLogger.Printf("CORS set using default: %s", CORS)
	} else {
		InfoLogger.Printf("CORS set using env: %s", CORS)
	}

	if SECRET_KEY = os.Getenv("SECRET_KEY"); SECRET_KEY == "" {
		WarningLogger.Println("SECRET_KEY is not set")
		// setting random key
		SECRET_KEY = RandomStringBytes(32)
		InfoLogger.Printf("SECRET_KEY is not set, using random key: %s", SECRET_KEY)
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

	if SLACK_TOKEN = os.Getenv("SLACK_TOKEN"); SLACK_TOKEN == "" {
		WarningLogger.Println("SLACK_TOKEN is not set")
		SLACK_TOKEN = ""
	}

	if BroadCastChannelID = os.Getenv("SLACK_BROADCAST_CHANNEL_ID"); BroadCastChannelID == "" {
		WarningLogger.Println("SLACK_BROADCAST_CHANNEL_ID is not set")
		BroadCastChannelID = ""
	}

	// get every env variable starting with STORAGE_COST_ and parse it to STORAGE_COST with the storage class name after STORAGE_COST_ as key
	storageClasses := []string{}
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
				Status = "Error: " + err.Error()
				return err
			}
			// add to tempStorageCost
			tempStorageCost[key[2]] = map[string]float64{"cost": value}
			InfoLogger.Printf("storage class %s set to cost value: %f", key[2], value)
			storageClasses = append(storageClasses, key[2])
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

	if QUOTA_CPU_LABEL = os.Getenv("QUOTA_CPU_LABEL"); QUOTA_CPU_LABEL == "" {
		WarningLogger.Println("QUOTA_CPU_LABEL is not set")
		QUOTA_CPU_LABEL = "natron.io/cpu-quota"
		InfoLogger.Printf("QUOTA_CPU_LABEL set using default: %s", QUOTA_CPU_LABEL)
	} else {
		InfoLogger.Printf("QUOTA_CPU_LABEL set using env: %s", QUOTA_CPU_LABEL)
	}

	if QUOTA_MEMORY_LABEL = os.Getenv("QUOTA_MEMORY_LABEL"); QUOTA_MEMORY_LABEL == "" {
		WarningLogger.Println("QUOTA_MEMORY_LABEL is not set")
		QUOTA_MEMORY_LABEL = "natron.io/memory-quota"
		InfoLogger.Printf("QUOTA_MEMORY_LABEL set using default: %s", QUOTA_MEMORY_LABEL)
	} else {
		InfoLogger.Printf("QUOTA_MEMORY_LABEL set using env: %s", QUOTA_MEMORY_LABEL)
	}

	if QUOTA_NAMESPACE_SUFFIX = os.Getenv("QUOTA_NAMESPACE_SUFFIX"); QUOTA_NAMESPACE_SUFFIX == "" {
		WarningLogger.Println("QUOTA_NAMESPACE_SUFFIX is not set")
		QUOTA_NAMESPACE_SUFFIX = "config"
		InfoLogger.Printf("QUOTA_NAMESPACE_SUFFIX set using default: %s", QUOTA_NAMESPACE_SUFFIX)
	} else {
		InfoLogger.Printf("QUOTA_NAMESPACE_SUFFIX set using env: %s", QUOTA_NAMESPACE_SUFFIX)
	}

	// for each storageclass get the quota label
	QUOTA_STORAGE_LABEL = make(map[string]string)
	// get storage class names from env
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "QUOTA_STORAGE_LABEL_") {
			// split env variable to key and value
			keyValue := strings.Split(env, "=")
			// split key to storage class name and cost
			key := strings.Split(keyValue[0], "_")
			// add to tempStorageCost

			// check if storage class already exists in storageClasses
			for _, storageClass := range storageClasses {
				if storageClass == key[3] {
					break
				}
			}
			storageClasses = append(storageClasses, key[3])
		}
	}

	for _, storageClass := range storageClasses {
		label := os.Getenv("QUOTA_STORAGE_LABEL_" + storageClass)
		if label == "" {
			WarningLogger.Printf("QUOTA_STORAGE_LABEL_%s is not set", storageClass)
			label = "natron.io/storage-quota-" + storageClass
			InfoLogger.Printf("QUOTA_STORAGE_LABEL_%s set using default: %s", storageClass, label)
			// add to QUOTA_STORAGE_LABEL
			QUOTA_STORAGE_LABEL[storageClass] = label
		} else {
			InfoLogger.Printf("QUOTA_STORAGE_LABEL_%s set using env: %s", storageClass, label)
			QUOTA_STORAGE_LABEL[storageClass] = label
		}
	}

	return nil
}
