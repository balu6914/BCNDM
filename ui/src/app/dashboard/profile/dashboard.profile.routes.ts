import { NgModule } from '@angular/core';
import { Routes, CanActivate, RouterModule } from '@angular/router';
import { DashboardProfileComponent } from './dashboard.profile.component';
import { AuthGuardService as AuthGuard } from '../..//auth/guardians/auth.guardian';

// Define our Auth Routes
const routes: Routes = [
       {
           path: '',
           component: DashboardProfileComponent ,
           canActivate: [AuthGuard],
       },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})

export class DashboardProfileRoutingModule { }
