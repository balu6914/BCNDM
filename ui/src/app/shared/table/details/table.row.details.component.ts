import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { Stream } from 'app/common/interfaces/stream.interface';
import { DpcPipe } from 'app/common/pipes/converter.pipe';
import { TableType } from 'app/shared/table/table';
import { DashboardSellEditComponent } from 'app/dashboard/sell/edit';
import { DashboardSellDeleteComponent } from 'app/dashboard/sell/delete';
import { DashboardBuyAddComponent } from 'app/dashboard/buy/add';

@Component({
  selector: 'dpc-table-row-details',
  templateUrl: './table.row.details.component.html',
  styleUrls: ['./table.row.details.component.scss']
})

export class TableRowDetailsComponent implements OnInit {
  types = TableType;
  bsModalRef: BsModalRef;

  @Input() row: any;
  @Input() rowType: TableType;
  @Output() backClicked = new EventEmitter<String>();
  @Output() deleteEvt: EventEmitter<any> = new EventEmitter();
  @Output() editEvt: EventEmitter<any> = new EventEmitter();

  constructor(
    private modalService: BsModalService,
    private dpcPipe: DpcPipe,

  ) {}

  ngOnInit() {
  }

  close() {
    this.backClicked.emit('trigger');
  }

  openModalEdit() {
    // Parameters editData and streamID are used in DashboardSellEditComponent
    const initialState = {
      editData: {
        name:        this.row.name,
        type:        this.row.type,
        description: this.row.description,
        url:         this.row.url,
        price:       this.dpcPipe.transform(this.row.price),
        long:        this.row.location.coordinates[0],
        lat:         this.row.location.coordinates[1],
        snippet:     this.row.snippet,
        terms:       this.row.terms,
      },
      streamID: this.row.id,
    };
    // Open DashboardSellEditComponent as Modal
    this.bsModalRef = this.modalService.show(DashboardSellEditComponent, {initialState})
      .content.streamEdited.subscribe(
        stream => {
          this.editEvt.emit(stream);
        }
      );
  }

  openModalDelete() {
    // Parameter stream is used in DashboardSellDeleteComponent
    const initialState = {
      stream: {
        id:          this.row.id,
        name:        this.row.name,
        type:        this.row.type,
        description: this.row.description,
        price:       this.dpcPipe.transform(this.row.price),
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

  openModalSubscription() {
    // Parameter stream is set on modal component
    const initialState = {
      stream: {
        id:    this.row.id,
        name:  this.row.name,
        price: this.row.price,
      },
    };
    // Open DashboardBuyAddComponent Modal
    this.bsModalRef = this.modalService.show(DashboardBuyAddComponent, {initialState});
  }

}
