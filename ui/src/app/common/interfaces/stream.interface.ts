export class Stream {
  constructor(
    public name: string,
    public type: string,
    public description: string,
    public url: string,
    public price: number,
    public location: object,
    public snippet?: string,
    public owner?: string,
    public id?: string,
    public bq?: boolean,
    public project?: string,
    public dataset?: string,
    public table?: string,
    public fields?: string
  ) { }
}
