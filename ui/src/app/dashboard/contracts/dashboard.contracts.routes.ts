import { NgModule } from '@angular/core';
import { Routes, CanActivate, RouterModule } from '@angular/router';
import { DashboardContractsComponent } from './dashboard.contracts.component';
import { DashboardContractsAddComponent } from './add/dashboard.contracts.add.component';
import { AuthGuardService as AuthGuard } from 'app/auth/guardians/auth.guardian';

// Define our Auth Routes
const routes: Routes = [
    {
        path: '',
        component: DashboardContractsComponent ,
        canActivate: [AuthGuard],
    },
    {
        path: 'add',
        component: DashboardContractsAddComponent ,
        canActivate: [AuthGuard]
    }
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})

export class DashboardContractsRoutingModule { }
