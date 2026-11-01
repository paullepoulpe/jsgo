// Sample program for the behavioural oracle's TypeScript workload (tools/behave/workloads.py).
import { Shape, Circle, largest, Kind } from "./shapes";

type Listener<E> = (event: E) => void;
type Events = { add: Shape; remove: { id: number }; clear: undefined };
type Handlers = { [K in keyof Events as `on${Capitalize<K>}`]?: Listener<Events[K]> };
type Awaited2<T> = T extends PromiseLike<infer U> ? Awaited2<U> : T;
type DeepReadonly<T> = { readonly [K in keyof T]: T[K] extends object ? DeepReadonly<T[K]> : T[K] };

export class Store<S extends Shape = Shape> implements Iterable<S> {
  private items = new Map<number, S>();
  private listeners: { [K in keyof Events]?: Listener<Events[K]>[] } = {};
  constructor(private readonly handlers: Handlers = {}) {}

  on<K extends keyof Events>(name: K, fn: Listener<Events[K]>): () => void {
    (this.listeners[name] ??= [] as any).push(fn);
    return () => { this.listeners[name] = (this.listeners[name] as Listener<Events[K]>[]).filter(f => f !== fn) as any; };
  }
  private emit<K extends keyof Events>(name: K, e: Events[K]) {
    for (const fn of (this.listeners[name] ?? []) as Listener<Events[K]>[]) fn(e);
  }
  add(...shapes: S[]): this {
    for (const s of shapes) { this.items.set(this.items.size + 1, s); this.emit("add", s); this.handlers.onAdd?.(s); }
    return this;
  }
  *[Symbol.iterator](): Iterator<S> { yield* this.items.values(); }
  async loadAll(sources: Array<() => Promise<S>>): Promise<number> {
    const loaded = await Promise.all(sources.map(async (f) => f()));
    this.add(...loaded);
    return loaded.length;
  }
  snapshot(): DeepReadonly<{ count: number; biggest?: S }> {
    return { count: this.items.size, biggest: largest(this) };
  }
}

export function byKind(store: Store, k: Kind) {
  const out: Shape[] = [];
  for (const s of store) if (s.kind === k) out.push(s);
  return out;
}

declare const ready: Awaited2<Promise<Promise<number>>>;
export const n: number = ready;
const s = new Store<Circle>().add(new Circle(2));
s.add(42);
export default s;
