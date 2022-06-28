package main

import (
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/agflow/tools/log"
	"github.com/agflow/tools/types"
)

func Hash(vs interface{}) (string, error) {
	h := sha512.New()
	r, err := json.Marshal(vs)
	if err != nil {
		return "", err
	}
	h.Write(r)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func getGroupHash(v reflect.Value, cols []string) (string, error) {
	grouped := make([]interface{}, len(cols))
	for j := range cols {
		grouped[j] = v.FieldByName(cols[j]).Interface()
	}
	return Hash(grouped)
}

type AggrFunc func([]interface{}) interface{}

var wg sync.WaitGroup //nolint: gochecknoglobals

// GroupBy groups a slice of structs with an aggregation function
func GroupBy(on, dest interface{}, cols []string, funcs map[string]AggrFunc) error { //nolint: deadcode
	groupMap := make(map[interface{}]chan interface{})
	if reflect.TypeOf(on).Kind() != reflect.Slice {
		return errors.New("on needs to be slice")
	}

	destVal := reflect.ValueOf(dest)
	direct := reflect.Indirect(destVal)

	sliceType, err := types.BaseType(destVal.Type(), reflect.Slice)
	if err != nil {
		return err
	}

	baseType := types.DeRef(sliceType.Elem())

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

func processFun(v chan interface{}, funcs map[string]AggrFunc, dest reflect.Value, finish chan bool) {
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

type FoldFunc func(interface{}, interface{}) interface{}

// GroupByFold groups a slice of structs with a fold function
func GroupByFold(on, dest interface{}, cols []string, funcs map[string]FoldFunc) error { //nolint: deadcode
	groupMap := make(map[interface{}]int)
	if reflect.TypeOf(on).Kind() != reflect.Slice {
		return errors.New("on needs to be slice")
	}

	destVal := reflect.ValueOf(dest)
	direct := reflect.Indirect(destVal)

	sliceType, err := types.BaseType(destVal.Type(), reflect.Slice)
	if err != nil {
		return err
	}

	baseType := types.DeRef(sliceType.Elem())

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
