import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ExecutionResultGraphComponent } from './graph/execution.result..graph.component';


@NgModule({
  imports: [
    CommonModule,
  ],
  declarations: [
    ExecutionResultGraphComponent,
  ],
  exports: [
    ExecutionResultGraphComponent,
  ]
})
export class ExecutionResultModule { }
