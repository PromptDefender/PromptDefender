package wall

import (
	"encoding/json"
	"fmt"
	"github.com/safetorun/PromptDefender/cache"
	"github.com/safetorun/PromptDefender/utils"
)

func checkCache(cache *cache.Cache, prompt PromptToCheck) (bool, *CheckResult, error) {
	if cache == nil {
		println("Cache is nil")
		return false, nil, nil
	}

	b, err := json.Marshal(prompt)

	if err != nil {
		println("Error marshalling cache: ", err)
		return false, nil, err
	}

	cachedResult, err := (*cache).Get(utils.HashString(string(b)))

	if err != nil {
		println(fmt.Sprintf("Error getting cache: %v", err))
		return false, nil, err
	}

	if cachedResult != nil {
		var cachedResultReturn *CheckResult
		err := json.Unmarshal([]byte(*cachedResult), &cachedResultReturn)

		if err != nil {
			println("Error unmarshalling cache: ", err)
			return false, nil, err
		}

		return true, cachedResultReturn, nil
	}

	return false, nil, nil
}

func storeCache(cache *cache.Cache, prompt PromptToCheck, result *CheckResult) error {
	if cache == nil {
		return nil
	}

	b, err := json.Marshal(prompt)

	if err != nil {
		return err
	}

	bResult, err := json.Marshal(result)

	if err != nil {
		return err
	}

	return (*cache).Set(utils.HashString(string(b)), string(bResult))
}
