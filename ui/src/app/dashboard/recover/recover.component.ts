import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormControl, ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import { ValidationErrors, ValidatorFn, AbstractControl } from '@angular/forms';
import {ActivatedRoute, Router} from '@angular/router';
import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/throw';

import { User } from 'app/common/interfaces/user.interface';
import { UserService } from 'app/common/services/user.service';
import { AuthService } from 'app/auth/services/auth.service';
import { CustomValidators } from 'app/common/validators/customvalidators';

@Component({
  selector: 'dpc-recovery-form',
  templateUrl: './recover.component.html',
  styleUrls: ['./recover.component.scss'],
  providers: [
  ],
})
export class RecoverComponent implements OnInit {
  public recoveryForm: FormGroup;
  public passwordForm: FormGroup;
  public errorMsg: String;
  emailSubmitted = false;
  passwordSubmitted = false;
  tokenSent = false;
  tokenReceived = false;
  token: string;
  id: string;
  tokenChecking = false;
  tokenValid = false;
  tokenChecked = false;
  passwordChanged = false;

  @Output() userCreated: EventEmitter<any> = new EventEmitter();
  constructor(
    private fb: FormBuilder,
    private router: Router,
    private userService: UserService,
    private authService: AuthService,
    private activatedRoute: ActivatedRoute
  ) {
  }

  ngOnInit() {
    this.errorMsg = null;
    this.recoveryForm = this.fb.group({
      email:      ['', [Validators.required, Validators.email, Validators.maxLength(32)]],
    });
    this.passwordForm = this.fb.group({
      password:      ['', [Validators.required, Validators.minLength(8), Validators.maxLength(32)]],
    });
    this.activatedRoute.queryParams.subscribe(params => {
      this.token = params['token'];
      this.id = params['id'];
      if (this.token && this.token.length && this.id && this.id.length) {
        this.tokenReceived = true;
        this.validateRecoveryToken();
      }
    });
  }

   validateRecoveryToken () {
     this.tokenChecking = true;
     this.userService.validateRecoveryToken(this.token , this.id).subscribe(
       response => {
         console.log(response);
         this.tokenValid = true;
         this.tokenChecking = false;
         this.tokenChecked = true;
       },
       err => {
         this.errorMsg = 'Invalid token!';
         this.tokenChecking = false;
         this.tokenChecked = true;
         console.log(err);
       }
     );
  }

  onRecoverySubmit() {
    this.emailSubmitted = true;

    if (this.recoveryForm.valid) {
      const email = this.recoveryForm.value.email;

      this.userService.sendRecoveryToken(email).subscribe(
        response => {
          console.log(response);
          this.tokenSent = true;
        },
        err => {
          console.log(err);
        }
      );
    }
  }

  onPasswordSubmit() {
    this.passwordSubmitted = true;

    if (this.passwordForm.valid) {
      const password = this.passwordForm.value.password;

      this.userService.setNewPassword(password, this.id, this.token).subscribe(
        response => {
          console.log(response);
          this.passwordChanged = true;
        },
        err => {
          console.log(err);
        }
      );
    }
  }

  resetPage() {
    this.errorMsg = '';
    this.emailSubmitted = false;
    this.passwordSubmitted = false;
    this.tokenSent = false;
    this.tokenReceived = false;
    this.token = '';
    this.tokenChecking = false;
    this.tokenValid = false;
    this.tokenChecked = false;
  }
}
