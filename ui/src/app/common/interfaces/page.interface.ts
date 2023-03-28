export class Page<T> {
  constructor(
    public page: number,
    public limit: number,
    public total: number,
    public content: T[]
  ) { }
}
