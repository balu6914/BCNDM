import { Component, OnInit, Input } from '@angular/core';
import { Balance } from '../balance';
@Component({
  selector: 'dpc-balance-widget',
  templateUrl: './balance.component.html',
  styleUrls: ['./balance.component.scss']
})
export class BalanceComponent implements OnInit {

  @Input() balance: Balance = new Balance();
  constructor() { }

  ngOnInit() {
  }
}
