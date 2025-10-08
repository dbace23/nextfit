package helper

import "database/sql"

func NullStringToInterface(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}

func NullTimeToInterface(nt sql.NullTime) interface{} {
	if nt.Valid {
		return nt.Time
	}
	return nil
}

func BoolToTinyInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
