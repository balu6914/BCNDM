import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { NoContentComponent } from './no-content';
import { DashboardModule } from './dashboard/dashboard.module';

// Define our Application Routes
const routes: Routes = [
  // Not found page
  {path: '', redirectTo: "login", pathMatch: 'full'},
  { path: '**',    component: NoContentComponent },
];

@NgModule({
    imports: [
        DashboardModule,
        RouterModule.forRoot(routes, {useHash: false})],
    exports: [RouterModule],
})

export class AppRoutingModule { }
