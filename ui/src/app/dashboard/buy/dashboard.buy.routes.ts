import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuardService as AuthGuard } from '../../auth/guardians/auth.guardian';
import { DashboardBuyComponent } from './main/dashboard.buy.component';

// Define our Auth Routes
const routes: Routes = [
  {
    path: '',
    component: DashboardBuyComponent,
    canActivate: [AuthGuard]
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})

export class DashboardBuyRoutingModule { }
