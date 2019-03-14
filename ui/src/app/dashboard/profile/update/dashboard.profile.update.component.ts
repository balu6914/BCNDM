import { Component, Input, OnInit } from '@angular/core';
import { FormGroup, FormControl, ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/throw';

import { UserService } from 'app/common/services/user.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { User } from 'app/common/interfaces/user.interface';

@Component({
  selector: 'dpc-user-profile-update',
  templateUrl: 'dashboard.profile.update.component.html',
})
export class DashboardProfileUpdateComponent implements OnInit {
  @Input() user: User;
  public form: FormGroup;
  submitted = false;

  constructor(
    private fb: FormBuilder,
    private router: Router,
    private userService: UserService,
    private alertService: AlertService,
  ) { }

  ngOnInit() {
    this.form = this.fb.group({
      contact_email: ['', [Validators.required, Validators.email]],
      first_name: ['', [Validators.required, Validators.maxLength(32)]],
      last_name: ['', [Validators.required, Validators.maxLength(32)]],
      phone: ['', [Validators.required, Validators.maxLength(32)]],
      address: ['', [Validators.required, Validators.maxLength(128)]],
      company: ['', [Validators.required, Validators.maxLength(32)]]
    });
    this.form.setValue({
      contact_email: this.user.contact_email || '',
      first_name: this.user.first_name || '',
      last_name: this.user.last_name || '',
      phone: this.user.phone || '',
      address: this.user.address || '',
      company: this.user.company || ''
    });
  }

  onSubmit() {
    this.submitted = true;

    if (this.form.valid) {
      const user = {
        contact_email: this.form.value.contact_email,
        first_name: this.form.value.first_name,
        last_name: this.form.value.last_name,
        phone: this.form.value.phone,
        address: this.form.value.address,
        company: this.form.value.company
      };
      this.userService.updateUser(user).subscribe(
        response => {
          this.alertService.success('You profile is succesfully updated.');
        },
        err => {
            // Handle tmp case when user already exists and we don't have error msg on API side yet.
            this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
        });
    }
  }

}
