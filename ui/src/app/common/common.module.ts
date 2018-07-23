import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
// Pipes
import { MitasPipe, TasPipe} from './pipes/converter.pipe';

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
  ],
  providers: [
      TasPipe,
      MitasPipe
  ],
  exports: [
      TasPipe,
      MitasPipe
  ]
})
export class CommonAppModule { }
