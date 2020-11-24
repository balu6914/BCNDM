import { Component, Output, EventEmitter, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap';

import { AlertService } from 'app/shared/alerts/services/alert.service';
import { User } from 'app/common/interfaces/user.interface';
import { UserService } from 'app/common/services/user.service';
import { CustomValidators } from 'app/common/validators/customvalidators';

@Component({
  selector: 'dpc-dashboard-admin-edit',
  templateUrl: './dashboard.admin.edit.component.html',
  styleUrls: ['./dashboard.admin.edit.component.scss']
})
export class DashboardAdminEditComponent implements OnInit {
  form: FormGroup;
  user: User = {};
  submitted = false;

  @Output() userEdited: EventEmitter<any> = new EventEmitter();
  constructor(
    private modalEditUser: BsModalRef,
    private alertService: AlertService,
    private formBuilder: FormBuilder,
    private userService: UserService,
  ) {
  }

  ngOnInit() {
    this.form = this.formBuilder.group({
      email:      [this.user.email, [Validators.required, Validators.email, Validators.maxLength(256)]],
      first_name: [this.user.first_name, [Validators.maxLength(32)]],
      last_name:  [this.user.last_name, [Validators.maxLength(32)]],
      company:    [this.user.company, [Validators.maxLength(32)]],
      address:    [this.user.address, [Validators.maxLength(128)]],
      phone:      [this.user.phone, [Validators.maxLength(32)]],
      password:   ['', [Validators.minLength(8), Validators.maxLength(32),
        // 2. check whether the entered password has a number
        CustomValidators.patternValidator(/\d/, { hasNumber: true }),
        // 3. check whether the entered password has upper case letter
        CustomValidators.patternValidator(/[A-Z]/, { hasCapitalCase: true }),
        // 4. check whether the entered password has a lower-case letter
        CustomValidators.patternValidator(/[a-z]/, { hasSmallCase: true }),
        // 5. check whether the entered password has a special character
        CustomValidators.patternValidator(/[ !@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/, { hasSpecialCharacters: true })
      ]]
    });
  }

  onSubmit() {
    this.submitted = true;

    if (this.form.valid) {
        const user: User = {
          id: this.user.id,
          email: this.form.value.email,
          first_name: this.form.value.first_name,
          last_name: this.form.value.last_name,
          company: this.form.value.company,
          address: this.form.value.address,
          phone: this.form.value.phone,
          password: this.form.value.password,
        };

      this.userService.updateUser(user).subscribe(
        response => {
          this.userEdited.emit(user);
          this.modalEditUser.hide();
          this.alertService.success(`User successfully edited.`);
        },
        err => {
          this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
        }
      );
    }
  }
}
