import { NgModule } from '@angular/core';
import { Routes, CanActivate, RouterModule } from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { DashboardMainComponent } from './main';
import { AuthGuardService as AuthGuard } from '../auth/guardians/auth.guardian';

// Define our Dashboard Routes
const routes: Routes = [
   {
       path: 'dashboard',
       component: DashboardComponent,
       canActivate: [AuthGuard],
       children: [
           {
               path: '',
               component: DashboardMainComponent,
           },
           {
               path: 'sell',
               loadChildren: 'app/dashboard/sell/dashboard.sell.module#DashboardSellModule',
           },
           {
               path: 'buy',
               loadChildren: 'app/dashboard/buy/dashboard.buy.module#DashboardBuyModule',
           },
           {
               path: 'wallet',
               loadChildren: 'app/dashboard/wallet/wallet.module#WalletModule',
           },
           {
               path: 'contract',
               loadChildren: 'app/dashboard/contracts/dashboard.contracts.module#DashboardContractsModule',
           },

       ]
   },
];

@NgModule({
    imports: [RouterModule.forRoot(routes, {useHash: false})],
    exports: [RouterModule],
})

export class DashboardRoutingModule { }
