export class Subscription {
  constructor(
    public id: string,
    public user_id: string,
    public stream_id: string,
    public stream_owner: string,
    public hours: string,
    public price: number,
    public start_date: number,
    public end_date: number,
    public url: string,
    public type?: string
  ) { }
}
