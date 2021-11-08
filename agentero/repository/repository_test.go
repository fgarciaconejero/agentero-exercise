package repository_test

import (
	"database/sql"
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agentero-exercise/agentero/repository"
	"github.com/agentero-exercise/agentero/resources/constants"
	"github.com/agentero-exercise/agentero/resources/protos"
	"github.com/stretchr/testify/assert"
)

func NewMockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

var getByIdTestingParameters = []struct {
	name            string
	id              string
	sqlExpectations func(sqlmock.Sqlmock)
	expectedResult  []*protos.PolicyHolder
	err             error
}{
	{
		"successful",
		"some-agent-id",
		func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"name", "ph_mobile_number"}).AddRow("some-name", "000000001")
			mock.ExpectQuery(regexp.QuoteMeta(constants.GetPolicyHoldersSQL)).WillReturnRows(rows)

			rows = sqlmock.NewRows([]string{"ip_id", "ip_mobile_number", "premium", "type", "agent_id"}).
				AddRow("some-ip-id", "000000001", 0, "some-type", "some-agent-id")
			mock.ExpectQuery(regexp.QuoteMeta(constants.GetInsurancePoliciesByIdSQL)).WillReturnRows(rows)
		},
		[]*protos.PolicyHolder{
			{
				Name:         "some-name",
				MobileNumber: "000000001",
				InsurancePolicy: []*protos.InsurancePolicy{
					{
						MobileNumber: "000000001",
						Premium:      0,
						Type:         "some-type",
						AgentId:      "some-agent-id",
					},
				},
			},
		},
		nil,
	},
	{
		"getPolicyHoldersSQL will return error",
		"some-agent-id",
		func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery(regexp.QuoteMeta(constants.GetPolicyHoldersSQL)).WillReturnError(errors.New("there was a problem while trying to get policy holders"))
		},
		nil,
		errors.New("there was a problem while trying to get policy holders"),
	},
	{
		"getInsurancePoliciesByIdSQL will return error",
		"some-agent-id",
		func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"name", "ph_mobile_number"}).AddRow("some-name", "000000001")
			mock.ExpectQuery(regexp.QuoteMeta(constants.GetPolicyHoldersSQL)).WillReturnRows(rows)

			mock.ExpectQuery(regexp.QuoteMeta(constants.GetInsurancePoliciesByIdSQL)).WillReturnError(errors.New("there was a problem while trying to get insurance policies"))
		},
		nil,
		errors.New("there was a problem while trying to get insurance policies"),
	},
}

func TestGetById(t *testing.T) {
	db, mock := NewMockDB()
	r := &repository.Repository{Db: *db}
	defer r.Db.Close()
	for _, tt := range getByIdTestingParameters {
		tt.sqlExpectations(mock)
		res, err := r.GetById(tt.id)
		if tt.err != nil {
			assert.EqualError(t, err, tt.err.Error())
		}
		assert.Equal(t, tt.expectedResult, res)
	}
}

var getByMobileNumberTestingParameters = []struct {
	name            string
	mobileNumber    string
	sqlExpectations func(sqlmock.Sqlmock)
	expectedResult  *protos.PolicyHolder
	err             error
}{
	{
		"successful",
		"000000001",
		func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"name", "ph_mobile_number"}).AddRow("some-name", "000000001")
			mock.ExpectQuery(regexp.QuoteMeta(constants.GetPolicyHoldersSQL)).WillReturnRows(rows)

			rows = sqlmock.NewRows([]string{"ip_id", "ip_mobile_number", "premium", "type", "agent_id"}).
				AddRow("some-ip-id", "000000001", 0, "some-type", "some-agent-id")
			mock.ExpectQuery(regexp.QuoteMeta(constants.GetInsurancePoliciesByMobileNumberSQL)).WillReturnRows(rows)
		},
		&protos.PolicyHolder{
			Name:         "some-name",
			MobileNumber: "000000001",
			InsurancePolicy: []*protos.InsurancePolicy{
				{
					MobileNumber: "000000001",
					Premium:      0,
					Type:         "some-type",
					AgentId:      "some-agent-id",
				},
			},
		},
		nil,
	},
	{
		"getPolicyHoldersSQL will return error",
		"some-agent-id",
		func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery(regexp.QuoteMeta(constants.GetPolicyHoldersSQL)).WillReturnError(errors.New("there was a problem while trying to get policy holders"))
		},
		nil,
		errors.New("there was a problem while trying to get policy holders"),
	},
	{
		"getInsurancePoliciesByMobileNumberSQL will return error",
		"some-agent-id",
		func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"name", "ph_mobile_number"}).AddRow("some-name", "000000001")
			mock.ExpectQuery(regexp.QuoteMeta(constants.GetPolicyHoldersSQL)).WillReturnRows(rows)

			mock.ExpectQuery(regexp.QuoteMeta(constants.GetInsurancePoliciesByMobileNumberSQL)).WillReturnError(errors.New("there was a problem while trying to get insurance policies"))
		},
		nil,
		errors.New("there was a problem while trying to get insurance policies"),
	},
}

