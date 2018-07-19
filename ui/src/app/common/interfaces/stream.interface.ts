export class Stream {
  constructor(
    public name: string,
    public description: string,
    public type: string,
    public url: string,
    public price: number,
    public tags: [string]
  ) { }
}
