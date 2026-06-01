package store

import (
	"log"
	"time"
)

func expireSample(s *Store) float32 {
	var limit int = 20
	var expiredCount int = 0

	// iterate over 20 keys that have expiration set
	for key, entry := range s.data {
		if entry.HasExpiry {
			limit --

			if time.Now().After(entry.ExpiresAt){
				// log.Printf("Key %s has expired and will be deleted\n", key)
				s.Delete(key)
				expiredCount ++
			}
		}
		
		if (limit == 0){
			break
		}
	}

	return  float32(expiredCount) / float32(20.0)
}

func DeleteExpiredKeys(s *Store){
	for {
		fraction := expireSample(s)
		if fraction < 0.25 {
			break
		}
	}

	log.Println("After deleting expired keys, total keys :", len(s.data))
}