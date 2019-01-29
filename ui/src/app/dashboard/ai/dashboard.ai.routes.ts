import { NgModule } from '@angular/core';
import { Routes, CanActivate, RouterModule } from '@angular/router';

import { DashboardAiComponent } from './dashboard.ai.component';
import { DashboardAiExecuteComponent } from './execute/dashboard.ai.execute.component';
import { AuthGuardService as AuthGuard } from 'app/auth/guardians/auth.guardian';

// Define our Auth Routes
const routes: Routes = [
    {
        path: '',
        component: DashboardAiComponent ,
        canActivate: [AuthGuard]
    },
    {
        path: 'execute',
        component: DashboardAiExecuteComponent ,
        canActivate: [AuthGuard],
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})

export class DashboardAiRoutingModule { }
