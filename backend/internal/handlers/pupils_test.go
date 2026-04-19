package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple_project/internal/models"

	"github.com/stretchr/testify/require"
)

func TestApplicationGetPupils(t *testing.T) {
	tests := []struct {
		name       string
		db         *fakeDB
		wantStatus int
		assertBody func(t *testing.T, body *bytes.Buffer)
	}{
		{
			name: "returns pupils",
			db: &fakeDB{pupils: []models.Pupil{{
				ID: 3, Name: "Petr", Surname: "Sidorov", ParentName: "Anna", ParentPhone: "+79990000000",
			}}},
			wantStatus: http.StatusOK,
			assertBody: func(t *testing.T, body *bytes.Buffer) {
				var response []models.Pupil
				require.NoError(t, json.NewDecoder(body).Decode(&response))
				require.Len(t, response, 1)
				require.Equal(t, "Sidorov", response[0].Surname)
			},
		},
		{
			name:       "returns pupils db error",
			db:         &fakeDB{getPupilsErr: errDB},
			wantStatus: http.StatusInternalServerError,
			assertBody: func(t *testing.T, body *bytes.Buffer) {
				require.Contains(t, body.String(), errDB.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{DB: tt.db}
			req := httptest.NewRequest(http.MethodGet, "/pupils", nil)
			rr := httptest.NewRecorder()

			app.GetPupils(rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)
			tt.assertBody(t, rr.Body)
		})
	}
}

func TestApplicationPupilMutations(t *testing.T) {
	pupil := models.Pupil{ID: 8, Name: "Ivan", Surname: "Smirnov", ParentName: "Olga", ParentPhone: "+79991112233"}

	tests := []struct {
		name       string
		method     string
		body       string
		db         *fakeDB
		call       func(app *Application, w http.ResponseWriter, req *http.Request)
		wantStatus int
		assertDB   func(t *testing.T, db *fakeDB)
	}{
		{
			name:       "add pupil",
			method:     http.MethodPost,
			body:       mustJSON(t, pupil),
			db:         &fakeDB{},
			call:       (*Application).AddPupil,
			wantStatus: http.StatusOK,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.NotNil(t, db.addedPupil)
				require.Equal(t, pupil.Surname, db.addedPupil.Surname)
			},
		},
		{
			name:       "update pupil",
			method:     http.MethodPut,
			body:       mustJSON(t, pupil),
			db:         &fakeDB{},
			call:       (*Application).UpdatePupil,
			wantStatus: http.StatusOK,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.NotNil(t, db.updatedPupil)
				require.Equal(t, pupil.ID, db.updatedPupil.ID)
			},
		},
		{
			name:       "delete pupil",
			method:     http.MethodDelete,
			body:       mustJSON(t, pupil),
			db:         &fakeDB{},
			call:       (*Application).DeletePupil,
			wantStatus: http.StatusOK,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.NotNil(t, db.deletedPupil)
				require.Equal(t, pupil.ID, db.deletedPupil.ID)
			},
		},
		{
			name:       "bad json does not call db",
			method:     http.MethodPost,
			body:       `{"name":`,
			db:         &fakeDB{},
			call:       (*Application).AddPupil,
			wantStatus: http.StatusBadRequest,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.Nil(t, db.addedPupil)
			},
		},
		{
			name:       "db error",
			method:     http.MethodDelete,
			body:       mustJSON(t, pupil),
			db:         &fakeDB{deletePupilErr: errDB},
			call:       (*Application).DeletePupil,
			wantStatus: http.StatusInternalServerError,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.NotNil(t, db.deletedPupil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{DB: tt.db}
			req := httptest.NewRequest(tt.method, "/pupils", bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()

			tt.call(app, rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)
			tt.assertDB(t, tt.db)
		})
	}
}
