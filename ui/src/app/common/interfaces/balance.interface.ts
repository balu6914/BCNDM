export class Balance {
  constructor(
    public amount: number = 0,
    public symbol: string = 'TAS',
    public fiatAmount: number = 0,
    public fiatSymbol: string = 'EUR'
  ) { }
}
