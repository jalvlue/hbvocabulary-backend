package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testStore *Store

func init() {
	testDSN := "user_course_design:123456@tcp(127.0.0.1:3306)/db_course_design?charset=utf8mb4&parseTime=True&loc=Local"
	InitDB(testDSN)
	testStore = NewStore(DB)
}

func TestUser(t *testing.T) {

	u := &User{
		Username:  "jm",
		Password:  "123456",
		MaxScore:  0,
		TestCount: 0,
	}

	err := testStore.CreateUser(u)
	require.NoError(t, err)
	getUser, err := testStore.GetUserByUsername(u.Username)
	require.NoError(t, err)

	require.Equal(t, u.Username, getUser.Username)
	require.Equal(t, u.Password, getUser.Password)
	require.Equal(t, u.MaxScore, getUser.MaxScore)
	require.Equal(t, u.TestCount, getUser.TestCount)
}
