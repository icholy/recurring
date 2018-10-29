package recurring

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestIncludes(t *testing.T) {

	// yyyy/mm/dd
	layout := "2006/01/02"

	tests := []struct {
		name    string
		expr    TemporalExpression
		match   []string
		nomatch []string
	}{
		{
			name:    "day",
			expr:    Day(1),
			match:   []string{"2012/01/01", "2012/12/01", "2014/02/01"},
			nomatch: []string{"2012/01/02", "2016/07/20", "2018/08/08"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, st := range tt.match {
				mt, err := time.Parse(layout, st)
				assert.NilError(t, err)
				assert.Assert(t, tt.expr.Includes(mt))
			}
			for _, st := range tt.nomatch {
				mt, err := time.Parse(layout, st)
				assert.NilError(t, err)
				assert.Assert(t, !tt.expr.Includes(mt))
			}
		})
	}
}
