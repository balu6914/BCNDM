import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { SignupComponent } from './signup';
import { LoginComponent } from './login';
import { AuthLoggedinGuardService as AuthGuard } from './guardians/auth.loggedin.guardin';

// Define our Auth Routes
const routes: Routes = [
   { path: 'login',  component: LoginComponent, canActivate: [AuthGuard]},
   { path: 'signup',  component: SignupComponent, canActivate: [AuthGuard] },
];

@NgModule({
    imports: [RouterModule.forRoot(routes, {useHash: false})],
    exports: [RouterModule],
})

export class AuthRoutingModule { }
