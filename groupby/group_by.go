package groupby

import (
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/agflow/tools/log"
	"github.com/agflow/tools/security"
	"github.com/agflow/tools/typing"
)

func getGroupHash(v reflect.Value, cols []string) (string, error) {
	grouped := make([]interface{}, len(cols))
	for j := range cols {
		grouped[j] = v.FieldByName(cols[j]).Interface()
	}
	return security.Hash(grouped)
}

// AggrFunc is an aggregation function
type AggrFunc func([]interface{}) interface{}

var wg sync.WaitGroup //nolint: gochecknoglobals

// Agg groups by a slice of structs with an aggregation function
func Agg(on, dest interface{}, cols []string, funcs map[string]AggrFunc) error {
	groupMap := make(map[interface{}]chan interface{})
	if reflect.TypeOf(on).Kind() != reflect.Slice {
		return errors.New("on needs to be slice")
	}

	destVal := reflect.ValueOf(dest)
	direct := reflect.Indirect(destVal)

	sliceType, err := typing.Base(destVal.Type(), reflect.Slice)
	if err != nil {
		return err
	}

	baseType := typing.DeRef(sliceType.Elem())

	s := reflect.ValueOf(on)
	finishChan := make(chan bool)
	for i := 0; i < s.Len(); i++ {
		v := reflect.Indirect(s.Index(i))

		key, err := getGroupHash(v, cols)
		if err != nil {
			return err
		}

		if _, ok := groupMap[key]; !ok {
			emptyAggVal := reflect.New(baseType)
			for j := range cols {
				newField := reflect.Indirect(emptyAggVal).FieldByName(cols[j])
				if newField.CanSet() {
					newField.Set(v.FieldByName(cols[j]))
				}
			}

			wg.Add(1)

			valChan := make(chan interface{})
			go processFun(valChan, funcs, emptyAggVal, finishChan)

			groupMap[key] = valChan
			direct.Set(reflect.Append(direct, emptyAggVal))
		}

		if valChan, ok := groupMap[key]; ok {
			valChan <- v.Interface()
		}
	}

	for i := 0; i < len(groupMap); i++ {
		finishChan <- true
	}

	wg.Wait()

	return nil
}

func processFun(
	v chan interface{},
	funcs map[string]AggrFunc,
	dest reflect.Value,
	finish chan bool,
) {
	defer wg.Done()
	var isFinished bool
	grouped := make([]interface{}, 0)
	for !isFinished {
		select {
		case res := <-v:
			grouped = append(grouped, res)
		case isFinished = <-finish:
		case <-time.After(3 * time.Second):
			log.Error("groupby process timeout")
			return
		}
	}
	for col, f := range funcs {
		field := reflect.Indirect(dest).FieldByName(col)
		field.Set(reflect.ValueOf(f(grouped)))
	}
}

// FoldFunc is a fold function
type FoldFunc func(interface{}, interface{}) interface{}

// Fold groups by a slice of structs with a fold function
func Fold(on, dest interface{}, cols []string, funcs map[string]FoldFunc) error {
	groupMap := make(map[interface{}]int)
	if reflect.TypeOf(on).Kind() != reflect.Slice {
		return errors.New("on needs to be slice")
	}

	destVal := reflect.ValueOf(dest)
	direct := reflect.Indirect(destVal)

	sliceType, err := typing.Base(destVal.Type(), reflect.Slice)
	if err != nil {
		return err
	}

	baseType := typing.DeRef(sliceType.Elem())

	s := reflect.ValueOf(on)
	for i := 0; i < s.Len(); i++ {
		v := reflect.Indirect(s.Index(i))

		key, err := getGroupHash(v, cols)
		if err != nil {
			return err
		}

		if _, ok := groupMap[key]; !ok {
			emptyAggVal := reflect.New(baseType)
			for j := range cols {
				newField := reflect.Indirect(emptyAggVal).FieldByName(cols[j])
				if newField.CanSet() {
					newField.Set(v.FieldByName(cols[j]))
				}
			}
			direct.Set(reflect.Append(direct, emptyAggVal))
			groupMap[key] = direct.Len() - 1
		}

		if ix, ok := groupMap[key]; ok {
			dest := direct.Index(ix)
			for col, f := range funcs {
				field := reflect.Indirect(dest).FieldByName(col)
				result := f(field.Interface(), v.Addr().Interface())
				field.Set(reflect.ValueOf(result))
			}
		}
	}
	return nil
}
