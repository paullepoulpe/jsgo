package rhino

import (
	"fmt"
	"sync"
	"testing"
)

type parallelRaceNoopReporter struct{}

func (parallelRaceNoopReporter) Warning(string, string, int, int) {}
func (parallelRaceNoopReporter) Error(string, string, int, int)   {}

func TestDetectInheritanceCycleConcurrentReadOnly(t *testing.T) {
	reg := NewJSTypeRegistry(parallelRaceNoopReporter{})

	// 1. Deep acyclic prototype chain: a[0] -> a[1] -> ... -> a[31] -> Object.prototype.
	const chainLen = 32
	acyclic := make([]*PrototypeObjectTypeBase, chainLen)
	prev := reg.GetNativeObjectType(JSTypeNative_OBJECT_TYPE)
	for i := chainLen - 1; i >= 0; i-- {
		b := NewPrototypeObjectType_Builder(reg)
		b.SetName(fmt.Sprintf("Acyclic_%d", i))
		b.SetImplicitPrototype(prev)
		acyclic[i] = NewPrototypeObjectType(b)
		prev = acyclic[i]
	}

	// 2. Direct self-cycle: selfCyc -> selfCyc.
	bSelf := NewPrototypeObjectType_Builder(reg)
	bSelf.SetName("SelfCycle")
	selfCyc := NewPrototypeObjectType(bSelf)
	selfCyc.SetImplicitPrototype(selfCyc)

	// 3. 3-hop cycle: cycA -> cycB -> cycC -> cycA.
	bA := NewPrototypeObjectType_Builder(reg)
	bA.SetName("CycleA")
	cycA := NewPrototypeObjectType(bA)
	bB := NewPrototypeObjectType_Builder(reg)
	bB.SetName("CycleB")
	cycB := NewPrototypeObjectType(bB)
	bC := NewPrototypeObjectType_Builder(reg)
	bC.SetName("CycleC")
	cycC := NewPrototypeObjectType(bC)
	cycA.SetImplicitPrototype(cycB)
	cycB.SetImplicitPrototype(cycC)
	cycC.SetImplicitPrototype(cycA)

	// 4. Lasso (rho-shaped) cycle: tail0 -> tail1 -> loop0 -> loop1 -> loop0.
	bLoop0 := NewPrototypeObjectType_Builder(reg)
	bLoop0.SetName("Loop0")
	loop0 := NewPrototypeObjectType(bLoop0)
	bLoop1 := NewPrototypeObjectType_Builder(reg)
	bLoop1.SetName("Loop1")
	loop1 := NewPrototypeObjectType(bLoop1)
	loop0.SetImplicitPrototype(loop1)
	loop1.SetImplicitPrototype(loop0)

	bTail1 := NewPrototypeObjectType_Builder(reg)
	bTail1.SetName("Tail1")
	bTail1.SetImplicitPrototype(loop0)
	tail1 := NewPrototypeObjectType(bTail1)

	bTail0 := NewPrototypeObjectType_Builder(reg)
	bTail0.SetName("Tail0")
	bTail0.SetImplicitPrototype(tail1)
	tail0 := NewPrototypeObjectType(bTail0)

	reg.SetParallelReadOnly(true)
	defer reg.SetParallelReadOnly(false)

	allNodes := make([]*PrototypeObjectTypeBase, 0, chainLen+8)
	allNodes = append(allNodes, acyclic...)
	allNodes = append(allNodes, selfCyc, cycA, cycB, cycC, loop0, loop1, tail1, tail0)
	for _, n := range allNodes {
		if n.hasFlag(jsTypeFlagVisited) {
			t.Fatalf("unexpected jsTypeFlagVisited before walk on %v", n)
		}
	}

	const workers = 16
	const iters = 100
	start := make(chan struct{})
	var wg sync.WaitGroup
	errCh := make(chan error, workers)

	for g := 0; g < workers; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-start
			for it := 0; it < iters; it++ {
				if acyclic[(id+it)%chainLen].DetectInheritanceCycle() {
					errCh <- fmt.Errorf("acyclic node reported cycle")
					return
				}
				if !selfCyc.DetectInheritanceCycle() {
					errCh <- fmt.Errorf("selfCyc missed cycle")
					return
				}
				if !cycA.DetectInheritanceCycle() || !cycB.DetectInheritanceCycle() || !cycC.DetectInheritanceCycle() {
					errCh <- fmt.Errorf("3-cycle missed cycle")
					return
				}
				if !tail0.DetectInheritanceCycle() || !tail1.DetectInheritanceCycle() || !loop0.DetectInheritanceCycle() {
					errCh <- fmt.Errorf("lasso cycle missed cycle")
					return
				}
			}
		}(g)
	}
	close(start)
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatal(err)
	}

	for _, n := range allNodes {
		if n.hasFlag(jsTypeFlagVisited) {
			t.Fatalf("detectImplicitPrototypeCycle mutated jsTypeFlagVisited on %v", n)
		}
	}
}

