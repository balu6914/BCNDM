import { NgModule } from '@angular/core';
import { Routes, CanActivate, RouterModule } from '@angular/router';
import { DashboardAdminComponent } from './dashboard.admin.component';
import { AuthGuardService as AuthGuard } from 'app/auth/guardians/auth.guardian';

// Define our Auth Routes
const routes: Routes = [
    {
        path: '',
        component: DashboardAdminComponent ,
        canActivate: [AuthGuard]
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})

export class DashboardAdminRoutingModule { }
