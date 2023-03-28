import { NgModule } from '@angular/core';

import { KubeflowComponent } from './kubeflow.component';
import { KubeflowRoutingModule } from './kubeflow.routes';

@NgModule({
  imports: [
    KubeflowRoutingModule,
  ],
  declarations: [
    KubeflowComponent,
  ],
})
export class KubeflowModule { }
