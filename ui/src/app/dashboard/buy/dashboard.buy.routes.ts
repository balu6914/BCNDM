import { NgModule } from '@angular/core';
import { Routes, CanActivate, RouterModule } from '@angular/router';
import { DashboardBuyMapComponent } from './map/dashboard.buy.map.component';
import { AuthGuardService as AuthGuard } from '../../auth/guardians/auth.guardian';
import { SubscriptionAddComponent } from '../../dashboard/subscription/add';

// Define our Auth Routes
const routes: Routes = [
       {
            path: 'map',
            component: DashboardBuyMapComponent ,
            canActivate: [AuthGuard]
        },
        {
             path: 'subscribe/:id',
             component: SubscriptionAddComponent ,
             canActivate: [AuthGuard]
         },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})

export class DashboardBuyRoutingModule { }
