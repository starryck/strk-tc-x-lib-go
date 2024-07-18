package xbflow

import (
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/starryck/x-lib-go/source/core/base/xbconst"
	"github.com/starryck/x-lib-go/source/core/toolkit/xbjson"
	"github.com/starryck/x-lib-go/source/core/toolkit/xbrand"
	"github.com/starryck/x-lib-go/source/core/utility/xblogger"
)

type Flow interface {
	GetID() string
	GetTrails() []string
}

type BaseFlow struct {
	storage *sync.Map
}

type Operate func(args ...any) error

func (flow *BaseFlow) Initiate() {
	flow.storage = &sync.Map{}
	flow.storage.Store(xbconst.FlowKeyFlowID, xbrand.MakeKSUID())
	flow.storage.Store(xbconst.FlowKeyFlowTrails, []string{})
	flow.storage.Store(xbconst.FlowKeyFlowError, nil)
	flow.storage.Store(xbconst.FlowKeyFlowOutcome, nil)
	return
}

func (flow *BaseFlow) Inherit(fore Flow) {
	flow.storage = &sync.Map{}
	flow.storage.Store(xbconst.FlowKeyFlowID, fore.GetID())
	flow.storage.Store(xbconst.FlowKeyFlowTrails, append(fore.GetTrails(), xbrand.MakeBase62String(8)))
	flow.storage.Store(xbconst.FlowKeyFlowError, nil)
	flow.storage.Store(xbconst.FlowKeyFlowOutcome, nil)
	return
}

func (flow *BaseFlow) GetID() string {
	id, _ := flow.storage.Load(xbconst.FlowKeyFlowID)
	return id.(string)
}

func (flow *BaseFlow) GetTrails() []string {
	trails, _ := flow.storage.Load(xbconst.FlowKeyFlowTrails)
	return trails.([]string)
}

func (flow *BaseFlow) HasError() bool {
	err, _ := flow.storage.Load(xbconst.FlowKeyFlowError)
	return err != nil
}

func (flow *BaseFlow) GetError() error {
	err, _ := flow.storage.Load(xbconst.FlowKeyFlowError)
	if err == nil {
		return nil
	} else {
		return err.(error)
	}
}

func (flow *BaseFlow) SetError(err error) {
	flow.storage.Store(xbconst.FlowKeyFlowError, err)
	return
}

func (flow *BaseFlow) GetOutcome() any {
	outcome, _ := flow.storage.Load(xbconst.FlowKeyFlowOutcome)
	return outcome
}

func (flow *BaseFlow) SetOutcome(outcome any) {
	flow.storage.Store(xbconst.FlowKeyFlowOutcome, outcome)
	return
}

func (flow *BaseFlow) GetStorage() *sync.Map {
	return flow.storage
}

func (flow *BaseFlow) SetStorage(storage *sync.Map) {
	flow.storage = storage
	return
}

func (flow *BaseFlow) GetLogger() *xblogger.Entry {
	entry := xblogger.WithFields(xblogger.Fields{
		"flowID":     flow.GetID(),
		"flowTrails": fmt.Sprintf("/%s", strings.Join(flow.GetTrails(), "/")),
	})
	return entry
}

func (flow *BaseFlow) Contain(key string) bool {
	_, ok := flow.storage.Load(key)
	return ok
}

func (flow *BaseFlow) Expose(key string, value any) {
	flow.storage.Store(key, value)
	return
}

func (flow *BaseFlow) Require(key string) any {
	value, ok := flow.storage.Load(key)
	if !ok {
		panic(fmt.Sprintf("The key `%s` must exist in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireBool(key string) bool {
	value, ok := flow.Require(key).(bool)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of bool type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireInt(key string) int {
	value, ok := flow.Require(key).(int)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of int type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireInt64(key string) int64 {
	value, ok := flow.Require(key).(int64)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of int64 type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireUint(key string) uint {
	value, ok := flow.Require(key).(uint)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of uint type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireUint64(key string) uint64 {
	value, ok := flow.Require(key).(uint64)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of uint64 type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireFloat64(key string) float64 {
	value, ok := flow.Require(key).(float64)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of float64 type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireString(key string) string {
	value, ok := flow.Require(key).(string)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of string type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireBools(key string) []bool {
	value, ok := flow.Require(key).([]bool)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of []bool type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireBytes(key string) []byte {
	value, ok := flow.Require(key).([]byte)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of []byte type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireInts(key string) []int {
	value, ok := flow.Require(key).([]int)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of []int type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireInt64s(key string) []int64 {
	value, ok := flow.Require(key).([]int64)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of []int64 type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireFloat64s(key string) []float64 {
	value, ok := flow.Require(key).([]float64)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of []float64 type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireStrings(key string) []string {
	value, ok := flow.Require(key).([]string)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of []string type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireIntMap(key string) map[int]any {
	value, ok := flow.Require(key).(map[int]any)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of map[int]any type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireIntBoolMap(key string) map[int]bool {
	value, ok := flow.Require(key).(map[int]bool)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of map[int]bool type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireIntIntMap(key string) map[int]int {
	value, ok := flow.Require(key).(map[int]int)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of map[int]int type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireIntStringMap(key string) map[int]string {
	value, ok := flow.Require(key).(map[int]string)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of map[int]string type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireStringMap(key string) map[string]any {
	value, ok := flow.Require(key).(map[string]any)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of map[string]any type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireStringBoolMap(key string) map[string]bool {
	value, ok := flow.Require(key).(map[string]bool)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of map[string]bool type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireStringIntMap(key string) map[string]int {
	value, ok := flow.Require(key).(map[string]int)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of map[string]int type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireStringStringMap(key string) map[string]string {
	value, ok := flow.Require(key).(map[string]string)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of map[string]string type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireError(key string) error {
	value, ok := flow.Require(key).(error)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of error type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireTime(key string) time.Time {
	value, ok := flow.Require(key).(time.Time)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of time.Time type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireDuration(key string) time.Duration {
	value, ok := flow.Require(key).(time.Duration)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of time.Duration type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) RequireRawMessage(key string) xbjson.RawMessage {
	value, ok := flow.Require(key).(xbjson.RawMessage)
	if !ok {
		panic(fmt.Sprintf("The value of key `%s` must be of json.RawMessage type in flow storage.", key))
	}
	return value
}

func (flow *BaseFlow) Async(operate Operate, args ...any) {
	go func() {
		defer func() {
			if v := recover(); v != nil {
				flow.GetLogger().WithField(xblogger.PanicKey, xblogger.FormatPanic(v, debug.Stack())).Error("Flow async operation panicked.")
			}
		}()
		if err := operate(args...); err != nil {
			flow.GetLogger().WithError(err).Error("Flow async operation failed.")
		}
	}()
}
