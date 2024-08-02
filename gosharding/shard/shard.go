package shard

import (
	"github.com/google/uuid"
	"hash/crc32"
)

func (s *Shard) MustBe(wantShardIdx int, id string) (string, error) {
	idx := s.Index(id)

	if idx != uint32(wantShardIdx) {
		newID, err := uuid.NewV7()
		if err != nil {
			return "", err
		}

		return s.MustBe(wantShardIdx, newID.String())
	}
	return id, nil
}

func (s *Shard) Index(id string) uint32 {
	hash := fnvHash(id)
	return hash % uint32(s.Count)
}

func fnvHash(id string) uint32 {
	const prime = 16777619
	hash := uint32(2166136261)
	//hash := uint64(18446744073709551615)

	for i := 0; i < len(id); i++ {
		hash ^= uint32(id[i])
		hash *= prime
	}

	return hash
}

func getcrc32(userID string) uint32 {
	ff := crc32.ChecksumIEEE([]byte(userID))
	return ff % 3
}

type Shard struct {
	Count int
}

func NewShard(count int) *Shard {
	return &Shard{
		Count: count,
	}
}
