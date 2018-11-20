package ops

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"
	"time"

	"goprizm/log"
)

// retries given operation atmost n times if it returns error.
// Retry is performed by backingoff exponentially after each failure.
// If n is -1, retry is performed till op succeeds.
// https://en.wikipedia.org/wiki/Exponential_backoff
// https://developers.google.com/admin-sdk/directory/v1/limits#backoff
func Retry(opName string, op func() error, n int) (err error) {
	// random milli seconds 0-1
	randMsecs := func() (int, error) {
		nBig, err := crand.Int(crand.Reader, big.NewInt(1000))
		if err != nil {
			return 0, err
		}

		n := nBig.Int64()
		return int(n), nil
	}

	after := func(secs int) {
		mSecs, err := randMsecs()
		if err != nil {
			log.Errorf("%s(retry) after(%d secs): %v", opName, secs, err)
			return
		}

		duration := time.Duration(time.Duration(secs*1000+mSecs) * time.Millisecond)
		time.Sleep(duration)
	}

	i := 0
	for {
		// If n is +ve and num of retries reached n.
		if n > 0 && i == n {
			return err
		}

		// Execute op and on error sleep for random duration.
		if err = op(); err != nil {
			i += 1
			log.Errorf("%s(retry:%d) failed err:%v", opName, i, err)
			// limit back off duration to 5 mins
			after(int(math.Min(300, math.Exp2(float64(i)))))
			continue
		}
		return nil
	}
}

// Shard maps given key to one of the partitions.
func Shard(key string, numShards int) int {
	checkSum := sha256.Sum256([]byte(key))
	hashKey := checkSum[len(checkSum)-4:]
	return int(binary.LittleEndian.Uint32(hashKey) % uint32(numShards))
}
