import { Stream, Subscription, Contract, Access, Execution, User } from 'app/common/interfaces';
import { Page } from 'app/common/interfaces/page.interface';

export enum TableType {
  Buy,
  Sell,
  Contract,
  Access,
  Dashboard,
  Wallet,
  Ai,
  Executions,
  Users,
}

export class Table {
  title: string;
  headers: string[];
  page: Page<Stream | Subscription | Contract | Access | Execution | User>;
  tableType: TableType;
  hasDetails: Boolean = false;

  constructor() {
    this.page = new Page(0, 20, 0, []);
  }
}
