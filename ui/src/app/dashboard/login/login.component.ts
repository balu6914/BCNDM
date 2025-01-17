import { Component, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, Validators  } from '@angular/forms';
import {Router} from '@angular/router';

import { User } from 'app/common/interfaces/user.interface';
import { AuthService } from 'app/auth/services/auth.service';
import { TermsComponent } from './terms/terms.component';
import { BsModalService } from 'ngx-bootstrap/modal';

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
		private modalService: BsModalService,
		private authService: AuthService
  ) {}

  ngOnInit() {
    this.errorMsg = null;
    this.form = this.formBuilder.group({
      email:    ['', [<any>Validators.required, Validators.email]],
      password: ['', [<any>Validators.required, <any>Validators.minLength(5)]],
      acceptterms: ['', [<any>Validators.required]]
    });
  }

  onSubmit(user: User, isValid: boolean) {
    this.errorMsg = null;
		console.log(this.form.value.acceptterms)
		if (this.form.value.acceptterms !== true) {
			this.errorMsg = "Accept terms & conditions"
			return false
		}
    if (isValid) {
      this.authService.login(user.email, user.password).subscribe(
        (token: any) => {
          this.router.navigate(['/dashboard']);
        },
        err => {
          if ( err.status === 423) {
            this.errorMsg = 'User Locked Please contact administrator.';
          } else {
            this.errorMsg = 'Invalid Credentials.';
          }
        }
      );
    }
  }
	
	openModalAdd(e) {
		// Show DashboardSellAddComponent as Modal
		this.modalService.show(TermsComponent, { });
		e.preventDefault()
	}
}
