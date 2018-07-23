import { Stream, Subscription, Contract } from '../../common/interfaces';

export enum TableType {
  Buy,
  Sell,
  Contract,
  Transaction,
  Dashboard
}

export class Table {
  title: string;
  headers: string[];
  content: Stream[] | Subscription[] | Contract[]
  tableType: TableType

  constructor() {
  }
}
