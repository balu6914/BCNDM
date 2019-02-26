
export class BigQuery {
  constructor(
    public project: string,
    public dataset: string,
    public table: string,
    public fields: string
  ) { }
}

export class Stream {
  constructor(
    public visibility: string,
    public name: string,
    public type: string,
    public description: string,
    public price: number,
    public location: object,
    public url?: string,
    public snippet?: string,
    public owner?: string,
    public id?: string,
    public external?: boolean,
    public bq?: BigQuery
  ) { }
}
