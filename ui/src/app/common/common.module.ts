import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
// Pipes
import { MitasPipe, TasPipe} from './pipes/converter.pipe';
// services
import { StreamService } from './services/stream.service';
import { SubscriptionService } from './services/subscription.service';


// Google Map
import { MapComponent } from './map/map.component';

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
    MapComponent,
  ],
  providers: [
      TasPipe,
      MitasPipe,
      StreamService,
      SubscriptionService
  ],
  exports: [
      TasPipe,
      MitasPipe,
      MapComponent,
  ]
})
export class CommonAppModule { }
