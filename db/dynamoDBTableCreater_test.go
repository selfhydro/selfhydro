package db

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func Test_ShouldCreateTableNameWithDate(t *testing.T) {
	time := time.Date(2019, 8, 27, 0, 0, 0, 0, time.UTC)
	tableName := getTableName(time)
	assert.Equal(t, tableName, "selfhydro-state-2019-08-27")

}
