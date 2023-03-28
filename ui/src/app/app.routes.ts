import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { NoContentComponent } from './no-content/no-content.component';
import { DashboardModule } from './dashboard/dashboard.module';

import { SignupComponent } from 'app/dashboard/signup/signup.component';
import { LoginComponent } from 'app/dashboard/login/login.component';
import { RecoverComponent } from 'app/dashboard/recover/recover.component';
import { AuthLoggedinGuardService as AuthGuard } from 'app/auth/guardians/auth.loggedin.guardian';

// Define our Application Routes
const routes: Routes = [
  {
    path: 'login',
    component: LoginComponent,
    canActivate: [AuthGuard]
  },
  {
    path: 'signup',
    component: SignupComponent,
    canActivate: [AuthGuard]
  },
  {
    path: 'recover',
    component: RecoverComponent,
    canActivate: [AuthGuard]
  },
  // Not found page
  {
    path: '',
    redirectTo: 'login',
    pathMatch: 'full',
  },
  {
    path: '**',
    component: NoContentComponent,
  },
];

@NgModule({
    imports: [
        DashboardModule,
        RouterModule.forRoot(routes, {useHash: false})],
    exports: [RouterModule],
})

export class AppRoutingModule { }
