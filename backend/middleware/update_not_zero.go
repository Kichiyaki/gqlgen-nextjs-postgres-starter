package middleware

import (
	"reflect"

	"github.com/go-pg/pg/v9/orm"
)

type allowedZeroValues []string

func (zv allowedZeroValues) IsAllowed(name string) bool {
	for _, zvn := range zv {
		if zvn == name {
			return true
		}
	}
	return false
}

func updateNotZero(q *orm.Query, azv allowedZeroValues) (*orm.Query, error) {
	tableModel := q.TableModel()
	columns := []string{}
	if !tableModel.IsNil() {
		fields := tableModel.Table().DataFields
		if tableModel.Kind() == reflect.Struct {
			strct := tableModel.Value()
			for _, f := range fields {
				kind := f.Type.Kind()
				if kind == reflect.Bool || isNumber(kind) || !f.HasZeroValue(strct) || azv.IsAllowed(f.GoName) {
					columns = append(columns, string(f.SQLName))
				}
			}
		}
	}

	return q.Column(columns...), nil
}

func isNumber(kind reflect.Kind) bool {
	return kind == reflect.Int ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64
}

func UpdateNotZero(q *orm.Query) (*orm.Query, error) {
	return updateNotZero(q, allowedZeroValues{})
}

func UpdateNotZeroWithExceptions(azv []string) func(*orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		return updateNotZero(q, allowedZeroValues(azv))
	}
}
