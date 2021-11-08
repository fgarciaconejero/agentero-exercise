package repository_test

import (
	"database/sql"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agentero-exercise/agentero/repository"
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
			getPolicyHoldersSQL := `SELECT * FROM policy_holders`
			rows := sqlmock.NewRows([]string{"name", "ph_mobile_number"}).AddRow("some-name", "000000001")

			mock.ExpectQuery(regexp.QuoteMeta(getPolicyHoldersSQL)).WillReturnRows(rows)

			getInsurancePoliciesByIdSQL := `SELECT * FROM insurance_policies WHERE agent_id = ?`
			rows = sqlmock.NewRows([]string{"ip_id", "ip_mobile_number", "premium", "type", "agent_id"}).
				AddRow("some-ip-id", "000000001", 0, "some-type", "some-agent-id")

			mock.ExpectQuery(regexp.QuoteMeta(getInsurancePoliciesByIdSQL)).WillReturnRows(rows)
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
	// {
	// 	"getPolicyHoldersSQL will return error",
	// 	"some-agent-id",
	// 	nil,
	// 	errors.New("there was a problem while trying to get policy holders"),
	// },
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
		assert.Equal(t, res, tt.expectedResult)
	}

}
