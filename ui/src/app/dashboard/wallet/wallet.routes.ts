import { NgModule } from '@angular/core';
import { Routes, CanActivate, RouterModule } from '@angular/router';
import { WalletBalanceComponent } from './balance';
import { AuthGuardService as AuthGuard } from '../../auth/guardians/auth.guardian';

// Define our Auth Routes
const routes: Routes = [
       {
           path: 'balance',
           component: WalletBalanceComponent ,
           canActivate: [AuthGuard],
       },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})

export class DashboardWalletRoutingModule { }
