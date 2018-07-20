import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MdlModule } from '@angular-mdl/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
// Pipes
import { MitasPipe, TasPipe} from './pipes/converter.pipe';
// Google Map
import { MapComponent } from './map/map.component';

@NgModule({
  imports: [
    CommonModule,
    MdlModule,
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
      MitasPipe
  ],
  exports: [
      TasPipe,
      MitasPipe,
      MapComponent,
  ]
})
export class CommonAppModule { }
