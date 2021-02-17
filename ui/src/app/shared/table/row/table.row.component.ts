import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { ngCopy } from 'angular-6-clipboard';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { AuthService } from 'app/auth/services/auth.service';
import { User } from 'app/common/interfaces/user.interface';

import { DashboardSellEditComponent } from 'app/dashboard/sell/edit/dashboard.sell.edit.component';
import { DashboardSellDeleteComponent } from 'app/dashboard/sell/delete/dashboard.sell.delete.component';
import { DashboardAiEditComponent } from 'app/dashboard/ai/edit/dashboard.ai.edit.component';
import { DashboardAiDeleteComponent } from 'app/dashboard/ai/delete/dashboard.ai.delete.component';
import { DashboardBuyAddComponent } from 'app/dashboard/buy/add/dashboard.buy.add.component';
import { DashboardContractsSignComponent } from 'app/dashboard/contracts/sign/dashboard.contracts.sign.component';
import { Stream, Subscription } from 'app/common/interfaces';
import { DpcPipe } from 'app/common/pipes/converter.pipe';
import { TableType } from 'app/shared/table/table';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { Access } from 'app/common/interfaces/access.interface';
import { Execution } from 'app/common/interfaces/execution.interface';
import { BalanceComponent } from 'app/shared/balance/balance.component';

@Component({
  selector: 'dpc-table-row',
  templateUrl: './table.row.component.html',
  styleUrls: ['./table.row.component.scss']
})

export class TableRowComponent implements OnInit {
  user: User;
  types = TableType;
  bsModalRef: BsModalRef;
  currentDate: string;
  execFetched = false;

  @Input() row: any;
  @Input() rowType: TableType;
  @Input() rowMarked: any;
  @Output() deleteEvt: EventEmitter<any> = new EventEmitter();
  @Output() editEvt: EventEmitter<any> = new EventEmitter();
  @Output() walletEvt: EventEmitter<any> = new EventEmitter();
  @Output() lockEvt: EventEmitter<any> = new EventEmitter();
  @Output() rowSelected = new EventEmitter<Stream | Subscription>();
  @Output() contractSigned: EventEmitter<any> = new EventEmitter();
  @Output() checkboxChangedEvt: EventEmitter<any> = new EventEmitter();
  @Output() accessApproved: EventEmitter<any> = new EventEmitter();
  @Output() accessRevoked: EventEmitter<any> = new EventEmitter();
  @Output() fetchExecResult: EventEmitter<any> = new EventEmitter();
  @Output() editEvent: EventEmitter<any> = new EventEmitter();

  constructor(
    private authService: AuthService,
    private modalService: BsModalService,
    private dpcPipe: DpcPipe,
    private alertService: AlertService
  ) { }

  private isStream(row: Stream | Subscription): row is Stream {
    return (<Stream>row).url !== undefined;
  }

  ngOnInit() {
    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      data => {
        this.user = data;
      },
      err => {
        console.log(err);
      }
    );

