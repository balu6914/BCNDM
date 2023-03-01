import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { BsModalRef } from 'ngx-bootstrap';

@Component({
	selector: 'dpc-terms',
	templateUrl: './terms.component.html',
  styleUrls: []
})
export class TermsComponent implements OnInit {
  constructor(
    public modalRefStream: BsModalRef
  ) { }

  ngOnInit() {
  }

	closeModal() {
		this.modalRefStream.hide();
	}
}
