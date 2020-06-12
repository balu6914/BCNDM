import { Component, Output, EventEmitter } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap';
import { UserService } from 'app/common/services/user.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';

@Component({
  selector: 'dpc-dashboard-admin-delete',
  templateUrl: './dashboard.admin.delete.component.html',
  styleUrls: ['./dashboard.admin.delete.component.scss']
})
export class DashboardAdminDeleteComponent {
  user: any;

  @Output() userDeleted: EventEmitter<any> = new EventEmitter();
  constructor(
    private userService: UserService,
    public  modalDeleteUser: BsModalRef,
    public  alertService: AlertService,
  ) {}


  confirm(): void {
    /*this.userService.removeUser(this.user.id).subscribe(
      res => {
        this.userDeleted.emit(this.user.id);
        this.modalDeleteUser.hide();
        this.alertService.success(`User succesfully removed!`);
      },
      err => {
        this.modalDeleteUser.hide();
        this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
    });*/
  }

}
