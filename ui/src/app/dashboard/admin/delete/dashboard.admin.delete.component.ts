import { Component, Output, EventEmitter } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap';

import { User } from 'app/common/interfaces/user.interface';
import { UserService } from 'app/common/services/user.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';

@Component({
  selector: 'dpc-dashboard-admin-delete',
  templateUrl: './dashboard.admin.delete.component.html',
  styleUrls: ['./dashboard.admin.delete.component.scss']
})
export class DashboardAdminDeleteComponent {
  user: User = {};

  @Output() userDeleted: EventEmitter<any> = new EventEmitter();
  constructor(
    private userService: UserService,
    public  modalDeleteUser: BsModalRef,
    public  alertService: AlertService,
  ) {}

  confirm(): void {
    const userUpdateReq: User = {
      id: this.user.id,
      disabled: !this.user.disabled,
    };

    this.userService.updateUser(userUpdateReq).subscribe(
      response => {
        this.userDeleted.emit(userUpdateReq.disabled);
        this.modalDeleteUser.hide();
        const msg =  this.user.disabled ? 'User Login successfully disabled.' : 'User Login successfully enabled.';
        this.alertService.success(msg);
      },
      err => {
        this.modalDeleteUser.hide();
        this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
      }
    );
  }
}
