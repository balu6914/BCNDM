import { Component, OnInit, Input } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';

import { UserService } from 'app/common/services/user.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { User } from 'app/common/interfaces/user.interface';
import { CustomValidators } from 'app/common/validators/customvalidators';

@Component({
  selector: 'dpc-user-profile-password-update',
  templateUrl: 'dashboard.profile.change.password.component.html',
})
export class DashboardProfilePasswordUpdateComponent implements OnInit {
  form: FormGroup;
  submitted = false;

  @Input() user: User;
  constructor(
    private formBuilder: FormBuilder,
    private userService: UserService,
    private alertService: AlertService,
  ) {}

  ngOnInit() {
    this.form = this.formBuilder.group({
      new_password: ['', [Validators.required, Validators.minLength(8),
        // 2. check whether the entered password has a number
        CustomValidators.patternValidator(/\d/, { hasNumber: true }),
        // 3. check whether the entered password has upper case letter
        CustomValidators.patternValidator(/[A-Z]/, { hasCapitalCase: true }),
        // 4. check whether the entered password has a lower-case letter
        CustomValidators.patternValidator(/[a-z]/, { hasSmallCase: true }),
        // 5. check whether the entered password has a special character
        CustomValidators.patternValidator(/[ !@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/, { hasSpecialCharacters: true })
      ]],
      re_password: ['', [Validators.required, Validators.minLength(8)]]
    },
    {
      validator: this.passwordMatchValidator
    });
  }

  passwordMatchValidator(formGroup: FormGroup) {
    // Compare passwords only if minLength is valid
    if (formGroup.get('re_password').value.length) {
      if (formGroup.get('new_password').value === formGroup.get('re_password').value) {
        return null;
      }
      formGroup.controls.re_password.setErrors({'invalid': true});
    }

    // Set form confirm password missmatch
    return {'missmatch': true };
  }

  onSubmit() {
    this.submitted = true;

    if (this.form.valid) {
      const updateUserReq = {
        id: this.user.id,
        // change new_password to password feild so API call take intor account the new password to update
        password: this.form.value.new_password,
        re_password: this.form.value.re_password,
      };
      this.userService.updateUser(updateUserReq).subscribe(
        response => {
          this.alertService.success('You succesfully change your password.');
          this.form.reset();
        },
        err => {
          this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
          this.submitted = false;
        });
    }
  }

}
