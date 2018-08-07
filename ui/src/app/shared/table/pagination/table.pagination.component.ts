import { Component, EventEmitter, Input, OnChanges, Output, SimpleChanges } from '@angular/core';
import { provideRouterInitializer } from '@angular/router/src/router_module';
import { ThrowStmt } from '@angular/compiler';

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
    let first = 0;
    let last = 0;

    last = this.page + this.size;
    if (this.page > this.size) {
      first = this.page - this.size;
      if (totalPages - this.page <= this.size) {
        first = first - (this.size - (totalPages - this.page) + 1);
      }
    } else {
      last = 2 * this.size;
    }

    first = first < 0 ? 0 : first;
    last = last < totalPages ? last : totalPages - 1;
    this.last = this.page === last;

    for (let i = first; i <= last; i++) {
      this.pages.push(i);
    }
  }

  pageChange(page: number) {
    if (this.page === page) {
      return;
    }

    this.page = page;
    this.pageChanged.emit(this.page);
  }

}
