# Keep flatten_helpers.go as a shallow shared-helpers file

The Flatten data source and Flatten function both use `errorTitle` and `flattenedToMapValue` from `flatten_helpers.go`. We considered extracting a deeper module that owns the full "call Flattener → convert result → map errors" workflow, but decided against it.

The proposed deeper module would have an interface (`flattener + yaml string → types.Map + error`) nearly as complex as its implementation (three lines of orchestration). The two adapters are asymmetric — the data source owns file I/O, path validation, and mutual-exclusivity checks before flattening, while the function is a straight pass-through. A shared workflow module would only cover the simple tail end that both adapters share, adding indirection without adding depth.

`flatten_helpers.go` is shallow by design: it holds the two genuinely shared conversions and nothing else. That's the right level of extraction for ~25 lines of shared code behind two callers.
