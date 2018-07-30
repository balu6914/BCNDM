import { Component, Input, OnInit } from '@angular/core';
import { ngCopy } from 'angular-6-clipboard';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { DashboardSellEditComponent } from '../../../dashboard/sell/edit';
import { DashboardSellDeleteComponent } from '../../../dashboard/sell/delete';
import { SubscriptionAddComponent } from '../../../dashboard/buy/add';
import { Stream, Subscription } from '../../../common/interfaces';
import { TasPipe } from '../../../common/pipes/converter.pipe';
import { TableType } from '../table';

@Component({
  selector: 'dpc-table-row',
  templateUrl: './table.row.component.html',
  styleUrls: ['./table.row.component.scss']
})

export class TableRowComponent implements OnInit {
  types = TableType
  bsModalRef: BsModalRef;

  @Input() row: Stream | Subscription;
  @Input() rowType: TableType
  constructor(
    private modalService: BsModalService,
    private tasPipe: TasPipe,
  ) { }

  private isStream(row: Stream | Subscription): row is Stream {
    return (<Stream>row).url !== undefined;
  }

  ngOnInit() {
  }

  public copyToClipboard() {
    if (this.isStream(this.row)) {
      ngCopy(this.row.url, null)
    }
  }

  openModal(row: any) {
    // Parameter formEdit is set on modal component
    const initialState = {
      editData: {
        name:        row.name,
        type:        row.type,
        description: row.description,
        url:         row.url,
        price:       this.tasPipe.transform(row.price),
        long:        row.location.coordinates[0],
        lat:         row.location.coordinates[1],
      },
      streamID: row.id,
    };
    // Open DashboardSellAddComponent Modal
    this.bsModalRef = this.modalService.show(DashboardSellEditComponent, {initialState});
  }

  openModalDelete(row: any) {
    // Parameter stream is set on modal component
    const initialState = {
      stream: {
        id:          row.id,
        name:        row.name,
        type:        row.type,
        description: row.description,
        price:       this.tasPipe.transform(row.price),
      },
    };
    // Open DashboardSellAddComponent Modal
    this.bsModalRef = this.modalService.show(DashboardSellDeleteComponent, {initialState});
  }
}
