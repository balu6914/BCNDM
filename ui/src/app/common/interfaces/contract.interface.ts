export class Contract {
  constructor(
    public id: string,
    public stream: object,
    public creationDate: string,
    public expirationDate: string,
    public share: string,
    public signed: boolean,
    public expired: boolean,
  ) { }
}
