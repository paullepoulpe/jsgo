package jscomp

import "os"

// Go-only: opt-in concrete traversal experiment. Read once at process startup;
// all ordinary traversal and callback behavior remains the default oracle.
//
//go:generate python3 ../../tools/bench/gen_concrete_traversal.py
var concreteTraversalExperimentEnabled = os.Getenv("JSGO_EXPERIMENT_CONCRETE_TRAVERSAL") == "1"
