import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { ngCopy } from 'angular-6-clipboard';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { DashboardSellEditComponent } from '../../../dashboard/sell/edit';
import { DashboardSellDeleteComponent } from '../../../dashboard/sell/delete';
import { DashboardBuyAddComponent } from '../../../dashboard/buy/add';
import { Stream, Subscription } from '../../../common/interfaces';
import { TasPipe } from '../../../common/pipes/converter.pipe';
import { TableType } from '../table';

@Component({
  selector: 'dpc-table-row',
  templateUrl: './table.row.component.html',
  styleUrls: ['./table.row.component.scss']
})

export class TableRowComponent implements OnInit {
  types = TableType;
  bsModalRef: BsModalRef;

  @Input() row: any;
  @Input() rowType: TableType;
  @Output() deleteEvt: EventEmitter<any> = new EventEmitter();
  @Output() editEvt: EventEmitter<any> = new EventEmitter();
  @Output() rowSelected = new EventEmitter<Stream | Subscription>();

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
        name:        row.name,
        type:        row.type,
        description: row.description,
        url:         row.url,
        price:       this.tasPipe.transform(row.price),
        long:        row.location.coordinates[0],
        lat:         row.location.coordinates[1],
        snippet: row.snippet,
      },
      streamID: row.id,
    };
    // Open DashboardSellEditComponent as Modal
    this.bsModalRef = this.modalService.show(DashboardSellEditComponent, {initialState})
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
        price:       this.tasPipe.transform(row.price),
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

  openModalSubscription(row: any) {
    // Parameter stream is set on modal component
    const initialState = {
      stream: {
        id:         row.id,
        name:       row.name,
        price:      row.price,
      },
    };
    // Open DashboardSellAddComponent Modal
    this.bsModalRef = this.modalService.show(DashboardBuyAddComponent, {initialState});
  }

  // Select/Click on Row emits a selectedRow event and pass selected row data
  // In order to show row details.
  selectRow(row: Stream) {
    this.rowSelected.emit(row);
  }
}
