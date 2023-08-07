package taxyfare

import (
	"reflect"
	"testing"
)

func TestAddRecord(t *testing.T) {

	type args struct {
		record string
	}
	tests := []struct {
		name    string
		args    args
		prepare func(m *manager)
		wantErr bool
	}{
		{
			"Add record success",
			args{
				"00:02:00.125 1141.2",
			},
			nil,
			false,
		},
		{
			"Add record error - elapsed time is earlier",
			args{
				"00:01:00.125 1141.2",
			},
			func(m *manager) {
				err := m.AddRecord("00:01:00.127 1141.2")
				if err != nil {
					t.Errorf("Error on preparation: %s", err.Error())
				}
			},
			true,
		},
		{
			"Add record error - elapsed time aparted more than 5 min from last time",
			args{
				"10:10:00.125 2000.2",
			},
			func(m *manager) {
				err := m.AddRecord("00:02:00.125 1141.2")
				if err != nil {
					t.Errorf("Error on preparation: %s", err.Error())
				}
			},
			true,
		},
		{
			"Add record error - distance smaller than last distance",
			args{
				"00:03:00.125 800.2",
			},
			func(m *manager) {
				err := m.AddRecord("00:02:00.125 1141.2")
				if err != nil {
					t.Errorf("Error on preparation: %s", err.Error())
				}
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &manager{}

			if tt.prepare != nil {
				tt.prepare(manager)
			}

			if err := manager.AddRecord(tt.args.record); (err != nil) != tt.wantErr {
				t.Errorf("AddRecord() = %v, want %v", err != nil, tt.wantErr)
			}
		})

	}
}

func TestGetCurrentFare(t *testing.T) {
	type args struct {
		record string
	}
	tests := []struct {
		name    string
		prepare func(m *manager)
		want    int
		wantErr bool
	}{
		{
			"Get Current Fare Success",
			func(m *manager) {
				err := m.AddRecord("00:01:00.125 400.2")
				if err != nil {
					t.Errorf("Error on preparation: %s", err.Error())
				}

				err = m.AddRecord("00:02:00.125 1800.2")
				if err != nil {
					t.Errorf("Error on preparation: %s", err.Error())
				}

				err = m.AddRecord("00:06:00.125 9400.2")
				if err != nil {
					t.Errorf("Error on preparation: %s", err.Error())
				}

				err = m.AddRecord("00:10:30.125 13500.2")
				if err != nil {
					t.Errorf("Error on preparation: %s", err.Error())
				}
			},
			1700,
			false,
		},
		{
			"Get Current Fare error - records is below minimum",
			func(m *manager) {
				err := m.AddRecord("00:00:00.125 1141.2")
				if err != nil {
					t.Errorf("Error on preparation: %s", err.Error())
				}
			},
			0,
			true,
		},
		{
			"Get Current Fare error - total mileage is 0.0",
			func(m *manager) {
				err := m.AddRecord("00:02:20.125 0.0")
				if err != nil {
					t.Errorf("Error on preparation: %s", err.Error())
				}

				err = m.AddRecord("00:05:00.123 0.0")
				if err != nil {
					t.Errorf("Error on preparation: %s", err.Error())
				}
			},
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &manager{}

			if tt.prepare != nil {
				tt.prepare(manager)
			}

			got, err := manager.GetCurrentFare()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() = %v, want %v", err != nil, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCurrentFare() got = %v, want %v", got, tt.want)
			}
		})

	}
}
