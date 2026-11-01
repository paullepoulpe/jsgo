// Sample program for the behavioural oracle's TypeScript workload (tools/behave/workloads.py).
export enum Kind { Circle = 1, Square, Polygon = "poly" }
export const enum Flags { None = 0, Filled = 1 << 0, Stroked = 1 << 1, Both = Filled | Stroked }

export interface Point { readonly x: number; readonly y: number }
export interface Shape { kind: Kind; area(): number; centre?: Point; flags: Flags }

export abstract class Base implements Shape {
  static count = 0;
  #id: number;
  protected constructor(public kind: Kind, public flags: Flags = Flags.None) { this.#id = ++Base.count; }
  abstract area(): number;
  get id(): number { return this.#id; }
  describe(this: Base): string { return `${Kind[this.kind] ?? this.kind}#${this.id}: ${this.area().toFixed(2)}`; }
}

export class Circle extends Base {
  constructor(public r: number, public centre: Point = { x: 0, y: 0 }) { super(Kind.Circle, Flags.Filled); }
  area() { return Math.PI * this.r ** 2; }
}

export class Square extends Base {
  constructor(private side: number) { super(Kind.Square, Flags.Both); }
  area() { return this.side * this.side; }
  override describe(): string { return "square " + super.describe(); }
}

export class Polygon extends Base {
  constructor(private pts: readonly Point[]) { super(Kind.Polygon); }
  area() {
    let s = 0;
    for (let i = 0, j = this.pts.length - 1; i < this.pts.length; j = i++) {
      const { x: xi, y: yi } = this.pts[i], { x: xj, y: yj } = this.pts[j];
      s += (xj + xi) * (yj - yi);
    }
    return Math.abs(s / 2);
  }
}

export function isCircle(s: Shape): s is Circle { return s.kind === Kind.Circle; }

export function largest<T extends Shape>(shapes: Iterable<T>): T | undefined {
  let best: T | undefined;
  for (const s of shapes) if (!best || s.area() > best.area()) best = s;
  return best;
}

export function overload(x: string): string;
export function overload(x: number): number;
export function overload(x: string | number) { return typeof x === "string" ? x.toUpperCase() : x * 2; }

const wrong: number = "not a number";
export const total = [new Circle(1), new Square(2)].reduce((a, s) => a + s.area(), 0) + wrong;
