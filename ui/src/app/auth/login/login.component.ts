import { Component, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, Validators  } from '@angular/forms';
import {Router} from '@angular/router';

import { User } from 'app/common/interfaces/user.interface';
import { AuthService } from 'app/auth/services/auth.service';

@Component({
  selector: 'dpc-login-form',
  templateUrl: './login.component.html',
  styleUrls: [ './login.component.scss' ],
  providers: [
  ],
})

export class LoginComponent implements OnInit {
  public form: FormGroup;
  public errorMsg: String;

  constructor(
    private formBuilder: FormBuilder,
    private router: Router,
    private authService: AuthService
  ) {}

  ngOnInit() {
    this.errorMsg = null;
    this.form = this.formBuilder.group({
      email:    ['', [<any>Validators.required, Validators.email]],
      password: ['', [<any>Validators.required, <any>Validators.minLength(5)]]
    });
  }

  onSubmit(user: User, isValid: boolean) {
    this.errorMsg = null;
    if (isValid) {
      this.authService.login(user.email, user.password).subscribe(
        (token: any) => {
          this.router.navigate(['/dashboard']);
        },
        err => {
          this.errorMsg = 'Invalid Credentials.';
        }
      );
    }
  }
}
