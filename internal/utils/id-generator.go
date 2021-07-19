package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/chilts/sid"

	guuid "github.com/google/uuid"
	"github.com/kjk/betterguid"
	"github.com/lithammer/shortuuid"
	"github.com/oklog/ulid"
	"github.com/rs/xid"

	"github.com/segmentio/ksuid"
	"github.com/sony/sonyflake"
)

func GenShortUUID() string {
	id := shortuuid.New()
	return id
}

func GenUUID() string {
	id := guuid.New()
	return id.String()
}

func GenXid() string {
	id := xid.New()
	return id.String()
}

func GenKsuid() string {
	id := ksuid.New()
	return id.String()
}

func GenBetterGUID() string {
	id := betterguid.New()
	return id
}

func GenUlid() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	return id.String()
}

func GenSonyflake() (uint64, error) {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		return id, err
	}
	return id, nil
}

func GenSid() {
	id := sid.Id()
	fmt.Printf("github.com/chilts/sid:          %s\n", id)
}

func UniqueDigits() uint64 {
	value := uint64(10000000000 + rand.Intn(99999999999-10000000000))
	return value
}
