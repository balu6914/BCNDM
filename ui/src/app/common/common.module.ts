import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MdlModule } from '@angular-mdl/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { LayoutModule } from '../layout'
import { HttpModule } from '@angular/http';
// Interfaces
import { Stream } from './interfaces/stream.interface';
import { Subscription } from './interfaces/subscription.interface';
import { User } from './interfaces/user.interface';
// Pipes
import { MitasPipe, TasPipe} from './pipes/converter.pipe';

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
