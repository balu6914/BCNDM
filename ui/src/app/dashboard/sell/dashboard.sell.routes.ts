import { NgModule } from '@angular/core';
import { Routes, CanActivate, RouterModule } from '@angular/router';
import { DashboardSellAddComponent } from './add/dashboard.sell.add.component';
import { DashboardSellEditComponent } from './edit/dashboard.sell.edit.component';
import { DashboardSellComponent } from './dashboard.sell.component';
import { AuthGuardService as AuthGuard } from 'app/auth/guardians/auth.guardian';

// Define our Auth Routes
const routes: Routes = [
    {
        path: '',
        component: DashboardSellComponent ,
        canActivate: [AuthGuard]
    },
    {
        path: 'add',
        component: DashboardSellAddComponent ,
        canActivate: [AuthGuard],
    },
    {
        path: 'edit/:id',
        component: DashboardSellEditComponent ,
        canActivate: [AuthGuard],
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})

export class DashboardSellRoutingModule { }
