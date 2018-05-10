import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MdlModule } from '@angular-mdl/core';
import {
    FormsModule,
    ReactiveFormsModule
} from '@angular/forms';
// Add wallet routes
import { DashboardWalletRoutingModule } from './wallet.routes';
// Add components
import { WalletBalanceComponent } from './balance';
import { WalletAddComponent } from './add';
import { CommonAppModule } from '../../common/common.module';

@NgModule({
  imports: [
    CommonModule,
    MdlModule,
    FormsModule,
    ReactiveFormsModule,
    CommonAppModule,
    DashboardWalletRoutingModule,
  ],
  declarations: [
      WalletBalanceComponent,
      WalletAddComponent,
  ],
  entryComponents:[WalletAddComponent]
})
export class WalletModule { }
