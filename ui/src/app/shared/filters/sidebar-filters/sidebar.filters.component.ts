import { Component } from '@angular/core';

@Component({
  selector: 'dpc-sidebar-filters',
  templateUrl: './sidebar.filters.component.html',
  styleUrls: [ './sidebar.filters.component.scss' ]
})
export class SidebarFiltersComponent {
    _opened: boolean = false

    constructor(
    ) { }

    _toggleSidebar() {
      this._opened = !this._opened;
    }

}
