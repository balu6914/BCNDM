import { Component, Output, EventEmitter } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap';

import { User } from 'app/common/interfaces/user.interface';
import { UserService } from 'app/common/services/user.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';

@Component({
  selector: 'dpc-dashboard-admin-lock',
  templateUrl: './dashboard.admin.lock.component.html',
  styleUrls: ['./dashboard.admin.lock.component.scss']
})
export class DashboardAdminLockComponent {
  user: User = {};

  @Output() userLocked: EventEmitter<any> = new EventEmitter();
  constructor(
    private userService: UserService,
    public  modalLockedUser: BsModalRef,
    public  alertService: AlertService,
  ) {}

  confirm(): void {
    const userUpdateReq: User = {
      id: this.user.id,
      locked: !this.user.locked,
    };

    this.userService.updateUser(userUpdateReq).subscribe(
      response => {
        this.userLocked.emit(userUpdateReq.locked);
        this.modalLockedUser.hide();
        const msg = userUpdateReq.locked ? 'User successfully Locked.' : 'User successfully unlocked.';
        this.alertService.success(msg);
      },
      err => {
        this.modalLockedUser.hide();
        this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
      }
    );
  }
}
