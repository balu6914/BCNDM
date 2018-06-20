import { NgModule } from '@angular/core';
import { Routes, CanActivate, RouterModule } from '@angular/router';
import { DashboardSellAddComponent } from './add/dashboard.sell.add.component';
import { DashboardSellEditComponent } from './edit/dashboard.sell.edit.component';
import { DashboardSellMapComponent } from './map/dashboard.sell.map.component';
import { DashboardSellComponent } from './dashboard.sell.component';
import { SubscriptionAddComponent } from '../../dashboard/subscription/add';
import { AuthGuardService as AuthGuard } from '../../auth/guardians/auth.guardian';

// Define our Auth Routes
const routes: Routes = [
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
    {
        path: 'map',
        component: DashboardSellMapComponent ,
        canActivate: [AuthGuard]
    }
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})

export class DashboardSellRoutingModule { }
