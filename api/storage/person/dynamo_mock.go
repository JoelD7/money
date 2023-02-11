package person

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	ForceNotFound   = false
	ForceUserExists = false

	ErrForceNotFound   = errors.New("force not found")
	ErrForceUserExists = errors.New("force user exists")
)

type MockDynamo struct{}

func InitDynamoMock() {
	Dynamo.Db = &MockDynamo{}
}

func (d *MockDynamo) GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	if ForceNotFound {
		return &dynamodb.GetItemOutput{}, ErrForceNotFound
	}

	return &dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String("test@gmail.com"),
			},
			"password": {
				// bcrypt hash for "1234"
				S: aws.String("$2a$10$.THF8QG33va8JTSIBz3lPuULaO6NiDb6yRmew63OtzujhVHbnZMFe"),
			},
			"full_name": {
				S: aws.String("Joel"),
			},
		},
	}, nil
}

func (d *MockDynamo) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	if ForceUserExists {
		return &dynamodb.PutItemOutput{}, ErrExistingUser
	}

	return &dynamodb.PutItemOutput{}, nil
}

func (d *MockDynamo) BatchExecuteStatement(input *dynamodb.BatchExecuteStatementInput) (*dynamodb.BatchExecuteStatementOutput, error) {
	return &dynamodb.BatchExecuteStatementOutput{}, nil
}

func (d *MockDynamo) BatchExecuteStatementWithContext(context aws.Context, input *dynamodb.BatchExecuteStatementInput, option ...request.Option) (*dynamodb.BatchExecuteStatementOutput, error) {
	return &dynamodb.BatchExecuteStatementOutput{}, nil
}

func (d *MockDynamo) BatchExecuteStatementRequest(input *dynamodb.BatchExecuteStatementInput) (*request.Request, *dynamodb.BatchExecuteStatementOutput) {
	return &request.Request{}, &dynamodb.BatchExecuteStatementOutput{}
}

func (d *MockDynamo) BatchGetItem(input *dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error) {
	return &dynamodb.BatchGetItemOutput{}, nil
}

func (d *MockDynamo) BatchGetItemWithContext(context aws.Context, input *dynamodb.BatchGetItemInput, option ...request.Option) (*dynamodb.BatchGetItemOutput, error) {
	return &dynamodb.BatchGetItemOutput{}, nil
}

