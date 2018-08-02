import { Component, EventEmitter, Input, OnChanges, Output, SimpleChanges } from '@angular/core';

@Component({
  selector: 'dpc-table-pagination',
  templateUrl: './table.pagination.component.html',
  styleUrls: ['./table.pagination.component.css']
})
export class TablePaginationComponent implements OnChanges {
  readonly size = 3;

  public last: boolean;
  public pages: number[];

  @Input() page: number;
  @Input() total: number;
  @Input() limit: number;

  @Output() pageChanged = new EventEmitter<number>();

  constructor() { }

  ngOnChanges(changes: SimpleChanges) {
    this.pages = [];
    const totalPages = Math.ceil(this.total / this.limit);
    this.last = this.page === totalPages;

    const start = this.size >= this.page ? 1 : this.page - this.size + 1;
    let end = start === 1 ? 2 * this.size : this.page + this.size;
    end = end < totalPages ? end + 1 : totalPages;

    for (let i = start; i <= end; i++) {
      this.pages.push(i);
    }
  }

  pageChange(page: number) {
    if (this.page === page - 1) {
      return;
    }

    this.page = page - 1;
    this.pageChanged.emit(this.page);
  }

}
