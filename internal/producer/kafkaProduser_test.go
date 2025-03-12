package producer

/*
import (
	"encoding/json"
	"testing"

	"github.com/mcu_service/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestParseData(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		expected    models.Record
		errExpected bool
		errValue    string
	}{
		{
			name:     "success parse data",
			data:     "1;temperature: 25, humidity: 50",
			expected: models.Record{ID: 1, Data: "temperature: 25, humidity: 50"},
		},
		{
			name:        "failed to parse data",
			data:        "1",
			expected:    models.Record{},
			errExpected: true,
			errValue:    "producer.parseData failed to parse data",
		},
		{
			name:        "failed to parse data, invalid id",
			data:        "ssa;",
			expected:    models.Record{},
			errExpected: true,
			errValue:    "producer.parseData strconv.Atoi: parsing \"ssa\": invalid syntax",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp, err := parseData(tc.data)
			if err != nil && tc.errExpected {
				assert.Equal(t, tc.errValue, err.Error())
			} else {
				assert.NoError(t, err)
				var record models.Record
				err = json.Unmarshal([]byte(resp), &record)
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, record)
			}
		})
	}

}
*/
