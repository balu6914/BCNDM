import { Component, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormControl, ReactiveFormsModule, FormBuilder, Validators  } from '@angular/forms';
import { Router } from '@angular/router';
import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/throw';

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
                email:            ['', emailValidator],
                password:         ['', [Validators.required, Validators.minLength(5)]],
                passwordConfirm:  [''],
                first_name:       ['', [Validators.required, Validators.minLength(2)]],
                last_name:        ['', [Validators.required, Validators.minLength(2)]]
        },
        {
                validator: this.passwordMatchValidator
        });
    }

    passwordMatchValidator(fg: FormGroup) {
      // Compare passwords only if minLength is valid
      if (fg.get('passwordConfirm').value.length >= 5) {
        if (fg.get('password').value === fg.get('passwordConfirm').value) {
          return null;
        }
        fg.controls.passwordConfirm.setErrors({'invalid': true});
      }
      return {'mismatch': true };
    }

    onSubmit(model: User, isValid: boolean) {
        this.errorMsg = null;
        if(isValid) {
            this.UserService.addUser(model).subscribe(
                response => {
                    this.AuthService.login(model.email, model.password)
                      .subscribe(
                        (token: any) => this.router.navigate(['/dashboard']),
                        err => { console.log("Response error: ", err) }
                      );
                },
                err => {
                  // Handle tmp case when user already exists and we don't have error msg on API side yet.
                  if (err.status === 409) {
                    this.errorMsg = "User with this email already exists."
                  } else {
                    this.errorMsg = err;
                  }
                });
        }
    }

}
