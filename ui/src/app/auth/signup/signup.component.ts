import { Component, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormControl, ReactiveFormsModule, FormBuilder, Validators  } from '@angular/forms';
import {Router} from '@angular/router';

import { User } from '../../common/interfaces/user.interface';
import { UserService } from '../services/user.service';
import { AuthService } from '../services/auth.service';

const emailValidator = Validators.pattern('^[a-z]+[a-z0-9._]+@[a-z]+\.[a-z.]{2,5}$');

@Component({
  selector: 'signup-form',
  templateUrl: './signup.component.html',
  styleUrls: [ './signup.component.scss' ],
  providers: [
  ],
})

export class SignupComponent {
    public form: FormGroup;
    public errorMsg: String;
    @Output()
    userCreated: EventEmitter<any> = new EventEmitter();
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
                password: ['', [<any>Validators.required, <any>Validators.minLength(5)]],
                first_name: ['', []],
                last_name: ['', []]
        });
    }

    onSubmit(model: User, isValid: boolean) {
        this.errorMsg = null;
        if(isValid) {
            this.UserService.addUser(model).subscribe(
                response => {
                    console.log("here is a response", response)
                    this.AuthService.login(model.email, model.password)
                      .subscribe(
                        (token: any) => this.router.navigate(['/dashboard/me']),
                        err => { console.log(err) }
                      );
                },
                err => {
                    this.errorMsg = err;
                });
        }
    }

}
