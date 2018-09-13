import { Component } from '@angular/core';
import { FormGroup, FormControl, ReactiveFormsModule, FormBuilder, Validators  } from '@angular/forms';

import {Router} from '@angular/router';
import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/throw';

import { User } from 'app/common/interfaces/user.interface';
import { UserService } from 'app/common/services/user.service';
import { AuthService } from '../services/auth.service';

const emailValidator = Validators.pattern('^[a-z]+[a-z0-9._]+@[a-z]+\.[a-z.]{2,5}$');

@Component({
  selector: 'dpc-login-form',
  templateUrl: './login.component.html',
  styleUrls: [ './login.component.scss' ],
  providers: [
  ],
})

export class LoginComponent {
    public form: FormGroup;
    private subscription;
    public errorMsg: String;
    constructor(
        private fb: FormBuilder,
        private router: Router,
        private UserService: UserService,
        private AuthService: AuthService

    ) {

    }

    ngOnInit() {
        this.errorMsg = null;
        this.form = this.fb.group({
                email: ['', [<any>Validators.required, emailValidator]],
                password: ['', [<any>Validators.required, <any>Validators.minLength(5)]]
        });
    }

    onSubmit(model: User, isValid: boolean) {
        this.errorMsg = null;
        if (isValid) {
            this.AuthService.login(model.email, model.password)
              .subscribe(
                (token: any) => {
                    this.router.navigate(['/dashboard'])
                },
                err => {
                    this.errorMsg = 'Invalid Credentials.'
                 }
              );
        }
    }

}