func TestFunctionTypeGetPrototypeAndGetSlotConcurrentParallelReadOnly(t *testing.T) {
	reg := NewJSTypeRegistry(parallelRaceNoopReporter{})

	const numFuncs = 24
	namedCtors := make([]FunctionType, numFuncs)
	namedFns := make([]FunctionType, numFuncs)
	structuralFns := make([]FunctionType, numFuncs)

	closer := reg.GetResolver().OpenForDefinition()
	for i := 0; i < numFuncs; i++ {
		var cb FunctionType_Builder
		cb.Init(reg)
		cb.WithName(fmt.Sprintf("ParallelCtor_%d", i)).
			ForConstructor().
			WithReturnType(reg.GetNativeType(JSTypeNative_VOID_TYPE))
		namedCtors[i] = cb.Build()

		var fb FunctionType_Builder
		fb.Init(reg)
		fb.WithName(fmt.Sprintf("ParallelFn_%d", i)).
			WithReturnType(reg.GetNativeType(JSTypeNative_NUMBER_TYPE))
		namedFns[i] = fb.Build()

		var sb FunctionType_Builder
		sb.Init(reg)
		sb.WithReturnType(reg.GetNativeType(JSTypeNative_STRING_TYPE))
		structuralFns[i] = sb.Build()
	}
	closer.Close()

	// closer.Close() resolves each constructor's InstanceObjectType, which eagerly
	// populates prototypeSlot on constructors; reset constructor prototypeSlots back
	// to nil so lazy prototype initialization is exercised concurrently under -race
	// for constructors, named ordinary functions, and structural functions alike.
	for i := 0; i < numFuncs; i++ {
		if slot := namedCtors[i].prototypeSlot(); slot != nil {
			if obj, ok := slot.GetType().(ObjectType); ok && obj != nil {
				obj.SetOwnerFunction(nil)
			}
			namedCtors[i].setPrototypeSlot(nil)
		}
		if namedCtors[i].prototypeSlot() != nil || namedFns[i].prototypeSlot() != nil || structuralFns[i].prototypeSlot() != nil {
			t.Fatalf("expected lazy nil prototypeSlot before parallel phase at %d", i)
		}
	}

	reg.SetParallelReadOnly(true)
	defer reg.SetParallelReadOnly(false)

	const workers = 16
	ctorProtos := make([][numFuncs]ObjectType, workers)
	ctorSlots := make([][numFuncs]*Property, workers)
	fnProtos := make([][numFuncs]ObjectType, workers)
	structProtos := make([][numFuncs]ObjectType, workers)

	start := make(chan struct{})
	var wg sync.WaitGroup
	errCh := make(chan error, workers)

	for g := 0; g < workers; g++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			<-start
			for i := 0; i < numFuncs; i++ {
				idx := (i + w) % numFuncs
				ctor := namedCtors[idx]
				fn := namedFns[idx]
				sfn := structuralFns[idx]

				if w%2 == 0 {
					ctorProtos[w][idx] = ctor.GetPrototype()
					ctorSlots[w][idx] = ctor.GetSlot("prototype")
				} else {
					ctorSlots[w][idx] = ctor.GetSlot("prototype")
					ctorProtos[w][idx] = ctor.GetPrototype()
				}
				if owner := ctorProtos[w][idx].GetOwnerFunction(); owner != ctor {
					errCh <- fmt.Errorf("worker %d idx %d: owner=%v, want %v", w, idx, owner, ctor)
					return
				}
				_ = ctor.IsAmbiguousConstructor()
				_ = ctor.GetOwnPropertyNames()
				ctor.ForEachOwnPropertyName(func(string) {})
				_ = ctor.GetPropertyType("call")
				_ = ctor.GetPropertyType("bind")
				_ = ctor.GetPropertyType("apply")

				if w%2 == 0 {
					fnProtos[w][idx] = fn.GetPrototype()
				} else {
					fnProtos[w][idx] = fn.GetSlot("prototype").GetType().(ObjectType)
				}
				if owner := fnProtos[w][idx].GetOwnerFunction(); owner != fn {
					errCh <- fmt.Errorf("worker %d fn idx %d: owner=%v, want %v", w, idx, owner, fn)
					return
				}

				structProtos[w][idx] = sfn.GetPrototype()
				if sfn.GetSlot("prototype") == nil {
					errCh <- fmt.Errorf("worker %d structural idx %d: nil prototype slot", w, idx)
					return
				}
			}
		}(g)
	}
	close(start)
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatal(err)
	}

	unknownObj := reg.GetNativeObjectType(JSTypeNative_UNKNOWN_TYPE)
	for i := 0; i < numFuncs; i++ {
		for w := 1; w < workers; w++ {
			if ctorProtos[w][i] != ctorProtos[0][i] {
				t.Fatalf("ctor[%d] double-created prototype across workers: %p vs %p", i, ctorProtos[0][i], ctorProtos[w][i])
			}
			if ctorSlots[w][i] != ctorSlots[0][i] {
				t.Fatalf("ctor[%d] double-created prototypeSlot across workers: %p vs %p", i, ctorSlots[0][i], ctorSlots[w][i])
			}
			if fnProtos[w][i] != fnProtos[0][i] {
				t.Fatalf("fn[%d] double-created prototype across workers: %p vs %p", i, fnProtos[0][i], fnProtos[w][i])
			}
			if structProtos[w][i] != unknownObj {
				t.Fatalf("structural[%d] expected UNKNOWN_TYPE prototype, got %v", i, structProtos[w][i])
			}
		}
	}
}
