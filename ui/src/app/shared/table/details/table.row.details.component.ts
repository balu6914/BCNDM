import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { Stream } from '../../../common/interfaces/stream.interface';
import { TasPipe } from '../../../common/pipes/converter.pipe';
import { TableType } from '../table';
import { DashboardSellEditComponent } from '../../../dashboard/sell/edit';
import { DashboardSellDeleteComponent } from '../../../dashboard/sell/delete';
import { DashboardBuyAddComponent } from '../../../dashboard/buy/add';

@Component({
  selector: 'dpc-table-row-details',
  templateUrl: './table.row.details.component.html',
  styleUrls: ['./table.row.details.component.scss']
})

export class TableRowDetailsComponent implements OnInit {
  types = TableType;
  bsModalRef: BsModalRef;

  @Input() stream: any;
  @Input() rowType: TableType;
  @Output() backClicked = new EventEmitter<String>();
  @Output() deleteEvt: EventEmitter<any> = new EventEmitter();
  @Output() editEvt: EventEmitter<any> = new EventEmitter();

  constructor(
    private modalService: BsModalService,
    private tasPipe: TasPipe,

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
        name:        this.stream.name,
        type:        this.stream.type,
        description: this.stream.description,
        url:         this.stream.url,
        price:       this.tasPipe.transform(this.stream.price),
        long:        this.stream.location.coordinates[0],
        lat:         this.stream.location.coordinates[1],
        snippet:     this.stream.snippet,
      },
      streamID: this.stream.id,
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
        id:          this.stream.id,
        name:        this.stream.name,
        type:        this.stream.type,
        description: this.stream.description,
        price:       this.tasPipe.transform(this.stream.price),
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
        id:    this.stream.id,
        name:  this.stream.name,
        price: this.stream.price,
      },
    };
    // Open DashboardSellAddComponent Modal
    this.bsModalRef = this.modalService.show(DashboardBuyAddComponent, {initialState});
  }

}
