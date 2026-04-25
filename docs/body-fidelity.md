# Body Fidelity Report

Generated: auto-computed by tools/bodyfidelity.

## Summary

| Metric | Value |
|--------|-------|
| Total matched methods | 888 |
| Average score | 73.6 |
| Methods below 50 | 189 |

## Methods by Body-Similarity Score (worst first)

| Rank | Score | Go File | Method | Java File | Java Line | Go LoC | Java LoC | Notes |
|------|-------|---------|--------|-----------|-----------|--------|----------|-------|
| 1 | 0 | `internal/jscomp/checkjsdoc.go` | `validateClassLevelJsDoc` | `com/google/javascript/jscomp/CheckJSDoc.java` | 315 | 17 | 2 | much larger (17 LoC vs Java's 2) |
| 2 | 0 | `internal/rhino/jstype/jstype.go` | `IsSubtype` | `com/google/javascript/rhino/jstype/JSType.java` | 1464 | 20 | 1 | much larger (20 LoC vs Java's 1) |
| 3 | 0 | `internal/jscomp/nodeutil.go` | `IsLhsByDestructuring` | `com/google/javascript/jscomp/NodeUtil.java` | 3187 | 15 | 4 | much larger (15 LoC vs Java's 4) |
| 4 | 0 | `internal/rhino/jstype/jstype.go` | `MatchesObjectContext` | `com/google/javascript/rhino/jstype/JSType.java` | 828 | 13 | 1 | much larger (13 LoC vs Java's 1) |
| 5 | 0 | `internal/rhino/jstype/jstype.go` | `AutoboxesTo` | `com/google/javascript/rhino/jstype/JSType.java` | 914 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 6 | 0 | `internal/rhino/jstype/jstype.go` | `RestrictByNotNullOrUndefined` | `com/google/javascript/rhino/jstype/JSType.java` | 1390 | 7 | 1 | much larger (7 LoC vs Java's 1) |
| 7 | 0 | `internal/rhino/jstype/jstype.go` | `CollapseUnion` | `com/google/javascript/rhino/jstype/JSType.java` | 1097 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 8 | 0 | `internal/rhino/jstype/jstype.go` | `testForEquality` | `com/google/javascript/rhino/jstype/JSType.java` | 976 | 11 | 1 | much larger (11 LoC vs Java's 1) |
| 9 | 0 | `internal/rhino/jstype/jstype.go` | `IsNullable` | `com/google/javascript/rhino/jstype/JSType.java` | 1075 | 7 | 1 | much larger (7 LoC vs Java's 1) |
| 10 | 0 | `internal/rhino/jstype/jstype.go` | `IsVoidable` | `com/google/javascript/rhino/jstype/JSType.java` | 1082 | 7 | 1 | much larger (7 LoC vs Java's 1) |
| 11 | 0 | `internal/rhino/jstype/jstype.go` | `IsObject` | `com/google/javascript/rhino/jstype/JSType.java` | 622 | 8 | 1 | much larger (8 LoC vs Java's 1) |
| 12 | 0 | `internal/rhino/jstype/property.go` | `HashCode` | `com/google/javascript/rhino/jstype/Property.java` | 277 | 11 | 1 | much larger (11 LoC vs Java's 1) |
| 13 | 0 | `internal/rhino/jstype/property_map.go` | `GetOwnProperty` | `com/google/javascript/rhino/jstype/PropertyMap.java` | 185 | 9 | 1 | much larger (9 LoC vs Java's 1) |
| 14 | 0 | `internal/jscomp/nodeutil.go` | `GetVarsDeclaredInBranch` | `com/google/javascript/jscomp/NodeUtil.java` | 3698 | 15 | 3 | much larger (15 LoC vs Java's 3) |
| 15 | 0 | `internal/jscomp/nodeutil.go` | `GetEnclosingBlockScopeRoot` | `com/google/javascript/jscomp/NodeUtil.java` | 1725 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 16 | 0 | `internal/jscomp/nodeutil.go` | `GetEnclosingHoistScopeRoot` | `com/google/javascript/jscomp/NodeUtil.java` | 1737 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 17 | 0 | `internal/rhino/jstype/union_type.go` | `RestrictByNotNullOrUndefined` | `com/google/javascript/rhino/jstype/UnionType.java` | 248 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 18 | 0 | `internal/rhino/jstype/union_type.go` | `RestrictByNotUndefined` | `com/google/javascript/rhino/jstype/UnionType.java` | 253 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 19 | 0 | `internal/rhino/jstype/union_type.go` | `RestrictByNotNull` | `com/google/javascript/rhino/jstype/UnionType.java` | 258 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 20 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordWizcallback` | `com/google/javascript/rhino/JSDocInfo.java` | 2700 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 21 | 10 | `internal/rhino/jstype/union_type.go` | `MatchesNumberContext` | `com/google/javascript/rhino/jstype/UnionType.java` | 157 | 7 | 1 | much larger (7 LoC vs Java's 1) |
| 22 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordWizaction` | `com/google/javascript/rhino/JSDocInfo.java` | 2690 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 23 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordCustomElement` | `com/google/javascript/rhino/JSDocInfo.java` | 2730 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 24 | 10 | `internal/rhino/jstype/function_type.go` | `HasEqualCallType` | `com/google/javascript/rhino/jstype/FunctionType.java` | 915 | 8 | 1 | much larger (8 LoC vs Java's 1) |
| 25 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordDescription` | `com/google/javascript/rhino/JSDocInfo.java` | 2247 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 26 | 10 | `internal/jscomp/typedscopecreator.go` | `isLValueRootedInGlobalScope` | `com/google/javascript/jscomp/TypedScopeCreator.java` | 2294 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 27 | 10 | `internal/rhino/jstype/function_type.go` | `checkExtendsLoop` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1231 | 19 | 1 | much larger (19 LoC vs Java's 1) |
| 28 | 10 | `internal/rhino/jstype/jstype.go` | `IsRecordType` | `com/google/javascript/rhino/jstype/JSType.java` | 511 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 29 | 10 | `internal/rhino/jsdoc.go` | `ContainsDeclaration` | `com/google/javascript/rhino/JSDocInfo.java` | 978 | 12 | 1 | much larger (12 LoC vs Java's 1) |
| 30 | 10 | `internal/jscomp/nodeutil.go` | `GetEnclosingScopeRoot` | `com/google/javascript/jscomp/NodeUtil.java` | 1729 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 31 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordNgInject` | `com/google/javascript/rhino/JSDocInfo.java` | 2680 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 32 | 10 | `internal/jscomp/nodeutil.go` | `GetEnclosingStatement` | `com/google/javascript/jscomp/NodeUtil.java` | 1750 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 33 | 10 | `internal/rhino/jstype/jstype.go` | `ToMaybeObjectType` | `com/google/javascript/rhino/jstype/JSType.java` | 1634 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 34 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordMeaning` | `com/google/javascript/rhino/JSDocInfo.java` | 2263 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 35 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordPolymer` | `com/google/javascript/rhino/JSDocInfo.java` | 2720 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 36 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordMixinFunction` | `com/google/javascript/rhino/JSDocInfo.java` | 2750 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 37 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordMixinClass` | `com/google/javascript/rhino/JSDocInfo.java` | 2740 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 38 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordProvideAlreadyProvided` | `com/google/javascript/rhino/JSDocInfo.java` | 2470 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 39 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordProvideGoog` | `com/google/javascript/rhino/JSDocInfo.java` | 2466 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 40 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordRequireInlining` | `com/google/javascript/rhino/JSDocInfo.java` | 2413 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 41 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordEncourageInlining` | `com/google/javascript/rhino/JSDocInfo.java` | 2401 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 42 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordNoDts` | `com/google/javascript/rhino/JSDocInfo.java` | 2367 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 43 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordNoCompile` | `com/google/javascript/rhino/JSDocInfo.java` | 2356 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 44 | 10 | `internal/rhino/jstype/union_type.go` | `HasAnyTemplateTypesInternal` | `com/google/javascript/rhino/jstype/UnionType.java` | 611 | 8 | 1 | much larger (8 LoC vs Java's 1) |
| 45 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordCollapsibleOrBreakMyCode` | `com/google/javascript/rhino/JSDocInfo.java` | 2435 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 46 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordPureOrBreakMyCode` | `com/google/javascript/rhino/JSDocInfo.java` | 2424 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 47 | 10 | `internal/rhino/jstype/union_type.go` | `MatchesSymbolContext` | `com/google/javascript/rhino/jstype/UnionType.java` | 185 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 48 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordOverride` | `com/google/javascript/rhino/JSDocInfo.java` | 2551 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 49 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordPolymerBehavior` | `com/google/javascript/rhino/JSDocInfo.java` | 2710 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 50 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordDeprecated` | `com/google/javascript/rhino/JSDocInfo.java` | 2559 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 51 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordTypeSummary` | `com/google/javascript/rhino/JSDocInfo.java` | 2630 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 52 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordNoCoverage` | `com/google/javascript/rhino/JSDocInfo.java` | 2622 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 53 | 10 | `internal/rhino/jstype/property_map.go` | `GetOwnKnownSymbols` | `com/google/javascript/rhino/jstype/PropertyMap.java` | 209 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 54 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordDeprecationReason` | `com/google/javascript/rhino/JSDocInfo.java` | 2053 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 55 | 10 | `internal/rhino/jstype/union_type.go` | `IsDict` | `com/google/javascript/rhino/jstype/UnionType.java` | 315 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 56 | 10 | `internal/rhino/jstype/union_type.go` | `IsStruct` | `com/google/javascript/rhino/jstype/UnionType.java` | 310 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 57 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordExterns` | `com/google/javascript/rhino/JSDocInfo.java` | 2613 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 58 | 10 | `internal/rhino/jstype/union_type.go` | `IsObject` | `com/google/javascript/rhino/jstype/UnionType.java` | 415 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 59 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordConstancy` | `com/google/javascript/rhino/JSDocInfo.java` | 2216 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 60 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordFinality` | `com/google/javascript/rhino/JSDocInfo.java` | 2237 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 61 | 10 | `internal/rhino/jstype/union_type.go` | `IsVoidable` | `com/google/javascript/rhino/jstype/UnionType.java` | 294 | 7 | 1 | much larger (7 LoC vs Java's 1) |
| 62 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordImplicitCast` | `com/google/javascript/rhino/JSDocInfo.java` | 2597 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 63 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordNoSideEffects` | `com/google/javascript/rhino/JSDocInfo.java` | 2605 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 64 | 10 | `internal/jscomp/checkjsdoc.go` | `validateImplicitCast` | `com/google/javascript/jscomp/CheckJSDoc.java` | 703 | 7 | 2 | much larger (7 LoC vs Java's 2) |
| 65 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RemoveExport` | `com/google/javascript/rhino/JSDocInfo.java` | 2588 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 66 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordExport` | `com/google/javascript/rhino/JSDocInfo.java` | 2580 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 67 | 10 | `internal/rhino/jstype/union_type.go` | `Contains` | `com/google/javascript/rhino/jstype/UnionType.java` | 427 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 68 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordNoInline` | `com/google/javascript/rhino/JSDocInfo.java` | 2389 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 69 | 10 | `internal/rhino/jstype/registry.go` | `RegisterTemplateTypeNamesInScope` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 2418 | 8 | 2 | much larger (8 LoC vs Java's 2) |
| 70 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordNoCollapse` | `com/google/javascript/rhino/JSDocInfo.java` | 2378 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 71 | 10 | `internal/rhino/jsdoc_info_builder.go` | `RecordAlternateMessageId` | `com/google/javascript/rhino/JSDocInfo.java` | 2272 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 72 | 13 | `internal/jscomp/nodeutil.go` | `IsAnyFor` | `com/google/javascript/jscomp/NodeUtil.java` | 2072 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 73 | 13 | `internal/jscomp/parsing/jsdocinfoparser.go` | `GetFileOverviewJSDocInfo` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2824 | 3 | 1 |  |
| 74 | 13 | `internal/rhino/qualified_name.go` | `GetOwner` | `com/google/javascript/rhino/QualifiedName.java` | 221 | 3 | 1 |  |
| 75 | 13 | `internal/jscomp/nodeutil.go` | `IsStatementBlock` | `com/google/javascript/jscomp/NodeUtil.java` | 2170 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 76 | 13 | `internal/jscomp/checkunreachablecode.go` | `enterScopeWithCfg` | `com/google/javascript/jscomp/CheckUnreachableCode.java` | 41 | 8 | 2 | much larger (8 LoC vs Java's 2) |
| 77 | 13 | `internal/rhino/qualified_name.go` | `GetOwner` | `com/google/javascript/rhino/QualifiedName.java` | 221 | 3 | 1 |  |
| 78 | 14 | `internal/jscomp/nodeutil.go` | `IsCallOrNew` | `com/google/javascript/jscomp/NodeUtil.java` | 2541 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 79 | 14 | `internal/rhino/jstype/object_type_ext.go` | `GetPropertyMap` | `com/google/javascript/rhino/jstype/ObjectType.java` | 99 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 80 | 14 | `internal/jscomp/checkunreachablecode.go` | `initScope` | `com/google/javascript/jscomp/CheckUnreachableCode.java` | 72 | 14 | 2 | much larger (14 LoC vs Java's 2) |
| 81 | 17 | `internal/rhino/jsdoc_info_builder.go` | `RecordLends` | `com/google/javascript/rhino/JSDocInfo.java` | 2670 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 82 | 17 | `internal/rhino/jsdoc_info_builder.go` | `RecordThisType` | `com/google/javascript/rhino/JSDocInfo.java` | 2181 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 83 | 17 | `internal/jscomp/nodeutil.go` | `IsFunctionDeclaration` | `com/google/javascript/jscomp/NodeUtil.java` | 2604 | 10 | 1 | much larger (10 LoC vs Java's 1) |
| 84 | 17 | `internal/jscomp/nodeutil.go` | `IsInvocation` | `com/google/javascript/jscomp/NodeUtil.java` | 4676 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 85 | 17 | `internal/rhino/jsdoc_info_builder.go` | `RecordClosurePrimitiveId` | `com/google/javascript/rhino/JSDocInfo.java` | 2281 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 86 | 17 | `internal/jscomp/nodeutil.go` | `IsStatementParent` | `com/google/javascript/jscomp/NodeUtil.java` | 2238 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 87 | 17 | `internal/rhino/jsdoc_info_builder.go` | `RecordType` | `com/google/javascript/rhino/JSDocInfo.java` | 2101 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 88 | 17 | `internal/rhino/jsdoc_info_builder.go` | `RecordReturnType` | `com/google/javascript/rhino/JSDocInfo.java` | 2126 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 89 | 17 | `internal/jscomp/nodeutil.go` | `isStatementParent` | `com/google/javascript/jscomp/NodeUtil.java` | 2238 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 90 | 18 | `internal/jscomp/nodeutil.go` | `IsStatement` | `com/google/javascript/jscomp/NodeUtil.java` | 2225 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 91 | 19 | `internal/rhino/jsdoc_info_builder.go` | `RecordModifies` | `com/google/javascript/rhino/JSDocInfo.java` | 2091 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 92 | 19 | `internal/jscomp/checkclosureimports.go` | `enterModule` | `com/google/javascript/jscomp/CheckClosureImports.java` | 242 | 7 | 1 | much larger (7 LoC vs Java's 1) |
| 93 | 19 | `internal/rhino/jsdoc_info_builder.go` | `RecordParameter` | `com/google/javascript/rhino/JSDocInfo.java` | 1901 | 7 | 2 | much larger (7 LoC vs Java's 2) |
| 94 | 19 | `internal/jscomp/nodeutil.go` | `IsPrototypePropertyDeclaration` | `com/google/javascript/jscomp/NodeUtil.java` | 3829 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 95 | 20 | `internal/rhino/jsdoc_info_builder.go` | `RecordBaseType` | `com/google/javascript/rhino/JSDocInfo.java` | 2191 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 96 | 20 | `internal/rhino/jsdoc_info_builder.go` | `RecordStableIdGenerator` | `com/google/javascript/rhino/JSDocInfo.java` | 2015 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 97 | 20 | `internal/jscomp/checkclosureimports.go` | `callName` | `com/google/javascript/jscomp/CheckClosureImports.java` | 92 | 10 | 1 | much larger (10 LoC vs Java's 1) |
| 98 | 20 | `internal/rhino/jstype/union_type.go` | `IsNullable` | `com/google/javascript/rhino/jstype/UnionType.java` | 286 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 99 | 20 | `internal/jscomp/checkclosureimports.go` | `allowDestructuring` | `com/google/javascript/jscomp/CheckClosureImports.java` | 77 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 100 | 20 | `internal/rhino/jstype/union_type.go` | `CanBeCalled` | `com/google/javascript/rhino/jstype/UnionType.java` | 238 | 9 | 1 | much larger (9 LoC vs Java's 1) |
| 101 | 20 | `internal/rhino/jstype/union_type.go` | `IsExplicitlyVoidable` | `com/google/javascript/rhino/jstype/UnionType.java` | 300 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 102 | 20 | `internal/rhino/jstype/union_type.go` | `IsUnknownType` | `com/google/javascript/rhino/jstype/UnionType.java` | 305 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 103 | 20 | `internal/jscomp/nodeutil.go` | `GetEnclosingType` | `com/google/javascript/jscomp/NodeUtil.java` | 1693 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 104 | 20 | `internal/rhino/jsdoc_info_builder.go` | `RecordIdGenerator` | `com/google/javascript/rhino/JSDocInfo.java` | 2048 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 105 | 20 | `internal/rhino/jsdoc_info_builder.go` | `RecordConsistentIdGenerator` | `com/google/javascript/rhino/JSDocInfo.java` | 2004 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 106 | 20 | `internal/rhino/jsdoc_info_builder.go` | `RecordXidGenerator` | `com/google/javascript/rhino/JSDocInfo.java` | 2026 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 107 | 20 | `internal/jscomp/nodeutil.go` | `IsExprAssign` | `com/google/javascript/jscomp/NodeUtil.java` | 2050 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 108 | 20 | `internal/jscomp/checkmissingrequires.go` | `isIIFE` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 177 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 109 | 20 | `internal/rhino/jstype/jstype.go` | `ToMaybeFunctionType` | `com/google/javascript/rhino/jstype/JSType.java` | 448 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 110 | 20 | `internal/rhino/jstype/jstype.go` | `ToMaybeEnumType` | `com/google/javascript/rhino/jstype/JSType.java` | 499 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 111 | 20 | `internal/rhino/jsdoc_info_builder.go` | `RecordExtendedInterface` | `com/google/javascript/rhino/JSDocInfo.java` | 2654 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 112 | 20 | `internal/rhino/jstype_expr.go` | `IsEquivalentTo` | `com/google/javascript/rhino/JSTypeExpression.java` | 186 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 113 | 20 | `internal/jscomp/nodeutil.go` | `IsObjectLitKey` | `com/google/javascript/jscomp/NodeUtil.java` | 3242 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 114 | 20 | `internal/rhino/jsdoc_info_builder.go` | `RecordImplementedInterface` | `com/google/javascript/rhino/JSDocInfo.java` | 2649 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 115 | 20 | `internal/jscomp/checkclosureimports.go` | `aliasMustBeConstant` | `com/google/javascript/jscomp/CheckClosureImports.java` | 82 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 116 | 20 | `internal/rhino/jstype/contains_upper_bound_super_type_visitor.go` | `CaseNamedType` | `com/google/javascript/rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java` | 82 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 117 | 20 | `internal/rhino/jsdoc_info_builder.go` | `RecordMappedIdGenerator` | `com/google/javascript/rhino/JSDocInfo.java` | 2037 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 118 | 20 | `internal/jscomp/nodeutil.go` | `IsExprCall` | `com/google/javascript/jscomp/NodeUtil.java` | 2060 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 119 | 20 | `internal/rhino/jstype/jstype.go` | `ToMaybeUnionType` | `com/google/javascript/rhino/jstype/JSType.java` | 423 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 120 | 20 | `internal/rhino/jstype/registry.go` | `GetNativeObjectType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1712 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 121 | 20 | `internal/rhino/jstype/registry.go` | `GetNativeFunctionType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1716 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 122 | 23 | `internal/rhino/jsdoc.go` | `GetMods` | `com/google/javascript/rhino/JSDocInfo.java` | 1450 | 3 | 1 |  |
| 123 | 23 | `internal/rhino/jsdoc.go` | `GetAuthors` | `com/google/javascript/rhino/JSDocInfo.java` | 1405 | 3 | 1 |  |
| 124 | 23 | `internal/rhino/jsdoc.go` | `GetLicense` | `com/google/javascript/rhino/JSDocInfo.java` | 1291 | 3 | 1 |  |
| 125 | 23 | `internal/rhino/jsdoc.go` | `GetBlockDescription` | `com/google/javascript/rhino/JSDocInfo.java` | 1420 | 3 | 1 |  |
| 126 | 23 | `internal/rhino/jstype/registry.go` | `RegisterLegacyClosureNamespace` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 2523 | 3 | 1 |  |
| 127 | 23 | `internal/rhino/jsdoc.go` | `GetEnhance` | `com/google/javascript/rhino/JSDocInfo.java` | 1440 | 3 | 1 |  |
| 128 | 23 | `internal/rhino/jsdoc.go` | `GetSuppressionsAndTheirDescription` | `com/google/javascript/rhino/JSDocInfo.java` | 1381 | 6 | 2 |  |
| 129 | 23 | `internal/rhino/jsdoc.go` | `GetTemplateTypes` | `com/google/javascript/rhino/JSDocInfo.java` | 1470 | 6 | 2 |  |
| 130 | 23 | `internal/rhino/jstype/function_type.go` | `HasProperties` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1158 | 3 | 1 |  |
| 131 | 23 | `internal/rhino/jsdoc.go` | `GetFileOverview` | `com/google/javascript/rhino/JSDocInfo.java` | 1430 | 3 | 1 |  |
| 132 | 23 | `internal/rhino/jsdoc.go` | `GetExtendedInterfaces` | `com/google/javascript/rhino/JSDocInfo.java` | 1348 | 3 | 1 |  |
| 133 | 23 | `internal/rhino/jsdoc.go` | `GetOriginalCommentString` | `com/google/javascript/rhino/JSDocInfo.java` | 1560 | 3 | 1 |  |
| 134 | 23 | `internal/rhino/jsdoc.go` | `HasParameterType` | `com/google/javascript/rhino/JSDocInfo.java` | 1047 | 3 | 1 |  |
| 135 | 23 | `internal/rhino/jsdoc.go` | `GetTypeTransformations` | `com/google/javascript/rhino/JSDocInfo.java` | 1476 | 6 | 2 |  |
| 136 | 23 | `internal/rhino/jsdoc.go` | `GetAlternateMessageId` | `com/google/javascript/rhino/JSDocInfo.java` | 1195 | 3 | 1 |  |
| 137 | 23 | `internal/rhino/jsdoc.go` | `GetClosurePrimitiveId` | `com/google/javascript/rhino/JSDocInfo.java` | 1215 | 3 | 1 |  |
| 138 | 23 | `internal/rhino/jstype/named_type.go` | `GetReferencedType` | `com/google/javascript/rhino/jstype/NamedType.java` | 208 | 3 | 1 |  |
| 139 | 23 | `internal/rhino/jsdoc.go` | `GetMeaning` | `com/google/javascript/rhino/JSDocInfo.java` | 1182 | 3 | 1 |  |
| 140 | 23 | `internal/rhino/jsdoc.go` | `GetDeprecationReason` | `com/google/javascript/rhino/JSDocInfo.java` | 1359 | 3 | 1 |  |
| 141 | 23 | `internal/rhino/jsdoc.go` | `GetReturnDescription` | `com/google/javascript/rhino/JSDocInfo.java` | 1415 | 3 | 1 |  |
| 142 | 23 | `internal/jscomp/typedscope.go` | `GetDeclarativelyUnboundVarsWithoutTypes` | `com/google/javascript/jscomp/TypedScope.java` | 257 | 3 | 1 |  |
| 143 | 23 | `internal/rhino/jsdoc.go` | `GetThisType` | `com/google/javascript/rhino/JSDocInfo.java` | 1147 | 3 | 1 |  |
| 144 | 23 | `internal/jscomp/parsing/jsdocinfoparser.go` | `RetrieveAndResetParsedJSDocInfo` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2819 | 3 | 1 |  |
| 145 | 23 | `internal/rhino/jstype/function_type.go` | `IsInstanceType` | `com/google/javascript/rhino/jstype/FunctionType.java` | 241 | 3 | 1 |  |
| 146 | 23 | `internal/rhino/jsdoc.go` | `GetReferences` | `com/google/javascript/rhino/JSDocInfo.java` | 1410 | 3 | 1 |  |
| 147 | 23 | `internal/rhino/jsdoc_info_builder.go` | `Build` | `com/google/javascript/rhino/JSDocInfo.java` | 1748 | 3 | 1 |  |
| 148 | 23 | `internal/rhino/jstype/property.go` | `SetType` | `com/google/javascript/rhino/jstype/Property.java` | 246 | 3 | 1 |  |
| 149 | 23 | `internal/rhino/jsdoc.go` | `GetDescription` | `com/google/javascript/rhino/JSDocInfo.java` | 1162 | 3 | 1 |  |
| 150 | 24 | `internal/jscomp/checkmissingrequires.go` | `visitJsDocExpr` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 255 | 10 | 3 | much larger (10 LoC vs Java's 3) |
| 151 | 24 | `internal/rhino/jstype/property_map.go` | `GetOwnPropertyNames` | `com/google/javascript/rhino/jstype/PropertyMap.java` | 205 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 152 | 24 | `internal/rhino/jstype/property_map.go` | `Values` | `com/google/javascript/rhino/jstype/PropertyMap.java` | 308 | 4 | 1 | much larger (4 LoC vs Java's 1) |
| 153 | 25 | `internal/jscomp/checkduplicatemember.go` | `checkObjectLiteralOrClass` | `com/google/javascript/jscomp/StrictModeCheck.java` | 163 | 67 | 21 | much larger (67 LoC vs Java's 21) |
| 154 | 27 | `internal/jscomp/checkjsdoc.go` | `checkJsDocInBlockComments` | `com/google/javascript/jscomp/CheckJSDoc.java` | 118 | 35 | 11 | much larger (35 LoC vs Java's 11) |
| 155 | 28 | `internal/jscomp/parsing/jsdocinfoparser.go` | `parseRecordType` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2568 | 28 | 10 |  |
| 156 | 29 | `internal/jscomp/nodeutil.go` | `IsNameReferenced` | `com/google/javascript/jscomp/NodeUtil.java` | 4056 | 3 | 1 |  |
| 157 | 31 | `internal/jscomp/nodeutil.go` | `GetFunctionParameters` | `com/google/javascript/jscomp/NodeUtil.java` | 4308 | 6 | 2 |  |
| 158 | 31 | `internal/rhino/jstype/jstype.go` | `GetLeastSupertype` | `com/google/javascript/rhino/jstype/JSType.java` | 1116 | 15 | 5 |  |
| 159 | 31 | `internal/jscomp/checksuspiciouscode.go` | `reportIfNaN` | `com/google/javascript/jscomp/CheckSuspiciousCode.java` | 126 | 6 | 2 |  |
| 160 | 31 | `internal/rhino/jstype_expr.go` | `Copy` | `com/google/javascript/rhino/JSTypeExpression.java` | 207 | 6 | 1 | much larger (6 LoC vs Java's 1) |
| 161 | 33 | `internal/jscomp/typedscope.go` | `GetParent` | `com/google/javascript/jscomp/TypedScope.java` | 149 | 3 | 1 |  |
| 162 | 33 | `internal/jscomp/typedscope.go` | `GetDepth` | `com/google/javascript/jscomp/TypedScope.java` | 144 | 3 | 1 |  |
| 163 | 33 | `internal/jscomp/nodeutil.go` | `And` | `com/google/javascript/jscomp/NodeUtil.java` | 1507 | 3 | 1 |  |
| 164 | 33 | `internal/jscomp/nodeutil.go` | `GetNameReferenceCount` | `com/google/javascript/jscomp/NodeUtil.java` | 4066 | 3 | 1 |  |
| 165 | 33 | `internal/jscomp/typedscope.go` | `IsBottom` | `com/google/javascript/jscomp/TypedScope.java` | 134 | 3 | 1 |  |
| 166 | 33 | `internal/rhino/jstype/function_type.go` | `SetStruct` | `com/google/javascript/rhino/jstype/FunctionType.java` | 327 | 3 | 1 |  |
| 167 | 33 | `internal/rhino/jsdoc.go` | `IsConstructorOrInterface` | `com/google/javascript/rhino/JSDocInfo.java` | 833 | 3 | 1 |  |
| 168 | 33 | `internal/rhino/jstype/function_type.go` | `SetDict` | `com/google/javascript/rhino/jstype/FunctionType.java` | 331 | 3 | 1 |  |
| 169 | 33 | `internal/rhino/jsdoc_info_builder.go` | `IsModsRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2312 | 3 | 1 |  |
| 170 | 34 | `internal/rhino/jstype/named_type.go` | `Build` | `com/google/javascript/rhino/jstype/NamedType.java` | 561 | 14 | 1 | much larger (14 LoC vs Java's 1) |
| 171 | 37 | `internal/rhino/jstype/function_type.go` | `GetConstructorOnlyTemplateParameters` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1307 | 17 | 7 |  |
| 172 | 39 | `internal/jscomp/checkmissingreturn.go` | `fastAllPathsReturnCheck` | `com/google/javascript/jscomp/CheckMissingReturn.java` | 114 | 13 | 5 |  |
| 173 | 39 | `internal/rhino/jstype/property_map.go` | `HashCode` | `com/google/javascript/rhino/jstype/PropertyMap.java` | 313 | 8 | 3 |  |
| 174 | 40 | `internal/jscomp/typedscope.go` | `CreateLatticeBottom` | `com/google/javascript/jscomp/TypedScope.java` | 111 | 5 | 1 | much larger (5 LoC vs Java's 1) |
| 175 | 42 | `internal/rhino/jstype_expr.go` | `GetAllTypeNames` | `com/google/javascript/rhino/JSTypeExpression.java` | 138 | 7 | 3 |  |
| 176 | 43 | `internal/jscomp/checkjsdoc.go` | `visit` | `com/google/javascript/jscomp/CheckJSDoc.java` | 761 | 19 | 8 |  |
| 177 | 44 | `internal/jscomp/nodeutil.go` | `RemoveName` | `com/google/javascript/jscomp/NodeUtil.java` | 584 | 10 | 4 |  |
| 178 | 45 | `internal/jscomp/checkmissingrequires.go` | `isGoogScopeBody` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 268 | 10 | 4 |  |
| 179 | 45 | `internal/jscomp/checkduplicatemember.go` | `report` | `com/google/javascript/jscomp/StrictModeCheck.java` | 249 | 8 | 2 | much larger (8 LoC vs Java's 2) |
| 180 | 46 | `internal/rhino/jsdoc_info_builder.go` | `RecordUnrestricted` | `com/google/javascript/rhino/JSDocInfo.java` | 2489 | 7 | 3 |  |
| 181 | 47 | `internal/jscomp/checkdebuggerstatement.go` | `visit` | `com/google/javascript/jscomp/CheckDebuggerStatement.java` | 52 | 14 | 2 | much larger (14 LoC vs Java's 2) |
| 182 | 48 | `internal/rhino/jstype/registry.go` | `RegisterPropertyOnType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1316 | 22 | 10 |  |
| 183 | 48 | `internal/jscomp/checksuper.go` | `visitSuperPropertyAccess` | `com/google/javascript/jscomp/CheckSuper.java` | 460 | 2 | 1 |  |
| 184 | 48 | `internal/rhino/jsdoc.go` | `HasDescriptionForParameter` | `com/google/javascript/rhino/JSDocInfo.java` | 1393 | 4 | 2 |  |
| 185 | 48 | `internal/rhino/jsdoc.go` | `HasParameter` | `com/google/javascript/rhino/JSDocInfo.java` | 1036 | 4 | 2 |  |
| 186 | 48 | `internal/rhino/jsdoc.go` | `HasSideEffectsArgumentsAnnotation` | `com/google/javascript/rhino/JSDocInfo.java` | 1575 | 4 | 2 |  |
| 187 | 48 | `internal/rhino/jstype/equality_checker.go` | `checkHasNotRun` | `com/google/javascript/rhino/jstype/EqualityChecker.java` | 116 | 2 | 1 |  |
| 188 | 48 | `internal/rhino/jstype/subtype_checker.go` | `checkHasNotRun` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 125 | 2 | 1 |  |
| 189 | 48 | `internal/rhino/jsdoc_info_builder.go` | `HasParameter` | `com/google/javascript/rhino/JSDocInfo.java` | 1036 | 4 | 2 |  |
| 190 | 50 | `internal/rhino/jstype_expr.go` | `ReplaceNamesWithUnknownType` | `com/google/javascript/rhino/JSTypeExpression.java` | 87 | 9 | 4 |  |
| 191 | 50 | `internal/rhino/jsdoc_info_builder.go` | `IsPopulatedWithFileOverview` | `com/google/javascript/rhino/JSDocInfo.java` | 1704 | 7 | 10 |  |
| 192 | 50 | `internal/rhino/jsdoc.go` | `ContainsDeclarationExcludingTypelessConst` | `com/google/javascript/rhino/JSDocInfo.java` | 936 | 17 | 22 |  |
| 193 | 50 | `internal/jscomp/checkjsdoc.go` | `report` | `com/google/javascript/jscomp/CheckJSDoc.java` | 608 | 8 | 1 | much larger (8 LoC vs Java's 1) |
| 194 | 51 | `internal/jscomp/checkjsdoc.go` | `containsOnlySuppressionFor` | `com/google/javascript/jscomp/CheckJSDoc.java` | 211 | 6 | 3 |  |
| 195 | 52 | `internal/rhino/jsdoc.go` | `GetTypeExpressions` | `com/google/javascript/rhino/JSDocInfo.java` | 1504 | 20 | 11 |  |
| 196 | 53 | `internal/rhino/jstype/function_param_builder.go` | `Build` | `com/google/javascript/rhino/jstype/FunctionParamBuilder.java` | 142 | 3 | 1 |  |
| 197 | 53 | `internal/jscomp/typemismatch.go` | `bothAreNotTemplateTypes` | `com/google/javascript/jscomp/TypeMismatch.java` | 133 | 3 | 1 |  |
| 198 | 53 | `internal/jscomp/typemismatch.go` | `GetMismatches` | `com/google/javascript/jscomp/TypeMismatch.java` | 114 | 3 | 1 |  |
| 199 | 54 | `internal/jscomp/parsing/jsdoctokenstream.go` | `isJSDocString` | `com/google/javascript/jscomp/parsing/JsDocTokenStream.java` | 276 | 6 | 22 | likely stub (6 LoC vs Java's 22) |
| 200 | 54 | `internal/jscomp/nodeutil.go` | `MayBeObjectLitKey` | `com/google/javascript/jscomp/NodeUtil.java` | 3229 | 6 | 4 |  |
| 201 | 55 | `internal/jscomp/typedscope.go` | `ValidateCompletelyBuilt` | `com/google/javascript/jscomp/TypedScope.java` | 121 | 2 | 7 | likely stub (2 LoC vs Java's 7) |
| 202 | 58 | `internal/jscomp/nodeutil.go` | `IsDeclarationLValue` | `com/google/javascript/jscomp/NodeUtil.java` | 2958 | 20 | 10 |  |
| 203 | 58 | `internal/rhino/jstype/jstype.go` | `Autobox` | `com/google/javascript/rhino/jstype/JSType.java` | 937 | 6 | 3 |  |
| 204 | 58 | `internal/rhino/jstype/subtype_checker.go` | `SetSubtype` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 102 | 7 | 4 |  |
| 205 | 58 | `internal/rhino/jstype/subtype_checker.go` | `SetSupertype` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 94 | 7 | 4 |  |
| 206 | 58 | `internal/rhino/jstype/function_type.go` | `Build` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1593 | 26 | 11 |  |
| 207 | 58 | `internal/jscomp/checkjsdoc.go` | `validateArrowFunction` | `com/google/javascript/jscomp/CheckJSDoc.java` | 613 | 8 | 4 |  |
| 208 | 58 | `internal/jscomp/parsing/jsdoctokenstream.go` | `ungetChar` | `com/google/javascript/jscomp/parsing/JsDocTokenStream.java` | 250 | 6 | 3 |  |
| 209 | 58 | `internal/jscomp/nodeutil.go` | `IsImportedName` | `com/google/javascript/jscomp/NodeUtil.java` | 2980 | 4 | 2 |  |
| 210 | 59 | `internal/jscomp/checkregexp.go` | `Visit` | `com/google/javascript/jscomp/CheckRegExp.java` | 91 | 1 | 26 | likely stub (1 LoC vs Java's 26) |
| 211 | 59 | `internal/jscomp/checkclosureimports.go` | `checkValidImportCodeReference` | `com/google/javascript/jscomp/CheckClosureImports.java` | 289 | 1 | 18 | likely stub (1 LoC vs Java's 18) |
| 212 | 59 | `internal/jscomp/checksuspiciouscode.go` | `checkLeftOperandOfLogicalOperator` | `com/google/javascript/jscomp/CheckSuspiciousCode.java` | 166 | 15 | 6 |  |
| 213 | 59 | `internal/jscomp/checksuspiciouscode.go` | `reportIfNonObject` | `com/google/javascript/jscomp/CheckSuspiciousCode.java` | 144 | 9 | 4 |  |
| 214 | 59 | `internal/jscomp/nodeutil.go` | `Precedence` | `com/google/javascript/jscomp/NodeUtil.java` | 1219 | 57 | 94 |  |
| 215 | 59 | `internal/jscomp/checkclosureimports.go` | `checkShortName` | `com/google/javascript/jscomp/CheckClosureImports.java` | 420 | 19 | 11 |  |
| 216 | 59 | `internal/jscomp/nodeutil.go` | `GetKnownValueType` | `com/google/javascript/jscomp/NodeUtil.java` | 1355 | 1 | 98 | likely stub (1 LoC vs Java's 98) |
| 217 | 59 | `internal/jscomp/checksuspiciouscode.go` | `checkNegatedLeftOperandOfInOperator` | `com/google/javascript/jscomp/CheckSuspiciousCode.java` | 152 | 5 | 2 |  |
| 218 | 60 | `internal/rhino/jstype/function_type.go` | `GetTypeOfThis` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1028 | 3 | 3 |  |
| 219 | 60 | `internal/rhino/jsdoc.go` | `GetExtendedInterfacesCount` | `com/google/javascript/rhino/JSDocInfo.java` | 1353 | 3 | 2 |  |
| 220 | 60 | `internal/rhino/jstype/function_type.go` | `SetSource` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1046 | 1 | 9 | likely stub (1 LoC vs Java's 9) |
| 221 | 60 | `internal/rhino/jstype_expr.go` | `GetAllTypeNodes` | `com/google/javascript/rhino/JSTypeExpression.java` | 131 | 5 | 3 |  |
| 222 | 60 | `internal/rhino/jsdoc.go` | `GetMarkers` | `com/google/javascript/rhino/JSDocInfo.java` | 1455 | 3 | 2 |  |
| 223 | 60 | `internal/rhino/jsdoc.go` | `GetTsTypes` | `com/google/javascript/rhino/JSDocInfo.java` | 1167 | 3 | 2 |  |
| 224 | 60 | `internal/rhino/jsdoc.go` | `GetThrowsAnnotations` | `com/google/javascript/rhino/JSDocInfo.java` | 1082 | 3 | 2 |  |
| 225 | 60 | `internal/jscomp/checksuper.go` | `visitContextNode` | `com/google/javascript/jscomp/CheckSuper.java` | 410 | 1 | 10 | likely stub (1 LoC vs Java's 10) |
| 226 | 60 | `internal/jscomp/typedscopecreator.go` | `getLValueRootScope` | `com/google/javascript/jscomp/TypedScopeCreator.java` | 2299 | 1 | 21 | likely stub (1 LoC vs Java's 21) |
| 227 | 60 | `internal/jscomp/checksuper.go` | `visitContextNode` | `com/google/javascript/jscomp/CheckSuper.java` | 410 | 1 | 10 | likely stub (1 LoC vs Java's 10) |
| 228 | 60 | `internal/rhino/jstype/registry.go` | `OverwriteDeclaredType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1493 | 3 | 2 |  |
| 229 | 60 | `internal/jscomp/checksuper.go` | `visitContextNode` | `com/google/javascript/jscomp/CheckSuper.java` | 410 | 1 | 10 | likely stub (1 LoC vs Java's 10) |
| 230 | 60 | `internal/rhino/jsdoc.go` | `GetImplementedInterfaceCount` | `com/google/javascript/rhino/JSDocInfo.java` | 1338 | 3 | 2 |  |
| 231 | 60 | `internal/rhino/jsdoc.go` | `GetParameterCount` | `com/google/javascript/rhino/JSDocInfo.java` | 1076 | 3 | 2 |  |
| 232 | 60 | `internal/rhino/jstype_expr.go` | `GetRecordPropertyNames` | `com/google/javascript/rhino/JSTypeExpression.java` | 222 | 5 | 3 |  |
| 233 | 60 | `internal/rhino/jsdoc.go` | `GetOriginalCommentPosition` | `com/google/javascript/rhino/JSDocInfo.java` | 1564 | 3 | 2 |  |
| 234 | 60 | `internal/rhino/jsdoc.go` | `GetVisibility` | `com/google/javascript/rhino/JSDocInfo.java` | 1018 | 3 | 2 |  |
| 235 | 60 | `internal/jscomp/checksuper.go` | `visitContextNode` | `com/google/javascript/jscomp/CheckSuper.java` | 410 | 1 | 10 | likely stub (1 LoC vs Java's 10) |
| 236 | 60 | `internal/rhino/jsdoc.go` | `GetModifies` | `com/google/javascript/rhino/JSDocInfo.java` | 1387 | 3 | 2 |  |
| 237 | 61 | `internal/jscomp/checkmissingrequires.go` | `checkMissingRequireThroughShortNameInGoogScope` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 471 | 5 | 41 | likely stub (5 LoC vs Java's 41) |
| 238 | 61 | `internal/jscomp/checkunreachablecode.go` | `shouldTraverse` | `com/google/javascript/jscomp/CheckUnreachableCode.java` | 48 | 19 | 9 |  |
| 239 | 61 | `internal/jscomp/checkconformance.go` | `validateBehaviorSettingOfConfigs` | `com/google/javascript/jscomp/CheckConformance.java` | 244 | 2 | 31 | likely stub (2 LoC vs Java's 31) |
| 240 | 62 | `internal/jscomp/nodeutil.go` | `NewQName` | `com/google/javascript/jscomp/NodeUtil.java` | 3449 | 1 | 22 | likely stub (1 LoC vs Java's 22) |
| 241 | 62 | `internal/jscomp/checkconformance.go` | `shouldTraverse` | `com/google/javascript/jscomp/CheckConformance.java` | 175 | 3 | 3 |  |
| 242 | 62 | `internal/rhino/jstype/named_type.go` | `resolveViaRegistry` | `com/google/javascript/rhino/jstype/NamedType.java` | 311 | 11 | 6 |  |
| 243 | 63 | `internal/rhino/jsdoc_info_builder.go` | `RecordInterface` | `com/google/javascript/rhino/JSDocInfo.java` | 2570 | 5 | 3 |  |
| 244 | 63 | `internal/jscomp/checksuspiciouscode.go` | `reportIfWasEmpty` | `com/google/javascript/jscomp/CheckSuspiciousCode.java` | 105 | 5 | 3 |  |
| 245 | 63 | `internal/jscomp/checkmissingrequires.go` | `process` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 130 | 3 | 1 |  |
| 246 | 63 | `internal/rhino/jstype/subtype_checker.go` | `isEnumSubtype` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 553 | 3 | 3 |  |
| 247 | 63 | `internal/rhino/jsdoc_info_builder.go` | `RecordDict` | `com/google/javascript/rhino/JSDocInfo.java` | 2537 | 5 | 3 |  |
| 248 | 63 | `internal/rhino/jsdoc_info_builder.go` | `RecordStruct` | `com/google/javascript/rhino/JSDocInfo.java` | 2520 | 5 | 3 |  |
| 249 | 63 | `internal/rhino/jsdoc.go` | `SetType` | `com/google/javascript/rhino/JSDocInfo.java` | 2809 | 1 | 5 | likely stub (1 LoC vs Java's 5) |
| 250 | 63 | `internal/rhino/qualified_name.go` | `Matches` | `com/google/javascript/rhino/QualifiedName.java` | 187 | 1 | 15 | likely stub (1 LoC vs Java's 15) |
| 251 | 64 | `internal/jscomp/checkmissingreturn.go` | `isVoidOrUnknown` | `com/google/javascript/jscomp/CheckMissingReturn.java` | 199 | 3 | 2 |  |
| 252 | 64 | `internal/rhino/jstype/function_type.go` | `SetImplementedInterfaces` | `com/google/javascript/rhino/jstype/FunctionType.java` | 617 | 1 | 4 | likely stub (1 LoC vs Java's 4) |
| 253 | 64 | `internal/rhino/jstype/function_type.go` | `SetExtendedInterfaces` | `com/google/javascript/rhino/jstype/FunctionType.java` | 636 | 1 | 4 | likely stub (1 LoC vs Java's 4) |
| 254 | 64 | `internal/rhino/jsdoc_info_builder.go` | `RecordTypedef` | `com/google/javascript/rhino/JSDocInfo.java` | 2113 | 5 | 4 |  |
| 255 | 64 | `internal/rhino/qualified_name.go` | `Join` | `com/google/javascript/rhino/QualifiedName.java` | 125 | 4 | 3 |  |
| 256 | 64 | `internal/jscomp/nodeutil.go` | `IsHoistedFunctionDeclaration` | `com/google/javascript/jscomp/NodeUtil.java` | 2666 | 12 | 7 |  |
| 257 | 64 | `internal/rhino/jsdoc_info_builder.go` | `RecordAbstract` | `com/google/javascript/rhino/JSDocInfo.java` | 2506 | 6 | 4 |  |
| 258 | 65 | `internal/rhino/jsdoc_info_builder.go` | `RecordConstructor` | `com/google/javascript/rhino/JSDocInfo.java` | 2446 | 5 | 3 |  |
| 259 | 65 | `internal/jscomp/checkmissingrequires.go` | `visitMaybeDeclaration` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 301 | 4 | 38 | likely stub (4 LoC vs Java's 38) |
| 260 | 65 | `internal/jscomp/checkconformance.go` | `mergeRequirements` | `com/google/javascript/jscomp/CheckConformance.java` | 418 | 6 | 48 | likely stub (6 LoC vs Java's 48) |
| 261 | 66 | `internal/rhino/jstype/function_type.go` | `SetPrototypeBasedOn` | `com/google/javascript/rhino/jstype/FunctionType.java` | 476 | 2 | 14 | likely stub (2 LoC vs Java's 14) |
| 262 | 66 | `internal/jscomp/parsing/jsdoctokenstream.go` | `GetJsDocToken` | `com/google/javascript/jscomp/parsing/JsDocTokenStream.java` | 58 | 10 | 99 | likely stub (10 LoC vs Java's 99) |
| 263 | 66 | `internal/jscomp/nodeutil.go` | `Has` | `com/google/javascript/jscomp/NodeUtil.java` | 4073 | 1 | 8 | likely stub (1 LoC vs Java's 8) |
| 264 | 66 | `internal/rhino/jstype/function_type.go` | `GetPrototype` | `com/google/javascript/rhino/jstype/FunctionType.java` | 445 | 1 | 13 | likely stub (1 LoC vs Java's 13) |
| 265 | 66 | `internal/jscomp/checksuper.go` | `isSuperPropertyAccess` | `com/google/javascript/jscomp/CheckSuper.java` | 190 | 5 | 3 |  |
| 266 | 66 | `internal/jscomp/checkmissingrequires.go` | `checkMissingRequireThroughShortNameInModule` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 548 | 5 | 14 |  |
| 267 | 67 | `internal/jscomp/checksuper.go` | `visitSuperConstructorCall` | `com/google/javascript/jscomp/CheckSuper.java` | 318 | 2 | 8 | likely stub (2 LoC vs Java's 8) |
| 268 | 67 | `internal/rhino/jstype/function_type.go` | `ForgetParameterAndReturnTypes` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1293 | 15 | 10 |  |
| 269 | 67 | `internal/rhino/jsdoc_info_builder.go` | `RecordImplicitMatch` | `com/google/javascript/rhino/JSDocInfo.java` | 2459 | 6 | 4 |  |
| 270 | 67 | `internal/rhino/jsdoc_info_builder.go` | `RecordReturnDescription` | `com/google/javascript/rhino/JSDocInfo.java` | 2136 | 7 | 3 |  |
| 271 | 67 | `internal/jscomp/checksuper.go` | `visitSuperConstructorCall` | `com/google/javascript/jscomp/CheckSuper.java` | 318 | 2 | 8 | likely stub (2 LoC vs Java's 8) |
| 272 | 67 | `internal/rhino/jsdoc.go` | `IsAtSignCodePresent` | `com/google/javascript/rhino/JSDocInfo.java` | 1009 | 3 | 2 |  |
| 273 | 67 | `internal/jscomp/nodeutil.go` | `GetBooleanValue` | `com/google/javascript/jscomp/NodeUtil.java` | 113 | 58 | 56 |  |
| 274 | 67 | `internal/rhino/jstype_expr.go` | `baseName` | `com/google/javascript/rhino/JSTypeExpression.java` | 125 | 3 | 2 |  |
| 275 | 67 | `internal/rhino/jstype/subtype_checker.go` | `isSubtypeDispatching` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 191 | 11 | 11 |  |
| 276 | 68 | `internal/rhino/jstype/subtype_checker.go` | `isTemplateSubtype` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 589 | 1 | 7 | likely stub (1 LoC vs Java's 7) |
| 277 | 68 | `internal/rhino/qualified_name.go` | `Matches` | `com/google/javascript/rhino/QualifiedName.java` | 187 | 3 | 15 | likely stub (3 LoC vs Java's 15) |
| 278 | 68 | `internal/jscomp/checksuspiciouscode.go` | `checkMissingSemicolon` | `com/google/javascript/jscomp/CheckSuspiciousCode.java` | 89 | 16 | 10 |  |
| 279 | 68 | `internal/jscomp/checksuper.go` | `isSuperConstructorCall` | `com/google/javascript/jscomp/CheckSuper.java` | 184 | 4 | 3 |  |
| 280 | 68 | `internal/jscomp/checkclosureimports.go` | `checkGoogModuleGet` | `com/google/javascript/jscomp/CheckClosureImports.java` | 319 | 8 | 44 | likely stub (8 LoC vs Java's 44) |
| 281 | 69 | `internal/jscomp/nodeutil.go` | `GetFunctionBody` | `com/google/javascript/jscomp/NodeUtil.java` | 2546 | 3 | 2 |  |
| 282 | 69 | `internal/rhino/jstype/union_type.go` | `GetAlternates` | `com/google/javascript/rhino/jstype/UnionType.java` | 132 | 1 | 5 | likely stub (1 LoC vs Java's 5) |
| 283 | 69 | `internal/rhino/jstype/subtype_checker.go` | `isObjectSubtypeHelper` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 269 | 8 | 49 | likely stub (8 LoC vs Java's 49) |
| 284 | 69 | `internal/jscomp/ast_analyzer.go` | `NodeTypeMayHaveSideEffects` | `com/google/javascript/jscomp/AstAnalyzer.java` | 537 | 30 | 22 |  |
| 285 | 69 | `internal/jscomp/checkconformance.go` | `initRules` | `com/google/javascript/jscomp/CheckConformance.java` | 207 | 3 | 24 | likely stub (3 LoC vs Java's 24) |
| 286 | 69 | `internal/rhino/jstype/equality_checker.go` | `areFunctionEqual` | `com/google/javascript/rhino/jstype/EqualityChecker.java` | 308 | 14 | 13 |  |
| 287 | 69 | `internal/rhino/jstype/jstype.go` | `GetTypesUnderEquality` | `com/google/javascript/rhino/jstype/JSType.java` | 1296 | 12 | 7 |  |
| 288 | 69 | `internal/jscomp/checkglobalthis.go` | `Visit` | `com/google/javascript/jscomp/CheckGlobalThis.java` | 189 | 11 | 4 |  |
| 289 | 70 | `internal/rhino/jsdoc.go` | `ContainsTypeDeclaration` | `com/google/javascript/rhino/JSDocInfo.java` | 962 | 9 | 9 |  |
| 290 | 70 | `internal/rhino/jsdoc.go` | `HasTypeInformation` | `com/google/javascript/rhino/JSDocInfo.java` | 1115 | 3 | 3 |  |
| 291 | 70 | `internal/jscomp/checksuper.go` | `visitReturn` | `com/google/javascript/jscomp/CheckSuper.java` | 403 | 1 | 2 |  |
| 292 | 70 | `internal/jscomp/checksuper.go` | `visitReturn` | `com/google/javascript/jscomp/CheckSuper.java` | 403 | 1 | 2 |  |
| 293 | 70 | `internal/jscomp/checksuper.go` | `visitReturn` | `com/google/javascript/jscomp/CheckSuper.java` | 403 | 1 | 2 |  |
| 294 | 70 | `internal/jscomp/checksuper.go` | `visitReturn` | `com/google/javascript/jscomp/CheckSuper.java` | 403 | 1 | 2 |  |
| 295 | 70 | `internal/jscomp/checkjsdoc.go` | `validateTemplates` | `com/google/javascript/jscomp/CheckJSDoc.java` | 239 | 13 | 11 |  |
| 296 | 70 | `internal/jscomp/ast_analyzer.go` | `checkForStateChangeHelper` | `com/google/javascript/jscomp/AstAnalyzer.java` | 285 | 76 | 104 |  |
| 297 | 70 | `internal/jscomp/nodeutil.go` | `IsFromTypeSummary` | `com/google/javascript/jscomp/NodeUtil.java` | 1063 | 4 | 3 |  |
| 298 | 70 | `internal/jscomp/nodeutil.go` | `IsInvocationTarget` | `com/google/javascript/jscomp/NodeUtil.java` | 4668 | 4 | 4 |  |
| 299 | 70 | `internal/rhino/jsdoc.go` | `ContainsTypeDefinition` | `com/google/javascript/rhino/JSDocInfo.java` | 1003 | 3 | 2 |  |
| 300 | 70 | `internal/jscomp/nodeutil.go` | `IsControlStructure` | `com/google/javascript/jscomp/NodeUtil.java` | 2114 | 9 | 18 |  |
| 301 | 70 | `internal/rhino/token_util.go` | `IsStrWhiteSpaceChar` | `com/google/javascript/rhino/TokenUtil.java` | 77 | 17 | 15 |  |
| 302 | 70 | `internal/jscomp/nodeutil.go` | `IsControlStructureCodeBlock` | `com/google/javascript/jscomp/NodeUtil.java` | 2136 | 15 | 11 |  |
| 303 | 70 | `internal/rhino/jstype/property.go` | `String` | `com/google/javascript/rhino/jstype/Property.java` | 110 | 1 | 1 |  |
| 304 | 71 | `internal/jscomp/nodeutil.go` | `IsGoogModuleGetCall` | `com/google/javascript/jscomp/NodeUtil.java` | 5151 | 1 | 5 | likely stub (1 LoC vs Java's 5) |
| 305 | 71 | `internal/jscomp/checkjsdoc.go` | `validateFunctionJsDoc` | `com/google/javascript/jscomp/CheckJSDoc.java` | 439 | 10 | 7 |  |
| 306 | 72 | `internal/jscomp/checkjsdoc.go` | `validateReturnJsDoc` | `com/google/javascript/jscomp/CheckJSDoc.java` | 721 | 8 | 4 |  |
| 307 | 72 | `internal/jscomp/checkjsdoc.go` | `validateRestParameter` | `com/google/javascript/jscomp/CheckJSDoc.java` | 623 | 26 | 23 |  |
| 308 | 72 | `internal/jscomp/parsing/jsdocinfoparser.go` | `parseParamTypeExpressionAnnotation` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2064 | 2 | 9 | likely stub (2 LoC vs Java's 9) |
| 309 | 72 | `internal/rhino/qualified_name.go` | `GetComponent` | `com/google/javascript/rhino/QualifiedName.java` | 271 | 1 | 6 | likely stub (1 LoC vs Java's 6) |
| 310 | 72 | `internal/jscomp/checkclosureimports.go` | `verifyRequireOrder` | `com/google/javascript/jscomp/CheckClosureImports.java` | 523 | 4 | 19 | likely stub (4 LoC vs Java's 19) |
| 311 | 72 | `internal/rhino/qualified_name.go` | `GetComponentCount` | `com/google/javascript/rhino/QualifiedName.java` | 301 | 1 | 6 | likely stub (1 LoC vs Java's 6) |
| 312 | 72 | `internal/rhino/qualified_name.go` | `GetComponent` | `com/google/javascript/rhino/QualifiedName.java` | 271 | 1 | 6 | likely stub (1 LoC vs Java's 6) |
| 313 | 72 | `internal/jscomp/checkconformance.go` | `validateExtendingRequirement` | `com/google/javascript/jscomp/CheckConformance.java` | 566 | 3 | 18 | likely stub (3 LoC vs Java's 18) |
| 314 | 73 | `internal/jscomp/parsing/jsdocinfoparser.go` | `parseHelperLoop` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 337 | 11 | 40 | likely stub (11 LoC vs Java's 40) |
| 315 | 73 | `internal/rhino/qualified_name.go` | `appendTo` | `com/google/javascript/rhino/QualifiedName.java` | 167 | 1 | 4 | likely stub (1 LoC vs Java's 4) |
| 316 | 73 | `internal/rhino/jstype/function_type.go` | `ClearCachedValues` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1069 | 2 | 6 |  |
| 317 | 73 | `internal/rhino/jsdoc_info_builder.go` | `RecordSuppressions` | `com/google/javascript/rhino/JSDocInfo.java` | 2062 | 3 | 10 |  |
| 318 | 74 | `internal/jscomp/nodeutil.go` | `isDeclarationParent` | `com/google/javascript/jscomp/NodeUtil.java` | 2245 | 6 | 4 |  |
| 319 | 74 | `internal/rhino/token_stream.go` | `IsKeyword` | `com/google/javascript/rhino/TokenStream.java` | 68 | 229 | 125 |  |
| 320 | 74 | `internal/rhino/jstype/named_type.go` | `SetValidator` | `com/google/javascript/rhino/jstype/NamedType.java` | 434 | 2 | 5 |  |
| 321 | 74 | `internal/jscomp/parsing/jsdocinfoparser.go` | `Parse` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 289 | 3 | 12 | likely stub (3 LoC vs Java's 12) |
| 322 | 74 | `internal/rhino/jstype/named_type.go` | `warning` | `com/google/javascript/rhino/jstype/NamedType.java` | 446 | 0 | 1 |  |
| 323 | 74 | `internal/jscomp/checkmissingrequires.go` | `checkMissingRequireThroughShortName` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 430 | 7 | 19 |  |
| 324 | 74 | `internal/jscomp/parsing/jsdocinfoparser.go` | `ParseTypeString` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 244 | 5 | 2 |  |
| 325 | 74 | `internal/jscomp/checksuper.go` | `visitSuperConstructorCall` | `com/google/javascript/jscomp/CheckSuper.java` | 318 | 2 | 8 | likely stub (2 LoC vs Java's 8) |
| 326 | 74 | `internal/jscomp/checksuper.go` | `visitSuperConstructorCall` | `com/google/javascript/jscomp/CheckSuper.java` | 318 | 3 | 8 |  |
| 327 | 75 | `internal/rhino/qualified_name.go` | `Of` | `com/google/javascript/rhino/QualifiedName.java` | 58 | 2 | 11 | likely stub (2 LoC vs Java's 11) |
| 328 | 75 | `internal/jscomp/checkconformance.go` | `removeDuplicates` | `com/google/javascript/jscomp/CheckConformance.java` | 594 | 1 | 13 | likely stub (1 LoC vs Java's 13) |
| 329 | 75 | `internal/jscomp/nodeutil.go` | `NewName` | `com/google/javascript/jscomp/NodeUtil.java` | 3580 | 1 | 5 | likely stub (1 LoC vs Java's 5) |
| 330 | 75 | `internal/rhino/jstype/function_type.go` | `WithTemplateKeys` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1511 | 2 | 7 | likely stub (2 LoC vs Java's 7) |
| 331 | 76 | `internal/jscomp/nodeutil.go` | `IsEs6ConstructorMemberFunctionDef` | `com/google/javascript/jscomp/NodeUtil.java` | 5305 | 1 | 5 | likely stub (1 LoC vs Java's 5) |
| 332 | 76 | `internal/jscomp/checkmissingrequires.go` | `swapRootOfQualifiedName` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 760 | 6 | 7 |  |
| 333 | 77 | `internal/jscomp/parsing/jsdocinfoparser.go` | `parseTopLevelTypeExpression` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2106 | 14 | 8 |  |
| 334 | 77 | `internal/jscomp/checkconformance.go` | `mergeExtendingRequirementIntoExtended` | `com/google/javascript/jscomp/CheckConformance.java` | 508 | 2 | 14 | likely stub (2 LoC vs Java's 14) |
| 335 | 77 | `internal/rhino/jstype/subtype_checker.go` | `isStructuralSubtypeHelper` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 369 | 4 | 15 | likely stub (4 LoC vs Java's 15) |
| 336 | 77 | `internal/jscomp/nodeutil.go` | `NewPropertyAccess` | `com/google/javascript/jscomp/NodeUtil.java` | 3528 | 1 | 4 | likely stub (1 LoC vs Java's 4) |
| 337 | 77 | `internal/rhino/jstype/registry.go` | `CreateUnionType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1770 | 1 | 4 | likely stub (1 LoC vs Java's 4) |
| 338 | 77 | `internal/jscomp/nodeutil.go` | `GetRootOfQualifiedName` | `com/google/javascript/jscomp/NodeUtil.java` | 3562 | 8 | 4 |  |
| 339 | 77 | `internal/rhino/jstype/contains_upper_bound_super_type_visitor.go` | `caseDefault` | `com/google/javascript/rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java` | 68 | 5 | 3 |  |
| 340 | 78 | `internal/jscomp/parsing/jsdoctokenstream.go` | `getChar` | `com/google/javascript/jscomp/parsing/JsDocTokenStream.java` | 315 | 17 | 34 |  |
| 341 | 78 | `internal/rhino/jstype/jstype.go` | `IsNamedType` | `com/google/javascript/rhino/jstype/JSType.java` | 503 | 2 | 1 |  |
| 342 | 78 | `internal/rhino/jstype/jstype.go` | `IsTemplatizedType` | `com/google/javascript/rhino/jstype/JSType.java` | 531 | 2 | 1 |  |
| 343 | 78 | `internal/jscomp/syntacticscopecreator.go` | `declareLHS` | `com/google/javascript/jscomp/SyntacticScopeCreator.java` | 168 | 7 | 4 |  |
| 344 | 78 | `internal/rhino/jstype/subtype_checker.go` | `SetSubtypingMode` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 118 | 5 | 4 |  |
| 345 | 78 | `internal/jscomp/checkjsdoc.go` | `isFromTs` | `com/google/javascript/jscomp/CheckJSDoc.java` | 742 | 2 | 1 |  |
| 346 | 78 | `internal/rhino/jsdoc.go` | `ToStringVerbose` | `com/google/javascript/rhino/JSDocInfo.java` | 1301 | 3 | 13 | likely stub (3 LoC vs Java's 13) |
| 347 | 78 | `internal/rhino/qualified_name.go` | `GetComponent` | `com/google/javascript/rhino/QualifiedName.java` | 271 | 9 | 6 |  |
| 348 | 78 | `internal/jscomp/nodeutil.go` | `IsImmutableValue` | `com/google/javascript/jscomp/NodeUtil.java` | 648 | 15 | 15 |  |
| 349 | 78 | `internal/rhino/qualified_name.go` | `GetComponentCount` | `com/google/javascript/rhino/QualifiedName.java` | 301 | 1 | 6 | likely stub (1 LoC vs Java's 6) |
| 350 | 78 | `internal/jscomp/typedscopecreator.go` | `declareTemplatesInFunctionBody` | `com/google/javascript/jscomp/TypedScopeCreator.java` | 3310 | 6 | 26 | likely stub (6 LoC vs Java's 26) |
| 351 | 78 | `internal/rhino/jstype/function_type.go` | `WithName` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1450 | 2 | 1 |  |
| 352 | 79 | `internal/rhino/jstype/function_type.go` | `GetOwnPropertyNames` | `com/google/javascript/rhino/jstype/FunctionType.java` | 426 | 11 | 7 |  |
| 353 | 79 | `internal/jscomp/nodeutil.go` | `IsObjectDefinePropertiesDefinition` | `com/google/javascript/jscomp/NodeUtil.java` | 3770 | 13 | 7 |  |
| 354 | 80 | `internal/jscomp/checksuper.go` | `visitThis` | `com/google/javascript/jscomp/CheckSuper.java` | 396 | 1 | 2 |  |
| 355 | 80 | `internal/rhino/jstype/named_type.go` | `setReferencedType` | `com/google/javascript/rhino/jstype/NamedType.java` | 551 | 1 | 2 |  |
| 356 | 80 | `internal/rhino/jstype/template_type.go` | `GetBound` | `com/google/javascript/rhino/jstype/TemplateType.java` | 135 | 1 | 1 |  |
| 357 | 80 | `internal/rhino/jstype/property_map.go` | `GetAllKeys` | `com/google/javascript/rhino/jstype/PropertyMap.java` | 213 | 28 | 18 |  |
| 358 | 80 | `internal/jscomp/checkconformance.go` | `shouldCheck` | `com/google/javascript/jscomp/CheckConformance.java` | 128 | 1 | 1 |  |
| 359 | 80 | `internal/jscomp/nodeutil.go` | `IsInSyntheticScript` | `com/google/javascript/jscomp/NodeUtil.java` | 2311 | 1 | 4 | likely stub (1 LoC vs Java's 4) |
| 360 | 80 | `internal/jscomp/nodeutil.go` | `GetCallTargetResolvingIndirectCalls` | `com/google/javascript/jscomp/NodeUtil.java` | 2555 | 8 | 8 |  |
| 361 | 80 | `internal/jscomp/checksuper.go` | `visitThis` | `com/google/javascript/jscomp/CheckSuper.java` | 396 | 1 | 2 |  |
| 362 | 80 | `internal/rhino/jstype/function_type.go` | `SetGoogModuleId` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1461 | 1 | 2 |  |
| 363 | 80 | `internal/jscomp/checksuper.go` | `visitThis` | `com/google/javascript/jscomp/CheckSuper.java` | 396 | 1 | 2 |  |
| 364 | 80 | `internal/jscomp/parsing/jsdoctokenstream.go` | `isAlpha` | `com/google/javascript/jscomp/parsing/JsDocTokenStream.java` | 267 | 1 | 4 | likely stub (1 LoC vs Java's 4) |
| 365 | 80 | `internal/rhino/jstype/function_type.go` | `GetAncestorInterfaces` | `com/google/javascript/rhino/jstype/FunctionType.java` | 587 | 3 | 6 |  |
| 366 | 80 | `internal/jscomp/parsing/jsdocinfoparser.go` | `parseAnnotation` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 398 | 256 | 604 |  |
| 367 | 80 | `internal/rhino/jsdoc.go` | `ParseDocumentation` | `com/google/javascript/rhino/JSDocInfo.java` | 1672 | 1 | 2 |  |
| 368 | 80 | `internal/jscomp/checksuper.go` | `visitThis` | `com/google/javascript/jscomp/CheckSuper.java` | 396 | 1 | 2 |  |
| 369 | 81 | `internal/jscomp/checkmissingrequires.go` | `hasAcceptableRequire` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 800 | 11 | 9 |  |
| 370 | 81 | `internal/rhino/jstype/registry.go` | `CreateObjectType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1946 | 4 | 4 |  |
| 371 | 81 | `internal/rhino/jstype/subtype_checker.go` | `SetUsingStructuralSubtyping` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 110 | 5 | 4 |  |
| 372 | 81 | `internal/rhino/jstype/equality_checker.go` | `SetEqMethod` | `com/google/javascript/rhino/jstype/EqualityChecker.java` | 109 | 5 | 4 |  |
| 373 | 82 | `internal/rhino/jstype/property.go` | `IsOwnedByInterface` | `com/google/javascript/rhino/jstype/Property.java` | 147 | 1 | 3 |  |
| 374 | 82 | `internal/rhino/qualified_name.go` | `Components` | `com/google/javascript/rhino/QualifiedName.java` | 110 | 1 | 3 |  |
| 375 | 82 | `internal/jscomp/nodeutil.go` | `IsMethodDeclaration` | `com/google/javascript/jscomp/NodeUtil.java` | 2631 | 11 | 14 |  |
| 376 | 82 | `internal/rhino/jstype/function_type.go` | `HasAnyTemplateTypesInternal` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1152 | 1 | 3 |  |
| 377 | 83 | `internal/jscomp/typedscopecreator.go` | `CreateScope` | `com/google/javascript/jscomp/TypedScopeCreator.java` | 361 | 3 | 9 |  |
| 378 | 83 | `internal/jscomp/nodeutil.go` | `IsCallTo` | `com/google/javascript/jscomp/NodeUtil.java` | 5347 | 1 | 3 |  |
| 379 | 83 | `internal/jscomp/parsing/jsdocinfoparser.go` | `parseSuppressTag` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 1405 | 14 | 32 |  |
| 380 | 83 | `internal/rhino/jstype/function_type.go` | `ToBuilder` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1398 | 10 | 20 |  |
| 381 | 83 | `internal/jscomp/nodeutil.go` | `IsUndefined` | `com/google/javascript/jscomp/NodeUtil.java` | 1322 | 6 | 5 |  |
| 382 | 84 | `internal/jscomp/nodeutil.go` | `GetEnclosingNode` | `com/google/javascript/jscomp/NodeUtil.java` | 1754 | 6 | 4 |  |
| 383 | 84 | `internal/jscomp/checkmissingreturn.go` | `isEmptyFunction` | `com/google/javascript/jscomp/CheckMissingReturn.java` | 191 | 2 | 1 |  |
| 384 | 84 | `internal/jscomp/nodeutil.go` | `GetSourceName` | `com/google/javascript/jscomp/NodeUtil.java` | 4411 | 7 | 5 |  |
| 385 | 84 | `internal/rhino/jsdoc_info_builder.go` | `RecordParameterDescription` | `com/google/javascript/rhino/JSDocInfo.java` | 1912 | 7 | 4 |  |
| 386 | 84 | `internal/jscomp/nodeutil.go` | `GetNodeTypeReferenceCount` | `com/google/javascript/jscomp/NodeUtil.java` | 4050 | 1 | 1 |  |
| 387 | 84 | `internal/jscomp/parsing/jsdocinfoparser.go` | `parseTypeName` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2291 | 20 | 12 |  |
| 388 | 84 | `internal/rhino/jstype/registry.go` | `isObjectLiteralThatCanBeSkipped` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1282 | 3 | 2 |  |
| 389 | 84 | `internal/rhino/jsdoc_info_builder.go` | `RecordSuppression` | `com/google/javascript/rhino/JSDocInfo.java` | 2086 | 2 | 1 |  |
| 390 | 84 | `internal/rhino/qualified_name.go` | `appendTo` | `com/google/javascript/rhino/QualifiedName.java` | 167 | 3 | 4 |  |
| 391 | 84 | `internal/jscomp/checkjsdoc.go` | `reportMisplaced` | `com/google/javascript/jscomp/CheckJSDoc.java` | 604 | 2 | 1 |  |
| 392 | 84 | `internal/jscomp/syntacticscopecreator.go` | `isShadowingAllowed` | `com/google/javascript/jscomp/SyntacticScopeCreator.java` | 376 | 7 | 4 |  |
| 393 | 84 | `internal/jscomp/nodeutil.go` | `IsLiteralValue` | `com/google/javascript/jscomp/NodeUtil.java` | 735 | 20 | 43 |  |
| 394 | 85 | `internal/jscomp/checkconformance.go` | `validateExtendableRequirement` | `com/google/javascript/jscomp/CheckConformance.java` | 533 | 5 | 11 |  |
| 395 | 85 | `internal/jscomp/checkduplicatemember.go` | `visit` | `com/google/javascript/jscomp/StrictModeCheck.java` | 103 | 6 | 10 |  |
| 396 | 86 | `internal/rhino/jsdoc.go` | `ContainsFunctionDeclaration` | `com/google/javascript/rhino/JSDocInfo.java` | 989 | 7 | 7 |  |
| 397 | 86 | `internal/rhino/jstype/property_map.go` | `getSecondaryParentObjects` | `com/google/javascript/rhino/jstype/PropertyMap.java` | 146 | 3 | 5 |  |
| 398 | 86 | `internal/rhino/jstype/function_type.go` | `isDelegateProxy` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1383 | 2 | 1 |  |
| 399 | 86 | `internal/rhino/jstype/template_type.go` | `SetBound` | `com/google/javascript/rhino/jstype/TemplateType.java` | 139 | 1 | 2 |  |
| 400 | 86 | `internal/jscomp/nodeutil.go` | `MayBeString` | `com/google/javascript/jscomp/NodeUtil.java` | 1554 | 14 | 21 |  |
| 401 | 86 | `internal/rhino/jstype/function_type.go` | `GetInstanceType` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1015 | 1 | 2 |  |
| 402 | 86 | `internal/rhino/jstype/jstype.go` | `IsSymbol` | `com/google/javascript/rhino/jstype/JSType.java` | 287 | 1 | 2 |  |
| 403 | 86 | `internal/jscomp/nodeutil.go` | `IsEs6Constructor` | `com/google/javascript/jscomp/NodeUtil.java` | 5314 | 2 | 4 |  |
| 404 | 86 | `internal/rhino/jstype/function_type.go` | `SetImplicitMatch` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1167 | 1 | 2 |  |
| 405 | 86 | `internal/jscomp/nodeutil.go` | `IteratesImpureIterable` | `com/google/javascript/jscomp/NodeUtil.java` | 1089 | 16 | 21 |  |
| 406 | 87 | `internal/rhino/jsdoc_info_builder.go` | `RecordTemplateTypeName` | `com/google/javascript/rhino/JSDocInfo.java` | 1930 | 8 | 7 |  |
| 407 | 87 | `internal/rhino/qualified_name.go` | `Matches` | `com/google/javascript/rhino/QualifiedName.java` | 187 | 18 | 15 |  |
| 408 | 87 | `internal/rhino/jstype/jstype.go` | `GetTypesUnderInequality` | `com/google/javascript/rhino/jstype/JSType.java` | 1322 | 10 | 7 |  |
| 409 | 88 | `internal/jscomp/nodeutil.go` | `IsAssignmentOp` | `com/google/javascript/jscomp/NodeUtil.java` | 1636 | 8 | 20 |  |
| 410 | 88 | `internal/rhino/jsdoc_info_builder.go` | `RecordSassGeneratedCssTs` | `com/google/javascript/rhino/JSDocInfo.java` | 2760 | 2 | 1 |  |
| 411 | 88 | `internal/jscomp/checkmissingrequires.go` | `visitMaybeGoogModuleGet` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 283 | 15 | 10 |  |
| 412 | 88 | `internal/rhino/jsdoc_info_builder.go` | `RecordThrowsAnnotation` | `com/google/javascript/rhino/JSDocInfo.java` | 1957 | 5 | 9 |  |
| 413 | 88 | `internal/jscomp/checksuper.go` | `shouldTraverse` | `com/google/javascript/jscomp/CheckSuper.java` | 98 | 15 | 24 |  |
| 414 | 88 | `internal/jscomp/checksuper.go` | `visitSuperPropertyAccess` | `com/google/javascript/jscomp/CheckSuper.java` | 460 | 2 | 1 |  |
| 415 | 88 | `internal/jscomp/checkjsdoc.go` | `isNameInitializeWithClass` | `com/google/javascript/jscomp/CheckJSDoc.java` | 293 | 2 | 1 |  |
| 416 | 89 | `internal/rhino/jstype/subtype_checker.go` | `isFunctionSubtype` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 428 | 17 | 13 |  |
| 417 | 89 | `internal/jscomp/nodeutil.go` | `hasCatchHandler` | `com/google/javascript/jscomp/NodeUtil.java` | 4299 | 1 | 2 |  |
| 418 | 89 | `internal/rhino/jstype/union_type.go` | `TestForEquality` | `com/google/javascript/rhino/jstype/UnionType.java` | 263 | 13 | 9 |  |
| 419 | 89 | `internal/jscomp/checkmissingreturn.go` | `goesThroughTrueCondition` | `com/google/javascript/jscomp/CheckMissingReturn.java` | 44 | 12 | 10 |  |
| 420 | 89 | `internal/rhino/jstype_expr.go` | `Evaluate` | `com/google/javascript/rhino/JSTypeExpression.java` | 179 | 1 | 3 |  |
| 421 | 89 | `internal/rhino/jstype/function_type.go` | `AcceptsArguments` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1267 | 15 | 13 |  |
| 422 | 89 | `internal/jscomp/checkmissingrequires.go` | `isAllowedNamespace` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 773 | 3 | 5 |  |
| 423 | 89 | `internal/jscomp/checkjsdoc.go` | `validateSuppress` | `com/google/javascript/jscomp/CheckJSDoc.java` | 166 | 16 | 30 |  |
| 424 | 90 | `internal/jscomp/nodeutil.go` | `GetClassMembers` | `com/google/javascript/jscomp/NodeUtil.java` | 628 | 1 | 2 |  |
| 425 | 90 | `internal/rhino/jsdoc_info_builder.go` | `IsInterfaceRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2638 | 1 | 1 |  |
| 426 | 90 | `internal/jscomp/nodeutil.go` | `IsNamespaceDecl` | `com/google/javascript/jscomp/NodeUtil.java` | 1014 | 51 | 32 |  |
| 427 | 90 | `internal/rhino/jsdoc.go` | `HasBaseType` | `com/google/javascript/rhino/JSDocInfo.java` | 1324 | 1 | 1 |  |
| 428 | 90 | `internal/rhino/jsdoc.go` | `HasModifies` | `com/google/javascript/rhino/JSDocInfo.java` | 1552 | 1 | 1 |  |
| 429 | 90 | `internal/jscomp/typemismatch.go` | `GetRequired` | `com/google/javascript/jscomp/TypeMismatch.java` | 56 | 1 | 1 |  |
| 430 | 90 | `internal/rhino/jstype/union_type.go` | `MatchesStringContext` | `com/google/javascript/rhino/jstype/UnionType.java` | 174 | 1 | 1 |  |
| 431 | 90 | `internal/rhino/error_reporter.go` | `Error` | `com/google/javascript/rhino/ErrorReporter.java` | 101 | 1 | 1 |  |
| 432 | 90 | `internal/rhino/error_reporter.go` | `Error` | `com/google/javascript/rhino/ErrorReporter.java` | 101 | 1 | 1 |  |
| 433 | 90 | `internal/rhino/jstype/equality_checker.go` | `areRecordEqual` | `com/google/javascript/rhino/jstype/EqualityChecker.java` | 377 | 9 | 8 |  |
| 434 | 90 | `internal/rhino/jstype/function_type.go` | `WithConstructorTemplateKeys` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1533 | 2 | 2 |  |
| 435 | 90 | `internal/rhino/jstype/function_type.go` | `ForInterface` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1551 | 3 | 3 |  |
| 436 | 90 | `internal/rhino/jsdoc_info_builder.go` | `hasAnySingletonSideEffectTags` | `com/google/javascript/rhino/JSDocInfo.java` | 2856 | 1 | 1 |  |
| 437 | 90 | `internal/rhino/jsdoc_info_builder.go` | `hasAnySingletonTypeTags` | `com/google/javascript/rhino/JSDocInfo.java` | 2846 | 3 | 3 |  |
| 438 | 90 | `internal/rhino/jsdoc_info_builder.go` | `IsDeprecationReasonRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2058 | 1 | 1 |  |
| 439 | 90 | `internal/rhino/jstype/union_type.go` | `CollapseUnion` | `com/google/javascript/rhino/jstype/UnionType.java` | 571 | 20 | 20 |  |
| 440 | 90 | `internal/rhino/jstype_expr.go` | `IsOptionalArg` | `com/google/javascript/rhino/JSTypeExpression.java` | 169 | 1 | 1 |  |
| 441 | 90 | `internal/rhino/jstype/template_type.go` | `RecursionUnsafeHashCode` | `com/google/javascript/rhino/jstype/TemplateType.java` | 165 | 1 | 1 |  |
| 442 | 90 | `internal/rhino/jstype_expr.go` | `IsVarArgs` | `com/google/javascript/rhino/JSTypeExpression.java` | 174 | 1 | 1 |  |
| 443 | 90 | `internal/rhino/jsdoc.go` | `HasReturnType` | `com/google/javascript/rhino/JSDocInfo.java` | 1104 | 1 | 1 |  |
| 444 | 90 | `internal/rhino/jsdoc.go` | `MakesStructs` | `com/google/javascript/rhino/JSDocInfo.java` | 801 | 1 | 1 |  |
| 445 | 90 | `internal/rhino/jsdoc.go` | `OverwriteVisibility` | `com/google/javascript/rhino/JSDocInfo.java` | 1890 | 2 | 2 |  |
| 446 | 90 | `internal/rhino/jsdoc.go` | `HasTypedefType` | `com/google/javascript/rhino/JSDocInfo.java` | 1099 | 1 | 1 |  |
| 447 | 90 | `internal/rhino/jsdoc.go` | `HasEnumParameterType` | `com/google/javascript/rhino/JSDocInfo.java` | 1091 | 1 | 1 |  |
| 448 | 90 | `internal/rhino/jstype/function_type.go` | `SetPrototype` | `com/google/javascript/rhino/jstype/FunctionType.java` | 521 | 6 | 5 |  |
| 449 | 90 | `internal/rhino/jstype/function_type.go` | `WithParameters` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1473 | 2 | 2 |  |
| 450 | 90 | `internal/rhino/jsdoc.go` | `MakesDicts` | `com/google/javascript/rhino/JSDocInfo.java` | 806 | 1 | 1 |  |
| 451 | 90 | `internal/jscomp/typemismatch.go` | `GetLocation` | `com/google/javascript/jscomp/TypeMismatch.java` | 61 | 1 | 1 |  |
| 452 | 90 | `internal/rhino/jsdoc.go` | `IsAbstract` | `com/google/javascript/rhino/JSDocInfo.java` | 786 | 1 | 1 |  |
| 453 | 90 | `internal/jscomp/checksuper.go` | `getContextForArrowFunctionNode` | `com/google/javascript/jscomp/CheckSuper.java` | 374 | 1 | 1 |  |
| 454 | 90 | `internal/rhino/jstype/function_type.go` | `HasInstanceType` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1022 | 1 | 1 |  |
| 455 | 90 | `internal/jscomp/checksuper.go` | `getContextForArrowFunctionNode` | `com/google/javascript/jscomp/CheckSuper.java` | 374 | 1 | 1 |  |
| 456 | 90 | `internal/rhino/jsdoc.go` | `IsNoSideEffects` | `com/google/javascript/rhino/JSDocInfo.java` | 848 | 1 | 1 |  |
| 457 | 90 | `internal/rhino/jstype/function_type.go` | `GetExtendedInterfacesCount` | `com/google/javascript/rhino/jstype/FunctionType.java` | 632 | 1 | 1 |  |
| 458 | 90 | `internal/rhino/jstype/function_type.go` | `HasCachedValues` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1084 | 1 | 1 |  |
| 459 | 90 | `internal/rhino/jsdoc.go` | `HasConstAnnotation` | `com/google/javascript/rhino/JSDocInfo.java` | 771 | 1 | 1 |  |
| 460 | 90 | `internal/jscomp/checksuper.go` | `visitSuperPropertyAccess` | `com/google/javascript/jscomp/CheckSuper.java` | 460 | 1 | 1 |  |
| 461 | 90 | `internal/rhino/jstype/function_param_builder.go` | `NewParameterFrom` | `com/google/javascript/rhino/jstype/FunctionParamBuilder.java` | 110 | 1 | 1 |  |
| 462 | 90 | `internal/jscomp/typemismatch.go` | `GetFound` | `com/google/javascript/jscomp/TypeMismatch.java` | 51 | 1 | 1 |  |
| 463 | 90 | `internal/jscomp/checksuper.go` | `getContextForArrowFunctionNode` | `com/google/javascript/jscomp/CheckSuper.java` | 374 | 1 | 1 |  |
| 464 | 90 | `internal/rhino/jsdoc.go` | `IsMappedIdGenerator` | `com/google/javascript/rhino/JSDocInfo.java` | 753 | 1 | 1 |  |
| 465 | 90 | `internal/rhino/jsdoc.go` | `IsXidGenerator` | `com/google/javascript/rhino/JSDocInfo.java` | 748 | 1 | 1 |  |
| 466 | 90 | `internal/jscomp/nodeutil.go` | `IsSomeCompileTimeConstStringValue` | `com/google/javascript/jscomp/NodeUtil.java` | 860 | 19 | 12 |  |
| 467 | 90 | `internal/rhino/jstype/function_type.go` | `GetMaxArity` | `com/google/javascript/rhino/jstype/FunctionType.java` | 384 | 7 | 7 |  |
| 468 | 90 | `internal/rhino/jstype/function_type.go` | `GetMinArity` | `com/google/javascript/rhino/jstype/FunctionType.java` | 365 | 7 | 7 |  |
| 469 | 90 | `internal/rhino/jsdoc.go` | `IsStableIdGenerator` | `com/google/javascript/rhino/JSDocInfo.java` | 743 | 1 | 1 |  |
| 470 | 90 | `internal/rhino/jsdoc.go` | `IsConsistentIdGenerator` | `com/google/javascript/rhino/JSDocInfo.java` | 738 | 1 | 1 |  |
| 471 | 90 | `internal/rhino/jsdoc.go` | `IsIdGenerator` | `com/google/javascript/rhino/JSDocInfo.java` | 758 | 1 | 1 |  |
| 472 | 90 | `internal/rhino/jsdoc.go` | `RecordVisibility` | `com/google/javascript/rhino/JSDocInfo.java` | 1881 | 5 | 5 |  |
| 473 | 90 | `internal/rhino/jstype/registry.go` | `CreateTemplateType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 2013 | 1 | 1 |  |
| 474 | 90 | `internal/rhino/jsdoc.go` | `IsConstructor` | `com/google/javascript/rhino/JSDocInfo.java` | 781 | 1 | 1 |  |
| 475 | 90 | `internal/rhino/jstype/registry.go` | `CreateOptionalNullableType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1755 | 1 | 2 |  |
| 476 | 90 | `internal/jscomp/checksuper.go` | `getContextForArrowFunctionNode` | `com/google/javascript/jscomp/CheckSuper.java` | 374 | 1 | 1 |  |
| 477 | 90 | `internal/jscomp/checkconformance.go` | `reportInvalidRequirement` | `com/google/javascript/jscomp/CheckConformance.java` | 665 | 2 | 3 |  |
| 478 | 90 | `internal/jscomp/nodeutil.go` | `getCatchBlock` | `com/google/javascript/jscomp/NodeUtil.java` | 4290 | 1 | 2 |  |
| 479 | 90 | `internal/rhino/jsdoc.go` | `IsDeprecated` | `com/google/javascript/rhino/JSDocInfo.java` | 824 | 1 | 1 |  |
| 480 | 90 | `internal/rhino/jsdoc_info_builder.go` | `RecordBlockDescription` | `com/google/javascript/rhino/JSDocInfo.java` | 1867 | 6 | 4 |  |
| 481 | 90 | `internal/jscomp/checksuper.go` | `visitSuperPropertyAccess` | `com/google/javascript/jscomp/CheckSuper.java` | 460 | 1 | 1 |  |
| 482 | 90 | `internal/rhino/jsdoc_info_builder.go` | `OverwriteVisibility` | `com/google/javascript/rhino/JSDocInfo.java` | 1890 | 3 | 2 |  |
| 483 | 90 | `internal/rhino/jsdoc.go` | `IsFinal` | `com/google/javascript/rhino/JSDocInfo.java` | 776 | 1 | 1 |  |
| 484 | 90 | `internal/jscomp/checksuper.go` | `getContextForArrowFunctionNode` | `com/google/javascript/jscomp/CheckSuper.java` | 374 | 1 | 1 |  |
| 485 | 90 | `internal/rhino/jsdoc.go` | `HasType` | `com/google/javascript/rhino/JSDocInfo.java` | 1111 | 1 | 1 |  |
| 486 | 90 | `internal/rhino/jstype/function_type.go` | `IsStructuralInterface` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1173 | 1 | 1 |  |
| 487 | 90 | `internal/rhino/jsdoc_info_builder.go` | `IsDescriptionRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 1718 | 1 | 1 |  |
| 488 | 90 | `internal/rhino/jstype/registry.go` | `GetNativeType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1708 | 1 | 1 |  |
| 489 | 90 | `internal/rhino/jstype/registry.go` | `RegisterNativeType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1275 | 1 | 1 |  |
| 490 | 90 | `internal/jscomp/checksuper.go` | `createMemberFieldDefContext` | `com/google/javascript/jscomp/CheckSuper.java` | 254 | 1 | 1 |  |
| 491 | 90 | `internal/rhino/jstype/property.go` | `HumanReadableName` | `com/google/javascript/rhino/jstype/Property.java` | 120 | 1 | 1 |  |
| 492 | 90 | `internal/rhino/jsdoc_info_builder.go` | `RecordLicense` | `com/google/javascript/rhino/JSDocInfo.java` | 2326 | 3 | 3 |  |
| 493 | 90 | `internal/rhino/jsdoc.go` | `HasThisType` | `com/google/javascript/rhino/JSDocInfo.java` | 1152 | 1 | 1 |  |
| 494 | 90 | `internal/rhino/jsdoc.go` | `IsOverride` | `com/google/javascript/rhino/JSDocInfo.java` | 819 | 1 | 1 |  |
| 495 | 90 | `internal/rhino/jstype/registry.go` | `IsForwardDeclaredType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1499 | 1 | 1 |  |
| 496 | 90 | `internal/rhino/jstype/function_param_builder.go` | `hasOptionalOrVarArgs` | `com/google/javascript/rhino/jstype/FunctionParamBuilder.java` | 130 | 4 | 4 |  |
| 497 | 90 | `internal/rhino/jstype/function_type.go` | `GetReturnType` | `com/google/javascript/rhino/jstype/FunctionType.java` | 397 | 1 | 1 |  |
| 498 | 90 | `internal/rhino/jstype/property_map.go` | `incrementCachedKeySetCounter` | `com/google/javascript/rhino/jstype/PropertyMap.java` | 325 | 2 | 3 |  |
| 499 | 90 | `internal/rhino/qualified_name.go` | `IsSimple` | `com/google/javascript/rhino/QualifiedName.java` | 281 | 1 | 1 |  |
| 500 | 90 | `internal/rhino/jstype/function_type.go` | `GetParameters` | `com/google/javascript/rhino/jstype/FunctionType.java` | 360 | 1 | 1 |  |
| 501 | 90 | `internal/rhino/jstype/function_param_builder.go` | `HasVarArgs` | `com/google/javascript/rhino/jstype/FunctionParamBuilder.java` | 138 | 1 | 1 |  |
| 502 | 90 | `internal/rhino/qualified_name.go` | `IsSimple` | `com/google/javascript/rhino/QualifiedName.java` | 281 | 1 | 1 |  |
| 503 | 90 | `internal/rhino/jstype/property.go` | `Matches` | `com/google/javascript/rhino/jstype/Property.java` | 76 | 1 | 1 |  |
| 504 | 90 | `internal/rhino/error_reporter.go` | `Warning` | `com/google/javascript/rhino/ErrorReporter.java` | 96 | 1 | 1 |  |
| 505 | 90 | `internal/jscomp/checksuper.go` | `createStaticBlockContext` | `com/google/javascript/jscomp/CheckSuper.java` | 258 | 1 | 1 |  |
| 506 | 90 | `internal/rhino/qualified_name.go` | `IsSimple` | `com/google/javascript/rhino/QualifiedName.java` | 281 | 1 | 1 |  |
| 507 | 90 | `internal/rhino/jsdoc.go` | `ModifiesThis` | `com/google/javascript/rhino/JSDocInfo.java` | 1570 | 1 | 1 |  |
| 508 | 90 | `internal/rhino/jsdoc_info_builder.go` | `IsConstructorRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2478 | 1 | 1 |  |
| 509 | 90 | `internal/rhino/error_reporter.go` | `Warning` | `com/google/javascript/rhino/ErrorReporter.java` | 96 | 1 | 1 |  |
| 510 | 90 | `internal/rhino/jsdoc_info_builder.go` | `IsStructRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2526 | 1 | 1 |  |
| 511 | 90 | `internal/rhino/jsdoc_info_builder.go` | `IsDictRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2543 | 1 | 1 |  |
| 512 | 90 | `internal/rhino/jsdoc.go` | `HasClosurePrimitiveId` | `com/google/javascript/rhino/JSDocInfo.java` | 1220 | 1 | 1 |  |
| 513 | 90 | `internal/rhino/jsdoc.go` | `HasLendsName` | `com/google/javascript/rhino/JSDocInfo.java` | 1210 | 1 | 1 |  |
| 514 | 90 | `internal/rhino/jstype/property.go` | `GetScope` | `com/google/javascript/rhino/jstype/Property.java` | 272 | 1 | 1 |  |
| 515 | 90 | `internal/rhino/jstype/property.go` | `Matches` | `com/google/javascript/rhino/jstype/Property.java` | 76 | 1 | 1 |  |
| 516 | 90 | `internal/rhino/jstype/property.go` | `GetOwnerInstanceType` | `com/google/javascript/rhino/jstype/Property.java` | 143 | 1 | 1 |  |
| 517 | 90 | `internal/rhino/jstype/property.go` | `IsFromExterns` | `com/google/javascript/rhino/jstype/Property.java` | 242 | 1 | 1 |  |
| 518 | 91 | `internal/rhino/jsdoc_info_builder.go` | `RecordTypeTransformation` | `com/google/javascript/rhino/JSDocInfo.java` | 1943 | 6 | 4 |  |
| 519 | 91 | `internal/rhino/jstype/subtype_checker.go` | `isSubtypeHelper` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 215 | 27 | 24 |  |
| 520 | 91 | `internal/jscomp/checkconformance.go` | `visit` | `com/google/javascript/jscomp/CheckConformance.java` | 187 | 6 | 7 |  |
| 521 | 91 | `internal/rhino/jstype/equality_checker.go` | `areUnionEqual` | `com/google/javascript/rhino/jstype/EqualityChecker.java` | 270 | 11 | 13 |  |
| 522 | 91 | `internal/jscomp/checkjsdoc.go` | `validateTypeAnnotations` | `com/google/javascript/jscomp/CheckJSDoc.java` | 544 | 33 | 28 |  |
| 523 | 91 | `internal/jscomp/nodeutil.go` | `IsLValue` | `com/google/javascript/jscomp/NodeUtil.java` | 2930 | 22 | 15 |  |
| 524 | 91 | `internal/rhino/jstype/function_type.go` | `addRelatedInterfaces` | `com/google/javascript/rhino/jstype/FunctionType.java` | 570 | 13 | 8 |  |
| 525 | 91 | `internal/rhino/jstype/function_type.go` | `getCallOrBindSignature` | `com/google/javascript/rhino/jstype/FunctionType.java` | 709 | 25 | 25 |  |
| 526 | 91 | `internal/jscomp/nodeutil.go` | `IsPrototypeProperty` | `com/google/javascript/jscomp/NodeUtil.java` | 3836 | 6 | 4 |  |
| 527 | 92 | `internal/rhino/jstype/subtype_checker.go` | `isEnumElementSubtype` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 559 | 5 | 8 |  |
| 528 | 92 | `internal/rhino/jsdoc_info_builder.go` | `hasAnyTypeRelatedTags` | `com/google/javascript/rhino/JSDocInfo.java` | 2821 | 8 | 8 |  |
| 529 | 92 | `internal/jscomp/typedscopecreator.go` | `inferTypeForDestructuredTarget` | `com/google/javascript/jscomp/TypedScopeCreator.java` | 1270 | 11 | 13 |  |
| 530 | 92 | `internal/jscomp/checkjsdoc.go` | `validateNoCollapse` | `com/google/javascript/jscomp/CheckJSDoc.java` | 411 | 23 | 17 |  |
| 531 | 92 | `internal/rhino/jstype/equality_checker.go` | `areEqualInternal` | `com/google/javascript/rhino/jstype/EqualityChecker.java` | 170 | 38 | 51 |  |
| 532 | 92 | `internal/rhino/jstype/named_type.go` | `GetBangType` | `com/google/javascript/rhino/jstype/NamedType.java` | 157 | 7 | 7 |  |
| 533 | 93 | `internal/jscomp/nodeutil.go` | `IsFunctionBlock` | `com/google/javascript/jscomp/NodeUtil.java` | 2699 | 1 | 1 |  |
| 534 | 93 | `internal/rhino/jstype_expr.go` | `MakeOptionalArg` | `com/google/javascript/rhino/JSTypeExpression.java` | 158 | 10 | 6 |  |
| 535 | 93 | `internal/jscomp/parsing/jsdocinfoparser.go` | `lookAheadFor` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2848 | 6 | 13 |  |
| 536 | 93 | `internal/rhino/jstype/contains_upper_bound_super_type_visitor.go` | `caseForwardingType` | `com/google/javascript/rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java` | 107 | 6 | 7 |  |
| 537 | 93 | `internal/rhino/jstype/function_type.go` | `SupAndInfHelper` | `com/google/javascript/rhino/jstype/FunctionType.java` | 777 | 21 | 24 |  |
| 538 | 93 | `internal/rhino/jstype/function_param_builder.go` | `newParameter` | `com/google/javascript/rhino/jstype/FunctionParamBuilder.java` | 125 | 1 | 2 |  |
| 539 | 93 | `internal/jscomp/checktypeimportcodereferences.go` | `visit` | `com/google/javascript/jscomp/CheckTypeImportCodeReferences.java` | 45 | 20 | 21 |  |
| 540 | 93 | `internal/rhino/jstype/union_type.go` | `GetPropertyKind` | `com/google/javascript/rhino/jstype/UnionType.java` | 372 | 17 | 17 |  |
| 541 | 93 | `internal/jscomp/nodeutil.go` | `IsNonlocalModuleExportName` | `com/google/javascript/jscomp/NodeUtil.java` | 2281 | 13 | 10 |  |
| 542 | 93 | `internal/jscomp/parsing/jsdocinfoparser.go` | `parseTypeExpression` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2152 | 26 | 32 |  |
| 543 | 93 | `internal/jscomp/checkclosureimports.go` | `isValidDestructuringImport` | `com/google/javascript/jscomp/CheckClosureImports.java` | 401 | 7 | 12 |  |
| 544 | 93 | `internal/rhino/jstype/jstype.go` | `GetRestrictedTypeGivenOutcome` | `com/google/javascript/rhino/jstype/JSType.java` | 1239 | 17 | 13 |  |
| 545 | 93 | `internal/jscomp/checksuper.go` | `process` | `com/google/javascript/jscomp/CheckSuper.java` | 89 | 2 | 2 |  |
| 546 | 93 | `internal/jscomp/typemismatch.go` | `removeNullUndefinedAndTemplates` | `com/google/javascript/jscomp/TypeMismatch.java` | 137 | 4 | 5 |  |
| 547 | 93 | `internal/rhino/jstype_expr.go` | `IsExplicitUnknownTemplateBound` | `com/google/javascript/rhino/JSTypeExpression.java` | 212 | 4 | 3 |  |
| 548 | 93 | `internal/jscomp/checkmissingrequires.go` | `visitQualifiedName` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 387 | 7 | 12 |  |
| 549 | 93 | `internal/rhino/jstype/registry.go` | `CreateFunctionType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1794 | 1 | 1 |  |
| 550 | 94 | `internal/jscomp/nodeutil.go` | `IsLogicalAssignmentOp` | `com/google/javascript/jscomp/NodeUtil.java` | 1660 | 4 | 4 |  |
| 551 | 94 | `internal/rhino/jstype/function_type.go` | `updatePropertyTypeMap` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1195 | 17 | 18 |  |
| 552 | 94 | `internal/rhino/jstype/function_type.go` | `GetImplementedInterfaces` | `com/google/javascript/rhino/jstype/FunctionType.java` | 598 | 9 | 9 |  |
| 553 | 94 | `internal/jscomp/checkmissingreturn.go` | `enterScopeWithCfg` | `com/google/javascript/jscomp/CheckMissingReturn.java` | 72 | 25 | 21 |  |
| 554 | 94 | `internal/jscomp/checkjsdoc.go` | `validateTypedefs` | `com/google/javascript/jscomp/CheckJSDoc.java` | 217 | 16 | 16 |  |
| 555 | 94 | `internal/rhino/jstype/function_param_builder.go` | `NewOptionalParameterFrom` | `com/google/javascript/rhino/jstype/FunctionParamBuilder.java` | 115 | 5 | 5 |  |
| 556 | 94 | `internal/jscomp/nodeutil.go` | `IsLoopStructure` | `com/google/javascript/jscomp/NodeUtil.java` | 2077 | 4 | 4 |  |
| 557 | 94 | `internal/rhino/jstype/property.go` | `GetName` | `com/google/javascript/rhino/jstype/Property.java` | 204 | 7 | 5 |  |
| 558 | 94 | `internal/jscomp/checkjsdoc.go` | `validateAbstractJsDoc` | `com/google/javascript/jscomp/CheckJSDoc.java` | 321 | 28 | 55 |  |
| 559 | 94 | `internal/jscomp/ast_analyzer.go` | `ConstructorCallHasSideEffects` | `com/google/javascript/jscomp/AstAnalyzer.java` | 509 | 8 | 9 |  |
| 560 | 94 | `internal/rhino/jstype_expr.go` | `getRecordPropertyNamesRecursive` | `com/google/javascript/rhino/JSTypeExpression.java` | 228 | 6 | 6 |  |
| 561 | 94 | `internal/jscomp/nodeutil.go` | `IsExpressionResultUsed` | `com/google/javascript/jscomp/NodeUtil.java` | 4973 | 15 | 24 |  |
| 562 | 94 | `internal/jscomp/checksuper.go` | `getPropertyName` | `com/google/javascript/jscomp/CheckSuper.java` | 337 | 4 | 5 |  |
| 563 | 94 | `internal/jscomp/nodeutil.go` | `ReplaceDeclarationChild` | `com/google/javascript/jscomp/NodeUtil.java` | 2449 | 22 | 20 |  |
| 564 | 94 | `internal/jscomp/checkunreachablecode.go` | `isReachable` | `com/google/javascript/jscomp/CheckUnreachableCode.java` | 77 | 11 | 10 |  |
| 565 | 94 | `internal/jscomp/parsing/jsdocinfoparser.go` | `parseIdGeneratorTag` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 1553 | 23 | 32 |  |
| 566 | 94 | `internal/rhino/qualified_name.go` | `Components` | `com/google/javascript/rhino/QualifiedName.java` | 110 | 3 | 3 |  |
| 567 | 94 | `internal/jscomp/nodeutil.go` | `TryMergeBlock` | `com/google/javascript/jscomp/NodeUtil.java` | 2490 | 13 | 13 |  |
| 568 | 94 | `internal/jscomp/nodeutil.go` | `CanBeSideEffected` | `com/google/javascript/jscomp/NodeUtil.java` | 1170 | 21 | 17 |  |
| 569 | 94 | `internal/jscomp/nodeutil.go` | `IsObjLitProperty` | `com/google/javascript/jscomp/NodeUtil.java` | 5611 | 2 | 5 |  |
| 570 | 94 | `internal/rhino/jstype/function_type.go` | `BuildAndResolve` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1623 | 1 | 1 |  |
| 571 | 94 | `internal/jscomp/parsing/jsdocinfoparser.go` | `parseBasicTypeExpression` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2226 | 23 | 22 |  |
| 572 | 95 | `internal/rhino/jstype/property_map.go` | `GetPropertiesCount` | `com/google/javascript/rhino/jstype/PropertyMap.java` | 196 | 3 | 4 |  |
| 573 | 95 | `internal/rhino/jstype/function_type.go` | `GetBindReturnType` | `com/google/javascript/rhino/jstype/FunctionType.java` | 684 | 11 | 14 |  |
| 574 | 95 | `internal/jscomp/nodeutil.go` | `hasFinally` | `com/google/javascript/jscomp/NodeUtil.java` | 4282 | 1 | 2 |  |
| 575 | 95 | `internal/jscomp/ast_analyzer.go` | `FunctionCallHasSideEffects` | `com/google/javascript/jscomp/AstAnalyzer.java` | 145 | 38 | 75 |  |
| 576 | 95 | `internal/rhino/jstype/function_type.go` | `calculateConstructorAmbiguity` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1337 | 12 | 18 |  |
| 577 | 95 | `internal/jscomp/syntacticscopecreator.go` | `populate` | `com/google/javascript/jscomp/SyntacticScopeCreator.java` | 104 | 37 | 36 |  |
| 578 | 95 | `internal/jscomp/typemismatch.go` | `RegisterMismatch` | `com/google/javascript/jscomp/TypeMismatch.java` | 81 | 19 | 19 |  |
| 579 | 95 | `internal/jscomp/nodeutil.go` | `IsPrototypeMethod` | `com/google/javascript/jscomp/NodeUtil.java` | 3847 | 8 | 6 |  |
| 580 | 95 | `internal/rhino/jstype/subtype_checker.go` | `isSubtypeCaching` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 145 | 20 | 21 |  |
| 581 | 95 | `internal/jscomp/nodeutil.go` | `isExportFrom` | `com/google/javascript/jscomp/NodeUtil.java` | 4274 | 1 | 2 |  |
| 582 | 96 | `internal/jscomp/syntacticscopecreator.go` | `declareVar` | `com/google/javascript/jscomp/SyntacticScopeCreator.java` | 313 | 18 | 17 |  |
| 583 | 96 | `internal/jscomp/parsing/jsdoctokenstream.go` | `GetRemainingJSDocLine` | `com/google/javascript/jscomp/parsing/JsDocTokenStream.java` | 202 | 12 | 17 |  |
| 584 | 96 | `internal/rhino/jsdoc_info_builder.go` | `RecordFileOverview` | `com/google/javascript/rhino/JSDocInfo.java` | 2291 | 7 | 5 |  |
| 585 | 96 | `internal/jscomp/syntacticscopecreator.go` | `scanVars` | `com/google/javascript/jscomp/SyntacticScopeCreator.java` | 186 | 57 | 69 |  |
| 586 | 96 | `internal/jscomp/checksuper.go` | `visitContextNode` | `com/google/javascript/jscomp/CheckSuper.java` | 410 | 16 | 10 |  |
| 587 | 96 | `internal/jscomp/nodeutil.go` | `IsDefinedValue` | `com/google/javascript/jscomp/NodeUtil.java` | 4565 | 27 | 49 |  |
| 588 | 96 | `internal/jscomp/ast_analyzer.go` | `checkAssignmentForChangeState` | `com/google/javascript/jscomp/AstAnalyzer.java` | 460 | 19 | 15 |  |
| 589 | 96 | `internal/jscomp/checksuper.go` | `visit` | `com/google/javascript/jscomp/CheckSuper.java` | 150 | 23 | 19 |  |
| 590 | 96 | `internal/jscomp/nodeutil.go` | `IsShallowStatementTree` | `com/google/javascript/jscomp/NodeUtil.java` | 4033 | 2 | 4 |  |
| 591 | 96 | `internal/jscomp/checksuper.go` | `getContextForFunctionNode` | `com/google/javascript/jscomp/CheckSuper.java` | 133 | 7 | 9 |  |
| 592 | 96 | `internal/jscomp/nodeutil.go` | `IsConstantDeclaration` | `com/google/javascript/jscomp/NodeUtil.java` | 4369 | 12 | 14 |  |
| 593 | 96 | `internal/rhino/jsdoc.go` | `IsConstant` | `com/google/javascript/rhino/JSDocInfo.java` | 763 | 1 | 2 |  |
| 594 | 96 | `internal/rhino/jstype/property_map.go` | `collectAllAncestors` | `com/google/javascript/rhino/jstype/PropertyMap.java` | 255 | 15 | 9 |  |
| 595 | 97 | `internal/rhino/jstype/registry.go` | `CreateNullableType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1746 | 1 | 1 |  |
| 596 | 97 | `internal/rhino/jstype/union_type.go` | `GetPossibleToBooleanOutcomes` | `com/google/javascript/rhino/jstype/UnionType.java` | 480 | 6 | 7 |  |
| 597 | 97 | `internal/rhino/jstype/function_type.go` | `GetAllImplementedInterfaces` | `com/google/javascript/rhino/jstype/FunctionType.java` | 559 | 5 | 4 |  |
| 598 | 97 | `internal/rhino/jstype/equality_checker.go` | `areEqualCaching` | `com/google/javascript/rhino/jstype/EqualityChecker.java` | 135 | 18 | 19 |  |
| 599 | 97 | `internal/jscomp/parsing/jsdocinfoparser.go` | `next` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2779 | 4 | 4 |  |
| 600 | 97 | `internal/rhino/jstype/contains_upper_bound_super_type_visitor.go` | `CaseTemplatizedType` | `com/google/javascript/rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java` | 87 | 1 | 1 |  |
| 601 | 97 | `internal/rhino/jstype/contains_upper_bound_super_type_visitor.go` | `CaseTemplateType` | `com/google/javascript/rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java` | 77 | 1 | 1 |  |
| 602 | 97 | `internal/rhino/jsdoc_info_builder.go` | `RecordVisibility` | `com/google/javascript/rhino/JSDocInfo.java` | 1881 | 4 | 5 |  |
| 603 | 97 | `internal/jscomp/checkjsdoc.go` | `validateMsgJsDoc` | `com/google/javascript/jscomp/CheckJSDoc.java` | 504 | 21 | 15 |  |
| 604 | 97 | `internal/jscomp/checkjsdoc.go` | `validateIsUsedViaDotConstructor` | `com/google/javascript/jscomp/CheckJSDoc.java` | 790 | 5 | 5 |  |
| 605 | 97 | `internal/rhino/qualified_name.go` | `GetComponentCount` | `com/google/javascript/rhino/QualifiedName.java` | 301 | 6 | 6 |  |
| 606 | 97 | `internal/rhino/jstype/function_type.go` | `HasImplementedInterfaces` | `com/google/javascript/rhino/jstype/FunctionType.java` | 349 | 6 | 6 |  |
| 607 | 97 | `internal/rhino/jstype/registry.go` | `createNamedType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 2074 | 3 | 5 |  |
| 608 | 97 | `internal/jscomp/checkmissingrequires.go` | `normalizeQualifiedNamePlusImport` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 743 | 4 | 4 |  |
| 609 | 97 | `internal/jscomp/checkmissingrequires.go` | `visit` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 182 | 4 | 4 |  |
| 610 | 97 | `internal/jscomp/checkmissingrequires.go` | `visitJsDocInfo` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 216 | 17 | 21 |  |
| 611 | 97 | `internal/rhino/jstype/property.go` | `HumanReadableName` | `com/google/javascript/rhino/jstype/Property.java` | 120 | 1 | 1 |  |
| 612 | 98 | `internal/rhino/jstype/named_type.go` | `setReferencedAndResolvedType` | `com/google/javascript/rhino/jstype/NamedType.java` | 359 | 5 | 7 |  |
| 613 | 98 | `internal/rhino/jstype/subtype_checker.go` | `checkPropertyPresence` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 405 | 7 | 8 |  |
| 614 | 98 | `internal/rhino/jstype_expr.go` | `visitAllTypeNodes` | `com/google/javascript/rhino/JSTypeExpression.java` | 145 | 6 | 6 |  |
| 615 | 98 | `internal/rhino/jstype/union_type.go` | `GetLeastSupertype` | `com/google/javascript/rhino/jstype/UnionType.java` | 320 | 7 | 6 |  |
| 616 | 98 | `internal/rhino/jstype/function_type.go` | `MakesDicts` | `com/google/javascript/rhino/jstype/FunctionType.java` | 300 | 16 | 16 |  |
| 617 | 99 | `internal/jscomp/nodeutil.go` | `CanMergeBlock` | `com/google/javascript/jscomp/NodeUtil.java` | 2516 | 8 | 12 |  |
| 618 | 99 | `internal/rhino/qualified_name.go` | `Join` | `com/google/javascript/rhino/QualifiedName.java` | 125 | 3 | 3 |  |
| 619 | 99 | `internal/rhino/qualified_name.go` | `appendTo` | `com/google/javascript/rhino/QualifiedName.java` | 167 | 4 | 4 |  |
| 620 | 99 | `internal/jscomp/checkmissingrequires.go` | `checkNormalizedShortNameImport` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 583 | 34 | 39 |  |
| 621 | 99 | `internal/rhino/jsdoc_info_builder.go` | `RecordAuthor` | `com/google/javascript/rhino/JSDocInfo.java` | 1978 | 4 | 4 |  |
| 622 | 99 | `internal/rhino/jsdoc_info_builder.go` | `RecordReference` | `com/google/javascript/rhino/JSDocInfo.java` | 1988 | 4 | 4 |  |
| 623 | 99 | `internal/jscomp/checkmissingrequires.go` | `shouldTraverse` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 135 | 5 | 6 |  |
| 624 | 99 | `internal/jscomp/parsing/jsdocinfoparser.go` | `parseClosurePrimitiveTag` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 1481 | 12 | 18 |  |
| 625 | 99 | `internal/jscomp/parsing/jsdocinfoparser.go` | `parseModifiesTag` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 1511 | 20 | 23 |  |
| 626 | 99 | `internal/jscomp/parsing/jsdocinfoparser.go` | `parseTypeExpressionAnnotation` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 1998 | 11 | 16 |  |
| 627 | 99 | `internal/rhino/token_util.go` | `IsJSSpace` | `com/google/javascript/rhino/TokenUtil.java` | 59 | 3 | 5 |  |
| 628 | 99 | `internal/jscomp/nodeutil.go` | `isReferenceName` | `com/google/javascript/jscomp/NodeUtil.java` | 2261 | 1 | 1 |  |
| 629 | 99 | `internal/jscomp/checkglobalthis.go` | `ShouldTraverse` | `com/google/javascript/jscomp/CheckGlobalThis.java` | 82 | 43 | 49 |  |
| 630 | 99 | `internal/jscomp/nodeutil.go` | `IsGoogModuleCall` | `com/google/javascript/jscomp/NodeUtil.java` | 5143 | 4 | 4 |  |
| 631 | 99 | `internal/rhino/jstype/function_param_builder.go` | `AddOptionalParams` | `com/google/javascript/rhino/jstype/FunctionParamBuilder.java` | 84 | 5 | 6 |  |
| 632 | 99 | `internal/jscomp/nodeutil.go` | `IsGoogProvideCall` | `com/google/javascript/jscomp/NodeUtil.java` | 5135 | 4 | 4 |  |
| 633 | 99 | `internal/rhino/jstype/contains_upper_bound_super_type_visitor.go` | `CaseUnionType` | `com/google/javascript/rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java` | 92 | 7 | 7 |  |
| 634 | 99 | `internal/jscomp/nodeutil.go` | `IsFunctionExpression` | `com/google/javascript/jscomp/NodeUtil.java` | 2754 | 1 | 1 |  |
| 635 | 100 | `internal/rhino/jstype/function_type.go` | `ToMaybeFunctionType` | `com/google/javascript/rhino/jstype/FunctionType.java` | 340 | 1 | 1 |  |
| 636 | 100 | `internal/rhino/jstype/function_type.go` | `IsAmbiguousConstructor` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1330 | 3 | 3 |  |
| 637 | 100 | `internal/jscomp/nodeutil.go` | `IsNameDeclaration` | `com/google/javascript/jscomp/NodeUtil.java` | 2008 | 1 | 1 |  |
| 638 | 100 | `internal/rhino/jstype/function_type.go` | `CheckExtendsLoop` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1231 | 1 | 1 |  |
| 639 | 100 | `internal/rhino/jstype/function_type.go` | `GetClosurePrimitive` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1389 | 1 | 1 |  |
| 640 | 100 | `internal/rhino/jstype/function_type.go` | `GetCanonicalRepresentation` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1628 | 1 | 1 |  |
| 641 | 100 | `internal/jscomp/nodeutil.go` | `IsBleedingFunctionName` | `com/google/javascript/jscomp/NodeUtil.java` | 2812 | 4 | 4 |  |
| 642 | 100 | `internal/rhino/jstype/function_type.go` | `GetPrototypeProperty` | `com/google/javascript/rhino/jstype/FunctionType.java` | 437 | 1 | 1 |  |
| 643 | 100 | `internal/rhino/jstype/function_type.go` | `CanBeCalled` | `com/google/javascript/rhino/jstype/FunctionType.java` | 345 | 1 | 1 |  |
| 644 | 100 | `internal/rhino/jstype/function_type.go` | `IsInterface` | `com/google/javascript/rhino/jstype/FunctionType.java` | 252 | 1 | 1 |  |
| 645 | 100 | `internal/rhino/jstype/function_type.go` | `SetExplicitUnrestricted` | `com/google/javascript/rhino/jstype/FunctionType.java` | 335 | 1 | 1 |  |
| 646 | 100 | `internal/rhino/jstype/function_type.go` | `IsOrdinaryFunction` | `com/google/javascript/rhino/jstype/FunctionType.java` | 257 | 1 | 1 |  |
| 647 | 100 | `internal/rhino/jstype/function_type.go` | `GetKind` | `com/google/javascript/rhino/jstype/FunctionType.java` | 261 | 1 | 1 |  |
| 648 | 100 | `internal/jscomp/nodeutil.go` | `IsInFunction` | `com/google/javascript/jscomp/NodeUtil.java` | 1746 | 1 | 1 |  |
| 649 | 100 | `internal/rhino/jstype/function_type.go` | `IsAbstract` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1177 | 1 | 1 |  |
| 650 | 100 | `internal/rhino/jstype/function_type.go` | `GetSource` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1036 | 1 | 1 |  |
| 651 | 100 | `internal/rhino/jstype/function_type.go` | `IsReturnTypeInferred` | `com/google/javascript/rhino/jstype/FunctionType.java` | 401 | 1 | 1 |  |
| 652 | 100 | `internal/rhino/jstype/function_type.go` | `GetExtendedInterfaces` | `com/google/javascript/rhino/jstype/FunctionType.java` | 627 | 1 | 1 |  |
| 653 | 100 | `internal/rhino/jstype/function_param_builder.go` | `AddVarArgs` | `com/google/javascript/rhino/jstype/FunctionParamBuilder.java` | 100 | 4 | 4 |  |
| 654 | 100 | `internal/rhino/jstype/function_param_builder.go` | `AddRequiredParams` | `com/google/javascript/rhino/jstype/FunctionParamBuilder.java` | 67 | 5 | 5 |  |
| 655 | 100 | `internal/rhino/jstype/equality_checker.go` | `Check` | `com/google/javascript/rhino/jstype/EqualityChecker.java` | 122 | 3 | 3 |  |
| 656 | 100 | `internal/rhino/jstype/enum_element_type.go` | `ToMaybeEnumElementType` | `com/google/javascript/rhino/jstype/EnumElementType.java` | 102 | 1 | 1 |  |
| 657 | 100 | `internal/rhino/jstype/function_type.go` | `GetOwnImplementedInterfaces` | `com/google/javascript/rhino/jstype/FunctionType.java` | 613 | 1 | 1 |  |
| 658 | 100 | `internal/rhino/jstype/enum_element_type.go` | `GetReferenceName` | `com/google/javascript/rhino/jstype/EnumElementType.java` | 181 | 1 | 1 |  |
| 659 | 100 | `internal/rhino/jstype/enum_element_type.go` | `GetEnumType` | `com/google/javascript/rhino/jstype/EnumElementType.java` | 85 | 1 | 1 |  |
| 660 | 100 | `internal/rhino/jstype/enum_element_type.go` | `GetPrimitiveType` | `com/google/javascript/rhino/jstype/EnumElementType.java` | 226 | 1 | 1 |  |
| 661 | 100 | `internal/jscomp/nodeutil.go` | `IsFunctionObjectCall` | `com/google/javascript/jscomp/NodeUtil.java` | 2894 | 1 | 1 |  |
| 662 | 100 | `internal/jscomp/nodeutil.go` | `IsFunctionObjectApply` | `com/google/javascript/jscomp/NodeUtil.java` | 2902 | 1 | 1 |  |
| 663 | 100 | `internal/rhino/jstype/function_type.go` | `GetPropertyTypeMap` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1187 | 4 | 3 |  |
| 664 | 100 | `internal/rhino/jstype/arrow_type.go` | `RecursionUnsafeHashCode` | `com/google/javascript/rhino/jstype/ArrowType.java` | 90 | 4 | 4 |  |
| 665 | 100 | `internal/jscomp/nodeutil.go` | `GetEnclosingScript` | `com/google/javascript/jscomp/NodeUtil.java` | 1716 | 1 | 1 |  |
| 666 | 100 | `internal/jscomp/nodeutil.go` | `GetEnclosingFunction` | `com/google/javascript/jscomp/NodeUtil.java` | 1711 | 1 | 1 |  |
| 667 | 100 | `internal/jscomp/ast_analyzer.go` | `MayEffectMutableState` | `com/google/javascript/jscomp/AstAnalyzer.java` | 125 | 1 | 1 |  |
| 668 | 100 | `internal/rhino/jstype/function_type.go` | `MakesStructs` | `com/google/javascript/rhino/jstype/FunctionType.java` | 269 | 16 | 16 |  |
| 669 | 100 | `internal/rhino/jstype/arrow_type.go` | `hasTemplatedParameterType` | `com/google/javascript/rhino/jstype/ArrowType.java` | 166 | 4 | 5 |  |
| 670 | 100 | `internal/rhino/jstype/arrow_type.go` | `HasUnknownParamsOrReturn` | `com/google/javascript/rhino/jstype/ArrowType.java` | 145 | 4 | 5 |  |
| 671 | 100 | `internal/rhino/jstype/arrow_type.go` | `GetParameterList` | `com/google/javascript/rhino/jstype/ArrowType.java` | 102 | 1 | 1 |  |
| 672 | 100 | `internal/jscomp/checksuspiciouscode.go` | `getBooleanValueWithTypes` | `com/google/javascript/jscomp/CheckSuspiciousCode.java` | 217 | 27 | 27 |  |
| 673 | 100 | `internal/jscomp/checksuspiciouscode.go` | `checkNonObjectInstanceOf` | `com/google/javascript/jscomp/CheckSuspiciousCode.java` | 138 | 3 | 2 |  |
| 674 | 100 | `internal/rhino/jstype/function_type.go` | `WithSourceNode` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1455 | 2 | 2 |  |
| 675 | 100 | `internal/rhino/jstype/function_type.go` | `SetGoogModuleId` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1461 | 2 | 2 |  |
| 676 | 100 | `internal/rhino/jstype/function_type.go` | `WithReturnType` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1485 | 2 | 3 |  |
| 677 | 100 | `internal/rhino/jstype/function_type.go` | `WithInferredReturnType` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1498 | 3 | 3 |  |
| 678 | 100 | `internal/rhino/jstype/function_type.go` | `WithTypeOfThis` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1505 | 2 | 2 |  |
| 679 | 100 | `internal/jscomp/checksuspiciouscode.go` | `checkInvalidIn` | `com/google/javascript/jscomp/CheckSuspiciousCode.java` | 132 | 3 | 2 |  |
| 680 | 100 | `internal/rhino/jstype/arrow_type.go` | `GetReturnType` | `com/google/javascript/rhino/jstype/ArrowType.java` | 98 | 1 | 1 |  |
| 681 | 100 | `internal/rhino/jstype/function_type.go` | `WithClosurePrimitiveId` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1580 | 2 | 2 |  |
| 682 | 100 | `internal/rhino/jstype/function_type.go` | `WithKind` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1539 | 2 | 2 |  |
| 683 | 100 | `internal/rhino/jstype/function_type.go` | `ForConstructor` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1545 | 2 | 2 |  |
| 684 | 100 | `internal/rhino/jsdoc_info_builder.go` | `isConstructorOrInterface` | `com/google/javascript/rhino/JSDocInfo.java` | 833 | 1 | 1 |  |
| 685 | 100 | `internal/rhino/jstype/function_type.go` | `WithIsAbstract` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1563 | 2 | 2 |  |
| 686 | 100 | `internal/rhino/jstype/function_type.go` | `WithCanonicalRepresentation` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1586 | 2 | 2 |  |
| 687 | 100 | `internal/rhino/jstype/function_type.go` | `WithPrototypeBasedOn` | `com/google/javascript/rhino/jstype/FunctionType.java` | 1574 | 2 | 2 |  |
| 688 | 100 | `internal/rhino/jsdoc_info_builder.go` | `IsSassGeneratedCssTsRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2755 | 1 | 1 |  |
| 689 | 100 | `internal/jscomp/checksuspiciouscode.go` | `checkNaN` | `com/google/javascript/jscomp/CheckSuspiciousCode.java` | 116 | 4 | 5 |  |
| 690 | 100 | `internal/rhino/jsdoc_info_builder.go` | `IsWizcallbackRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2695 | 1 | 1 |  |
| 691 | 100 | `internal/jscomp/checksuspiciouscode.go` | `visit` | `com/google/javascript/jscomp/CheckSuspiciousCode.java` | 79 | 6 | 7 |  |
| 692 | 100 | `internal/rhino/jsdoc_info_builder.go` | `IsMixinFunctionRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2745 | 1 | 1 |  |
| 693 | 100 | `internal/rhino/jsdoc_info_builder.go` | `IsMixinClassRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2735 | 1 | 1 |  |
| 694 | 100 | `internal/jscomp/checksuper.go` | `visitSuperPropertyAccess` | `com/google/javascript/jscomp/CheckSuper.java` | 460 | 1 | 1 |  |
| 695 | 100 | `internal/rhino/jsdoc_info_builder.go` | `IsCustomElementRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2725 | 1 | 1 |  |
| 696 | 100 | `internal/jscomp/checksuper.go` | `getContextNode` | `com/google/javascript/jscomp/CheckSuper.java` | 205 | 1 | 1 |  |
| 697 | 100 | `internal/rhino/jsdoc_info_builder.go` | `IsPolymerRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2715 | 1 | 1 |  |
| 698 | 100 | `internal/jscomp/checksuper.go` | `visitReturn` | `com/google/javascript/jscomp/CheckSuper.java` | 403 | 2 | 2 |  |
| 699 | 100 | `internal/jscomp/checksuper.go` | `visitThis` | `com/google/javascript/jscomp/CheckSuper.java` | 396 | 2 | 2 |  |
| 700 | 100 | `internal/rhino/jsdoc_info_builder.go` | `IsPolymerBehaviorRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2705 | 1 | 1 |  |
| 701 | 100 | `internal/rhino/jsdoc_info_builder.go` | `IsEncourageInlining` | `com/google/javascript/rhino/JSDocInfo.java` | 1738 | 1 | 1 |  |
| 702 | 100 | `internal/jscomp/checksuper.go` | `getContextNode` | `com/google/javascript/jscomp/CheckSuper.java` | 205 | 1 | 1 |  |
| 703 | 100 | `internal/rhino/jsdoc_info_builder.go` | `IsRequireInlining` | `com/google/javascript/rhino/JSDocInfo.java` | 1730 | 1 | 1 |  |
| 704 | 100 | `internal/rhino/jsdoc_info_builder.go` | `IsNoInline` | `com/google/javascript/rhino/JSDocInfo.java` | 1723 | 1 | 1 |  |
| 705 | 100 | `internal/rhino/jsdoc_info_builder.go` | `IsUnrestrictedRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2495 | 1 | 1 |  |
| 706 | 100 | `internal/jscomp/checksuper.go` | `visitSuperConstructorCall` | `com/google/javascript/jscomp/CheckSuper.java` | 318 | 11 | 8 |  |
| 707 | 100 | `internal/jscomp/nodeutil.go` | `IsGoogRequireCall` | `com/google/javascript/jscomp/NodeUtil.java` | 5160 | 3 | 4 |  |
| 708 | 100 | `internal/jscomp/nodeutil.go` | `MayBeUndefined` | `com/google/javascript/jscomp/NodeUtil.java` | 4556 | 1 | 1 |  |
| 709 | 100 | `internal/jscomp/checksuper.go` | `getContextNode` | `com/google/javascript/jscomp/CheckSuper.java` | 205 | 1 | 1 |  |
| 710 | 100 | `internal/jscomp/nodeutil.go` | `isNamedFunction` | `com/google/javascript/jscomp/NodeUtil.java` | 2794 | 1 | 1 |  |
| 711 | 100 | `internal/jscomp/nodeutil.go` | `isNamedClass` | `com/google/javascript/jscomp/NodeUtil.java` | 2804 | 1 | 1 |  |
| 712 | 100 | `internal/jscomp/nodeutil.go` | `isSwitchCase` | `com/google/javascript/jscomp/NodeUtil.java` | 2253 | 1 | 1 |  |
| 713 | 100 | `internal/rhino/jstype/named_type.go` | `SetResolutionKind` | `com/google/javascript/rhino/jstype/NamedType.java` | 527 | 2 | 2 |  |
| 714 | 100 | `internal/rhino/jstype/named_type.go` | `SetErrorReportingLocation` | `com/google/javascript/rhino/jstype/NamedType.java` | 532 | 4 | 4 |  |
| 715 | 100 | `internal/rhino/jstype/named_type.go` | `SetReferencedType` | `com/google/javascript/rhino/jstype/NamedType.java` | 551 | 2 | 2 |  |
| 716 | 100 | `internal/rhino/jstype/named_type.go` | `SetRestrictByNull` | `com/google/javascript/rhino/jstype/NamedType.java` | 556 | 2 | 2 |  |
| 717 | 100 | `internal/rhino/jstype/named_type.go` | `SetTemplateTypes` | `com/google/javascript/rhino/jstype/NamedType.java` | 546 | 2 | 2 |  |
| 718 | 100 | `internal/jscomp/checksuper.go` | `getContextNode` | `com/google/javascript/jscomp/CheckSuper.java` | 205 | 1 | 1 |  |
| 719 | 100 | `internal/jscomp/nodeutil.go` | `isTryFinallyNode` | `com/google/javascript/jscomp/NodeUtil.java` | 2300 | 1 | 1 |  |
| 720 | 100 | `internal/rhino/jstype/named_type.go` | `GetReferenceName` | `com/google/javascript/rhino/jstype/NamedType.java` | 213 | 1 | 1 |  |
| 721 | 100 | `internal/jscomp/nodeutil.go` | `isTryCatchNodeContainer` | `com/google/javascript/jscomp/NodeUtil.java` | 2305 | 2 | 2 |  |
| 722 | 100 | `internal/jscomp/nodeutil.go` | `IsClassDeclaration` | `com/google/javascript/jscomp/NodeUtil.java` | 2658 | 1 | 1 |  |
| 723 | 100 | `internal/jscomp/nodeutil.go` | `IsClassExpression` | `com/google/javascript/jscomp/NodeUtil.java` | 2771 | 1 | 1 |  |
| 724 | 100 | `internal/rhino/jstype/named_type.go` | `GetTemplateTypes` | `com/google/javascript/rhino/jstype/NamedType.java` | 172 | 1 | 1 |  |
| 725 | 100 | `internal/jscomp/nodeutil.go` | `IsBlockScopedFunctionDeclaration` | `com/google/javascript/jscomp/NodeUtil.java` | 2677 | 13 | 13 |  |
| 726 | 100 | `internal/rhino/jstype/named_type.go` | `resolveTypeof` | `com/google/javascript/rhino/jstype/NamedType.java` | 321 | 23 | 22 |  |
| 727 | 100 | `internal/jscomp/nodeutil.go` | `RemoveChild` | `com/google/javascript/jscomp/NodeUtil.java` | 2362 | 42 | 50 |  |
| 728 | 100 | `internal/jscomp/checksuper.go` | `getContextNode` | `com/google/javascript/jscomp/CheckSuper.java` | 205 | 1 | 1 |  |
| 729 | 100 | `internal/rhino/jsdoc_info_builder.go` | `ParseDocumentation` | `com/google/javascript/rhino/JSDocInfo.java` | 1672 | 2 | 2 |  |
| 730 | 100 | `internal/rhino/jstype/object_type_ext.go` | `getCtorExtendedInterfaces` | `com/google/javascript/rhino/jstype/ObjectType.java` | 911 | 1 | 1 |  |
| 731 | 100 | `internal/rhino/jstype/property.go` | `Kind` | `com/google/javascript/rhino/jstype/Property.java` | 93 | 1 | 1 |  |
| 732 | 100 | `internal/rhino/jsdoc_info_builder.go` | `ShouldParseDocumentation` | `com/google/javascript/rhino/JSDocInfo.java` | 1677 | 1 | 1 |  |
| 733 | 100 | `internal/jscomp/nodeutil.go` | `IsNullOrUndefined` | `com/google/javascript/jscomp/NodeUtil.java` | 1330 | 1 | 1 |  |
| 734 | 100 | `internal/rhino/jstype/property.go` | `Kind` | `com/google/javascript/rhino/jstype/Property.java` | 93 | 1 | 1 |  |
| 735 | 100 | `internal/jscomp/nodeutil.go` | `IsEmptyBlock` | `com/google/javascript/jscomp/NodeUtil.java` | 883 | 6 | 6 |  |
| 736 | 100 | `internal/rhino/jsdoc_info_builder.go` | `RecordOriginalCommentPosition` | `com/google/javascript/rhino/JSDocInfo.java` | 1693 | 3 | 3 |  |
| 737 | 100 | `internal/rhino/jstype/property.go` | `GetOwner` | `com/google/javascript/rhino/jstype/Property.java` | 135 | 1 | 1 |  |
| 738 | 100 | `internal/rhino/jstype/property.go` | `GetValue` | `com/google/javascript/rhino/jstype/Property.java` | 139 | 1 | 1 |  |
| 739 | 100 | `internal/rhino/jsdoc_info_builder.go` | `RecordOriginalCommentString` | `com/google/javascript/rhino/JSDocInfo.java` | 1685 | 3 | 3 |  |
| 740 | 100 | `internal/jscomp/nodeutil.go` | `IsCompoundAssignmentOp` | `com/google/javascript/jscomp/NodeUtil.java` | 1667 | 1 | 1 |  |
| 741 | 100 | `internal/jscomp/nodeutil.go` | `IsHoistScopeRoot` | `com/google/javascript/jscomp/NodeUtil.java` | 1741 | 1 | 1 |  |
| 742 | 100 | `internal/rhino/jstype/property.go` | `GetNode` | `com/google/javascript/rhino/jstype/Property.java` | 213 | 1 | 1 |  |
| 743 | 100 | `internal/rhino/jstype/property.go` | `GetSymbol` | `com/google/javascript/rhino/jstype/Property.java` | 223 | 1 | 1 |  |
| 744 | 100 | `internal/rhino/jstype/property.go` | `GetType` | `com/google/javascript/rhino/jstype/Property.java` | 233 | 1 | 1 |  |
| 745 | 100 | `internal/rhino/jstype/property.go` | `IsTypeInferred` | `com/google/javascript/rhino/jstype/Property.java` | 238 | 1 | 1 |  |
| 746 | 100 | `internal/jscomp/nodeutil.go` | `IsEnhancedFor` | `com/google/javascript/jscomp/NodeUtil.java` | 2068 | 1 | 1 |  |
| 747 | 100 | `internal/jscomp/nodeutil.go` | `GetEnclosingBlock` | `com/google/javascript/jscomp/NodeUtil.java` | 1721 | 1 | 1 |  |
| 748 | 100 | `internal/rhino/jstype/property.go` | `SetJSDocInfo` | `com/google/javascript/rhino/jstype/Property.java` | 254 | 1 | 1 |  |
| 749 | 100 | `internal/rhino/jstype/property.go` | `SetNode` | `com/google/javascript/rhino/jstype/Property.java` | 258 | 1 | 1 |  |
| 750 | 100 | `internal/jscomp/nodeutil.go` | `GetEnclosingClass` | `com/google/javascript/jscomp/NodeUtil.java` | 1702 | 1 | 1 |  |
| 751 | 100 | `internal/jscomp/nodeutil.go` | `GetEnclosingModuleIfPresent` | `com/google/javascript/jscomp/NodeUtil.java` | 1706 | 1 | 1 |  |
| 752 | 100 | `internal/jscomp/nodeutil.go` | `IsNormalGet` | `com/google/javascript/jscomp/NodeUtil.java` | 1854 | 1 | 1 |  |
| 753 | 100 | `internal/rhino/jstype/property_map.go` | `ImmutableEmptyMap` | `com/google/javascript/rhino/jstype/PropertyMap.java` | 114 | 1 | 1 |  |
| 754 | 100 | `internal/rhino/jstype/property_map.go` | `SetParentSource` | `com/google/javascript/rhino/jstype/PropertyMap.java` | 118 | 4 | 4 |  |
| 755 | 100 | `internal/rhino/jstype/property_map.go` | `SetParentForTesting` | `com/google/javascript/rhino/jstype/PropertyMap.java` | 128 | 2 | 2 |  |
| 756 | 100 | `internal/jscomp/nodeutil.go` | `IsOptChainGet` | `com/google/javascript/jscomp/NodeUtil.java` | 1859 | 1 | 1 |  |
| 757 | 100 | `internal/jscomp/nodeutil.go` | `IsOptChainNode` | `com/google/javascript/jscomp/NodeUtil.java` | 1864 | 1 | 1 |  |
| 758 | 100 | `internal/jscomp/nodeutil.go` | `IsNormalOrOptChainGet` | `com/google/javascript/jscomp/NodeUtil.java` | 1839 | 1 | 1 |  |
| 759 | 100 | `internal/jscomp/nodeutil.go` | `IsNormalOrOptChainGetProp` | `com/google/javascript/jscomp/NodeUtil.java` | 1844 | 1 | 1 |  |
| 760 | 100 | `internal/jscomp/checkregexp.go` | `IsGlobalRegExpPropertiesUsed` | `com/google/javascript/jscomp/CheckRegExp.java` | 71 | 3 | 3 |  |
| 761 | 100 | `internal/jscomp/nodeutil.go` | `IsNormalOrOptChainCall` | `com/google/javascript/jscomp/NodeUtil.java` | 1849 | 1 | 1 |  |
| 762 | 100 | `internal/jscomp/nodeutil.go` | `IsNumericResult` | `com/google/javascript/jscomp/NodeUtil.java` | 1512 | 1 | 1 |  |
| 763 | 100 | `internal/rhino/jstype/property_map.go` | `PutProperty` | `com/google/javascript/rhino/jstype/PropertyMap.java` | 273 | 5 | 6 |  |
| 764 | 100 | `internal/jscomp/nodeutil.go` | `IsBigIntResult` | `com/google/javascript/jscomp/NodeUtil.java` | 1517 | 1 | 1 |  |
| 765 | 100 | `internal/jscomp/nodeutil.go` | `IsBooleanResult` | `com/google/javascript/jscomp/NodeUtil.java` | 1524 | 1 | 1 |  |
| 766 | 100 | `internal/jscomp/nodeutil.go` | `IsStringResult` | `com/google/javascript/jscomp/NodeUtil.java` | 1531 | 1 | 1 |  |
| 767 | 100 | `internal/jscomp/nodeutil.go` | `IsObjectResult` | `com/google/javascript/jscomp/NodeUtil.java` | 1538 | 1 | 1 |  |
| 768 | 100 | `internal/rhino/jstype/registry.go` | `GetResolver` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1720 | 1 | 1 |  |
| 769 | 100 | `internal/jscomp/nodeutil.go` | `IsLhsOfAssign` | `com/google/javascript/jscomp/NodeUtil.java` | 2975 | 2 | 2 |  |
| 770 | 100 | `internal/jscomp/nodeutil.go` | `IsNameDeclOrSimpleAssignLhs` | `com/google/javascript/jscomp/NodeUtil.java` | 2915 | 1 | 1 |  |
| 771 | 100 | `internal/jscomp/nodeutil.go` | `IsTopLevel` | `com/google/javascript/jscomp/NodeUtil.java` | 5267 | 1 | 1 |  |
| 772 | 100 | `internal/rhino/jsdoc_info_builder.go` | `RecordEnumParameterType` | `com/google/javascript/rhino/JSDocInfo.java` | 2162 | 5 | 5 |  |
| 773 | 100 | `internal/jscomp/nodeutil.go` | `IsBlocklessArrowFunctionResult` | `com/google/javascript/jscomp/NodeUtil.java` | 5622 | 2 | 2 |  |
| 774 | 100 | `internal/jscomp/checkmissingrequires.go` | `isTypeScriptSource` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 383 | 1 | 1 |  |
| 775 | 100 | `internal/jscomp/checkmissingrequires.go` | `isInGlobalControlFlowScope` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 153 | 13 | 12 |  |
| 776 | 100 | `internal/jscomp/nodeutil.go` | `IsNamedExportsLiteral` | `com/google/javascript/jscomp/NodeUtil.java` | 2584 | 9 | 8 |  |
| 777 | 100 | `internal/jscomp/nodeutil.go` | `NewDeclaration` | `com/google/javascript/jscomp/NodeUtil.java` | 3436 | 3 | 3 |  |
| 778 | 100 | `internal/rhino/jstype/registry.go` | `CreateOptionalType` | `com/google/javascript/rhino/jstype/JSTypeRegistry.java` | 1733 | 3 | 4 |  |
| 779 | 100 | `internal/jscomp/nodeutil.go` | `NewUndefinedNode` | `com/google/javascript/jscomp/NodeUtil.java` | 3940 | 4 | 4 |  |
| 780 | 100 | `internal/rhino/jsdoc.go` | `HasMods` | `com/google/javascript/rhino/JSDocInfo.java` | 1445 | 1 | 1 |  |
| 781 | 100 | `internal/jscomp/nodeutil.go` | `EmptyFunction` | `com/google/javascript/jscomp/NodeUtil.java` | 3976 | 1 | 1 |  |
| 782 | 100 | `internal/rhino/jsdoc.go` | `HasEnhance` | `com/google/javascript/rhino/JSDocInfo.java` | 1435 | 1 | 1 |  |
| 783 | 100 | `internal/rhino/jsdoc.go` | `GetLogTypeInCompiler` | `com/google/javascript/rhino/JSDocInfo.java` | 1583 | 1 | 1 |  |
| 784 | 100 | `internal/rhino/jsdoc.go` | `IsClosureUnawareCode` | `com/google/javascript/rhino/JSDocInfo.java` | 2775 | 1 | 1 |  |
| 785 | 100 | `internal/jscomp/checkmissingrequires.go` | `checkMissingRequireThroughFullyQualifiedName` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 687 | 25 | 22 |  |
| 786 | 100 | `internal/rhino/jstype/subtype_checker.go` | `isOptional` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 746 | 1 | 1 |  |
| 787 | 100 | `internal/rhino/jsdoc.go` | `IsSassGeneratedCssTs` | `com/google/javascript/rhino/JSDocInfo.java` | 1264 | 1 | 1 |  |
| 788 | 100 | `internal/rhino/jsdoc.go` | `IsWizcallback` | `com/google/javascript/rhino/JSDocInfo.java` | 1235 | 1 | 1 |  |
| 789 | 100 | `internal/rhino/jsdoc.go` | `RecordInlineType` | `com/google/javascript/rhino/JSDocInfo.java` | 2105 | 1 | 1 |  |
| 790 | 100 | `internal/rhino/jsdoc.go` | `IsDocumentationIncluded` | `com/google/javascript/rhino/JSDocInfo.java` | 728 | 1 | 1 |  |
| 791 | 100 | `internal/jscomp/nodeutil.go` | `VisitPreOrder` | `com/google/javascript/jscomp/NodeUtil.java` | 4149 | 4 | 4 |  |
| 792 | 100 | `internal/rhino/jstype/subtype_checker.go` | `Check` | `com/google/javascript/rhino/jstype/SubtypeChecker.java` | 133 | 3 | 3 |  |
| 793 | 100 | `internal/jscomp/nodeutil.go` | `VisitPostOrder` | `com/google/javascript/jscomp/NodeUtil.java` | 4169 | 6 | 6 |  |
| 794 | 100 | `internal/jscomp/checkmissingrequires.go` | `visitNode` | `com/google/javascript/jscomp/CheckMissingRequires.java` | 193 | 12 | 13 |  |
| 795 | 100 | `internal/jscomp/nodeutil.go` | `GetCount` | `com/google/javascript/jscomp/NodeUtil.java` | 4115 | 7 | 7 |  |
| 796 | 100 | `internal/jscomp/nodeutil.go` | `EvaluatesToLocalValue` | `com/google/javascript/jscomp/NodeUtil.java` | 4474 | 42 | 40 |  |
| 797 | 100 | `internal/jscomp/nodeutil.go` | `AllArgsUnescapedLocal` | `com/google/javascript/jscomp/NodeUtil.java` | 1143 | 4 | 4 |  |
| 798 | 100 | `internal/rhino/jsdoc.go` | `SetDescription` | `com/google/javascript/rhino/JSDocInfo.java` | 625 | 1 | 1 |  |
| 799 | 100 | `internal/rhino/jsdoc.go` | `IsProvideAlreadyProvided` | `com/google/javascript/rhino/JSDocInfo.java` | 925 | 1 | 1 |  |
| 800 | 100 | `internal/rhino/jsdoc.go` | `IsProvideGoog` | `com/google/javascript/rhino/JSDocInfo.java` | 917 | 1 | 1 |  |
| 801 | 100 | `internal/jscomp/parsing/jsdocinfoparser.go` | `ParseJsdoc` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 253 | 3 | 3 |  |
| 802 | 100 | `internal/jscomp/checkjsdoc.go` | `isClassDecl` | `com/google/javascript/jscomp/CheckJSDoc.java` | 286 | 5 | 4 |  |
| 803 | 100 | `internal/rhino/jstype/template_type.go` | `GetTypeClass` | `com/google/javascript/rhino/jstype/TemplateType.java` | 89 | 1 | 1 |  |
| 804 | 100 | `internal/rhino/jstype/template_type.go` | `GetReferenceName` | `com/google/javascript/rhino/jstype/TemplateType.java` | 94 | 1 | 1 |  |
| 805 | 100 | `internal/rhino/jstype/template_type.go` | `AppendTo` | `com/google/javascript/rhino/jstype/TemplateType.java` | 99 | 5 | 4 |  |
| 806 | 100 | `internal/rhino/jstype/template_type.go` | `ToMaybeTemplateType` | `com/google/javascript/rhino/jstype/TemplateType.java` | 109 | 1 | 1 |  |
| 807 | 100 | `internal/rhino/jstype/template_type.go` | `HasAnyTemplateTypesInternal` | `com/google/javascript/rhino/jstype/TemplateType.java` | 114 | 1 | 1 |  |
| 808 | 100 | `internal/rhino/jstype/template_type.go` | `Visit` | `com/google/javascript/rhino/jstype/TemplateType.java` | 119 | 1 | 1 |  |
| 809 | 100 | `internal/rhino/jstype/template_type.go` | `IsTypeTransformation` | `com/google/javascript/rhino/jstype/TemplateType.java` | 127 | 1 | 1 |  |
| 810 | 100 | `internal/rhino/jstype/template_type.go` | `GetTypeTransformation` | `com/google/javascript/rhino/jstype/TemplateType.java` | 131 | 1 | 1 |  |
| 811 | 100 | `internal/jscomp/parsing/jsdocinfoparser.go` | `lookAheadForType` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2835 | 1 | 1 |  |
| 812 | 100 | `internal/rhino/jsdoc.go` | `IsNoDts` | `com/google/javascript/rhino/JSDocInfo.java` | 873 | 1 | 1 |  |
| 813 | 100 | `internal/rhino/jstype/template_type.go` | `SetValidator` | `com/google/javascript/rhino/jstype/TemplateType.java` | 145 | 1 | 1 |  |
| 814 | 100 | `internal/rhino/jstype/template_type.go` | `ContainsCycle` | `com/google/javascript/rhino/jstype/TemplateType.java` | 157 | 3 | 2 |  |
| 815 | 100 | `internal/rhino/jsdoc.go` | `IsEncourageInlining` | `com/google/javascript/rhino/JSDocInfo.java` | 1738 | 1 | 1 |  |
| 816 | 100 | `internal/rhino/jsdoc.go` | `IsRequireInlining` | `com/google/javascript/rhino/JSDocInfo.java` | 1730 | 1 | 1 |  |
| 817 | 100 | `internal/jscomp/checkjsdoc.go` | `validateTsType` | `com/google/javascript/jscomp/CheckJSDoc.java` | 732 | 4 | 4 |  |
| 818 | 100 | `internal/rhino/jsdoc.go` | `IsWizactionRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2685 | 1 | 1 |  |
| 819 | 100 | `internal/rhino/jstype/union_type.go` | `MatchesObjectContext` | `com/google/javascript/rhino/jstype/UnionType.java` | 207 | 1 | 1 |  |
| 820 | 100 | `internal/rhino/jsdoc.go` | `IsNgInjectRecorded` | `com/google/javascript/rhino/JSDocInfo.java` | 2675 | 1 | 1 |  |
| 821 | 100 | `internal/rhino/jstype/union_type.go` | `ToMaybeUnionType` | `com/google/javascript/rhino/jstype/UnionType.java` | 410 | 1 | 1 |  |
| 822 | 100 | `internal/rhino/jsdoc.go` | `IsInlineType` | `com/google/javascript/rhino/JSDocInfo.java` | 1122 | 1 | 1 |  |
| 823 | 100 | `internal/rhino/jstype/union_type.go` | `GetRestrictedUnion` | `com/google/javascript/rhino/jstype/UnionType.java` | 445 | 5 | 6 |  |
| 824 | 100 | `internal/rhino/jsdoc.go` | `MakesUnrestricted` | `com/google/javascript/rhino/JSDocInfo.java` | 796 | 1 | 1 |  |
| 825 | 100 | `internal/jscomp/checkjsdoc.go` | `validateClosurePrimitive` | `com/google/javascript/jscomp/CheckJSDoc.java` | 710 | 4 | 4 |  |
| 826 | 100 | `internal/jscomp/checkjsdoc.go` | `validateDefinesDeclaration` | `com/google/javascript/jscomp/CheckJSDoc.java` | 696 | 3 | 2 |  |
| 827 | 100 | `internal/jscomp/checkjsdoc.go` | `validateNoSideEffects` | `com/google/javascript/jscomp/CheckJSDoc.java` | 680 | 7 | 6 |  |
| 828 | 100 | `internal/rhino/jsdoc.go` | `IsCollapsibleOrBreakMyCode` | `com/google/javascript/rhino/JSDocInfo.java` | 905 | 1 | 1 |  |
| 829 | 100 | `internal/rhino/jsdoc.go` | `IsPureOrBreakMyCode` | `com/google/javascript/rhino/JSDocInfo.java` | 912 | 1 | 1 |  |
| 830 | 100 | `internal/jscomp/checkjsdoc.go` | `isValidMsgName` | `com/google/javascript/jscomp/CheckJSDoc.java` | 533 | 8 | 6 |  |
| 831 | 100 | `internal/rhino/jsdoc.go` | `IsTypeSummary` | `com/google/javascript/rhino/JSDocInfo.java` | 863 | 1 | 1 |  |
| 832 | 100 | `internal/jscomp/checkjsdoc.go` | `isJSDocOnFunctionNode` | `com/google/javascript/jscomp/CheckJSDoc.java` | 460 | 16 | 23 |  |
| 833 | 100 | `internal/rhino/jsdoc.go` | `IsNoCoverage` | `com/google/javascript/rhino/JSDocInfo.java` | 858 | 1 | 1 |  |
| 834 | 100 | `internal/jscomp/checkjsdoc.go` | `hasClassLevelJsDoc` | `com/google/javascript/jscomp/CheckJSDoc.java` | 400 | 4 | 4 |  |
| 835 | 100 | `internal/rhino/jsdoc.go` | `IsNoCompile` | `com/google/javascript/rhino/JSDocInfo.java` | 868 | 1 | 1 |  |
| 836 | 100 | `internal/rhino/jsdoc.go` | `IsExterns` | `com/google/javascript/rhino/JSDocInfo.java` | 853 | 1 | 1 |  |
| 837 | 100 | `internal/rhino/jsdoc.go` | `HasFileOverview` | `com/google/javascript/rhino/JSDocInfo.java` | 1425 | 1 | 1 |  |
| 838 | 100 | `internal/jscomp/checkglobalthis.go` | `shouldReportThis` | `com/google/javascript/jscomp/CheckGlobalThis.java` | 198 | 4 | 4 |  |
| 839 | 100 | `internal/jscomp/parsing/jsdocinfoparser.go` | `lookAheadForAnnotation` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2839 | 1 | 1 |  |
| 840 | 100 | `internal/jscomp/parsing/jsdocinfoparser.go` | `match` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2761 | 3 | 2 |  |
| 841 | 100 | `internal/rhino/jstype_expr.go` | `GetRoot` | `com/google/javascript/rhino/JSTypeExpression.java` | 194 | 1 | 1 |  |
| 842 | 100 | `internal/rhino/jstype_expr.go` | `GetSourceName` | `com/google/javascript/rhino/JSTypeExpression.java` | 198 | 1 | 1 |  |
| 843 | 100 | `internal/jscomp/parsing/jsdocinfoparser.go` | `skipEOLs` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2798 | 4 | 4 |  |
| 844 | 100 | `internal/jscomp/parsing/jsdocinfoparser.go` | `eatUntilEOLIfNotAnnotation` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2713 | 5 | 4 |  |
| 845 | 100 | `internal/jscomp/parsing/jsdocinfoparser.go` | `eatTokensUntilEOL` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2733 | 5 | 6 |  |
| 846 | 100 | `internal/jscomp/parsing/jsdocinfoparser.go` | `eatIfMatch` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 2767 | 4 | 4 |  |
| 847 | 100 | `internal/jscomp/parsing/jsdocinfoparser.go` | `parseModsTag` | `com/google/javascript/jscomp/parsing/JsDocInfoParser.java` | 1456 | 11 | 14 |  |
| 848 | 100 | `internal/jscomp/parsing/jsdoctokenstream.go` | `peekChar` | `com/google/javascript/jscomp/parsing/JsDocTokenStream.java` | 309 | 3 | 3 |  |
| 849 | 100 | `internal/rhino/jsdoc.go` | `IsAnyIdGenerator` | `com/google/javascript/rhino/JSDocInfo.java` | 733 | 1 | 1 |  |
| 850 | 100 | `internal/rhino/jsdoc.go` | `IsUsedViaDotConstructor` | `com/google/javascript/rhino/JSDocInfo.java` | 1286 | 1 | 1 |  |
| 851 | 100 | `internal/rhino/jsdoc.go` | `IsDefine` | `com/google/javascript/rhino/JSDocInfo.java` | 814 | 1 | 1 |  |
| 852 | 100 | `internal/rhino/jsdoc.go` | `IsImplicitCast` | `com/google/javascript/rhino/JSDocInfo.java` | 843 | 1 | 1 |  |
| 853 | 100 | `internal/rhino/jsdoc.go` | `IsNoInline` | `com/google/javascript/rhino/JSDocInfo.java` | 1723 | 1 | 1 |  |
| 854 | 100 | `internal/rhino/jsdoc.go` | `IsNoCollapse` | `com/google/javascript/rhino/JSDocInfo.java` | 878 | 1 | 1 |  |
| 855 | 100 | `internal/jscomp/checkconformance.go` | `shouldCheck` | `com/google/javascript/jscomp/CheckConformance.java` | 128 | 1 | 1 |  |
| 856 | 100 | `internal/rhino/jsdoc.go` | `IsExport` | `com/google/javascript/rhino/JSDocInfo.java` | 838 | 1 | 1 |  |
| 857 | 100 | `internal/rhino/qualified_name.go` | `GetRoot` | `com/google/javascript/rhino/QualifiedName.java` | 98 | 4 | 4 |  |
| 858 | 100 | `internal/rhino/jsdoc.go` | `IsMixinFunction` | `com/google/javascript/rhino/JSDocInfo.java` | 1260 | 1 | 1 |  |
| 859 | 100 | `internal/rhino/qualified_name.go` | `buildComponents` | `com/google/javascript/rhino/QualifiedName.java` | 116 | 4 | 4 |  |
| 860 | 100 | `internal/rhino/jsdoc.go` | `IsMixinClass` | `com/google/javascript/rhino/JSDocInfo.java` | 1255 | 1 | 1 |  |
| 861 | 100 | `internal/rhino/qualified_name.go` | `Getprop` | `com/google/javascript/rhino/QualifiedName.java` | 137 | 1 | 1 |  |
| 862 | 100 | `internal/rhino/jsdoc.go` | `IsCustomElement` | `com/google/javascript/rhino/JSDocInfo.java` | 1250 | 1 | 1 |  |
| 863 | 100 | `internal/jscomp/checkclosureimports.go` | `checkImport` | `com/google/javascript/jscomp/CheckClosureImports.java` | 437 | 53 | 47 |  |
| 864 | 100 | `internal/jscomp/checkclosureimports.go` | `mustBeOrdered` | `com/google/javascript/jscomp/CheckClosureImports.java` | 87 | 1 | 1 |  |
| 865 | 100 | `internal/rhino/jsdoc.go` | `IsPolymer` | `com/google/javascript/rhino/JSDocInfo.java` | 1245 | 1 | 1 |  |
| 866 | 100 | `internal/rhino/jsdoc.go` | `IsPolymerBehavior` | `com/google/javascript/rhino/JSDocInfo.java` | 1240 | 1 | 1 |  |
| 867 | 100 | `internal/rhino/jsdoc.go` | `IsWizaction` | `com/google/javascript/rhino/JSDocInfo.java` | 1230 | 1 | 1 |  |
| 868 | 100 | `internal/jscomp/basicblock.go` | `GetRoot` | `com/google/javascript/jscomp/BasicBlock.java` | 103 | 1 | 1 |  |
| 869 | 100 | `internal/rhino/jsdoc.go` | `IsNgInject` | `com/google/javascript/rhino/JSDocInfo.java` | 1225 | 1 | 1 |  |
| 870 | 100 | `internal/rhino/qualified_name.go` | `GetOwner` | `com/google/javascript/rhino/QualifiedName.java` | 221 | 1 | 1 |  |
| 871 | 100 | `internal/jscomp/basicblock.go` | `IsLoop` | `com/google/javascript/jscomp/BasicBlock.java` | 99 | 1 | 1 |  |
| 872 | 100 | `internal/rhino/jsdoc.go` | `SetBit` | `com/google/javascript/rhino/JSDocInfo.java` | 2898 | 4 | 4 |  |
| 873 | 100 | `internal/rhino/jsdoc.go` | `CheckBit` | `com/google/javascript/rhino/JSDocInfo.java` | 407 | 1 | 1 |  |
| 874 | 100 | `internal/jscomp/basicblock.go` | `IsFunction` | `com/google/javascript/jscomp/BasicBlock.java` | 95 | 1 | 1 |  |
| 875 | 100 | `internal/jscomp/basicblock.go` | `ProvablyExecutesBefore` | `com/google/javascript/jscomp/BasicBlock.java` | 81 | 5 | 7 |  |
| 876 | 100 | `internal/jscomp/basicblock.go` | `IsGlobalScopeBlock` | `com/google/javascript/jscomp/BasicBlock.java` | 76 | 1 | 1 |  |
| 877 | 100 | `internal/jscomp/basicblock.go` | `GetParent` | `com/google/javascript/jscomp/BasicBlock.java` | 65 | 1 | 1 |  |
| 878 | 100 | `internal/jscomp/parsing/jsdoctokenstream.go` | `matchChar` | `com/google/javascript/jscomp/parsing/JsDocTokenStream.java` | 257 | 4 | 6 |  |
| 879 | 100 | `internal/jscomp/syntacticscopecreator.go` | `OnRedeclaration` | `com/google/javascript/jscomp/SyntacticScopeCreator.java` | 405 | 1 | 1 |  |
| 880 | 100 | `internal/jscomp/ast_analyzer.go` | `checkForChangeStateInChildren` | `com/google/javascript/jscomp/AstAnalyzer.java` | 450 | 4 | 4 |  |
| 881 | 100 | `internal/jscomp/typedscope.go` | `Typed` | `com/google/javascript/jscomp/TypedScope.java` | 116 | 1 | 1 |  |
| 882 | 100 | `internal/rhino/jsdoc.go` | `IsInterface` | `com/google/javascript/rhino/JSDocInfo.java` | 829 | 1 | 1 |  |
| 883 | 100 | `internal/jscomp/typemismatch.go` | `registerIfMismatch` | `com/google/javascript/jscomp/TypeMismatch.java` | 118 | 2 | 2 |  |
| 884 | 100 | `internal/rhino/token_stream.go` | `IsJSIdentifier` | `com/google/javascript/rhino/TokenStream.java` | 203 | 1 | 1 |  |
| 885 | 100 | `internal/jscomp/typemismatch.go` | `CreateForTesting` | `com/google/javascript/jscomp/TypeMismatch.java` | 70 | 1 | 1 |  |
| 886 | 100 | `internal/rhino/token_util.go` | `IsJSFormatChar` | `com/google/javascript/rhino/TokenUtil.java` | 68 | 1 | 1 |  |
| 887 | 100 | `internal/rhino/token_util.go` | `IsWhitespace` | `com/google/javascript/rhino/TokenUtil.java` | 72 | 1 | 1 |  |
| 888 | 100 | `internal/jscomp/ast_analyzer.go` | `MayHaveSideEffects` | `com/google/javascript/jscomp/AstAnalyzer.java` | 133 | 1 | 1 |  |
