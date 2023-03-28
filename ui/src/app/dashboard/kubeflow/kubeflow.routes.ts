import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuardService as AuthGuard } from 'app/auth/guardians/auth.guardian';
import { KubeflowComponent } from './kubeflow.component';

// Define our Auth Routes
const routes: Routes = [
  {
    path: '',
    component: KubeflowComponent,
    canActivate: [AuthGuard]
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})

export class KubeflowRoutingModule { }
