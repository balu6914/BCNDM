import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { FormArray, FormGroup, FormBuilder, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap';

import { AccessService } from 'app/dashboard/access/access.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { AuthService } from 'app/auth/services/auth.service';
import { UserService } from 'app/common/services/user.service';
import { User } from 'app/common/interfaces/user.interface';

@Component({
  selector: 'dpc-dashboard-access-add',
  templateUrl: './dashboard.access.add.component.html',
  styleUrls: [ './dashboard.access.add.component.scss' ]
})
export class DashboardAccessAddComponent implements OnInit {
  user = <User>{};
  form: FormGroup;
  submitted = false;
  users: User[] = [];

  @Output() accessCreated: EventEmitter<any> = new EventEmitter();
  constructor(
    private accessService: AccessService,
    private authService: AuthService,
    private userService: UserService,
    private formBuilder: FormBuilder,
    public modalNewAccess: BsModalRef,
    public alertService: AlertService,
  ) {

    this.form = this.formBuilder.group({
      receiver: ['', [Validators.required]],
    }, {
      validator: [this.userValidator.bind(this)]
    });

  }

  userValidator(fg: FormGroup) {
    if (fg.value.receiver.id === this.user.id) {
      // Create a custom error field used as *ngIf condition for style
      fg.controls.receiver.setErrors({
        'ownerID': true
      });
    }
  }

  ngOnInit () {
    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      data => {
        this.user = data;
        // Fetch all registered users
        this.userService.getAllUsers().subscribe(
          (result: any) => {
            this.users = result.users;
          },
          err => {
            this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
          }
        );
      },
      err => {
        console.log(err);
      }
    );
  }

  onSubmit() {
    this.submitted = true;

    if (this.form.valid) {
      const requestAccessReq = {
        receiver: this.form.value.receiver.id,
      };

      // Send Accesss Request
      this.accessService.requestAccess(requestAccessReq).subscribe(
        resp => {
          this.accessCreated.emit(this.form.value);
        },
        err => {
          this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
        }
      );

      this.modalNewAccess.hide();
    }
  }
}
