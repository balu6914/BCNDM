import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { CommonAppModule } from '../../common/common.module';
import { BalanceComponent } from './main/balance.component';
import { BalanceAddComponent } from './add/balance.add.component';
import { BalanceWithdrawComponent } from './withdraw/balance.withdraw.component';


@NgModule({
  imports: [
    CommonModule,
    CommonAppModule
  ],
  declarations: [
    BalanceAddComponent,
    BalanceWithdrawComponent,
    BalanceComponent,
  ],
  exports: [
    BalanceComponent
  ]
})
export class BalanceModule { }
