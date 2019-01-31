import { NgModule } from '@angular/core';
import { Routes, CanActivate, RouterModule } from '@angular/router';
import { DashboardAccessComponent } from './dashboard.access.component';
import { DashboardAccessAddComponent } from './add/dashboard.access.add.component';
import { AuthGuardService as AuthGuard } from 'app/auth/guardians/auth.guardian';

// Define our Auth Routes
const routes: Routes = [
    {
        path: '',
        component: DashboardAccessComponent ,
        canActivate: [AuthGuard],
    },
    {
        path: 'add',
        component: DashboardAccessAddComponent ,
        canActivate: [AuthGuard]
    }
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})

export class DashboardAccessRoutingModule { }