func (d *MockDynamo) BatchGetItemRequest(input *dynamodb.BatchGetItemInput) (*request.Request, *dynamodb.BatchGetItemOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) BatchGetItemPages(input *dynamodb.BatchGetItemInput, f func(*dynamodb.BatchGetItemOutput, bool) bool) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) BatchGetItemPagesWithContext(context aws.Context, input *dynamodb.BatchGetItemInput, f func(*dynamodb.BatchGetItemOutput, bool) bool, option ...request.Option) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) BatchWriteItem(input *dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) BatchWriteItemWithContext(context aws.Context, input *dynamodb.BatchWriteItemInput, option ...request.Option) (*dynamodb.BatchWriteItemOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) BatchWriteItemRequest(input *dynamodb.BatchWriteItemInput) (*request.Request, *dynamodb.BatchWriteItemOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) CreateBackup(input *dynamodb.CreateBackupInput) (*dynamodb.CreateBackupOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) CreateBackupWithContext(context aws.Context, input *dynamodb.CreateBackupInput, option ...request.Option) (*dynamodb.CreateBackupOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) CreateBackupRequest(input *dynamodb.CreateBackupInput) (*request.Request, *dynamodb.CreateBackupOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) CreateGlobalTable(input *dynamodb.CreateGlobalTableInput) (*dynamodb.CreateGlobalTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) CreateGlobalTableWithContext(context aws.Context, input *dynamodb.CreateGlobalTableInput, option ...request.Option) (*dynamodb.CreateGlobalTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) CreateGlobalTableRequest(input *dynamodb.CreateGlobalTableInput) (*request.Request, *dynamodb.CreateGlobalTableOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) CreateTable(input *dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) CreateTableWithContext(context aws.Context, input *dynamodb.CreateTableInput, option ...request.Option) (*dynamodb.CreateTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) CreateTableRequest(input *dynamodb.CreateTableInput) (*request.Request, *dynamodb.CreateTableOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DeleteBackup(input *dynamodb.DeleteBackupInput) (*dynamodb.DeleteBackupOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DeleteBackupWithContext(context aws.Context, input *dynamodb.DeleteBackupInput, option ...request.Option) (*dynamodb.DeleteBackupOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DeleteBackupRequest(input *dynamodb.DeleteBackupInput) (*request.Request, *dynamodb.DeleteBackupOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DeleteItemWithContext(context aws.Context, input *dynamodb.DeleteItemInput, option ...request.Option) (*dynamodb.DeleteItemOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DeleteItemRequest(input *dynamodb.DeleteItemInput) (*request.Request, *dynamodb.DeleteItemOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DeleteTable(input *dynamodb.DeleteTableInput) (*dynamodb.DeleteTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DeleteTableWithContext(context aws.Context, input *dynamodb.DeleteTableInput, option ...request.Option) (*dynamodb.DeleteTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DeleteTableRequest(input *dynamodb.DeleteTableInput) (*request.Request, *dynamodb.DeleteTableOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeBackup(input *dynamodb.DescribeBackupInput) (*dynamodb.DescribeBackupOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeBackupWithContext(context aws.Context, input *dynamodb.DescribeBackupInput, option ...request.Option) (*dynamodb.DescribeBackupOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeBackupRequest(input *dynamodb.DescribeBackupInput) (*request.Request, *dynamodb.DescribeBackupOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeContinuousBackups(input *dynamodb.DescribeContinuousBackupsInput) (*dynamodb.DescribeContinuousBackupsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeContinuousBackupsWithContext(context aws.Context, input *dynamodb.DescribeContinuousBackupsInput, option ...request.Option) (*dynamodb.DescribeContinuousBackupsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeContinuousBackupsRequest(input *dynamodb.DescribeContinuousBackupsInput) (*request.Request, *dynamodb.DescribeContinuousBackupsOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeContributorInsights(input *dynamodb.DescribeContributorInsightsInput) (*dynamodb.DescribeContributorInsightsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeContributorInsightsWithContext(context aws.Context, input *dynamodb.DescribeContributorInsightsInput, option ...request.Option) (*dynamodb.DescribeContributorInsightsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeContributorInsightsRequest(input *dynamodb.DescribeContributorInsightsInput) (*request.Request, *dynamodb.DescribeContributorInsightsOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeEndpoints(input *dynamodb.DescribeEndpointsInput) (*dynamodb.DescribeEndpointsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeEndpointsWithContext(context aws.Context, input *dynamodb.DescribeEndpointsInput, option ...request.Option) (*dynamodb.DescribeEndpointsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeEndpointsRequest(input *dynamodb.DescribeEndpointsInput) (*request.Request, *dynamodb.DescribeEndpointsOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeExport(input *dynamodb.DescribeExportInput) (*dynamodb.DescribeExportOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeExportWithContext(context aws.Context, input *dynamodb.DescribeExportInput, option ...request.Option) (*dynamodb.DescribeExportOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeExportRequest(input *dynamodb.DescribeExportInput) (*request.Request, *dynamodb.DescribeExportOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeGlobalTable(input *dynamodb.DescribeGlobalTableInput) (*dynamodb.DescribeGlobalTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeGlobalTableWithContext(context aws.Context, input *dynamodb.DescribeGlobalTableInput, option ...request.Option) (*dynamodb.DescribeGlobalTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeGlobalTableRequest(input *dynamodb.DescribeGlobalTableInput) (*request.Request, *dynamodb.DescribeGlobalTableOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeGlobalTableSettings(input *dynamodb.DescribeGlobalTableSettingsInput) (*dynamodb.DescribeGlobalTableSettingsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeGlobalTableSettingsWithContext(context aws.Context, input *dynamodb.DescribeGlobalTableSettingsInput, option ...request.Option) (*dynamodb.DescribeGlobalTableSettingsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeGlobalTableSettingsRequest(input *dynamodb.DescribeGlobalTableSettingsInput) (*request.Request, *dynamodb.DescribeGlobalTableSettingsOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeImport(input *dynamodb.DescribeImportInput) (*dynamodb.DescribeImportOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeImportWithContext(context aws.Context, input *dynamodb.DescribeImportInput, option ...request.Option) (*dynamodb.DescribeImportOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeImportRequest(input *dynamodb.DescribeImportInput) (*request.Request, *dynamodb.DescribeImportOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeKinesisStreamingDestination(input *dynamodb.DescribeKinesisStreamingDestinationInput) (*dynamodb.DescribeKinesisStreamingDestinationOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeKinesisStreamingDestinationWithContext(context aws.Context, input *dynamodb.DescribeKinesisStreamingDestinationInput, option ...request.Option) (*dynamodb.DescribeKinesisStreamingDestinationOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeKinesisStreamingDestinationRequest(input *dynamodb.DescribeKinesisStreamingDestinationInput) (*request.Request, *dynamodb.DescribeKinesisStreamingDestinationOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeLimits(input *dynamodb.DescribeLimitsInput) (*dynamodb.DescribeLimitsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeLimitsWithContext(context aws.Context, input *dynamodb.DescribeLimitsInput, option ...request.Option) (*dynamodb.DescribeLimitsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeLimitsRequest(input *dynamodb.DescribeLimitsInput) (*request.Request, *dynamodb.DescribeLimitsOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeTable(input *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeTableWithContext(context aws.Context, input *dynamodb.DescribeTableInput, option ...request.Option) (*dynamodb.DescribeTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeTableRequest(input *dynamodb.DescribeTableInput) (*request.Request, *dynamodb.DescribeTableOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeTableReplicaAutoScaling(input *dynamodb.DescribeTableReplicaAutoScalingInput) (*dynamodb.DescribeTableReplicaAutoScalingOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeTableReplicaAutoScalingWithContext(context aws.Context, input *dynamodb.DescribeTableReplicaAutoScalingInput, option ...request.Option) (*dynamodb.DescribeTableReplicaAutoScalingOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeTableReplicaAutoScalingRequest(input *dynamodb.DescribeTableReplicaAutoScalingInput) (*request.Request, *dynamodb.DescribeTableReplicaAutoScalingOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeTimeToLive(input *dynamodb.DescribeTimeToLiveInput) (*dynamodb.DescribeTimeToLiveOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeTimeToLiveWithContext(context aws.Context, input *dynamodb.DescribeTimeToLiveInput, option ...request.Option) (*dynamodb.DescribeTimeToLiveOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DescribeTimeToLiveRequest(input *dynamodb.DescribeTimeToLiveInput) (*request.Request, *dynamodb.DescribeTimeToLiveOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DisableKinesisStreamingDestination(input *dynamodb.DisableKinesisStreamingDestinationInput) (*dynamodb.DisableKinesisStreamingDestinationOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DisableKinesisStreamingDestinationWithContext(context aws.Context, input *dynamodb.DisableKinesisStreamingDestinationInput, option ...request.Option) (*dynamodb.DisableKinesisStreamingDestinationOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) DisableKinesisStreamingDestinationRequest(input *dynamodb.DisableKinesisStreamingDestinationInput) (*request.Request, *dynamodb.DisableKinesisStreamingDestinationOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) EnableKinesisStreamingDestination(input *dynamodb.EnableKinesisStreamingDestinationInput) (*dynamodb.EnableKinesisStreamingDestinationOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) EnableKinesisStreamingDestinationWithContext(context aws.Context, input *dynamodb.EnableKinesisStreamingDestinationInput, option ...request.Option) (*dynamodb.EnableKinesisStreamingDestinationOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) EnableKinesisStreamingDestinationRequest(input *dynamodb.EnableKinesisStreamingDestinationInput) (*request.Request, *dynamodb.EnableKinesisStreamingDestinationOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ExecuteStatement(input *dynamodb.ExecuteStatementInput) (*dynamodb.ExecuteStatementOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ExecuteStatementWithContext(context aws.Context, input *dynamodb.ExecuteStatementInput, option ...request.Option) (*dynamodb.ExecuteStatementOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ExecuteStatementRequest(input *dynamodb.ExecuteStatementInput) (*request.Request, *dynamodb.ExecuteStatementOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ExecuteTransaction(input *dynamodb.ExecuteTransactionInput) (*dynamodb.ExecuteTransactionOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ExecuteTransactionWithContext(context aws.Context, input *dynamodb.ExecuteTransactionInput, option ...request.Option) (*dynamodb.ExecuteTransactionOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ExecuteTransactionRequest(input *dynamodb.ExecuteTransactionInput) (*request.Request, *dynamodb.ExecuteTransactionOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ExportTableToPointInTime(input *dynamodb.ExportTableToPointInTimeInput) (*dynamodb.ExportTableToPointInTimeOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ExportTableToPointInTimeWithContext(context aws.Context, input *dynamodb.ExportTableToPointInTimeInput, option ...request.Option) (*dynamodb.ExportTableToPointInTimeOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ExportTableToPointInTimeRequest(input *dynamodb.ExportTableToPointInTimeInput) (*request.Request, *dynamodb.ExportTableToPointInTimeOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) GetItemWithContext(context aws.Context, input *dynamodb.GetItemInput, option ...request.Option) (*dynamodb.GetItemOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) GetItemRequest(input *dynamodb.GetItemInput) (*request.Request, *dynamodb.GetItemOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ImportTable(input *dynamodb.ImportTableInput) (*dynamodb.ImportTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ImportTableWithContext(context aws.Context, input *dynamodb.ImportTableInput, option ...request.Option) (*dynamodb.ImportTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ImportTableRequest(input *dynamodb.ImportTableInput) (*request.Request, *dynamodb.ImportTableOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListBackups(input *dynamodb.ListBackupsInput) (*dynamodb.ListBackupsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListBackupsWithContext(context aws.Context, input *dynamodb.ListBackupsInput, option ...request.Option) (*dynamodb.ListBackupsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListBackupsRequest(input *dynamodb.ListBackupsInput) (*request.Request, *dynamodb.ListBackupsOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListContributorInsights(input *dynamodb.ListContributorInsightsInput) (*dynamodb.ListContributorInsightsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListContributorInsightsWithContext(context aws.Context, input *dynamodb.ListContributorInsightsInput, option ...request.Option) (*dynamodb.ListContributorInsightsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListContributorInsightsRequest(input *dynamodb.ListContributorInsightsInput) (*request.Request, *dynamodb.ListContributorInsightsOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListContributorInsightsPages(input *dynamodb.ListContributorInsightsInput, f func(*dynamodb.ListContributorInsightsOutput, bool) bool) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListContributorInsightsPagesWithContext(context aws.Context, input *dynamodb.ListContributorInsightsInput, f func(*dynamodb.ListContributorInsightsOutput, bool) bool, option ...request.Option) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListExports(input *dynamodb.ListExportsInput) (*dynamodb.ListExportsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListExportsWithContext(context aws.Context, input *dynamodb.ListExportsInput, option ...request.Option) (*dynamodb.ListExportsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListExportsRequest(input *dynamodb.ListExportsInput) (*request.Request, *dynamodb.ListExportsOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListExportsPages(input *dynamodb.ListExportsInput, f func(*dynamodb.ListExportsOutput, bool) bool) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListExportsPagesWithContext(context aws.Context, input *dynamodb.ListExportsInput, f func(*dynamodb.ListExportsOutput, bool) bool, option ...request.Option) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListGlobalTables(input *dynamodb.ListGlobalTablesInput) (*dynamodb.ListGlobalTablesOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListGlobalTablesWithContext(context aws.Context, input *dynamodb.ListGlobalTablesInput, option ...request.Option) (*dynamodb.ListGlobalTablesOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListGlobalTablesRequest(input *dynamodb.ListGlobalTablesInput) (*request.Request, *dynamodb.ListGlobalTablesOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListImports(input *dynamodb.ListImportsInput) (*dynamodb.ListImportsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListImportsWithContext(context aws.Context, input *dynamodb.ListImportsInput, option ...request.Option) (*dynamodb.ListImportsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListImportsRequest(input *dynamodb.ListImportsInput) (*request.Request, *dynamodb.ListImportsOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListImportsPages(input *dynamodb.ListImportsInput, f func(*dynamodb.ListImportsOutput, bool) bool) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListImportsPagesWithContext(context aws.Context, input *dynamodb.ListImportsInput, f func(*dynamodb.ListImportsOutput, bool) bool, option ...request.Option) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListTables(input *dynamodb.ListTablesInput) (*dynamodb.ListTablesOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListTablesWithContext(context aws.Context, input *dynamodb.ListTablesInput, option ...request.Option) (*dynamodb.ListTablesOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListTablesRequest(input *dynamodb.ListTablesInput) (*request.Request, *dynamodb.ListTablesOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListTablesPages(input *dynamodb.ListTablesInput, f func(*dynamodb.ListTablesOutput, bool) bool) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListTablesPagesWithContext(context aws.Context, input *dynamodb.ListTablesInput, f func(*dynamodb.ListTablesOutput, bool) bool, option ...request.Option) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListTagsOfResource(input *dynamodb.ListTagsOfResourceInput) (*dynamodb.ListTagsOfResourceOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListTagsOfResourceWithContext(context aws.Context, input *dynamodb.ListTagsOfResourceInput, option ...request.Option) (*dynamodb.ListTagsOfResourceOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ListTagsOfResourceRequest(input *dynamodb.ListTagsOfResourceInput) (*request.Request, *dynamodb.ListTagsOfResourceOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) PutItemWithContext(context aws.Context, input *dynamodb.PutItemInput, option ...request.Option) (*dynamodb.PutItemOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) PutItemRequest(input *dynamodb.PutItemInput) (*request.Request, *dynamodb.PutItemOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) QueryWithContext(context aws.Context, input *dynamodb.QueryInput, option ...request.Option) (*dynamodb.QueryOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) QueryRequest(input *dynamodb.QueryInput) (*request.Request, *dynamodb.QueryOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) QueryPages(input *dynamodb.QueryInput, f func(*dynamodb.QueryOutput, bool) bool) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) QueryPagesWithContext(context aws.Context, input *dynamodb.QueryInput, f func(*dynamodb.QueryOutput, bool) bool, option ...request.Option) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) RestoreTableFromBackup(input *dynamodb.RestoreTableFromBackupInput) (*dynamodb.RestoreTableFromBackupOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) RestoreTableFromBackupWithContext(context aws.Context, input *dynamodb.RestoreTableFromBackupInput, option ...request.Option) (*dynamodb.RestoreTableFromBackupOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) RestoreTableFromBackupRequest(input *dynamodb.RestoreTableFromBackupInput) (*request.Request, *dynamodb.RestoreTableFromBackupOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) RestoreTableToPointInTime(input *dynamodb.RestoreTableToPointInTimeInput) (*dynamodb.RestoreTableToPointInTimeOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) RestoreTableToPointInTimeWithContext(context aws.Context, input *dynamodb.RestoreTableToPointInTimeInput, option ...request.Option) (*dynamodb.RestoreTableToPointInTimeOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) RestoreTableToPointInTimeRequest(input *dynamodb.RestoreTableToPointInTimeInput) (*request.Request, *dynamodb.RestoreTableToPointInTimeOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ScanWithContext(context aws.Context, input *dynamodb.ScanInput, option ...request.Option) (*dynamodb.ScanOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ScanRequest(input *dynamodb.ScanInput) (*request.Request, *dynamodb.ScanOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ScanPages(input *dynamodb.ScanInput, f func(*dynamodb.ScanOutput, bool) bool) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) ScanPagesWithContext(context aws.Context, input *dynamodb.ScanInput, f func(*dynamodb.ScanOutput, bool) bool, option ...request.Option) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) TagResource(input *dynamodb.TagResourceInput) (*dynamodb.TagResourceOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) TagResourceWithContext(context aws.Context, input *dynamodb.TagResourceInput, option ...request.Option) (*dynamodb.TagResourceOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) TagResourceRequest(input *dynamodb.TagResourceInput) (*request.Request, *dynamodb.TagResourceOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) TransactGetItems(input *dynamodb.TransactGetItemsInput) (*dynamodb.TransactGetItemsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) TransactGetItemsWithContext(context aws.Context, input *dynamodb.TransactGetItemsInput, option ...request.Option) (*dynamodb.TransactGetItemsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) TransactGetItemsRequest(input *dynamodb.TransactGetItemsInput) (*request.Request, *dynamodb.TransactGetItemsOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) TransactWriteItems(input *dynamodb.TransactWriteItemsInput) (*dynamodb.TransactWriteItemsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) TransactWriteItemsWithContext(context aws.Context, input *dynamodb.TransactWriteItemsInput, option ...request.Option) (*dynamodb.TransactWriteItemsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) TransactWriteItemsRequest(input *dynamodb.TransactWriteItemsInput) (*request.Request, *dynamodb.TransactWriteItemsOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UntagResource(input *dynamodb.UntagResourceInput) (*dynamodb.UntagResourceOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UntagResourceWithContext(context aws.Context, input *dynamodb.UntagResourceInput, option ...request.Option) (*dynamodb.UntagResourceOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UntagResourceRequest(input *dynamodb.UntagResourceInput) (*request.Request, *dynamodb.UntagResourceOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateContinuousBackups(input *dynamodb.UpdateContinuousBackupsInput) (*dynamodb.UpdateContinuousBackupsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateContinuousBackupsWithContext(context aws.Context, input *dynamodb.UpdateContinuousBackupsInput, option ...request.Option) (*dynamodb.UpdateContinuousBackupsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateContinuousBackupsRequest(input *dynamodb.UpdateContinuousBackupsInput) (*request.Request, *dynamodb.UpdateContinuousBackupsOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateContributorInsights(input *dynamodb.UpdateContributorInsightsInput) (*dynamodb.UpdateContributorInsightsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateContributorInsightsWithContext(context aws.Context, input *dynamodb.UpdateContributorInsightsInput, option ...request.Option) (*dynamodb.UpdateContributorInsightsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateContributorInsightsRequest(input *dynamodb.UpdateContributorInsightsInput) (*request.Request, *dynamodb.UpdateContributorInsightsOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateGlobalTable(input *dynamodb.UpdateGlobalTableInput) (*dynamodb.UpdateGlobalTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateGlobalTableWithContext(context aws.Context, input *dynamodb.UpdateGlobalTableInput, option ...request.Option) (*dynamodb.UpdateGlobalTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateGlobalTableRequest(input *dynamodb.UpdateGlobalTableInput) (*request.Request, *dynamodb.UpdateGlobalTableOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateGlobalTableSettings(input *dynamodb.UpdateGlobalTableSettingsInput) (*dynamodb.UpdateGlobalTableSettingsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateGlobalTableSettingsWithContext(context aws.Context, input *dynamodb.UpdateGlobalTableSettingsInput, option ...request.Option) (*dynamodb.UpdateGlobalTableSettingsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateGlobalTableSettingsRequest(input *dynamodb.UpdateGlobalTableSettingsInput) (*request.Request, *dynamodb.UpdateGlobalTableSettingsOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateItemWithContext(context aws.Context, input *dynamodb.UpdateItemInput, option ...request.Option) (*dynamodb.UpdateItemOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateItemRequest(input *dynamodb.UpdateItemInput) (*request.Request, *dynamodb.UpdateItemOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateTable(input *dynamodb.UpdateTableInput) (*dynamodb.UpdateTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateTableWithContext(context aws.Context, input *dynamodb.UpdateTableInput, option ...request.Option) (*dynamodb.UpdateTableOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateTableRequest(input *dynamodb.UpdateTableInput) (*request.Request, *dynamodb.UpdateTableOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateTableReplicaAutoScaling(input *dynamodb.UpdateTableReplicaAutoScalingInput) (*dynamodb.UpdateTableReplicaAutoScalingOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateTableReplicaAutoScalingWithContext(context aws.Context, input *dynamodb.UpdateTableReplicaAutoScalingInput, option ...request.Option) (*dynamodb.UpdateTableReplicaAutoScalingOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateTableReplicaAutoScalingRequest(input *dynamodb.UpdateTableReplicaAutoScalingInput) (*request.Request, *dynamodb.UpdateTableReplicaAutoScalingOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateTimeToLive(input *dynamodb.UpdateTimeToLiveInput) (*dynamodb.UpdateTimeToLiveOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateTimeToLiveWithContext(context aws.Context, input *dynamodb.UpdateTimeToLiveInput, option ...request.Option) (*dynamodb.UpdateTimeToLiveOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) UpdateTimeToLiveRequest(input *dynamodb.UpdateTimeToLiveInput) (*request.Request, *dynamodb.UpdateTimeToLiveOutput) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) WaitUntilTableExists(input *dynamodb.DescribeTableInput) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) WaitUntilTableExistsWithContext(context aws.Context, input *dynamodb.DescribeTableInput, option ...request.WaiterOption) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) WaitUntilTableNotExists(input *dynamodb.DescribeTableInput) error {
	//TODO implement me
	panic("implement me")
}

func (d *MockDynamo) WaitUntilTableNotExistsWithContext(context aws.Context, input *dynamodb.DescribeTableInput, option ...request.WaiterOption) error {
	//TODO implement me
	panic("implement me")
}
