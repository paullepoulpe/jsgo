// Sample program for the behavioural oracle's TypeScript workload (tools/behave/workloads.py).
export namespace Text {
  export const enum Align { Left, Right, Centre }
  export function pad(s: string, width: number, align: Align = Align.Left, fill = " "): string {
    const gap = Math.max(0, width - s.length);
    switch (align) {
      case Align.Left: return s + fill.repeat(gap);
      case Align.Right: return fill.repeat(gap) + s;
      default: { const l = gap >> 1; return fill.repeat(l) + s + fill.repeat(gap - l); }
    }
  }
  export namespace Table {
    export function render(rows: ReadonlyArray<readonly unknown[]>): string {
      const cols = Math.max(...rows.map(r => r.length));
      const widths = Array.from({ length: cols }, (_, i) => Math.max(...rows.map(r => String(r[i] ?? "").length)));
      return rows.map(r => r.map((c, i) => pad(String(c ?? ""), widths[i]!, typeof c === "number" ? Align.Right : Align.Left)).join(" | ")).join("\n");
    }
  }
}

export function memo<A extends unknown[], R>(fn: (...a: A) => R, key: (...a: A) => string = (...a) => JSON.stringify(a)) {
  const cache = new Map<string, R>();
  return (...a: A): R => {
    const k = key(...a);
    if (!cache.has(k)) cache.set(k, fn(...a));
    return cache.get(k)!;
  };
}

export const fib = memo((n: number): number => (n < 2 ? n : fib(n - 1) + fib(n - 2)));

export function assertNever(x: never): never { throw new Error("unexpected " + x); }

export type Result<T, E = Error> = { ok: true; value: T } | { ok: false; error: E };
export function tryParse(s: string): Result<unknown, SyntaxError> {
  try { return { ok: true, value: JSON.parse(s) }; }
  catch (e) { return { ok: false, error: e as SyntaxError }; }
}
export function unwrap<T>(r: Result<T>): T {
  if (r.ok) return r.value;
  throw r.error;
}

export function* range(start: number, end?: number, step = 1): Generator<number, void, undefined> {
  if (end === undefined) [start, end] = [0, start];
  for (let i = start; step > 0 ? i < end : i > end; i += step) yield i;
}

let unused = Text.pad("x", 3, Text.Align.Centre, 4);
