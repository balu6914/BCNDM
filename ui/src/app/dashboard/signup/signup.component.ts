import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormControl, ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import { ValidationErrors, ValidatorFn, AbstractControl } from '@angular/forms';
import { Router } from '@angular/router';
import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/throw';

import { User } from 'app/common/interfaces/user.interface';
import { UserService } from 'app/common/services/user.service';
import { AuthService } from 'app/auth/services/auth.service';
import { CustomValidators } from 'app/common/validators/customvalidators';

@Component({
  selector: 'dpc-signup-form',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.scss'],
  providers: [
  ],
})
export class SignupComponent implements OnInit {
  public form: FormGroup;
  public errorMsg: String;
  submitted = false;

  @Output() userCreated: EventEmitter<any> = new EventEmitter();
  constructor(
    private fb: FormBuilder,
    private router: Router,
    private userService: UserService,
    private authService: AuthService
  ) {
  }

  ngOnInit() {
    this.errorMsg = null;
    this.form = this.fb.group({
      email:      ['', [Validators.required, Validators.email, Validators.maxLength(32)]],
      password:   ['', [Validators.required, Validators.minLength(9), Validators.maxLength(32),
                // 2. check whether the entered password has a number
                CustomValidators.patternValidator(/\d/, { hasNumber: true }),
                // 3. check whether the entered password has upper case letter
                CustomValidators.patternValidator(/[A-Z]/, { hasCapitalCase: true }),
                // 4. check whether the entered password has a lower-case letter
                CustomValidators.patternValidator(/[a-z]/, { hasSmallCase: true }),
                // 5. check whether the entered password has a special character
                CustomValidators.patternValidator(/[ !@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/, { hasSpecialCharacters: true })
      ]],
      confirm:    ['', [Validators.required]],
      first_name: ['', [Validators.maxLength(32)]],
      last_name:  ['', [Validators.maxLength(32)]],
      company:    ['', [Validators.maxLength(32)]],
      address:    ['', [Validators.maxLength(128)]],
      phone:      ['', [Validators.maxLength(32)]]
    },
    {
      validator: this.passwordMatchValidator
    });
  }

  passwordMatchValidator(fg: FormGroup) {
    // Compare passwords only if minLength is valid
    if (fg.get('confirm').value.length > 0) {
      if (fg.get('password').value === fg.get('confirm').value) {
        return null;
      }
      fg.controls.confirm.setErrors({ 'invalid': true });
    }

    // Set form confirm password missmatch
    return { 'missmatch': true };
  }

  onSubmit() {
    this.submitted = true;

    if (this.form.valid) {
      const user = {
        email: this.form.value.email,
        password: this.form.value.password,
        first_name: this.form.value.first_name,
        last_name: this.form.value.last_name,
        company: this.form.value.company,
        address: this.form.value.address,
        phone: this.form.value.phone,
      };
      this.userService.addUser(user).subscribe(
        response => {
          this.authService.login(user.email, user.password)
            .subscribe(
              token => this.router.navigate(['/dashboard']),
              err => this.errorMsg = err.status
            );
        },
        err => {
          // Handle tmp case when user already exists and we don't have error msg on API side yet.
          if (err.status === 409) {
            this.errorMsg = 'User with this email already exists.';
          } else {
            this.errorMsg = err.status;
          }
        }
      );
    }
  }
}
