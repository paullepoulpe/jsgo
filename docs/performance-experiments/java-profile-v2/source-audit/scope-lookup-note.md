Scope declaration lookup audit (not implemented)

SyntacticScopeCreator.declareVar looks up the name, Scope.declare checks it again,
and AbstractScope.declareInternal calls hasOwnSlot(name) || canDeclare(name), with
canDeclare repeating hasOwnSlot. This is a real source-level duplicate lookup.
However, AbstractScope.hasOwnImplicitSlot is protected and overridable, and
TypedScope is public/nonfinal. Algebraically simplifying repeated calls would
change callbacks for custom subclasses. A private fast path restricted to final
Scope could be explored, but it changes more code than the accumulator or
identity-return candidates and needs proof around compiler hooks and implicit
slots. No runtime gain or patch is claimed. JFR showed SyntacticScopeCreator at
8.91% inclusive CPU samples; that is not attribution to these duplicate lookups.
