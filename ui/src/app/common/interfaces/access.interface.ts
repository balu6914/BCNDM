export class Access {
  constructor(
    public receiver: string,
    public state: string,
    public origin: string,
    public id?: string,
  ) { }
}
