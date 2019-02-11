import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
// Pipes
import { MidpcPipe, SubscriptionTypePipe, DpcPipe } from './pipes/converter.pipe';
import { WalletBalanceStatisticPipe } from './pipes/balance.income.pipe';
// services
import { StreamService } from './services/stream.service';
import { SubscriptionService } from './services/subscription.service';
import { UserService } from './services/user.service';
import { ExecutionsService } from './services/executions.service';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    HttpModule,
    ReactiveFormsModule,
  ],
  declarations: [
    // Pipes
    DpcPipe,
    MidpcPipe,
    SubscriptionTypePipe,
    WalletBalanceStatisticPipe,
  ],
  providers: [
    DpcPipe,
    MidpcPipe,
    SubscriptionTypePipe,
    WalletBalanceStatisticPipe,
    StreamService,
    SubscriptionService,
    UserService,
    ExecutionsService,
  ],
  exports: [
    DpcPipe,
    MidpcPipe,
    SubscriptionTypePipe,
    WalletBalanceStatisticPipe,
  ]
})
export class CommonAppModule { }
