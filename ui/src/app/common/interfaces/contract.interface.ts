export class Contract {
  constructor(
    public stream_id: string,
    public stream_name: string,
    public start_time: string,
    public end_time: string,
    public share: string,
    public signed: boolean
  ) { }
}
