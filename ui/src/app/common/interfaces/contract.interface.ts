export class Contract {
  constructor(
    public id: number,
    public stream: string,
    public creationDate: string,
    public expirationDate: string,
    public share: string,
    public signed: boolean,
    public expired: boolean,
  ) { }
}
