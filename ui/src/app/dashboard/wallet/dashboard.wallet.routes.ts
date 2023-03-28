import { NgModule } from '@angular/core';
import { Routes, CanActivate, RouterModule } from '@angular/router';
import { DashboardWalletComponent } from './dashboard.wallet.component';
import { AuthGuardService as AuthGuard } from 'app/auth/guardians/auth.guardian';

// Define our Auth Routes
const routes: Routes = [
       {
           path: '',
           component: DashboardWalletComponent ,
           canActivate: [AuthGuard],
       },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})

export class DashboardWalletRoutingModule { }
