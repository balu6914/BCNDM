import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { AppBootstrapModule } from '../../app-bootstrap/app-bootstrap.module';
import { CommonAppModule } from '../../common/common.module';
import { BalanceAddComponent } from './add/balance.add.component';
import { BalanceComponent } from './balance.component';
import { BalanceWithdrawComponent } from './withdraw/balance.withdraw.component';



@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    AppBootstrapModule,
    CommonAppModule
  ],
  declarations: [
    BalanceComponent,
    BalanceAddComponent,
    BalanceWithdrawComponent,
  ],
  exports: [
    BalanceComponent,
  ],
  entryComponents: [
    BalanceAddComponent,
    BalanceWithdrawComponent,
  ]

})
export class BalanceModule { }
