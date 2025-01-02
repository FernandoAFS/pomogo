package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Register a flag with an environment variable as default
func StringArg(flagName, envName, help, defValue string) *string {

	help_ := fmt.Sprintf(
		"%s, env. var: %s, default: %s",
		help,
		envName,
		defValue,
	)

	osVal, hasVal := os.LookupEnv(envName)

	if !hasVal {
		return flag.String(flagName, defValue, help_)
	}
	return flag.String(flagName, osVal, help_)
}

// Register a flag with an environment variable as default
func IntArg(flagName, envName, help string, defValue int) *int {

	help_ := fmt.Sprintf(
		"%s, env. var: %s, default: %s",
		help,
		envName,
		strconv.FormatInt(int64(defValue), 10),
	)

	osVal, hasVal := os.LookupEnv(envName)
	if !hasVal {
		return flag.Int(flagName, defValue, help_)
	}

	intVal, err := strconv.ParseInt(osVal, 10, strconv.IntSize)
	if err != nil {
		return nil
	}

	return flag.Int(flagName, int(intVal), help_)
}

func DurationArg(flagName, envName, help string, defValue time.Duration) *time.Duration {

	help_ := fmt.Sprintf(
		"%s, env. var: %s, default: %s",
		help,
		envName,
		defValue.String(),
	)

	osVal, hasVal := os.LookupEnv(envName)
	if !hasVal {
		return flag.Duration(flagName, defValue, help_)
	}

	intVal, err := time.ParseDuration(osVal)
	if err != nil {
		return nil
	}

	return flag.Duration(flagName, intVal, help_)
}
