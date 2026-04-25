package other

import "os"

func GetEnv(key string) string {
	a := os.Getenv(key)
	if a == "" {
		panic(".env中缺少" + key)
	}
	return a
}
