package functioncalling

import (
	"sync"
)

var (
	once     sync.Once
	instance *FunctionCallingRegistry
)

type FunctionCallingRegistry struct {
	functionCallingsMap map[string]FunctionCallingManager
}

func GetRegistry() *FunctionCallingRegistry {
	once.Do(func() {
		instance = &FunctionCallingRegistry{
			functionCallingsMap: make(map[string]FunctionCallingManager),
		}
	})
	return instance
}

func (r *FunctionCallingRegistry) RegisterFunctionCalling(functionCalling FunctionCallingManager) {
	r.functionCallingsMap[functionCalling.GetID()] = functionCalling
}

func (r *FunctionCallingRegistry) GetFunctionCalling(id string) (FunctionCallingManager, bool) {
	functionCalling, ok := r.functionCallingsMap[id]
	return functionCalling, ok
}

func RegisterFunctionCalling(functionCalling FunctionCallingManager) {
	if functionCalling.GetDisabled() {
		return
	}

	GetRegistry().RegisterFunctionCalling(functionCalling)
}

func GetFunctionCalling(id string) (FunctionCallingManager, bool) {
	return GetRegistry().GetFunctionCalling(id)
}
