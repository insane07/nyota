package store

import (
	"encoding/json"
	"errors"
	"nyota/backend/model/config"

	"github.com/lib/pq"

	gorp "gopkg.in/gorp.v2"
)

type PgStringArray []string

// gorpTypeConverter is used by Gorp to Encode/Decode custom types in DB.
type gorpTypeConverter struct{}

// ToDb converts a Prizm object to one suitable for the DB representation.
func (tc gorpTypeConverter) ToDb(val interface{}) (interface{}, error) {
	switch t := val.(type) {
	case PgStringArray:
		return pq.Array(val), nil

	case []string,
		config.ExtraParam, map[string]interface{}, map[string]string:

		js, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}
		return string(js), nil
	default:
		return val, nil
	}
}

// FromDb converts a DB representation back into a Prizm object.
func (tc gorpTypeConverter) FromDb(target interface{}) (gorp.CustomScanner, bool) {
	switch target.(type) {
	case *PgStringArray:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New("gorp(type converter) unable to convert *string")
			}

			var array pq.StringArray
			if err := array.Scan(*s); err != nil {
				return err
			}

			value := target.(*PgStringArray)
			*value = PgStringArray([]string(array))
			return nil
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true

	case *[]string,
		*config.ExtraParam, *map[string]interface{}, *map[string]string:
		binder := func(holder, target interface{}) error {
			js, ok := holder.(*string)
			if !ok {
				return errors.New("gorp(type converter) unable to convert *string")
			}
			return json.Unmarshal([]byte(*js), target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	default:
		return gorp.CustomScanner{}, false
	}
}
