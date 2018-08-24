import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
// Pipes
import { MitasPipe, TasPipe, SubscriptionTypePipe } from './pipes/converter.pipe';
// services
import { StreamService } from './services/stream.service';
import { SubscriptionService } from './services/subscription.service';

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
  ],
  providers: [
      TasPipe,
      MitasPipe,
      SubscriptionTypePipe,
      StreamService,
      SubscriptionService
  ],
  exports: [
      TasPipe,
      MitasPipe,
      SubscriptionTypePipe
  ]
})
export class CommonAppModule { }
