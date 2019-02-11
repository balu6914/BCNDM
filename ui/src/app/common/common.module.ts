import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
// Pipes
import { MitasPipe, SubscriptionTypePipe, TasPipe } from './pipes/converter.pipe';
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
    TasPipe,
    MitasPipe,
    SubscriptionTypePipe,
    WalletBalanceStatisticPipe,
  ],
  providers: [
    TasPipe,
    MitasPipe,
    SubscriptionTypePipe,
    WalletBalanceStatisticPipe,
    StreamService,
    SubscriptionService,
    UserService,
    ExecutionsService,
  ],
  exports: [
    TasPipe,
    MitasPipe,
    SubscriptionTypePipe,
    WalletBalanceStatisticPipe,
  ]
})
export class CommonAppModule { }
