import { Component } from '@angular/core';
import { AuthService } from '../../auth/services/auth.service';

@Component({
    selector: 'header-component',
    templateUrl: './header.component.html',
    styleUrls: [ './header.component.scss' ],
    providers: [
    ],
})

export class HeaderComponent {
    isLoggedin: Boolean;
    loggedSubscription: any;
    subscription: any;
    user: any;

    constructor(
        private AuthService: AuthService
    ) {}

    ngOnInit() {
        this.AuthService.loggedIn
            .subscribe(is => {
                this.isLoggedin = is;

                this.AuthService.getCurrentUser()
                    .subscribe(data => {
                        this.user = data;
                    });
            })
    }

    logout() {
        this.AuthService.logout();
        this.isLoggedin = false;
    }
}
