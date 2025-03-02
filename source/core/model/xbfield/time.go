package xbfield

import (
	"database/sql/driver"
	"time"

	"github.com/starryck/strk-tc-x-lib-go/source/core/toolkit/xbvalue"
	"github.com/starryck/strk-tc-x-lib-go/source/core/utility/xberror"
	"github.com/starryck/strk-tc-x-lib-go/source/core/utility/xbjson"
)

type UnixTime time.Time

func UnixTimeNow() *UnixTime {
	return (*UnixTime)(xbvalue.Refer(time.Now()))
}

func UnixTimeDate(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) *UnixTime {
	return (*UnixTime)(xbvalue.Refer(time.Date(year, month, day, hour, min, sec, nsec, loc)))
}

func UnixTimeItot(unix int64) *UnixTime {
	return (*UnixTime)(xbvalue.Refer(time.Unix(unix, 0)))
}

func UnixTimeTtoi(unixTime UnixTime) int64 {
	return time.Time(unixTime).Unix()
}

func (unixTime UnixTime) Value() (driver.Value, error) {
	zero := time.Time{}
	value := time.Time(unixTime)
	if value.UnixNano() == zero.UnixNano() {
		return nil, nil
	}
	return value, nil
}

func (unixTime *UnixTime) Scan(data any) error {
	var value time.Time
	switch data.(type) {
	case time.Time:
		value = data.(time.Time)
	case string, []byte:
		var date string
		if mDate, ok := data.(string); ok {
			date = mDate
		} else {
			date = string(data.([]byte))
		}
		if mValue, err := time.Parse(time.RFC3339, date); err != nil {
			return err
		} else {
			value = mValue
		}
	default:
		return xberror.Newf("Unsupported type of data `%#v`.", []any{data})
	}
	*unixTime = UnixTime(value)
	return nil
}

func (unixTime UnixTime) MarshalJSON() ([]byte, error) {
	return xbjson.Marshal(UnixTimeTtoi(unixTime))
}

func (unixTime *UnixTime) UnmarshalJSON(data []byte) error {
	var unix int64
	if err := xbjson.Unmarshal(data, &unix); err != nil {
		return err
	}
	*unixTime = *UnixTimeItot(unix)
	return nil
}

func (unixTime UnixTime) String() string {
	return time.Time(unixTime).String()
}
