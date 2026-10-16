/**
 * @template T := cond(isUnknown(T), 'Object', mapunion(T, (x) => rawTypeOf(x))) =:
 * @template R := record({a: 'number'}, typeExpr('string'), union('a', 'b')) =:
 * @template S := maprecord(T, (k, v) => printType('m', propType('p', v))) =:
 * @template U := templateTypeOf(T, 1.5) =:
 * @template V := bogus(T) =:
 * @template W := sub(T, eq(T, R)) =:
 * @template X := type('Array', typeOfVar('y'), instanceOf(T), none(), all(), unknown()) =:
 * @template Y := cond(!isCtor(T) && (isRecord(T) || isTemplatized(T)), 'a', 'b') =:
 * @template Z := cond(streq('a', ''), 'a', 'b') =:
 * @template Q := cond(isDefined(T), 'a', 'b') =:
 * @template P := 1 + 2 =:
 */
function f() {}