func TestGetByMobileNumber(t *testing.T) {
	db, mock := NewMockDB()
	r := &repository.Repository{Db: *db}
	defer r.Db.Close()
	for _, tt := range getByMobileNumberTestingParameters {
		tt.sqlExpectations(mock)
		res, err := r.GetByMobileNumber(tt.mobileNumber)
		if tt.err != nil {
			assert.EqualError(t, err, tt.err.Error())
		}
		assert.Equal(t, tt.expectedResult, res)
	}
}

var getAllInsuranceAgentsIdsTestingParameters = []struct {
	name            string
	sqlExpectations func(sqlmock.Sqlmock)
	expectedResult  []string
	err             error
}{
	{
		"success",
		func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"agent_id", "name"}).AddRow("some-agent-id", "some-name")
			mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllInsuranceAgentsSQL)).WillReturnRows(rows)
		},
		[]string{"some-agent-id"},
		nil,
	},
	{
		"getAllInsuranceAgentsSQL returns an error",
		func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllInsuranceAgentsSQL)).WillReturnError(errors.New("GetAllInsuranceAgentsSQL returned an error"))
		},
		nil,
		errors.New("GetAllInsuranceAgentsSQL returned an error"),
	},
}

func TestGetAllInsuranceAgentsIds(t *testing.T) {
	db, mock := NewMockDB()
	r := &repository.Repository{Db: *db}
	defer r.Db.Close()
	for _, tt := range getAllInsuranceAgentsIdsTestingParameters {
		tt.sqlExpectations(mock)
		res, err := r.GetAllInsuranceAgentsIds()
		if tt.err != nil {
			assert.EqualError(t, err, tt.err.Error())
		}
		assert.Equal(t, tt.expectedResult, res)
	}
}

var upsertPolicyHolderTestingParameters = []struct {
	name            string
	policyHolder    *protos.PolicyHolder
	sqlExpectations func(sqlmock.Sqlmock)
	err             error
}{
	{
		"success",
		&protos.PolicyHolder{
			MobileNumber: "000000001",
			Name:         "some-name",
		},
		func(mock sqlmock.Sqlmock) {
			result := sqlmock.NewResult(0, 0)
			mock.ExpectExec(regexp.QuoteMeta(constants.InsertPolicyHolderSQL)).WillReturnResult(result)
		},
		nil,
	},
	{
		"insertPolicyHolderSQL returns error",
		&protos.PolicyHolder{
			MobileNumber: "000000001",
			Name:         "some-name",
		},
		func(mock sqlmock.Sqlmock) {
			mock.ExpectExec(regexp.QuoteMeta(constants.InsertPolicyHolderSQL)).WillReturnError(errors.New("insertPolicyHolderSQL returned an error"))
		},
		errors.New("insertPolicyHolderSQL returned an error"),
	},
}

func TestUpsertPolicyHolder(t *testing.T) {
	db, mock := NewMockDB()
	r := &repository.Repository{Db: *db}
	defer r.Db.Close()

	for _, tt := range upsertPolicyHolderTestingParameters {
		tt.sqlExpectations(mock)
		err := r.UpsertPolicyHolder(tt.policyHolder)
		if tt.err != nil {
			assert.EqualError(t, err, tt.err.Error())
		}
	}
}

var upsertInsurancePolicyTestingParameters = []struct {
	name            string
	agentId         string
	insurancePolicy *protos.InsurancePolicy
	sqlExpectations func(sqlmock.Sqlmock)
	err             error
}{
	{
		"success",
		"some-agent-id",
		&protos.InsurancePolicy{
			MobileNumber: "000000001",
			Premium:      0,
			Type:         "some-type",
		},
		func(mock sqlmock.Sqlmock) {
			result := sqlmock.NewResult(0, 0)
			mock.ExpectExec(regexp.QuoteMeta(constants.InsertInsurancePolicySQL)).WillReturnResult(result)
		},
		nil,
	},
	{
		"insertPolicyHolderSQL returns error",
		"some-agent-id",
		&protos.InsurancePolicy{
			MobileNumber: "000000001",
			Premium:      0,
			Type:         "some-type",
		},
		func(mock sqlmock.Sqlmock) {
			mock.ExpectExec(regexp.QuoteMeta(constants.InsertInsurancePolicySQL)).WillReturnError(errors.New("insertInsurancePolicySQL returned an error"))
		},
		errors.New("insertInsurancePolicySQL returned an error"),
	},
}

func TestUpsertInsurancePolicy(t *testing.T) {
	db, mock := NewMockDB()
	r := &repository.Repository{Db: *db}
	defer r.Db.Close()

	for _, tt := range upsertInsurancePolicyTestingParameters {
		tt.sqlExpectations(mock)
		err := r.UpsertInsurancePolicy(tt.insurancePolicy, tt.agentId)
		if tt.err != nil {
			assert.EqualError(t, err, tt.err.Error())
		}
	}
}
