export class User {
  constructor(
    public id: string,
    public email: string,
    public password: string,
    public contact_email: string,
    public first_name: string,
    public last_name: string,
  ) { }
}
