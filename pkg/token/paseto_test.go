package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/workspace/evoting/ev-webservice/pkg/faker"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(faker.RandomString(32))
	require.NoError(t, err)

	username := faker.RandomUser()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, false, duration)
	// fmt.Println(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Data)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(faker.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(faker.RandomUser(), false, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
