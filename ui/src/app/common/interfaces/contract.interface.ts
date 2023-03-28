export class Contract {
  constructor(
    public stream_name: string,
    public start_time: string,
    public end_time: string,
    public partner_id: string,
    public share: string,
    public signed: boolean,
    public stream_id?: string,
  ) { }
}
