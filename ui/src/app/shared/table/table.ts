import { Stream, Subscription, Contract } from '../../common/interfaces';
import { Page } from '../../common/interfaces/page.interface';

export enum TableType {
  Buy,
  Sell,
  Contract,
  Dashboard,
  Wallet
}

export class Table {
  title: string;
  headers: string[];
  page: Page<Stream | Subscription | Contract>;
  tableType: TableType;
  hasDetails: Boolean = false;

  constructor() {
    this.page = new Page(0, 20, 0, []);
  }
}
