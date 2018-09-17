import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl, ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/throw';

import { UserService } from 'app/common/services/user.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';

@Component({
  selector: 'dpc-user-profile-password-update',
  templateUrl: 'dashboard.profile.change.password.component.html',
})
export class DashboardProfilePasswordUpdateComponent implements OnInit {
  public form: FormGroup;
  submitted = false;

  constructor(
    private fb: FormBuilder,
    private userService: UserService,
    private alertService: AlertService,
  ) {}

  ngOnInit() {
    this.form = this.fb.group({
      old_password: ['', [Validators.required, Validators.minLength(8)]],
      new_password: ['', [Validators.required, Validators.minLength(8)]],
      re_password: ['', [Validators.required, Validators.minLength(8)]]
    },
    {
      validator: this.passwordMatchValidator
    });
  }

    passwordMatchValidator(fg: FormGroup) {
      // Compare passwords only if minLength is valid
      if (fg.get('re_password').value.length) {
        if (fg.get('new_password').value === fg.get('re_password').value) {
            return null;
        }
        fg.controls.re_password.setErrors({'invalid': true});
      }

      // Set form confirm password missmatch
      return {'missmatch': true };
    }

  onSubmit() {
    this.submitted = true;
    console.log("submit try", this.form)

    if (this.form.valid) {
      console.log("SENNNNNNNNNND", this.form);
      const data = this.form.value;
      this.userService.updatePassword(data).subscribe(
        response => {
          this.alertService.success('Yousuccesfully change your password.');
        },
        err => {
            this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
        });
    }
  }

}
