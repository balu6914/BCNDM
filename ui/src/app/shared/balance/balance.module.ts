import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ClipboardModule } from 'ngx-clipboard';

import { CommonAppModule } from 'app/common/common.module';
import { BalanceAddComponent } from './add/balance.add.component';
import { BalanceComponent } from './balance.component';
import { BalanceWithdrawComponent } from './withdraw/balance.withdraw.component';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    CommonAppModule,
    ClipboardModule,
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
    BalanceComponent,
    BalanceAddComponent,
    BalanceWithdrawComponent,
  ]

})
export class BalanceModule { }