    const date = new Date();
    this.currentDate = date.toISOString();
  }

  public copyToClipboard() {
    if (this.row.url) {
      return ngCopy(this.row.url, null);
    }
    if (this.row.stream_url) {
      return ngCopy(this.row.stream_url, null);
    }
  }

  openModalEdit(row: any) {
    // Parameters editData and streamID are used in DashboardSellEditComponent
    const initialState = {
      editData: {
        visibility:  row.visibility,
        name:        row.name,
        type:        row.type,
        description: row.description,
        url:         row.url,
        price:       this.dpcPipe.transform(row.price),
        long:        row.location.coordinates[0],
        lat:         row.location.coordinates[1],
        snippet:     row.snippet,
        terms:       row.terms,
      },
      streamID: row.id,
      ownerID:  row.owner,
    };

    // Open DashboardSellEditComponent as Modal
    this.bsModalRef = this.modalService.show(DashboardSellEditComponent, {initialState})
      .content.streamEdited.subscribe(
        stream => {
          this.editEvt.emit(stream);
        }
      );
  }

  openModalEditAi(row: any) {
    // editData, streamID and ownerID are used in DashboardAiEditComponent
    const initialState = {
      editData: {
        visibility:  row.visibility,
        name:        row.name,
        type:        row.type,
        description: row.description,
        url:         row.url,
        price:       this.dpcPipe.transform(row.price),
        long:        row.location.coordinates[0],
        lat:         row.location.coordinates[1],
        snippet:     row.snippet,
        metadata:    JSON.stringify(row.metadata),
        terms:       row.terms,
      },
      streamID: row.id,
      ownerID:  row.owner,
    };

    // Open DashboardAiEditComponent as Modal
    this.bsModalRef = this.modalService.show(DashboardAiEditComponent, {initialState})
      .content.streamEdited.subscribe(
        stream => {
          this.editEvt.emit(stream);
        }
      );
  }

  openModalDelete(row: any) {
    // Parameter stream is used in DashboardSellDeleteComponent
    const initialState = {
      stream: {
        id:          row.id,
        name:        row.name,
        type:        row.type,
        description: row.description,
        price:       this.dpcPipe.transform(row.price),
      },
    };
    // Open DashboardSellDeleteComponent as Modal
    this.bsModalRef = this.modalService.show(DashboardSellDeleteComponent, {initialState})
      .content.streamDeleted.subscribe(
        id => {
          // Emit event to TableComponent
          this.deleteEvt.emit(id);
        }
      );
  }

  openModalDeleteAi(row: any) {
    // Parameter stream is used in DashboardSellDeleteComponent
    const initialState = {
      stream: {
        id:          row.id,
        name:        row.name,
        type:        row.type,
        description: row.description,
        price:       this.dpcPipe.transform(row.price),
      },
    };
    // Open DashboardSellDeleteComponent as Modal
    this.bsModalRef = this.modalService.show(DashboardAiDeleteComponent, {initialState})
      .content.streamDeleted.subscribe(
        id => {
          // Emit event to TableComponent
          this.deleteEvt.emit(id);
        }
    );
  }

  openModalSubscription(row: any) {
    if (row.external) {
      const gmailSuffix = '@gmail.com';
      if (!(this.user.email.toLowerCase().endsWith(gmailSuffix) ||
        this.user.contact_email.toLowerCase().endsWith(gmailSuffix))) {
          this.alertService.warning('Please use your Gmail account as a contact email address.');
          return;
      }
    }
    // Parameter stream is set on modal component
    const initialState = {
      stream: row,
    };
    // Open DashboardBuyAddComponent Modal
    this.bsModalRef = this.modalService.show(DashboardBuyAddComponent, {initialState});
  }

  openModalSignContract(row: any) {
    const initialState = {
      contract: row,
    };

    this.bsModalRef = this.modalService.show(DashboardContractsSignComponent, {initialState})
      .content.contractSigned.subscribe(
        contract => {
          this.contractSigned.emit(contract);
        }
      );
  }

  // Select/Click on Row emits a selectedRow event and pass selected row data
  // In order to show row details.
  selectRow(row: Stream) {
    this.rowSelected.emit(row);
  }

  onCheckboxChanged(row: Stream) {
    this.checkboxChangedEvt.emit(row);
  }

  approveAccess(row: Access) {
    this.accessApproved.emit(row);
  }

  revokeAccess(row: Access) {
    this.accessRevoked.emit(row);
  }

  onFetchExecResult(row: Execution) {
    this.execFetched = true;
    this.fetchExecResult.emit(row);
  }

  openModalEditUser(row: any) {
    this.editEvt.emit(row);
  }

  openModalWalletUser(row: any) {
    const initialState = {
      user: row,
    };

    this.bsModalRef = this.modalService.show(BalanceComponent, {initialState});
    this.bsModalRef.content.balanceUpdate.subscribe(
      balance => {
        this.row.balance = balance;
        this.bsModalRef.hide();
      }
    );
  }

  openModalDeleteUser(row: any) {
    this.deleteEvt.emit(row);
  }

  openModalLockedUser(row: any) {
    this.lockEvt.emit(row);
  }
}
