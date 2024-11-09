package income

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/tests/e2e/utils"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := env.LoadEnvTesting()
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestGetAllIncomePeriods(t *testing.T) {
	c := require.New(t)

	var (
		incomeTableName       = env.GetString("INCOME_TABLE_NAME", "")
		periodUserIncomeIndex = env.GetString("PERIOD_USER_INCOME_INDEX", "")
	)

	dynamoClient := utils.InitDynamoClient()
	ctx := context.Background()
	username := "e2e_test@mail.com"

	incomeRepo, err := income.NewDynamoRepository(dynamoClient, incomeTableName, periodUserIncomeIndex)
	c.Nil(err, "failed to create income repository")

	setupIncome(ctx, "samples/income.json", c, incomeRepo, t)

	incomePeriods, err := incomeRepo.GetAllIncomePeriods(ctx, username)
	c.Nil(err, "failed to get all income periods")
	c.Len(incomePeriods, 3, "unexpected number of income periods")
}

func setupIncome(ctx context.Context, file string, c *require.Assertions, incomeRepo income.Repository, t *testing.T) {
	data, err := os.ReadFile(file)
	c.Nil(err, "reading income sample file failed")

	var incomeList []*models.Income
	err = json.Unmarshal(data, &incomeList)

	c.Len(incomeList, 10, "unexpected number of income in the sample file")

	err = incomeRepo.BatchCreateIncome(ctx, incomeList)
	c.Nil(err, "batch creating income failed")

	defer t.Cleanup(func() {
		err = incomeRepo.BatchDeleteIncome(ctx, incomeList)
		c.Nil(err, "batch deleting income failed")
	})
}
