package scalar

import (
	"fmt"
	"io"
	"time"
)

type DateTime time.Time

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (x *DateTime) UnmarshalGQL(v any) error {
	t, ok := v.(string)
	if !ok {
		return fmt.Errorf("date_time must be a string")
	}
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return err
	}
	*x = DateTime(parsed)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (x DateTime) MarshalGQL(w io.Writer) {
	t := (time.Time)(x)
	w.Write([]byte(fmt.Sprintf(`"%s"`, t.Format(time.RFC3339))))
}
