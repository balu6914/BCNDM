import { NgModule } from '@angular/core';
import { Routes, CanActivate, RouterModule } from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { DashboardMainComponent } from './main';
import { AuthGuardService as AuthGuard } from 'app/auth/guardians/auth.guardian';

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
      }, {
        path: 'sell',
        loadChildren: 'app/dashboard/sell/dashboard.sell.module#DashboardSellModule',
      }, {
        path: 'buy',
        loadChildren: 'app/dashboard/buy/dashboard.buy.module#DashboardBuyModule',
      }, {
        path: 'wallet',
        loadChildren: 'app/dashboard/wallet/dashboard.wallet.module#DashboardWalletModule',
      }, {
        path: 'contracts',
        loadChildren: 'app/dashboard/contracts/dashboard.contracts.module#DashboardContractsModule',
      }, {
        path: 'access',
        loadChildren: 'app/dashboard/access/dashboard.access.module#DashboardAccessModule',
      }, {
        path: 'ai',
        loadChildren: 'app/dashboard/ai/dashboard.ai.module#DashboardAiModule',
      }, {
        path: 'kubeflow',
        loadChildren: 'app/dashboard/kubeflow/kubeflow.module#KubeflowModule',
      }, {
        path: 'profile',
        loadChildren: 'app/dashboard/profile/dashboard.profile.module#DashboardProfileModule',
      }, {
        path: 'admin',
        loadChildren: 'app/dashboard/admin/dashboard.admin.module#DashboardAdminModule',
      },
    ]
  },
];

@NgModule({
    imports: [RouterModule.forRoot(routes, {useHash: false})],
    exports: [RouterModule],
})

export class DashboardRoutingModule { }
