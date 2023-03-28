export class Balance {
  constructor(
    public amount: number = 0,
    public symbol: string = 'DPC',
    public fiatAmount: number = 0,
    public fiatSymbol: string = 'EUR'
  ) { }
}
